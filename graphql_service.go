package main

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
)

var customerType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Customer",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"churnProbability": &graphql.Field{
				Type: graphql.Float,
			},
			"lifetimeValue": &graphql.Field{
				Type: graphql.Float,
			},
			"recommendations": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)

var dashboardSummaryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "DashboardSummary",
		Fields: graphql.Fields{
			"totalCustomers":       &graphql.Field{Type: graphql.Int},
			"activeCustomers":      &graphql.Field{Type: graphql.Int},
			"churnRate":            &graphql.Field{Type: graphql.Float},
			"totalRevenue":         &graphql.Field{Type: graphql.Float},
			"averageTicket":        &graphql.Field{Type: graphql.Float},
			"customerSatisfaction": &graphql.Field{Type: graphql.Float},
			"topProducts": &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(
					graphql.ObjectConfig{
						Name: "TopProduct",
						Fields: graphql.Fields{
							"name":  &graphql.Field{Type: graphql.String},
							"sales": &graphql.Field{Type: graphql.Int},
						},
					},
				)),
			},
			"salesTrend": &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(
					graphql.ObjectConfig{
						Name: "SalesTrend",
						Fields: graphql.Fields{
							"date":  &graphql.Field{Type: graphql.String},
							"sales": &graphql.Field{Type: graphql.Float},
						},
					},
				)),
			},
		},
	},
)

var rootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"customer": &graphql.Field{
				Type: customerType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					if !ok {
						return nil, nil
					}
					
					customerData := getCustomerData(id)
					churnProbability := predictChurn(customerData)
					recommendations := recommendationEngine.GetRecommendations(id, 3)
					
					return map[string]interface{}{
						"id":               id,
						"name":             customerData["name"],
						"email":            customerData["email"],
						"churnProbability": churnProbability,
						"lifetimeValue":    calculateLifetimeValue(id),
						"recommendations":  recommendations,
					}, nil
				},
			},
			"dashboardSummary": &graphql.Field{
				Type: dashboardSummaryType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return getDashboardSummary(), nil
				},
			},
		},
	},
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: rootQuery,
	},
)

func setupGraphQLRoutes(r *gin.Engine) {
	r.POST("/graphql", func(c *gin.Context) {
		var request struct {
			Query string `json:"query"`
		}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: request.Query,
		})
		
		c.JSON(http.StatusOK, result)
	})
}


