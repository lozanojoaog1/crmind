package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/mat"
	"github.com/sirupsen/logrus"
	"database/sql"
	_ "github.com/lib/pq"
	"strconv"
	"your-project/logger"
	"your-project/validator"
	"your-project/auth"
	"github.com/google/uuid"
	"time"
	"math"
)

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

var db *sql.DB

func initDB() {
	connStr := "user=your_user password=your_password dbname=crmind sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Conexão com o banco de dados estabelecida com sucesso")

	// Criar tabelas se não existirem
	createTables()
}

func createTables() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS customers (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS sales (
			id SERIAL PRIMARY KEY,
			customer_id INTEGER REFERENCES customers(id),
			product_name VARCHAR(100) NOT NULL,
			amount DECIMAL(10, 2) NOT NULL,
			date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func setupCustomerRoutes(r *gin.Engine) {
	customerGroup := r.Group("/customers")
	customerGroup.Use(AuthMiddleware())
	{
		customerGroup.GET("", listCustomers)
		customerGroup.POST("", createCustomer)
		customerGroup.GET("/:id", getCustomer)
		customerGroup.PUT("/:id", updateCustomer)
		customerGroup.DELETE("/:id", deleteCustomer)
		customerGroup.GET("/:id/insights", getCustomerInsights)
		customerGroup.POST("/:id/interaction", recordCustomerInteraction)
	}
}

func listCustomers(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, email FROM customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar clientes"})
		return
	}
	defer rows.Close()

	var customers []gin.H
	for rows.Next() {
		var id, name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao ler dados do cliente"})
			return
		}
		customers = append(customers, gin.H{"id": id, "name": name, "email": email})
	}

	c.JSON(http.StatusOK, customers)
}

func createCustomer(c *gin.Context) {
	var newCustomer struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id string
	err := db.QueryRow("INSERT INTO customers (name, email) VALUES ($1, $2) RETURNING id", newCustomer.Name, newCustomer.Email).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar cliente"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "name": newCustomer.Name, "email": newCustomer.Email})
}

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func getCustomer(c *gin.Context) {
	id := c.Param("id")
	customer, err := fetchCustomerFromDB(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
		return
	}

	aiSuggestion, _ := c.Get("ai_suggestion")
	c.JSON(http.StatusOK, gin.H{
		"customer": customer,
		"ai_suggestion": aiSuggestion,
	})
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

func recordCustomerInteraction(c *gin.Context) {
	id := c.Param("id")
	var interaction struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}
	if err := c.BindJSON(&interaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Analisar sentimento da interação
	sentiment, err := analyzeSentiment(interaction.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao analisar sentimento"})
		return
	}
	
	// Armazenar interação e sentimento (simulado)
	storeInteraction(id, interaction.Type, interaction.Content, sentiment)
	
	c.JSON(http.StatusOK, gin.H{"message": "Interaction recorded successfully"})
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
	// Implementação simplificada
	return map[string]interface{}{
		"total_purchases": 50,
		"months_active": 24,
		"support_tickets": 3,
		"average_sentiment": 0.7,
		"days_since_last_purchase": 15,
		"average_purchase_value": 100.0,
	}
}

func getRecentInteractions(id string) []string {
	// Implementação simplificada
	return []string{
		"Great product, very satisfied!",
		"Had an issue but support resolved it quickly.",
		"Thinking about trying the new features.",
	}
}

func storeInteraction(id, interactionType, content string, sentiment float64) {
	// Simula o armazenamento da interação
	logrus.WithFields(logrus.Fields{
		"customer_id": id,
		"type": interactionType,
		"content": content,
		"sentiment": sentiment,
	}).Info("Customer interaction recorded")
}

func getCustomerSentimentHistory(c *gin.Context) {
	customerID := c.Param("id")
	
	// Obter histórico de interações
	interactions := getCustomerInteractions(customerID)
	
	// Analisar sentimento de cada interação
	sentimentHistory := make([]gin.H, len(interactions))
	for i, interaction := range interactions {
		sentiment := analyzeSentiment(interaction.Content)
		sentimentHistory[i] = gin.H{
			"date": interaction.Date,
			"type": interaction.Type,
			"sentiment": sentiment,
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"customer_id": customerID,
		"sentiment_history": sentimentHistory,
	})
}

func getCustomerInteractions(customerID string) []Interaction {
	// Implementação simplificada
	return []Interaction{
		{Date: "2023-05-01", Type: "support", Content: "Ótimo atendimento!"},
		{Date: "2023-05-15", Type: "purchase", Content: "Produto chegou com defeito."},
		{Date: "2023-05-20", Type: "support", Content: "Problema resolvido rapidamente."},
	}
}

type Interaction struct {
	Date    string
	Type    string
	Content string
}
