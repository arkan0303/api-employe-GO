package controllers

import (
	"api-rect-go/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetForm5(c *gin.Context){
	form5 , err := services.GetForm5()
	if err != nil{
		c.JSON(http.StatusInternalServerError , gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK , form5)
}