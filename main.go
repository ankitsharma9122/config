package main

import (
	"ankit/project/config"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var configuration config.Config

func main() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&configuration); err != nil {
		panic(err)
	}
	router := gin.Default()
	router.GET("/api/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, configuration)
	})
	router.POST("/api/config", func(c *gin.Context) {
		// Check if the uploaded file is in the request
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Open the configuration file for writing
		outputFile, err := os.Create("config.yaml")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer outputFile.Close()

		// Open the uploaded file for reading
		uploadedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer uploadedFile.Close()

		// Copy the uploaded file content to the configuration file
		_, err = io.Copy(outputFile, uploadedFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Re-read the configuration from the updated file
		if err := viper.ReadInConfig(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Unmarshal the updated configuration
		if err := viper.Unmarshal(&configuration); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully"})
	})

	port := configuration.Server.Port
	router.Run(":" + strconv.Itoa(port))
}
