package main

import (
	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"math"
	"net/http"
	"sort"
)

type RecommendationEngine struct {
	UserItemMatrix *mat.Dense
	Users          []string
	Items          []string
}

func NewRecommendationEngine() *RecommendationEngine {
	return &RecommendationEngine{
		UserItemMatrix: mat.NewDense(0, 0, nil),
		Users:          make([]string, 0),
		Items:          make([]string, 0),
	}
}

func (re *RecommendationEngine) AddRating(user, item string, rating float64) {
	userIndex := re.getUserIndex(user)
	itemIndex := re.getItemIndex(item)

	rows, cols := re.UserItemMatrix.Dims()
	if userIndex >= rows || itemIndex >= cols {
		newRows := max(rows, userIndex+1)
		newCols := max(cols, itemIndex+1)
		newMatrix := mat.NewDense(newRows, newCols, nil)
		newMatrix.Copy(re.UserItemMatrix)
		re.UserItemMatrix = newMatrix
	}

	re.UserItemMatrix.Set(userIndex, itemIndex, rating)
}

func (re *RecommendationEngine) GetRecommendations(user string, n int) []string {
	userIndex := re.getUserIndex(user)
	if userIndex >= re.UserItemMatrix.RawMatrix().Rows {
		return []string{}
	}

	userRatings := re.UserItemMatrix.RowView(userIndex)
	similarities := make([]float64, len(re.Items))

	for i := 0; i < len(re.Items); i++ {
		if userRatings.AtVec(i) == 0 {
			itemRatings := re.UserItemMatrix.ColView(i)
			similarities[i] = stat.Correlation(userRatings, itemRatings, nil)
		}
	}

	type itemScore struct {
		item  string
		score float64
	}

	var scores []itemScore
	for i, score := range similarities {
		if !math.IsNaN(score) {
			scores = append(scores, itemScore{re.Items[i], score})
		}
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	recommendations := make([]string, min(n, len(scores)))
	for i := 0; i < len(recommendations); i++ {
		recommendations[i] = scores[i].item
	}

	return recommendations
}

func (re *RecommendationEngine) getUserIndex(user string) int {
	for i, u := range re.Users {
		if u == user {
			return i
		}
	}
	re.Users = append(re.Users, user)
	return len(re.Users) - 1
}

func (re *RecommendationEngine) getItemIndex(item string) int {
	for i, it := range re.Items {
		if it == item {
			return i
		}
	}
	re.Items = append(re.Items, item)
	return len(re.Items) - 1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var recommendationEngine = NewRecommendationEngine()

func setupRecommendationEngineRoutes(r *gin.Engine) {
	rec := r.Group("/recommendations")
	{
		rec.POST("/rate", addRating)
		rec.GET("/user/:user", getUserRecommendations)
	}
}

func addRating(c *gin.Context) {
	var rating struct {
		User   string  `json:"user"`
		Item   string  `json:"item"`
		Rating float64 `json:"rating"`
	}
	if err := c.BindJSON(&rating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recommendationEngine.AddRating(rating.User, rating.Item, rating.Rating)
	c.JSON(http.StatusOK, gin.H{"message": "Rating added successfully"})
}

func getUserRecommendations(c *gin.Context) {
	user := c.Param("user")
	recommendations := recommendationEngine.GetRecommendations(user, 5)
	c.JSON(http.StatusOK, gin.H{"recommendations": recommendations})
}


