package main

import (
	"context"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/phinexdaz/ipapk-server/conf"
	"github.com/phinexdaz/ipapk-server/middleware"
	"github.com/phinexdaz/ipapk-server/models"
	"github.com/phinexdaz/ipapk-server/templates"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	models.InitDB()
	conf.Init("config.json")

	router := gin.Default()
	router.Use(location.New(location.Config{
		Scheme: "http",
		Base:   "/ipapk",
	}))
	router.SetFuncMap(templates.TplFuncMap)
	router.LoadHTMLGlob("public/views/*")

	v1 := router.Group("/ipapk")
	v1.Static("icon", "data/icon")
	v1.Static("static", "public/static")

	v1.POST("/upload", middleware.Upload)
	v1.GET("/qrcode/:uuid", middleware.QRCode)
	v1.GET("/plist/:uuid", middleware.Plist)
	v1.GET("/ipa/:uuid", middleware.DownloadIPA)
	v1.GET("/apk/:uuid", middleware.DownloadAPK)
	v1.GET("/bundles/:uuid", middleware.GetBundle)
	v1.GET("/bundles/:uuid/versions", middleware.GetVersions)
	v1.GET("/bundles/:uuid/versions/:ver", middleware.GetBuilds)

	srv := &http.Server{
		Addr:    conf.AppConfig.Addr(),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %v\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
