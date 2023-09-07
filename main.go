package main

import (
	"mongo-crud/configs"
	"mongo-crud/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	configs.ConnectDB()

}

func main() {
	router := gin.Default()

	routes.ContactRoute(router);
	router.Run("localhost:8000")
}