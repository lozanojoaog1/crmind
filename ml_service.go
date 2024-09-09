package main

import (
	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"math"
	"net/http"
)

type MLModel struct {
	Weights *mat.VecDense
}

func NewMLModel(features int) *MLModel {
	return &MLModel{
		Weights: mat.NewVecDense(features, nil),
	}
}

func (m *MLModel) Predict(features *mat.VecDense) float64 {
	return mat.Dot(m.Weights, features)
}

func (m *MLModel) Train(X *mat.Dense, y *mat.VecDense, learningRate float64, epochs int) {
	nSamples, nFeatures := X.Dims()
	for epoch := 0; epoch < epochs; epoch++ {
		for i := 0; i < nSamples; i++ {
			xi := X.RowView(i)
			yi := y.AtVec(i)
			
			prediction := m.Predict(xi.(*mat.VecDense))
			error := yi - prediction
			
			gradient := mat.NewVecDense(nFeatures, nil)
			gradient.ScaleVec(-2*error, xi.(*mat.VecDense))
			
			m.Weights.AddScaledVec(m.Weights, -learningRate, gradient)
		}
	}
}

var churnModel = NewMLModel(5) // 5 features for simplicity

func setupMLRoutes(r *gin.Engine) {
	ml := r.Group("/ml")
	{
		ml.POST("/train", trainModel)
		ml.POST("/predict", predictChurn)
	}
}

func trainModel(c *gin.Context) {
	// Simulated training data
	X := mat.NewDense(100, 5, generateRandomData(100*5))
	y := mat.NewVecDense(100, generateRandomLabels(100))
	
	churnModel.Train(X, y, 0.01, 1000)
	
	c.JSON(http.StatusOK, gin.H{"message": "Model trained successfully"})
}

func predictChurn(c *gin.Context) {
	var features struct {
		Data []float64 `json:"data"`
	}
	if err := c.BindJSON(&features); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if len(features.Data) != 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expected 5 features"})
		return
	}
	
	featureVec := mat.NewVecDense(5, features.Data)
	prediction := churnModel.Predict(featureVec)
	
	c.JSON(http.StatusOK, gin.H{
		"churn_probability": math.Max(0, math.Min(1, prediction)),
	})
}

func generateRandomData(n int) []float64 {
	data := make([]float64, n)
	for i := range data {
		data[i] = stat.NormFloat64()
	}
	return data
}

func generateRandomLabels(n int) []float64 {
	labels := make([]float64, n)
	for i := range labels {
		if stat.NormFloat64() > 0 {
			labels[i] = 1
		}
	}
	return labels
}
