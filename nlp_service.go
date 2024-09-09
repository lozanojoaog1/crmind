package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func setupNLPRoutes(r *gin.Engine) {
	nlp := r.Group("/nlp")
	{
		nlp.POST("/analyze-sentiment", analyzeSentiment)
		nlp.POST("/extract-entities", extractEntities)
	}
}

func analyzeText(c *gin.Context) {
	var text struct {
		Content string `json:"content"`
	}
	if err := c.BindJSON(&text); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc, err := prose.NewDocument(text.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entities := make([]string, 0)
	for _, ent := range doc.Entities() {
		entities = append(entities, ent.Text)
	}

	tokenizer, err := sentences.NewSentenceTokenizer(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sentences := tokenizer.Tokenize(text.Content)

	c.JSON(http.StatusOK, gin.H{
		"entities":  entities,
		"sentences": len(sentences),
		"tokens":    len(doc.Tokens()),
	})
}

func analyzeSentiment(c *gin.Context) {
	var text struct {
		Content string `json:"content"`
	}
	if err := c.BindJSON(&text); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Implementação simplificada de análise de sentimento
	sentiment := simpleSentimentAnalysis(text.Content)

	c.JSON(http.StatusOK, gin.H{
		"sentiment": sentiment,
	})
}

func extractEntities(c *gin.Context) {
	// Implementar extração de entidades usando Google Cloud NLP
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func simpleSentimentAnalysis(text string) float64 {
	positiveWords := []string{"bom", "ótimo", "excelente", "maravilhoso", "adorei"}
	negativeWords := []string{"ruim", "péssimo", "terrível", "horrível", "odiei"}

	words := strings.Fields(strings.ToLower(text))
	var score float64

	for _, word := range words {
		if contains(positiveWords, word) {
			score += 1
		} else if contains(negativeWords, word) {
			score -= 1
		}
	}

	return score / float64(len(words))
}
