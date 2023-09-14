package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusHandler struct{}

func AttachStatusHandler(g *gin.RouterGroup) {
	h := &StatusHandler{}
	g.GET("/ping", h.Status)
}

func (h *StatusHandler) Status(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
