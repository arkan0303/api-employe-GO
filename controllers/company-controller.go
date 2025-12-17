package controllers

import (
	"api-rect-go/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetComapny(c *gin.Context) {
	// Optional: force refresh external recruitment cache
	// if c.Query("refresh") == "1" {
	// 	services.ClearRecruitmentCache()
	// }

	masterData, err := services.GetData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, masterData)
}