package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AISuggestion struct {
	gorm.Model
	Intent     string
	Suggestion string
	UsefulnessScore float32
}

func recordAISuggestion(c *gin.Context) {
	var suggestion AISuggestion
	if err := c.ShouldBindJSON(&suggestion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&suggestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar sugestão"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sugestão registrada com sucesso"})
}

func storeFeedback(suggestion string, useful bool, feedback string) error {
	db := database.GetDB()
	return db.Create(&AIFeedback{
		Suggestion: suggestion,
		Useful:     useful,
		Feedback:   feedback,
	}).Error
}

func updateAIModel(suggestion string, useful bool, feedback string) {
	// Aqui você implementaria a lógica para atualizar seu modelo de IA
	// Por enquanto, vamos apenas registrar que o modelo seria atualizado
	logger.Infof("Atualizando modelo de IA com novo feedback: %s, útil: %v", suggestion, useful)
}
