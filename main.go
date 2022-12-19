package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"usercenter/config/config"
	"usercenter/config/database"
)

func main() {
	config.InitConfig()
	database.Init()
	r := gin.Default()
	r.Use(cors.Default())
	err := r.Run(":" + config.Config.Server.Port)
	if err != nil {
		log.Fatal("ServerStartFailed", err)
	}
}
