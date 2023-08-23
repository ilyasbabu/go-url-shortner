package main

import (
	"log"

	"github.com/gin-gonic/gin"
	c "github.com/ilyasbabu/go-url-shortner/controllers"
	m "github.com/ilyasbabu/go-url-shortner/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&m.URLs{})
	c.SetDB(db)
}

func main() {
	r := gin.Default()
	r.POST("/shorten", c.ShortenURL)
	r.GET("/:slug", c.ResolveURL)
	r.Run(":8080")
}
