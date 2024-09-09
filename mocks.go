package main

import (
    "github.com/gin-gonic/gin"
)

type MockAIMiddleware struct{}

func (m *MockAIMiddleware) GenerateSuggestion(c *gin.Context) string {
    return "Mock AI Suggestion"
}

func mockGetCustomerData(customerID string) map[string]interface{} {
    return map[string]interface{}{
        "name": "Mock Customer",
        "email": "mock@example.com",
    }
}

func mockPredictChurn(customerData map[string]interface{}) float64 {
    return 0.3
}

func mockGetCustomerRecommendations(customerID string, count int) []string {
    return []string{"Product A", "Product B"}
}
