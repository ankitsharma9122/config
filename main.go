package main

import (
	"ankit/project/config"
	"net/http"
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
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := c.SaveUploadedFile(file, "config.yaml"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := viper.ReadInConfig(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := viper.Unmarshal(&configuration); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully"})
	})
	port := configuration.Server.Port
	router.Run(":" + strconv.Itoa(port))
}
