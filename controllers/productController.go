package controllers

import (
	models "api-rect-go/modals"
	"api-rect-go/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	products, err := services.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, products)
}

func CreateProduct(c *gin.Context){
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest , gin.H {
			"error": err.Error(),
		})
		return
	}

	if err := services.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError , gin.H {
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK , product)
}