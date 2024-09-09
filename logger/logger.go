package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger *logrus.Logger

func setupLogger() {
	logger = logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		
		aiProcessed, _ := c.Get("ai_processed")
		aiSuggestion, _ := c.Get("ai_suggestion")

		logger.WithFields(logrus.Fields{
			"status":        c.Writer.Status(),
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"ip":            c.ClientIP(),
			"user_agent":    c.Request.UserAgent(),
			"latency":       latencyTime,
			"ai_processed":  aiProcessed,
			"ai_suggestion": aiSuggestion,
		}).Info("Request processed")
	}
}
