package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Platform   string
	Address    string
	ReleaseID  uint
}

type Release struct {
	gorm.Model
	Title        string
	Subtitle     string
	Artist 	     string
	PicturePath  string
	DownloadPath string
	Type         string
	Links        []Link
}

func main() {
	// err := godotenv.Load()

	// if err != nil {
	// 	panic("Failed to load environment variables")
	// }

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	url := fmt.Sprintf("host=%s  user=%s  password=%s  dbname=%s  port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to the database")
	}

	db.AutoMigrate(&Release{}, &Link{})

	router := gin.Default()
	
	router.GET("/releases", func (c *gin.Context) {
		var releases []Release
		result := db.Preload("Links").Order("created_at desc").Find(&releases)

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

	router.GET("/links", func (c *gin.Context) {
		var links []Link
		result := db.Find(&links)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch entries"})
			return
		}

		c.IndentedJSON(http.StatusOK, links)
	})

	router.GET("/links/:id", func (c *gin.Context) {
		id := c.Param("id")
		var link Link
		result := db.First(&link, id)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch entries"})
			return
		}

		c.IndentedJSON(http.StatusOK, link)
	})

	router.Run(":8080")
}
