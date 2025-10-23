package controllers

import (
	"api-rect-go/db"
	models "api-rect-go/modals"
	"api-rect-go/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMasterData(c *gin.Context) {
	masterData, err := services.GetMasterData()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, masterData)
}

func PostMasterDataExternal(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka"})
        return
    }

    if err := services.PostMasterDataExternal(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Data berhasil dikirim ke API eksternal",
    })
}

func CreateMasterData(c *gin.Context) {
	var masterData models.TbMasterDataDiri
	if err := c.ShouldBindJSON(&masterData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := services.CreateMasterData(&masterData); err != nil{
		c.JSON(http.StatusInternalServerError , gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK , masterData)
}

func EditMasterData(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka"})
		return
	}

	if err := services.EditMasterData(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updatedData models.TbMasterDataDiri
	if err := db.DB.First(&updatedData, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diupdate",
		"data":    updatedData,
	})
}

