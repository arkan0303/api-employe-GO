package controllers

import (
	"net/http"
	"strconv"

	"api-rect-go/services"

	"github.com/gin-gonic/gin"
)

func NewTimesheetController(service *services.TimesheetService) *TimesheetController {
	return &TimesheetController{Service: service}
}

func GetMergedDatass(ctx *gin.Context) {
	bulan, _ := strconv.Atoi(ctx.Param("bulan"))
	tahun, _ := strconv.Atoi(ctx.Param("tahun"))
	idCutOff, _ := strconv.Atoi(ctx.Param("periode"))

	data, err := services.GetMergedDatas(bulan, tahun, idCutOff)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data retrieved successfully",
		"data":    data,
	})
}
