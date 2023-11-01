package handlers

import (
	"fmt"
	"fortest/core"
	"fortest/models"
	"fortest/router"
	"log"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	Core *core.Core
}

func New(core *core.Core) *Controller {
	return &Controller{Core: core}
}

func (c *Controller) AddApartment(ctx *router.Context) {
	var apartmentsResponse models.ApartmentsResponse
	if err := ctx.ShouldBind(&apartmentsResponse); err != nil {
		ctx.BadRequest(err)
		return
	}

	objectID, err := strconv.Atoi(ctx.Param("object_id"))
	if err!=nil{
		ctx.BadRequest(err)
		return
	}

	for i := range apartmentsResponse.Apartments {
		apartmentsResponse.Apartments[i].ObjectID = uint(objectID)
	}

	app, err := c.Core.NewAppartmentService.AddApartment(apartmentsResponse)
	if err != nil {
		ctx.Internal(err)
		return
	}

	ctx.OK(app)
}

func (c *Controller) AddDeveloper(ctx *router.Context) {
	var developer models.Developer
	if err := ctx.ShouldBind(&developer); err != nil {
		ctx.BadRequest(err)
		return
	}

	dev, err := c.Core.NewBuilderService.AddDeveloper(developer)
	if err != nil {
		ctx.Internal(err)
		return
	}

	ctx.OK(dev)
}

func (c *Controller) GetObjectsByDeveloperID(ctx *router.Context) {
	devIDStr := ctx.Param("id")
	devID, err := strconv.Atoi(devIDStr)

	if err != nil {
		ctx.BadRequest("Invalid developer ID")
		return
	}

	response, err := c.Core.NewBuilderService.GetDeveloperWithObjectsByID(uint(devID))
	if err != nil {
		ctx.Internal(err)
		return
	}

	ctx.OK(response)
}

func (c *Controller) GetDeveloperByIDHandler(ctx *router.Context) {
	developerIDStr := ctx.Param("id")
	developerID, err := strconv.Atoi(developerIDStr)
	if err != nil {
		ctx.BadRequest("Invalid developer ID")
		return
	}

	developer, err := c.Core.NewBuilderService.GetDeveloperByID(uint(developerID))
	if err != nil {
		ctx.NotFound("Developer not found")
		return
	}

	ctx.OK(developer)
}

func (c *Controller) AddObject(ctx *router.Context) {
	var object models.Object
	if err := ctx.ShouldBind(&object); err != nil {
		ctx.BadRequest(err)
		return
	}

	// Получение developer_id из параметра пути
	developerID := ctx.Param("developer_id")
	intID, err := strconv.Atoi(developerID)
	if err != nil {
		ctx.BadRequest("Invalid developer ID")
		return
	}
	object.DeveloperID = uint(intID)

	if err := c.Core.NewObjectService.CreateObject(&object); err != nil {
		ctx.Internal(err)
		return
	}

	ctx.OK(object)
}

func (c *Controller) DeleteObjectByID(ctx *router.Context) {
	objectIDStr := ctx.Param("id")
	objectID, err := strconv.Atoi(objectIDStr)
	if err != nil {
		ctx.BadRequest(fmt.Errorf("invalid object ID format. Error: %v", err))
		return
	}

	if err := c.Core.NewObjectService.DeleteObjectByID(objectID); err != nil {
		ctx.Internal(fmt.Errorf("error deleting object: %v", err))
		return
	}

	ctx.OK("Object deleted successfully")
}

func (c *Controller) GetAllObjectsHandler(ctx *router.Context) {
	var listLength int64
	total := GeneratePaginationFromRequest(ctx.Request)
	TotalPages, err := TotalPageTasks(int64(total.Limit), listLength)
	if err != nil {
		log.Println(err)
		ctx.BadRequest(err)
		return
	}
	var filter models.ObjectFilter

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.BadRequest(err)
		return
	}

	objects, apartments, count, err := c.Core.NewObjectService.GetAllObjects(filter, &total)
	if err != nil {
		ctx.Internal(err)
		return
	}
	if count == 0 {
		ctx.OK("No ads found")
	}

	var response models.FinalResponse
	for _, obj := range objects {
		var objResp models.ObjectResponse
		objResp.Name = obj.Name
		objResp.Address = obj.Address
		objResp.PhoneNumb = obj.Phone
		objResp.Photos = strings.Join(obj.Photos, ",")
		objResp.Status = obj.Status

		var types []models.ApartmentType
		for _, apt := range apartments {
			if apt.ObjectID == obj.ID {
				types = append(types, models.ApartmentType{
					Type:      apt.Type,
					MinSquare: fmt.Sprintf("%.2f", apt.Area),
					MinPrice:  fmt.Sprintf("%.2f", apt.Price),
				})
			}
		}

		objResp.Types = types
		response.Objects = append(response.Objects, objResp)
	}

	total.Records = response
	total.TotalPages = TotalPages

	ctx.OK(total)
}

func (c *Controller) GetObject(ctx *router.Context) {
	objectID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.BadRequest(err)
		return
	}

	object, err := c.Core.NewObjectService.GetObjectByID(objectID)
	if err != nil {
		ctx.Internal(err)
		return
	}
	ctx.OK(object)
}

