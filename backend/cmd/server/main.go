package main

import (
	"log"
	"net/http"
	"novel-to-screenplay-ai/internal/ai"
	"novel-to-screenplay-ai/internal/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	ai.LoadEnv()
	aiStatus := ai.RuntimeStatusFromEnv()

	r := gin.Default()
	r.Use(withCORS())
	r.GET("/api/health", handlers.Health)
	r.POST("/api/generate", handlers.Generate)
	addr := ":8080"

	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("AI provider: %s", aiStatus.AIProvider)
	log.Printf("AI model: %s", aiStatus.AIModel)
	log.Printf("AI base URL configured: %t", aiStatus.AIBaseURLConfigured)
	log.Printf("AI API key configured: %t", aiStatus.AIAPIKeyConfigured)
	log.Printf("AI timeout seconds: %d", aiStatus.AITimeoutSeconds)
	log.Printf("novel-to-screenplay-ai backend listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func withCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}

}
