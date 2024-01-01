package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Platform   string `json:"platform"`
	Address    string `json:"address"`
	ReleaseID  uint `gorm:"foreignKey:ID"`
}

type Release struct {
	gorm.Model
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	Artist 	     string `json:"artist"`
	PicturePath  string `json:"picturePath"`
	DownloadPath string `json:"downloadPath"`
	Links        []Link `json:"links"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("m1841-music.db"), &gorm.Config{})

	if err != nil {
		panic("[Server]: Failed to connect to the database")
	}

	db.AutoMigrate(&Release{}, &Link{})

	router := gin.Default()
	
	router.GET("/releases", func (c *gin.Context) {
		var releases []Release
		result := db.Preload("Links").Find(&releases)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch entries"})
			return
		}

		c.IndentedJSON(http.StatusOK, releases)
	})

	router.GET("/releases/:id", func (c *gin.Context) {
		id := c.Param("id")
		var release Release
		result := db.Preload("Links").First(&release, id)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch entries"})
			return
		}

		c.IndentedJSON(http.StatusOK, release)
	})

	router.POST("/releases", func (c *gin.Context) {
		var newRelease Release
		err := c.BindJSON(&newRelease)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bind JSON"})
			return
		}

		result := db.Create(&newRelease)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new entry"})
			return
		}

		c.IndentedJSON(http.StatusCreated, newRelease)
	})

	router.Run(":8080")
}