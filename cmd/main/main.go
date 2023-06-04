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
	router.GET("/landlords/search", controllers.SearchLandlords)
	router.GET("/properties/search", controllers.SearchProperties)

	//router.GET("/albums", getAlbums)
	//router.GET("/albums/:id", getAlbumByID)
	//router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}
