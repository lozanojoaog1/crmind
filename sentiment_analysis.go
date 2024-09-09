package main

import (
	"github.com/cdipaolo/sentiment"
	"strings"
)

var model sentiment.Models

func initSentimentModel() {
	var err error
	model, err = sentiment.Restore()
	if err != nil {
		logger.Fatalf("Falha ao inicializar o modelo de sentimento: %v", err)
	}
}

func analyzeSentiment(text string) float64 {
	analysis := model.SentimentAnalysis(strings.ToLower(text), sentiment.English)
	return float64(analysis.Score) / 4.0 // Normaliza para o intervalo 0-1
}

func analyzeSentimentBatch(texts []string) float64 {
	var totalScore float64
	for _, text := range texts {
		totalScore += analyzeSentiment(text)
	}
	return totalScore / float64(len(texts))
}
