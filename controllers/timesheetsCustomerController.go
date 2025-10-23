package controllers

import (
	"api-rect-go/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TimesheetController struct {
	Service *services.TimesheetService
}

func  GetMergedData(ctx *gin.Context) {
	bulan := strings.TrimSpace(ctx.Query("bulan"))
	tahunStr := strings.TrimSpace(ctx.Query("tahun"))

	tahun, err := strconv.Atoi(tahunStr)
if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parameter tahun tidak valid"})
    return
}

	data, err := services.GetMergedData(bulan, int32(tahun))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diambil",
		"data":    data,
	})
}
