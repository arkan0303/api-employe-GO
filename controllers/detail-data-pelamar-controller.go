package controllers

import (
	// models "api-rect-go/modals"
	"api-rect-go/services"
	"net/http"

	"github.com/gin-gonic/gin"
)


func GetDataDetailPelamar(c *gin.Context){
	mobils , err := services.GetDataPelamar()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK , mobils)
}
