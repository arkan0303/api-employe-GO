package controllers

import (
	models "api-rect-go/modals/mysql"
	"api-rect-go/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceUserController struct {
	Service *services.ServiceUserService
}

func NewServiceUserController(s *services.ServiceUserService) *ServiceUserController {
	return &ServiceUserController{Service: s}
}

func (ctrl *ServiceUserController) GetByIwoTemplateID(c *gin.Context) {
	iwoTemplateIDParam := c.Param("iwo_template_id")

	iwoTemplateID, err := strconv.Atoi(iwoTemplateIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	users, err := ctrl.Service.GetByIwoTemplateID(int32(iwoTemplateID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diambil",
		"data":    users,
	})
}

func (ctrl *ServiceUserController) CreateUser(c *gin.Context) {
	var user models.ServiceUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Gagal membuat user", "error": err.Error()})
		return
	}

	if err := ctrl.Service.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dibuat"})
}


