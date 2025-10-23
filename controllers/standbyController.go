package controllers

import (
	"net/http"

	"api-rect-go/services"

	"github.com/gin-gonic/gin"
)

func GetMasterDataAvailableWithForms(c *gin.Context) {
	// Optional: force refresh external recruitment cache
	// if c.Query("refresh") == "1" {
	// 	services.ClearRecruitmentCache()
	// }

	masterData, err := services.GetMasterDataAvailableWithForms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, masterData)
}