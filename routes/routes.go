package routes

import (
	"fortest/handlers"
	"fortest/router"
)

func Init(c *handlers.Controller) *router.Router {
	r := router.NewRouter().Group("spec")

	r.SetValidator(c.Core.Validator)
	r.StrictSlash(true)

	dashboard := r.Group("/dashboard")
	{

		newBuilding := dashboard.Group("/newBuilding")
		{
			newBuilding.POST("/dev_logo", c.UploadDevLogo)
			newBuilding.POST("/object_photos", c.UploadObjectPhotos)
			newBuilding.POST("/layout", c.UploadLayout)
			newBuilding.GET("/announcements", c.GetAllObjectsHandler)
			newBuilding.GET("/popular_announcements", c.GetPopularObjectsHandler)
			newBuilding.GET("/object/{id}", c.GetObject)

			newBuildingObjects := newBuilding.Group("/objects")
			{
				newBuildingObjects.DELETE("/{id}", c.DeleteObjectByID)
				newBuildingObjects.PUT("/{id}", c.UpdateObjectWithApartmentsHandler)
				newBuildingObjects.POST("/{object_id}/apartments", c.AddApartment)
			}

			newBuildingDevelopers := newBuilding.Group("/developers")
			{
				newBuildingDevelopers.POST("", c.AddDeveloper)
				newBuildingDevelopers.POST("/{developer_id}/objects", c.AddObject)
				newBuildingDevelopers.GET("/developers/{id}", c.GetDeveloperByIDHandler)
			}

			newBuildingDeveloper := newBuilding.Group("/developer")
			{
				newBuildingDeveloper.GET("/{id}/objects", c.GetObjectsByDeveloperID)
				newBuildingDeveloper.GET("/{id}/objects_summary", c.GetDeveloperObjectsSummary)
			}
		}
	}
	return r
}
