package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"snoop-server/pkg/controllers"
	"snoop-server/pkg/database"
)

func main() {
	ctx := context.Background()
	defer database.UnsafeClose(ctx)

	router := gin.Default()

	router.Use(database.ErrorHandler)

	// Landlord APIs
	router.GET("/landlords", controllers.SearchLandlords)
	router.GET("/landlords/:id", controllers.GetLandlord)
	router.POST("/landlords", controllers.AddLandlord)

	// Properties APIs
	router.GET("/properties", controllers.SearchProperties)
	router.GET("/properties/:id", controllers.GetProperty)

	//router.GET("/albums", getAlbums)
	//router.GET("/albums/:id", getAlbumByID)
	//router.POST("/albums", postAlbums)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
