package controllers

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilyasbabu/go-url-shortner/models"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDB(d *gorm.DB) {
	db = d
}

func genrateSlug() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func isValidUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func Home(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func updateExpiredURLs() {
	dateBeforeSevenDays := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	db.Model(&models.URLs{}).Where("created_at < ? AND active = ?", dateBeforeSevenDays, true).Update("active", false)
}

func ShortenURL(c *gin.Context) {
	urlText := c.PostForm("value")
	if !isValidUrl(urlText) {
		c.JSON(400, gin.H{
			"data": "Invalid URL",
		})
		return
	}
	dateBeforeSevenDays := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	var slugs []string
	db.Model(&models.URLs{}).Pluck("slug", &slugs).Where("created_at > ?", dateBeforeSevenDays)
	slug := genrateSlug()
	for contains(slugs, slug) {
		slug = genrateSlug()
	}
	go updateExpiredURLs()
	db.Create(&models.URLs{Slug: slug, Url: urlText, Active: true})
	c.JSON(200, gin.H{
		"data": slug,
	})
}

func ResolveURL(c *gin.Context) {
	slug := c.Param("slug")
	var url string
	dateBeforeSevenDays := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	db.Model(&models.URLs{}).Where("slug = ? AND created_at > ?", slug, dateBeforeSevenDays).Pluck("url", &url)
	go updateExpiredURLs()
	if url == "" {
		c.JSON(400, gin.H{
			"msg": "Invalid URL",
		})
		return
	} else {
		c.Redirect(301, url)
	}
}
