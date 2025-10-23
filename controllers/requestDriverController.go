package controllers

import (
	models "api-rect-go/modals"
	"api-rect-go/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetRequestDrivers(c *gin.Context){
requestDriver , err := services.GetRequestDrivers()
if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
	return
}

c.JSON(http.StatusOK, gin.H{
	"message": "Data berhasil diambil",
	"status":  http.StatusOK,
	"data":    requestDriver,
})

}

func CreateRequestDriverController(c *gin.Context) {
	var req models.TbRequestDriver

	// ambil form field
	req.NamaCustomer = c.PostForm("nama_customer")
	req.NamaDriver = c.PostForm("nama_driver")
	req.NoDriver = c.PostForm("no_driver")
	req.LokasiMobil = c.PostForm("lokasi_mobil")
	req.Pic = c.PostForm("pic")
	req.NoPic = c.PostForm("no_pic")

	// parsing tanggal
	layout := "2006-01-02"
	tglKerja, _ := time.Parse(layout, c.PostForm("tgl_kerja"))
	tglSelesai, _ := time.Parse(layout, c.PostForm("tgl_selesai"))
	req.TglKerja = tglKerja
	req.TglSelesai = tglSelesai

	// ambil file foto
	file, header, err := c.Request.FormFile("foto_mobil")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Foto mobil wajib diupload"})
		return
	}

	err = services.CreateRequestDriver(&req, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request driver berhasil dibuat",
		"data":    req,
	})
}