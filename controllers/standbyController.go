package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"api-rect-go/db"
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

// UpdateMasterData handles the update of master data including photo upload
func UpdateMasterData(c *gin.Context) {
	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	// Inisialisasi map untuk menyimpan data update
	updateData := make(map[string]interface{})

	// Cek apakah request berupa form-data (untuk upload file)
	if c.ContentType() == "multipart/form-data" {
		// Handle file upload
		file, err := c.FormFile("foto")
		if err == nil && file != nil {
			// Upload file ke Cloudinary
			fotoURL, err := db.UploadFile(file, "standby_photos")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Gagal mengunggah foto: " + err.Error(),
				})
				return
			}
			updateData["foto"] = fotoURL
		}

		// Handle other form fields
		for key, values := range c.Request.PostForm {
			if key != "foto" && len(values) > 0 {
				updateData[key] = values[0]
			}
		}
	} else {
		// Handle JSON request
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Data tidak valid: " + err.Error(),
			})
			return
		}

		// Handle base64 photo if exists
		if fotoBase64, ok := updateData["foto"].(string); ok && fotoBase64 != "" {
			// Jika diawali dengan 'data:image', berarti base64
			if strings.HasPrefix(fotoBase64, "data:image") {
				fotoURL, err := db.UploadFromBase64(fotoBase64, "standby_photos", fmt.Sprintf("standby_%d", id))
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Gagal mengunggah foto: " + err.Error(),
					})
					return
				}
				updateData["foto"] = fotoURL
			}
		}
	}

	// Remove ID from update data to prevent changing the ID
	delete(updateData, "id")

	// Jika tidak ada yang diupdate
	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tidak ada data yang diupdate",
		})
		return
	}

	// Call service to update data
	err = services.UpdateMasterData(int32(id), updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memperbarui data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diperbarui",
		"id":      id,
	})
}