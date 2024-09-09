package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func setupChatbotRoutes(r *gin.Engine) {
	chatbot := r.Group("/chatbot")
	chatbot.Use(AIMiddleware())
	{
		chatbot.POST("/message", handleChatbotMessage)
	}
}

func handleChatbotMessage(c *gin.Context) {
	var message struct {
		Text string `json:"text"`
	}

	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	intent := analyzeChatIntent(message.Text)
	response := generateChatResponse(intent)

	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}

func analyzeChatIntent(message string) string {
	message = strings.ToLower(message)

	switch {
	case strings.Contains(message, "preço") || strings.Contains(message, "custo"):
		return "pricing_inquiry"
	case strings.Contains(message, "entrega") || strings.Contains(message, "envio"):
		return "shipping_inquiry"
	case strings.Contains(message, "produto") || strings.Contains(message, "item"):
		return "product_inquiry"
	case strings.Contains(message, "reclamação") || strings.Contains(message, "problema"):
		return "complaint"
	default:
		return "general_inquiry"
	}
}

func generateChatResponse(intent string) string {
	switch intent {
	case "pricing_inquiry":
		return "Nossos preços variam dependendo do produto e da quantidade. Posso te ajudar a encontrar informações específicas sobre algum item?"
	case "shipping_inquiry":
		return "Geralmente, nossas entregas são feitas em 3-5 dias úteis. Para informações mais precisas, precisaria saber seu CEP e os itens que você está interessado."
	case "product_inquiry":
		return "Temos uma ampla gama de produtos. Você está procurando algo específico? Posso te ajudar a encontrar o produto ideal para suas necessidades."
	case "complaint":
		return "Lamento ouvir que você está tendo problemas. Pode me dar mais detalhes sobre a situação? Farei o possível para resolver seu problema rapidamente."
	default:
		return "Como posso te ajudar hoje? Estou aqui para responder perguntas sobre nossos produtos, preços, entregas ou qualquer outra dúvida que você possa ter."
	}
}
