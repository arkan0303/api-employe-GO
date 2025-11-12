package controllers

import (
	"fmt"
	"api-rect-go/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReplacementData(c *gin.Context) {
	data, err := services.GetServiceReplacementData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   data,
	})
}

func DeleteReplacementData(c *gin.Context) {
	// Get ID from URL parameter
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	// Convert ID to int64
	var idInt64 int64
	_, err := fmt.Sscanf(id, "%d", &idInt64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Call service to delete
	err = services.DeleteServiceReplacement(idInt64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"message": "Data berhasil dihapus",
	})
}
