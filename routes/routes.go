package routes

import (
	controllers "api-rect-go/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
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
	r.GET("/jobholder", controllers.GetMasterDataAvailableWithFormss)
	r.GET("/timesheets-customer", controllers.GetMergedData)
	r.GET("/request-driver", controllers.GetRequestDrivers)
	r.GET("/cuti-driver", controllers.GetReplacementData)
	r.POST("/request-driver", controllers.CreateRequestDriverController)	
}
