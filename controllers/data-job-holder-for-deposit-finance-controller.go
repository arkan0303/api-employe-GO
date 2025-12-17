package controllers

import (
	"api-rect-go/services"
	"net/http"

	"github.com/gin-gonic/gin"
)


func GetDataJobHolderForDepositFinance(c *gin.Context) {
	data, err := services.DataJobHolderForDepositFinance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   data,
	})
}