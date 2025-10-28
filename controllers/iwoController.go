package controllers

import (
	"api-rect-go/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIwoByCompanyIDController(c *gin.Context) {
	masterCompaniesIDParam := c.Param("master_companies_id")
	masterCompaniesID, err := strconv.Atoi(masterCompaniesIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ID perusahaan tidak valid",
		})
		return
	}

	iwos, err := services.GetIwoDataByMasterCompaniesID(int32(masterCompaniesID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diambil",
		"data":    iwos,
	})
}
