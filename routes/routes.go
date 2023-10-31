package routes

import (
	"fortest/handlers"
	"fortest/router"
)

func Init(c *handlers.Controller) *router.Router {
	r := router.NewRouter().Group("spec")

	r.StrictSlash(true)

	dashboard := r.Group("/dashboard")
	{

		newBuilding := dashboard.Group("/new_building")
		{
			newBuilding.GET("/announcements", c.GetAllObjectsHandler)             //DONE
			newBuilding.GET("/popular_announcements", c.GetPopularObjectsHandler) //DONE

			newBuildingObjects := newBuilding.Group("/objects")
			{
				newBuildingObjects.POST("/object_photos", c.UploadObjectPhotos)      //DONE надо папку создать
				newBuildingObjects.POST("/layout", c.UploadLayout)                   //DONE надо папку создать
				newBuildingObjects.GET("/{id}", c.GetObject)                         //DONE
				newBuildingObjects.DELETE("/{id}", c.DeleteObjectByID)               //DONE
				newBuildingObjects.PUT("/{id}", c.UpdateObjectWithApartmentsHandler) //failed:id
				newBuildingObjects.POST("/{object_id}/apartments", c.AddApartment)   //failed:id
			}

			newBuildingDevelopers := newBuilding.Group("/developers")
			{
				newBuildingDevelopers.POST("/logo", c.UploadDevLogo)                             //DONE надо папку создать
				newBuildingDevelopers.POST("", c.AddDeveloper)                                   //DONE
				newBuildingDevelopers.POST("/{developer_id}/objects", c.AddObject)               //DONE
				newBuildingDevelopers.GET("/{id}", c.GetDeveloperByIDHandler)                    //надо починить (нет всех данных)
				newBuildingDevelopers.GET("/{id}/objects", c.GetObjectsByDeveloperID)            //надо починить (излишние данные)
				newBuildingDevelopers.GET("/{id}/objects_summary", c.GetDeveloperObjectsSummary) //DONE
			}
		}
	}
	return r
}
