package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilyasbabu/go-url-shortner/models"
	"github.com/ilyasbabu/go-url-shortner/utils"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDB(d *gorm.DB) {
	db = d
}

func Home(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func updateExpiredURLs() {
	dateBeforeSevenDays := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	db.Model(&models.URLs{}).Where("created_at < ? AND active = ?", dateBeforeSevenDays, true).Update("active", false)
}

func ShortenURL(c *gin.Context) {
	urlText := c.PostForm("value")
	if !utils.IsValidUrl(urlText) {
		c.JSON(400, gin.H{
			"data": "Invalid URL",
		})
		return
	}
	dateBeforeSevenDays := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	var slugs []string
	db.Model(&models.URLs{}).Pluck("slug", &slugs).Where("created_at > ?", dateBeforeSevenDays)
	slug := utils.GenrateSlug()
	for utils.Contains(slugs, slug) {
		slug = utils.GenrateSlug()
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
