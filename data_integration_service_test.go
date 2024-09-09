package main

import (
    "testing"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "net/http"
    "encoding/json"
    "github.com/stretchr/testify/mock"
)

func TestGetCachedCustomer360View(t *testing.T) {
    // Inicializar o cache
    customer360Cache = make(map[string]gin.H)
    customer360CacheExpiration = make(map[string]time.Time)

    // Adicionar um item ao cache
    customerID := "123"
    view := gin.H{"name": "Test Customer"}
    setCachedCustomer360View(customerID, view)

    // Testar a recuperação do cache
    cachedView := getCachedCustomer360View(customerID)
    assert.Equal(t, view, cachedView)

    // Testar a expiração do cache
    time.Sleep(16 * time.Minute)
    expiredView := getCachedCustomer360View(customerID)
    assert.Nil(t, expiredView)
}

func TestSetCachedCustomer360View(t *testing.T) {
    // Inicializar o cache
    customer360Cache = make(map[string]gin.H)
    customer360CacheExpiration = make(map[string]time.Time)

    customerID := "456"
    view := gin.H{"name": "Another Test Customer"}
    setCachedCustomer360View(customerID, view)

    assert.Equal(t, view, customer360Cache[customerID])
    assert.True(t, time.Now().Before(customer360CacheExpiration[customerID]))
}

func TestGetCustomerInteractionHistory(t *testing.T) {
    customerID := "789"
    history := getCustomerInteractionHistory(customerID)
    
    assert.NotNil(t, history)
    assert.Len(t, history, 2)
    assert.Equal(t, "support", history[0]["type"])
    assert.Equal(t, "sales", history[1]["type"])
}

func TestGetCustomerSentimentHistory(t *testing.T) {
    customerID := "101112"
    history := getCustomerSentimentHistory(customerID)
    
    assert.NotNil(t, history)
    assert.Len(t, history, 2)
    assert.InDelta(t, 0.8, history[0]["sentiment"], 0.01)
    assert.InDelta(t, 0.6, history[1]["sentiment"], 0.01)
}

func TestCalculateLifetimeValue(t *testing.T) {
    customerID := "131415"
    ltv := calculateLifetimeValue(customerID)
    
    assert.InDelta(t, 5000.0, ltv, 0.01)
}

func TestGetCustomer360View(t *testing.T) {
    // Configurar o router e o contexto de teste
    router := gin.New()
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    // Configurar o mock do middleware de IA
    c.Set("ai_suggestion", "Mock AI Suggestion")

    // Substituir funções reais por mocks
    getCustomerData = mockGetCustomerData
    predictChurn = mockPredictChurn
    getCustomerRecommendations = mockGetCustomerRecommendations

    // Chamar a função
    getCustomer360View(c)

    // Verificar o resultado
    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)

    assert.Equal(t, "Mock Customer", response["customer_data"].(map[string]interface{})["name"])
    assert.InDelta(t, 0.3, response["churn_probability"], 0.01)
    assert.Equal(t, []interface{}{"Product A", "Product B"}, response["recommendations"])
    assert.Equal(t, "Mock AI Suggestion", response["ai_suggestion"])
}
