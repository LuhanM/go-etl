package router

import (
	"log"
	"net/http"

	"github.com/LuhanM/go-etl/handler"
	"github.com/gin-gonic/gin"
)

func Initialize() {
	router := gin.Default()
	initializeRoutes(router)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initializeRoutes(router *gin.Engine) {
	basePath := "/api/v1"
	v1 := router.Group(basePath)
	{
		v1.POST("/file", handler.ImportFile)
	}
}
