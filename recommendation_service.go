package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func setupRecommendationRoutes(r *gin.Engine) {
	recommendations := r.Group("/recommendations")
	recommendations.Use(AIMiddleware())
	{
		recommendations.GET("/products/:customer_id", getProductRecommendations)
		recommendations.GET("/cross-sell/:product_id", getCrossSellRecommendations)
	}
}

func setupRecommendationEngineRoutes(r *gin.Engine) {
	rec := r.Group("/recommendations")
	rec.Use(AuthMiddleware())
	{
		rec.POST("/rate", addRating)
		rec.GET("/user/:user", getUserRecommendations)
		rec.GET("/customer/:id", getCustomerRecommendations)
	}
}

func getProductRecommendations(c *gin.Context) {
	customerID := c.Param("customer_id")
	
	customerData := getCustomerData(customerID)
	recentPurchases := getRecentPurchases(customerID)
	browsingHistory := getBrowsingHistory(customerID)
	
	recommendations := recommendationEngine.GetRecommendations(customerID, 5)
	
	// Analisa o sentimento das últimas interações do cliente
	recentInteractions := getRecentInteractions(customerID)
	sentimentScore := analyzeSentimentBatch(recentInteractions)
	
	aiSuggestion, _ := c.Get("ai_suggestion")
	c.JSON(http.StatusOK, gin.H{
		"customer_id": customerID,
		"recommendations": recommendations,
		"context": gin.H{
			"recent_purchases": recentPurchases,
			"browsing_history": browsingHistory,
			"sentiment_score": sentimentScore,
		},
		"ai_suggestion": aiSuggestion,
	})
}

func getCrossSellRecommendations(c *gin.Context) {
	productID := c.Param("product_id")
	
	c.JSON(http.StatusOK, gin.H{
		"product_id": productID,
		"cross_sell_products": []string{"Produto X", "Produto Y", "Produto Z"},
	})
}

func getCustomerRecommendations(c *gin.Context) {
	customerID := c.Param("id")
	recommendations := recommendationEngine.GetRecommendations(customerID, 5)
	c.JSON(http.StatusOK, gin.H{"recommendations": recommendations})
}

func getRecentPurchases(customerID string) []string {
	// Implementação simulada
	return []string{"Produto A", "Produto B"}
}

func getBrowsingHistory(customerID string) []string {
	// Implementação simulada
	return []string{"Categoria X", "Produto Y", "Categoria Z"}
}
