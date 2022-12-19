package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"usercenter/app/midwares"
	"usercenter/config/config"
	"usercenter/config/database"
	"usercenter/config/router"
)

func main() {
	config.InitConfig()
	database.Init()
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(midwares.ErrHandler())
	r.NoMethod(midwares.HandleNotFound)
	r.NoRoute(midwares.HandleNotFound)
	router.Init(r)
	err := r.Run(":" + config.Config.GetString("server.port"))
	if err != nil {
		log.Fatal("ServerStartFailed", err)
	}
}
