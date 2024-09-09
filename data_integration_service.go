package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"sync"
)

var (
	customer360Cache map[string]gin.H
	customer360CacheMutex sync.RWMutex
	customer360CacheExpiration map[string]time.Time
)

func init() {
	customer360Cache = make(map[string]gin.H)
	customer360CacheExpiration = make(map[string]time.Time)
}

func getCachedCustomer360View(customerID string) gin.H {
	customer360CacheMutex.RLock()
	defer customer360CacheMutex.RUnlock()

	if view, ok := customer360Cache[customerID]; ok {
		if time.Now().Before(customer360CacheExpiration[customerID]) {
			return view
		}
	}

	return nil
}

func setCachedCustomer360View(customerID string, view gin.H) {
	customer360CacheMutex.Lock()
	defer customer360CacheMutex.Unlock()

	customer360Cache[customerID] = view
	customer360CacheExpiration[customerID] = time.Now().Add(15 * time.Minute)
}

func setupDataIntegrationRoutes(r *gin.Engine) {
	integration := r.Group("/integration")
	integration.Use(AIMiddleware())
	{
		integration.GET("/customer-360/:customer_id", getCustomer360View)
		integration.POST("/sync-data", syncExternalData)
	}
}

func getCustomer360View(c *gin.Context) {
	customerID := c.Param("customer_id")
	
	if cachedView := getCachedCustomer360View(customerID); cachedView != nil {
		c.JSON(http.StatusOK, cachedView)
		return
	}

	aiSuggestion, _ := c.Get("ai_suggestion")

	customerData := getCustomerData(customerID)
	salesData := getCustomerSalesData(customerID)
	analyticsData := getCustomerAnalytics(customerID)
	recommendationData := getCustomerRecommendations(customerID)
	interactionHistory := getCustomerInteractionHistory(customerID)
	sentimentHistory := getCustomerSentimentHistory(customerID)
	churnProbability := predictChurn(customerData)
	lifetimeValue := calculateLifetimeValue(customerID)

	view := gin.H{
		"customer_id":        customerID,
		"customer_data":      customerData,
		"sales_data":         salesData,
		"analytics_data":     analyticsData,
		"recommendations":    recommendationData,
		"interaction_history": interactionHistory,
		"sentiment_history":   sentimentHistory,
		"churn_probability":   churnProbability,
		"lifetime_value":      lifetimeValue,
		"ai_suggestion":       aiSuggestion,
	}

	setCachedCustomer360View(customerID, view)
	c.JSON(http.StatusOK, view)
}

func syncExternalData(c *gin.Context) {
	// Simula sincronização de dados externos
	c.JSON(http.StatusOK, gin.H{
		"message": "Dados externos sincronizados com sucesso",
		"synced_sources": []string{"ERP", "E-commerce platform", "Social media"},
	})
}

// Funções auxiliares para simular a obtenção de dados de diferentes serviços

