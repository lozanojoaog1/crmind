package main

import (
    "CRMind/backend/auth"
    "CRMind/backend/database"
    "fmt"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "database/sql"
    _ "github.com/lib/pq"
)

var db *sql.DB

// Inicializar o banco de dados
func initDB() {
    connStr := "host=p-2hrz9m0lvu.pg.biganimal.io port=5432 user=edb_admin password=CRMind*2024*898 dbname=edb_admin sslmode=require"
    // Substitua 'your_user' e 'your_password' pelas suas credenciais do PostgreSQL
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    if err = db.Ping(); err != nil {
        log.Fatal("Falha ao conectar ao banco de dados:", err)
    }
    log.Println("Conexão com o banco de dados estabelecida com sucesso")

    // Criar tabelas se não existirem
    createTables()
}

// Criar tabelas se necessário
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

// Função principal que inicia o servidor
func main() {
    initDB()

    r := gin.Default()

    // Configurar rotas de autenticação
    setupAuthRoutes(r)

    // Configurar rotas de clientes
    setupCustomerRoutes(r)

    // Configurar rotas de vendas
    setupSalesRoutes(r)

    // Iniciar o servidor na porta 8080
    r.Run(":8080")
}

// Configurar rotas de autenticação
func setupAuthRoutes(r *gin.Engine) {
    authGroup := r.Group("/auth")
    {
        authGroup.POST("/login", login)
        authGroup.POST("/register", register)
    }
}

// Função de login
func login(c *gin.Context) {
    var loginData struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var userID string
    err := db.QueryRow("SELECT id FROM users WHERE email = $1 AND password = $2", loginData.Email, loginData.Password).Scan(&userID)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
        return
    }

    token, err := auth.GenerateToken(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

// Função de registro
func register(c *gin.Context) {
    var registerData struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
        Name     string `json:"name" binding:"required"`
    }

    if err := c.ShouldBindJSON(&registerData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Aqui você pode salvar o novo usuário no banco de dados
    c.JSON(http.StatusOK, gin.H{"message": "Usuário registrado com sucesso"})
}

// Middleware para verificar token de autenticação
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autenticação não fornecido"})
            c.Abort()
            return
        }

        claims, err := auth.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Next()
    }
}

// Configurar rotas de clientes
func setupCustomerRoutes(r *gin.Engine) {
    customerGroup := r.Group("/customers")
    customerGroup.Use(AuthMiddleware())
    {
        customerGroup.POST("", createCustomer)
        customerGroup.GET("/:id", getCustomer)
        customerGroup.PUT("/:id", updateCustomer)
        customerGroup.DELETE("/:id", deleteCustomer)
    }
}

// Configurar rotas de vendas
func setupSalesRoutes(r *gin.Engine) {
    salesGroup := r.Group("/sales")
    salesGroup.Use(AuthMiddleware())
    {
        salesGroup.POST("", createSale)
        salesGroup.GET("/:id", getSale)
        salesGroup.GET("", listSales)
    }
}

// Exemplos de funções para CRUD de clientes e vendas
// Aqui você pode implementar as funções createCustomer, getCustomer, updateCustomer, deleteCustomer, createSale, etc.
