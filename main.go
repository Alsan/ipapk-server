package main

import (
	"context"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
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
	v1.GET("/bundles/:uuid", middleware.GetBundle)
	v1.GET("/qrcode/:uuid", middleware.QRCode)
	v1.GET("/detail/:uuid", middleware.Detail)
	v1.GET("/plist/:uuid", middleware.Plist)
	v1.GET("/ipa/:uuid", middleware.DownloadIPA)
	v1.GET("/apk/:uuid", middleware.DownloadAPK)

	srv := &http.Server{
		Addr:    "127.0.0.1:8090",
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
