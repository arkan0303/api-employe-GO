package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"api-rect-go/services"
)

type UpdateByStatusDataDiriRequest struct {
	IDCustomer     int32 `json:"id_customer" binding:"required"`
	IDUsers        int32 `json:"id_users" binding:"required"`
	ServiceUsersID int32 `json:"service_users_id" binding:"required"`
}

func UpdateByStatusDataDiri(c *gin.Context) {
	idStatusStr := c.Param("id_status_data_diri")
	idStatus, err := strconv.Atoi(idStatusStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id_status_data_diri tidak valid",
		})
		return
	}

	var req UpdateByStatusDataDiriRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "body tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if err := services.UpdateByStatusDataDiri(
		int32(idStatus),
		req.IDCustomer,
		req.IDUsers,
		req.ServiceUsersID,
	); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update berhasil",
	})
}
