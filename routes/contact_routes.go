package routes

import (
	"mongo-crud/controllers"

	"github.com/gin-gonic/gin"
)

func ContactRoute(router *gin.Engine)  {
	router.GET("/", controllers.Index())
	router.GET("/contacts/:contactId", controllers.Show())
	router.POST("/contacts", controllers.Store())
	router.PUT("/contacts/:contactId", controllers.Update())

}