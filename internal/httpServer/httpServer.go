package httpServer

import (
	"fmt"
	"net/http"
	"time"

	"gateway.bojiu.com/internal/httpServer/routers"
	"gateway.bojiu.com/pkg/viper"
	"github.com/gin-gonic/gin"
)

func Run() {
	mode := viper.Vp.GetString("active")
	if mode == "pro" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", viper.Vp.GetInt("ser.gateway.httpPort")),
		Handler:        r,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
