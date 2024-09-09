package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("sua_chave_secreta_aqui") // Em produção, use uma variável de ambiente

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func generateToken(userID string, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Token de autorização não fornecido"})
			c.Abort()
			return
		}

		claims, err := validateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func roleAuthorization(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(401, gin.H{"error": "Role não encontrada"})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(500, gin.H{"error": "Erro interno do servidor"})
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{"error": "Acesso não autorizado"})
		c.Abort()
	}
}

func registerUser(c *gin.Context) {
	var user struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	
	if err != nil {
		c.JSON(500, gin.H{"error": "Falha ao criar usuário"})
		return
	}

	// Aqui você salvaria o usuário no banco de dados
	// Por enquanto, vamos apenas simular isso
	c.JSON(201, gin.H{"message": "Usuário criado com sucesso"})
}

func loginUser(c *gin.Context) {
	var credentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Aqui você buscaria o usuário no banco de dados e verificaria a senha
	// Por enquanto, vamos simular isso
	
	hashedPassword, err := hashPassword(credentials.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Falha ao processar senha"})
		return
	}

	if !checkPasswordHash(credentials.Password, hashedPassword) {
		c.JSON(401, gin.H{"error": "Credenciais inválidas"})
		return
	}

	token, err := generateToken(credentials.Username, "user")
	if err != nil {
		c.JSON(500, gin.H{"error": "Falha ao gerar token"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func setupAuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", login)
		auth.POST("/refresh", refreshToken)
	}
}

func login(c *gin.Context) {
	// Implementar lógica de login
}

func refreshToken(c *gin.Context) {
	// Implementar lógica de refresh do token
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementar verificação do token JWT
	}
}

func generateRefreshToken(userID string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func refreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		c.JSON(400, gin.H{"error": "Refresh token não fornecido"})
		return
	}

	claims, err := validateToken(refreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "Refresh token inválido"})
		return
	}

	// Aqui você buscaria o usuário no banco de dados para obter o role atual
	// Por enquanto, vamos simular isso
	role := "user"

	newToken, err := generateToken(claims.UserID, role)
	if err != nil {
		c.JSON(500, gin.H{"error": "Falha ao gerar novo token"})
		return
	}

	c.JSON(200, gin.H{"token": newToken})
}
