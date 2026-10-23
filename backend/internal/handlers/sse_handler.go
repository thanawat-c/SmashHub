package handlers

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"smashhub/internal/services"
)

type SSEHandler struct{}

func NewSSEHandler() *SSEHandler {
	return &SSEHandler{}
}

// StreamEvents handles GET /api/v1/events/stream
func (h *SSEHandler) StreamEvents(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	clientChan := make(chan string, 10)
	services.Hub.Register(clientChan)
	defer services.Hub.Unregister(clientChan)

	// Send initial greeting
	c.SSEvent("CONNECTED", gin.H{
		"status": "connected",
		"time":   time.Now().Format(time.RFC3339),
	})
	c.Writer.Flush()

	// Keep-alive ticker every 25 seconds
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			return false
		case msg, ok := <-clientChan:
			if !ok {
				return false
			}
			w.Write([]byte(msg))
			c.Writer.Flush()
			return true
		case <-ticker.C:
			// Ping comment to keep connection alive
			w.Write([]byte(": ping\n\n"))
			c.Writer.Flush()
			return true
		}
	})
}
