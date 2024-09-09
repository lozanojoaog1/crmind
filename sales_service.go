package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
	"database/sql"
	_ "github.com/lib/pq"
	"strconv"
	"your-project/logger"
	"your-project/validator"
	"your-project/auth"
	"math"
	"gonum.org/v1/gonum/mat"
)

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Sale struct {
	ID          string    `json:"id"`
	CustomerID  string    `json:"customer_id"`
	ProductName string    `json:"product_name"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
}

type ChurnModel struct{}

func (m *ChurnModel) Predict(features []float64) float64 {
	// Modelo simplificado de previsão de churn
	// Baseado em: recência da última compra, frequência de compras e valor monetário
	recency := features[0]
	frequency := features[1]
	monetary := features[2]

	score := 0.3*recency + 0.3*frequency + 0.4*monetary
	return 1 / (1 + math.Exp(-score))
}

var churnModel = &ChurnModel{}

func predictCustomerChurn(customerID string) float64 {
	// Obter dados do cliente
	customerData := getCustomerData(customerID)

	// Extrair features
	recency := float64(customerData["days_since_last_purchase"].(int))
	frequency := float64(customerData["total_purchases"].(int)) / float64(customerData["months_active"].(int))
	monetary := customerData["average_purchase_value"].(float64)

	// Normalizar features (simplificado)
	recency = recency / 365 // Assumindo que 365 dias é o máximo
	frequency = frequency / 10 // Assumindo que 10 compras por mês é o máximo
	monetary = monetary / 1000 // Assumindo que $1000 é o valor máximo de compra

	features := []float64{recency, frequency, monetary}
	return churnModel.Predict(features)
}

func setupCustomerRoutes(r *gin.Engine) {
	customerGroup := r.Group("/customers")
	customerGroup.Use(AuthMiddleware())
	{
		customerGroup.POST("", createCustomer)
		customerGroup.GET("/:id", getCustomer)
		customerGroup.PUT("/:id", updateCustomer)
		customerGroup.DELETE("/:id", deleteCustomer)
		customerGroup.GET("/:id/churn", getChurnPrediction)
	}
}

func setupSalesRoutes(r *gin.Engine) {
	salesGroup := r.Group("/sales")
	salesGroup.Use(AuthMiddleware())
	{
		salesGroup.POST("", createSale)
		salesGroup.GET("/:id", getSale)
		salesGroup.GET("", listSales)
	}
}

func createCustomer(c *gin.Context) {
	var newCustomer Customer
	if err := c.ShouldBindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !validator.IsValidEmail(newCustomer.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email inválido"})
		return
	}

	if !validator.IsValidPhone(newCustomer.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Telefone inválido"})
		return
	}

	newCustomer.ID = uuid.New().String()
	newCustomer.CreatedAt = time.Now()
	newCustomer.UpdatedAt = time.Now()

	// Aqui você deve salvar o cliente no banco de dados
	// Por enquanto, vamos apenas retornar o cliente criado
	c.JSON(http.StatusCreated, newCustomer)
}

func getCustomer(c *gin.Context) {
	customerID := c.Param("id")

	// Aqui você deve buscar o cliente no banco de dados
	// Por enquanto, vamos retornar um cliente de exemplo
	customer := Customer{
		ID:        customerID,
		Name:      "Cliente Exemplo",
		Email:     "cliente@exemplo.com",
		Phone:     "123456789",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	c.JSON(http.StatusOK, customer)
}

func updateCustomer(c *gin.Context) {
	customerID := c.Param("id")
	var updatedCustomer Customer
	if err := c.ShouldBindJSON(&updatedCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !validator.IsValidEmail(updatedCustomer.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email inválido"})
		return
	}

	if !validator.IsValidPhone(updatedCustomer.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Telefone inválido"})
		return
	}

	// Aqui você deve atualizar o cliente no banco de dados
	// Por enquanto, vamos apenas retornar o cliente atualizado
	updatedCustomer.ID = customerID
	updatedCustomer.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, updatedCustomer)
}

func deleteCustomer(c *gin.Context) {
	customerID := c.Param("id")

	// Aqui você deve deletar o cliente do banco de dados
	// Por enquanto, vamos apenas retornar uma mensagem de sucesso
	c.JSON(http.StatusOK, gin.H{"message": "Cliente deletado com sucesso"})
}

func createSale(c *gin.Context) {
	var newSale Sale
	if err := c.ShouldBindJSON(&newSale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newSale.ID = uuid.New().String()
	newSale.Date = time.Now()

	// Aqui você deve salvar a venda no banco de dados
	// Por enquanto, vamos apenas retornar a venda criada
	c.JSON(http.StatusCreated, newSale)
}

func getSale(c *gin.Context) {
	saleID := c.Param("id")

	// Aqui você deve buscar a venda no banco de dados
	// Por enquanto, vamos retornar uma venda de exemplo
	sale := Sale{
		ID:          saleID,
		CustomerID:  "customer123",
		ProductName: "Produto Exemplo",
		Amount:      100.00,
		Date:        time.Now(),
	}

	c.JSON(http.StatusOK, sale)
}

func listSales(c *gin.Context) {
	// Aqui você deve buscar as vendas no banco de dados
	// Por enquanto, vamos retornar uma lista de exemplo
	sales := []Sale{
		{
			ID:          "sale1",
			CustomerID:  "customer1",
			ProductName: "Produto A",
			Amount:      100.00,
			Date:        time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "sale2",
			CustomerID:  "customer2",
			ProductName: "Produto B",
			Amount:      150.00,
			Date:        time.Now(),
		},
	}

	c.JSON(http.StatusOK, sales)
}

func getChurnPrediction(c *gin.Context) {
	customerID := c.Param("id")
	churnProbability := predictCustomerChurn(customerID)

	c.JSON(http.StatusOK, gin.H{
		"customer_id":        customerID,
		"churn_probability":  churnProbability,
	})
}

func getCustomerData(customerID string) map[string]interface{} {
	// Implementação simplificada para obter dados do cliente
	// Aqui você deve buscar os dados do cliente no banco de dados
	// Por enquanto, vamos retornar dados de exemplo
	return map[string]interface{}{
		"days_since_last_purchase": 30,
		"total_purchases":          10,
		"months_active":            12,
		"average_purchase_value":  100.0,
	}
}

func getTotalCustomers() int {
	// Implementação real: consultar banco de dados
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM customers").Scan(&count)
	if err != nil {
		logger.ErrorLogger.Printf("Erro ao contar clientes: %v", err)
		return 0
	}
	return count
}

func getActiveCustomers() int {
	// Implementação real: consultar banco de dados
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM customers WHERE last_activity > $1", time.Now().AddDate(0, -1, 0)).Scan(&count)
	if err != nil {
		logger.ErrorLogger.Printf("Erro ao contar clientes ativos: %v", err)
		return 0
	}
	return count
}

func getTotalRevenue() float64 {
	// Implementação real: consultar banco de dados
	var total float64
	err := db.QueryRow("SELECT SUM(amount) FROM sales").Scan(&total)
	if err != nil {
		logger.ErrorLogger.Printf("Erro ao calcular receita total: %v", err)
		return 0
	}
	return total
}

func getSalesTrend() []gin.H {
	// Implementação real: consultar banco de dados
package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func setupSalesRoutes(r *gin.Engine) {
	sales := r.Group("/sales")
	{
		sales.GET("", listSales)
		sales.POST("", createSale)
		sales.GET("/:id", getSale)
		sales.PUT("/:id", updateSale)
		sales.DELETE("/:id", deleteSale)
	}
}

func listSales(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"sales": []gin.H{
			{"id": 1, "customer_id": 1, "amount": 1000},
			{"id": 2, "customer_id": 2, "amount": 1500},
		},
	})
}

func createSale(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Venda criada com sucesso",
		"id": 3,
	})
}

func getSale(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id": id,
		"customer_id": 1,
		"amount": 2000,
	})
}

func updateSale(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Venda " + id + " atualizada com sucesso",
	})
}

func deleteSale(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Venda " + id + " deletada com sucesso",
	})
}
