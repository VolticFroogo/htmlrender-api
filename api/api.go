package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init() {
	r := gin.Default()
	r.Use(recoveryMiddleware)

	r.Handle(http.MethodGet, "/render", render)

	err := r.Run()
	if err != nil {
		panic(err)
	}
}

func recoveryMiddleware(c *gin.Context) {
	// defer sentry.Recover()
	c.Next()
}
