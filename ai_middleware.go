package main

import (
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"fmt"
	"log"
)

func AIMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ai_processed", true)
		
		// Analisa o contexto da requisição
		intent := analyzeIntent(c)
		c.Set("ai_intent", intent)
		
		// Gera sugestões baseadas no contexto e intenção
		suggestion := generateSuggestion(c, intent)
		c.Set("ai_suggestion", suggestion)
		
		// Adiciona informações do usuário ao contexto
		userID, exists := c.Get("user_id")
		if exists {
			userContext := getUserContext(userID.(string))
			c.Set("user_context", userContext)
		}

		c.Next()
	}
}

func analyzeIntent(c *gin.Context) string {
	path := c.FullPath()
	method := c.Request.Method

	switch {
	case strings.Contains(path, "/customers") && method == "GET":
		return "view_customer"
	case strings.Contains(path, "/customers") && method == "POST":
		return "create_customer"
	case strings.Contains(path, "/sales") && method == "GET":
		return "view_sales"
	case strings.Contains(path, "/sales") && method == "POST":
		return "create_sale"
	case strings.Contains(path, "/analytics"):
		return "analyze_data"
	case strings.Contains(path, "/recommendations"):
		return "get_recommendations"
	case strings.Contains(path, "/integration"):
		return "integrate_data"
	default:
		return "general_inquiry"
	}
}

func generateSuggestion(c *gin.Context, intent string) string {
    userContext, _ := c.Get("user_context")
    path := c.FullPath()

    // Obter sugestões anteriores bem avaliadas
    previousSuggestions := getPreviousSuccessfulSuggestions(intent)

    // Gerar uma nova sugestão baseada no contexto atual e sugestões anteriores
    suggestion := generateNewSuggestion(intent, userContext.(map[string]interface{}), previousSuggestions)

    return suggestion
}

func getPreviousSuccessfulSuggestions(intent string) []string {
    // Implementar lógica para buscar sugestões bem avaliadas do banco de dados
    // Por enquanto, retornaremos algumas sugestões fixas
    return []string{
        "Considere oferecer um desconto personalizado baseado no histórico de compras do cliente.",
        "Analise o padrão de compras sazonais para otimizar o estoque.",
        "Utilize dados de interações recentes para personalizar a abordagem de vendas.",
    }
}

func generateNewSuggestion(intent string, userContext map[string]interface{}, previousSuggestions []string) string {
    // Aqui você implementaria a lógica para gerar uma nova sugestão
    // Por enquanto, vamos retornar uma sugestão baseada no intent e contexto
    switch intent {
    case "view_customer":
        return fmt.Sprintf("Analise o histórico de compras do cliente e considere oferecer produtos complementares. Última compra: %v", userContext["last_purchase"])
    case "create_sale":
        return "Verifique produtos frequentemente comprados juntos e sugira uma oferta combinada."
    default:
        if len(previousSuggestions) > 0 {
            return previousSuggestions[0] // Retorna a primeira sugestão bem-sucedida anterior
        }
        return "Como posso ajudar a otimizar suas operações hoje?"
    }
}

// Função para obter o contexto do usuário (implementar posteriormente)
func getUserContext(userID string) map[string]interface{} {
	// Implementação temporária
	return map[string]interface{}{
		"last_login": time.Now().Add(-24 * time.Hour),
		"role":       "manager",
	}
}

func recordAISuggestion(c *gin.Context) {
	var feedback struct {
		Suggestion string  `json:"suggestion"`
		Useful     bool    `json:"useful"`
		Feedback   string  `json:"feedback"`
	}
	if err := c.ShouldBindJSON(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Armazenar o feedback no banco de dados
	err := storeFeedback(feedback.Suggestion, feedback.Useful, feedback.Feedback)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao armazenar feedback"})
		return
	}

	// Atualizar o modelo de IA com o novo feedback
	updateAIModel(feedback.Suggestion, feedback.Useful, feedback.Feedback)

	c.JSON(http.StatusOK, gin.H{"message": "Feedback registrado e modelo atualizado com sucesso"})
}
