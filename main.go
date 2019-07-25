package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var router *gin.Engine

func main() {
	router = gin.Default()
	initDB()
	initializeRoutes()
	router.Run()
}

func initializeRoutes() {
	router.POST("/orders", createOrder)
}
