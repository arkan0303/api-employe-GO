package controllers

import (
	models "api-rect-go/modals"
	"api-rect-go/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMobils(c *gin.Context){
	mobils , err := services.GetAllMobils()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK , mobils)
}

func CreateMobil(c *gin.Context){
	var mobil models.Mobil
	if err := c.ShouldBindJSON(&mobil); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" :err.Error(),
		})
		return
	}

	if err := services.CreateMobil(&mobil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK , mobil)
}