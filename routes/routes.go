package routes

import (
	"api-rect-go/controllers"
	"api-rect-go/db"
	"api-rect-go/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	serviceUserService := services.NewServiceUserService(db.DBMySQL)
	serviceUserController := controllers.NewServiceUserController(serviceUserService)
	r.GET("/products", controllers.GetProducts)
	r.POST("/products", controllers.CreateProduct)
	r.GET("/mobils", controllers.GetMobils)
	r.POST("/mobils", controllers.CreateMobil)

	r.GET("/form5", controllers.GetForm5)

	r.GET("/master-data", controllers.GetMasterData)
	r.POST("/master-data", controllers.CreateMasterData)
	r.PUT("/master-data/:id", controllers.EditMasterData)
	r.POST("/master-data/:id/post-external", controllers.PostMasterDataExternal)
	r.GET("/standby", controllers.GetMasterDataAvailableWithForms)
	r.PUT("/standby/:id", controllers.UpdateMasterData)
	r.GET("/jobholder", controllers.GetMasterDataAvailableWithFormss)
	r.GET("/timesheets-customer", controllers.GetMergedData)
	r.GET("/request-driver", controllers.GetRequestDrivers)
	r.GET("/cuti-driver", controllers.GetReplacementData)
	r.POST("/request-driver", controllers.CreateRequestDriverController)	
	r.GET("/iwo/:master_companies_id", controllers.GetIwoByCompanyIDController)
	r.GET("/users/:iwo_template_id", serviceUserController.GetByIwoTemplateID)
	r.POST("/users", serviceUserController.CreateUser)
	r.GET("/timesheets/:bulan/:tahun/:periode", controllers.GetMergedDatass)
}
