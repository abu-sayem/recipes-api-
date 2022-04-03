// Recipe API
// Language: go
// This is a sample recipe API. You can use this as a starting point for your own recipe API.
//
//    Schemes: http
//    Host: localhost:8080
//    BasePath: /
//    Version: 0.0.1
//	  Contact: Abu sayem < abusaayem@gmail.com
//
//    Consumes:
//    - application/json
//
//    Produces:
//    - application/json
//swagger:meta

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []Recipe

type Recipe struct {
	//swagger:ignore
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

// swager:opertaion get /recipes listRecipes
// List all recipes
//
// produces:
// - application/json
// responses:
//   '200':
//     description: A list of recipes
func ListRecipedHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}


// swager: opeartion POST /recipes createRecipe
// Create a new recipe
//
// consumes:
// - application/json
//
// produces:
// - application/json
//
// responses:
//   '201':
//     description: Recipe created
//   '400':
//     description: Invalid recipe
func CreateRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()

	recipes = append(recipes, recipe)

	c.JSON(http.StatusCreated, recipe)

}


// swager: opration put /recipes/{id} updateRecipe
// Update existing recipe
//
// parmaeters:
// - name: id
//   in: path
//   required: true
//   description: id of recipe to update
//   tyoe: string
//
// consumes:
// - application/json
//
// produces:
// - application/json
//
// responses:
//   '200':
//     description: Recipe updated
//   '400':
//     description: Invalid recipe
func UpdateRecipeHandler(c *gin.Context) {
	recipeID := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, item := range recipes {
		if item.ID == recipeID {
			recipes[i] = recipe
			break
		}
	}

	c.JSON(http.StatusOK, recipe)
}

// swager: opration delete /recipes/{id} deleteRecipe
// Delete existing recipe
//
// parameters:
// - name: id
//   in: path
//   required: true
//   description: id of recipe to delete
//   type: string
//
// produces:
// - application/json
//
// responses:
//   '200':
//     description: Recipe deleted
func DeleteRecipeHandler(c *gin.Context) {
	recipeID := c.Param("id")

	for i, item := range recipes {
		if item.ID == recipeID {
			recipes = append(recipes[:i], recipes[i+1:]...)
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"id": recipeID})
}

func SearchRecipeHandler(c *gin.Context){
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)

	for _, recipe := range recipes {
		found := false

		for _, item := range recipe.Tags {
			if strings.EqualFold(item, tag) {
				found = true
				break
			}
		}

		if found {
			listOfRecipes = append(listOfRecipes, recipe)
		}
	}

		c.JSON(http.StatusOK, listOfRecipes)
}


func init() {
	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal(file, &recipes)
}

func main() {
	router := gin.Default()
	router.GET("/recipes", ListRecipedHandler)
	router.POST("/recipes", CreateRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	

	router.Run(":8080")	
}
