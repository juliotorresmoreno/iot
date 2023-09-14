package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/iot/etl/config"
	"github.com/juliotorresmoreno/iot/etl/handlers"
	"github.com/juliotorresmoreno/iot/etl/tasks"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	handlers.AttachStatusHandler(r.Group("/status"))
	handlers.AttachJobsHandler(r.Group("/jobs"))

	r.SetTrustedProxies([]string{})

	return r
}

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	switch conf.Env {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "testing":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	tasks.DefaultTaskManager.Subscribe()

	r := setupRouter()
	r.Run(conf.Addr)
}
