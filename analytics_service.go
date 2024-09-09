package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func setupAnalyticsRoutes(r *gin.Engine) {
	analytics := r.Group("/analytics")
	analytics.Use(AIMiddleware())
	{
		analytics.GET("/customer-insights", getCustomerInsights)
		analytics.GET("/sales-forecast", getSalesForecast)
		analytics.GET("/churn-prediction", getChurnPrediction)
	}
}

func getCustomerInsights(c *gin.Context) {
	customerID := c.Query("customer_id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_id is required"})
		return
	}

	customerData := getCustomerData(customerID)
	churnProbability := predictChurn(customerData)
	lifetimeValue := calculateLifetimeValue(customerID)
	recentInteractions := getRecentInteractions(customerID)
	sentimentScore := analyzeSentimentBatch(recentInteractions)
	recommendations := recommendationEngine.GetRecommendations(customerID, 3)

	insights := []string{
		fmt.Sprintf("Cliente tem %.2f%% de probabilidade de churn", churnProbability*100),
		fmt.Sprintf("O valor vitalício estimado do cliente é R$ %.2f", lifetimeValue),
		fmt.Sprintf("O sentimento médio das interações recentes é %.2f", sentimentScore),
	}

	aiSuggestion, _ := c.Get("ai_suggestion")
	c.JSON(http.StatusOK, gin.H{
		"customer_id":        customerID,
		"insights":           insights,
		"churn_probability":  churnProbability,
		"lifetime_value":     lifetimeValue,
		"average_sentiment":  sentimentScore,
		"recommendations":    recommendations,
		"ai_suggestion":      aiSuggestion,
		"recent_interactions": len(recentInteractions),
	})
}

func getSalesForecast(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Previsão de vendas para os próximos 3 meses",
		"forecast": []gin.H{
			{"month": "Junho", "predicted_sales": 50000},
			{"month": "Julho", "predicted_sales": 55000},
			{"month": "Agosto", "predicted_sales": 60000},
		},
	})
}

func getChurnPrediction(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Previsão de churn de clientes",
		"high_risk_customers": []int{3, 7, 12},
		"churn_probability": 0.15,
	})
}