func (c *Controller) GetPopularObjectsHandler(ctx *router.Context) {
	var listLenght int64
	total := GeneratePaginationFromRequest(ctx.Request)
	TotalPages, err := TotalPageTasks(int64(total.Limit), listLenght)
	if err != nil {
		log.Println(err)
		ctx.BadRequest(err)
		return
	}

	objects, apartments, count, err := c.Core.NewObjectService.GetPopularObjects(&total)
	if err != nil {
		ctx.Internal(err)
		return
	}

	if count == 0 {
		ctx.OK("No objects found")
		return
	}

	var response models.FinalResponse

	for _, obj := range objects {
		var objResp models.ObjectResponse
		objResp.Name = obj.Name
		objResp.Address = obj.Address
		objResp.PhoneNumb = obj.Phone
		objResp.Photos = strings.Join(obj.Photos, ",")
		objResp.Status = obj.Status

		var types []models.ApartmentType

		for _, apt := range apartments {
			if apt.ObjectID == obj.ID {
				types = append(types, models.ApartmentType{
					Type:      apt.Type,
					MinSquare: fmt.Sprintf("%.2f", apt.Area),
					MinPrice:  fmt.Sprintf("%.2f", apt.Price),
				})
			}
		}

		objResp.Types = types
		response.Objects = append(response.Objects, objResp)
	}

	total.Records = response
	total.TotalPages = TotalPages
	ctx.OK(total)
}

func (c *Controller) UpdateObjectWithApartmentsHandler(ctx *router.Context) {
	var updatedObject models.Object

	if err := ctx.ShouldBind(&updatedObject); err != nil {
		ctx.BadRequest(err)
		return
	}

	if err := c.Core.NewObjectService.UpdateObjectWithApartments(&updatedObject); err != nil {
		ctx.Internal(err.Error())
		return
	}

	ctx.OK("updated")
}

func (c *Controller) GetDeveloperObjectsSummary(ctx *router.Context) {
	devIDStr := ctx.Param("id")
	devID, err := strconv.Atoi(devIDStr)

	if err != nil {
		ctx.BadRequest("Invalid developer ID")
		return
	}

	response, err := c.Core.NewObjectService.GetDeveloperObjectsSummaryByID(uint(devID))
	if err != nil {
		ctx.Internal(err)
		return
	}

	ctx.OK(response)
}

func (c *Controller) UploadLayout(ctx *router.Context) {
	//err := ctx.Request.ParseMultipartForm(64 << 20)
	//if err != nil {
	//	ctx.BadRequest(err)
	//	return
	//}
	//
	//fileHeaders := ctx.Request.MultipartForm.File["upload[]"]
	//var filenames []string
	//
	//for _, fileHeader := range fileHeaders {
	//	newFilename := generateRandomFilename(fileHeader.Filename)
	//	file, err := fileHeader.Open()
	//	if err != nil {
	//		ctx.BadRequest(err)
	//		return
	//	}
	//	_, err = utils.UploadImage(file, newFilename, "."+constant.UploadDirNewBuildingLayouts)
	//	if err != nil {
	//		ctx.Internal(err)
	//		return
	//	}
	//	filenames = append(filenames, newFilename)
	//}
	//
	//ctx.OK(fmt.Sprintf("new_filename: %v", filenames))

}

func (c *Controller) UploadObjectPhotos(ctx *router.Context) {
	//err := ctx.Request.ParseMultipartForm(64 << 20)
	//if err != nil {
	//	ctx.BadRequest(err)
	//	return
	//}
	//
	//fileHeaders := ctx.Request.MultipartForm.File["upload[]"]
	//var filenames []string
	//
	//for _, fileHeader := range fileHeaders {
	//	newFilename := generateRandomFilename(fileHeader.Filename)
	//	file, err := fileHeader.Open()
	//	if err != nil {
	//		ctx.BadRequest(err)
	//		return
	//	}
	//	_, err = utils.UploadImage(file, newFilename, "."+constant.UploadDirNewBuildingObjects)
	//	if err != nil {
	//		ctx.Internal(err)
	//		return
	//	}
	//	filenames = append(filenames, newFilename)
	//}
	//
	//ctx.OK(fmt.Sprintf("new_filename: %v", filenames))

}

func (c *Controller) UploadDevLogo(ctx *router.Context) {
	//err := ctx.Request.ParseMultipartForm(64 << 20)
	//if err != nil {
	//	ctx.BadRequest(err)
	//	return
	//}
	//
	//file, fileHeader, err := ctx.Request.FormFile("file")
	//if err != nil {
	//	ctx.BadRequest(err)
	//	return
	//}
	//defer file.Close()
	//
	////	newFilename := generateRandomFilename(fileHeader.Filename)
	////filename, err := utils.UploadImage(file, newFilename, "."+constant.UploadDirNewBuildingLogo)
	//if err != nil {
	//	ctx.Internal(err)
	//	return
	//}
	////ctx.OK(fmt.Sprintf("new_filename: %v", filename))
}

func generateRandomFilename(filename string) string {
	extension := filepath.Ext(filename)
	rand.Seed(time.Now().UnixNano())
	randomName := fmt.Sprintf("%d%s", rand.Int(), extension)
	return randomName
}
