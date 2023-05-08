package handler

import (
	"collection-format/database"
	"collection-format/model"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// query all Collection
func GetAllCollection(c *fiber.Ctx) error {
	db := database.DB
	var collection []model.Info
	db.Find(&collection)
	return c.JSON(collection)
}

// CreateProduct new product
func CreateCollection(c *fiber.Ctx) error {
	db := database.DB

	NewCollectionInfo := model.Info{
		ID:   uuid.NewString(),
		Name: "New Collection",
	}
	err := db.Create(&NewCollectionInfo).Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create Collection", "data": nil})
	}
	return c.JSON(NewCollectionInfo)
}

// Add folder for collection
func AddFolder(c *fiber.Ctx) error {
	db := database.DB
	type NewFolderData struct {
		ID string `json:"id"`
	}
	var NF NewFolderData
	if err := c.BodyParser(&NF); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	var collection model.Collection
	err := db.Where("info_id = ?", NF.ID).First(&collection).Error
	if err != nil {
		NewCltItem := model.Collection{
			ID:     uuid.NewString(),
			InfoID: NF.ID,
		}
		err := db.Create(&NewCltItem).Error
		if err != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create Collection Item", "data": nil})
		}
		NewFolder := model.Item{
			ID:           uuid.NewString(),
			Name:         "New Folder",
			CollectionID: NewCltItem.ID,
			IsFolder:     true,
			Folders:      []model.Folder{},
			Examples:     []model.Example{},
		}

		erro := db.Create(&NewFolder).Error
		if erro != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create Folder", "data": nil})
		}
		return c.JSON(NewFolder)
	}

	NewFolder := model.Item{
		ID:           uuid.NewString(),
		Name:         "New Folder",
		CollectionID: collection.ID,
		IsFolder:     true,
		Folders:      []model.Folder{},
		Examples:     []model.Example{},
	}

	erro := db.Create(&NewFolder).Error
	if erro != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create Folder", "data": nil})
	}
	return c.JSON(NewFolder)
}

// Add request for collection + folder
func AddRequest(c *fiber.Ctx) error {
	db := database.DB
	type postID struct {
		ID string `json:"id"`
	}
	var pID postID
	if err := c.BodyParser(&pID); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	//level folder
	var folder model.Item
	erro := db.Where("id = ?", pID.ID).First(&folder).Error
	if erro != nil {
		//level item
		var collection model.Collection
		err := db.Where("info_id = ?", pID.ID).First(&collection).Error
		if err != nil {
			NewCltItem := model.Collection{
				ID:     uuid.NewString(),
				InfoID: pID.ID,
			}
			err := db.Create(&NewCltItem).Error
			if err != nil {
				return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create Collection Item", "data": nil})
			}
			//Create Request
			NewRequest := model.Item{
				ID:           uuid.NewString(),
				Name:         "New Request",
				CollectionID: NewCltItem.ID,
				Examples:     []model.Example{},
			}

			erro := db.Create(&NewRequest).Error
			if erro != nil {
				return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create Request", "data": nil})
			}
			//Create Detail for Rq
			NewRequestDt := model.Request{
				ID:       uuid.NewString(),
				ParentID: NewRequest.ID,
			}

			errors := db.Create(&NewRequestDt).Error
			if errors != nil {
				return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create RqData", "data": nil})
			}
			return c.JSON(NewRequest)
		}
		NewRequest := model.Item{
			ID:           uuid.NewString(),
			Name:         "New Request",
			CollectionID: collection.ID,
			Examples:     []model.Example{},
		}

		erro := db.Create(&NewRequest).Error
		if erro != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create Request(2)", "data": nil})
		}
		NewRequestDt := model.Request{
			ID:       uuid.NewString(),
			ParentID: NewRequest.ID,
		}

		errors := db.Create(&NewRequestDt).Error
		if errors != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create RqData(2)", "data": nil})
		}
		return c.JSON(NewRequest)
	}
	if !folder.IsFolder {
		return c.JSON(fiber.Map{"status": "error", "message": "Not a folder"})
	}
	NewRequest := model.Folder{
		ID:       uuid.NewString(),
		Name:     "New Request",
		ItemID:   pID.ID,
		Examples: []model.Example{},
	}

	error := db.Create(&NewRequest).Error
	if error != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create Request(3)", "data": nil})
	}
	NewRequestDt := model.Request{
		ID:       uuid.NewString(),
		ParentID: NewRequest.ID,
	}

	errors := db.Create(&NewRequestDt).Error
	if errors != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create RqData(3)", "data": nil})
	}
	return c.JSON(NewRequest)

}

// GetDetail of a Collection
func GetDetailCollection(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	fmt.Println(id)
	var collection model.Collection
	db.Where("info_id = ?", id).Preload("Items.Examples").Preload("Items.Folders.Examples").First(&collection)
	return c.JSON(collection)
}

// Update Collection name
func UpdateCollection(c *fiber.Ctx) error {
	db := database.DB

	type UpdateCollectionName struct {
		ID   string `json:"ID"`
		Name string `json:"Name"`
	}
	var ucn UpdateCollectionName
	if err := c.BodyParser(&ucn); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// fmt.Println(id)
	var info model.Info
	err := db.First(&info, "id = ?", ucn.ID).Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "No collection found with ID", "data": nil})
	}
	// Update the product with the new information
	erro := db.Table("infos").Where("ID = ?", ucn.ID).Updates(map[string]interface{}{"Name": ucn.Name}).Error
	if erro != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot update data", "data": info})
	}

	info.Name = ucn.Name

	return c.JSON(info)
}

// Update folder name
func UpdateFolder(c *fiber.Ctx) error {
	db := database.DB

	type UpdateFolderName struct {
		Id   string `json:"ID"`
		Name string `json:"Name"`
	}
	var ufn UpdateFolderName
	if err := c.BodyParser(&ufn); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	var item model.Item
	err := db.First(&item, "id = ?", ufn.Id).Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "No folder found with ID", "data": nil})
	}
	if !item.IsFolder {
		return c.JSON(fiber.Map{"status": "error", "message": "Not a folder"})
	}
	// Update the product with the new information
	erro := db.Table("items").Where("ID = ?", ufn.Id).Updates(map[string]interface{}{"Name": ufn.Name}).Error
	if erro != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot update data", "data": erro})
	}
	item.Name = ufn.Name
	return c.JSON(item)
}

// Update request name
func UpdateRequest(c *fiber.Ctx) error {
	db := database.DB

	type UpdateRequestName struct {
		Id   string `json:"ID"`
		Name string `json:"Name"`
	}
	var urn UpdateRequestName
	if err := c.BodyParser(&urn); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	var request model.Folder
	errorRQ := db.First(&request, "id = ?", urn.Id).Error
	if errorRQ != nil {
		//level Item
		var item model.Item
		err := db.First(&item, "id = ?", urn.Id).Error
		if err != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "No folder found with ID", "data": nil})
		}
		if item.IsFolder {
			return c.JSON(fiber.Map{"status": "error", "message": "Not a request"})
		}
		// Update the product with the new information
		erro := db.Table("items").Where("ID = ?", urn.Id).Updates(map[string]interface{}{"Name": urn.Name}).Error
		if erro != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot update data", "data": erro})
		}
		item.Name = urn.Name
		return c.JSON(item)
	}
	//level Folder
	erro := db.Table("folders").Where("ID = ?", urn.Id).Updates(map[string]interface{}{"Name": urn.Name}).Error
	if erro != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot update data", "data": erro})
	}
	request.Name = urn.Name
	return c.JSON(request)
}

func UpdateExample(c *fiber.Ctx) error {
	db := database.DB

	type UpdateExampleName struct {
		Id   string `json:"ID"`
		Name string `json:"Name"`
	}
	var uen UpdateExampleName
	if err := c.BodyParser(&uen); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	var Examp model.Example
	err := db.First(&Examp, "id = ?", uen.Id).Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "No folder found with ID", "data": nil})
	}
	// Update the product with the new information
	erro := db.Table("examples").Where("ID = ?", uen.Id).Updates(map[string]interface{}{"Name": uen.Name}).Error
	if erro != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot update data", "data": erro})
	}
	Examp.Name = uen.Name
	return c.JSON(Examp)
}

//DELETE

func DeleteExample(c *fiber.Ctx, ID string) (*model.Example, error) {
	db := database.DB
	var Examp model.Example
	if err := db.First(&Examp, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	if err := db.Where("id = ?", ID).Delete(&Examp).Error; err != nil {
		return nil, err
	}
	return &Examp, nil
}

func DeleteExamp(c *fiber.Ctx) error {
	type ExampID struct {
		ExampID string `json:"ID"`
	}
	var ExID ExampID
	if err := c.BodyParser(&ExID); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	examp, err := DeleteExample(c, ExID.ExampID)
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Error deleting example", "data": err})
	}
	DeleteResponseData(c, examp.ID)
	return c.JSON(examp)
}

func DeleteResponseData(c *fiber.Ctx, ID string) error {
	db := database.DB
	var RpDt model.Response
	if err := db.First(&RpDt, "example_id = ?", ID).Error; err != nil {
		return err
	}
	if err := db.Where("example_id = ?", ID).Delete(&RpDt).Error; err != nil {
		return err
	}
	errDetele := DeleteRequestData(c, RpDt.ID)
	if errDetele != nil {
		return errDetele
	}
	return nil
}

func DeleteRequestData(c *fiber.Ctx, ParentID string) error {
	db := database.DB
	var RqDt model.Request
	if err := db.Where("parent_id = ?", ParentID).Delete(&RqDt).Error; err != nil {
		return err
	}
	return nil
}

// Delete Request
func DeleteRequest(c *fiber.Ctx) error {
	db := database.DB

	type DelRequest struct {
		ID string `json:"ID"`
	}
	var DelRQ DelRequest
	if err := c.BodyParser(&DelRQ); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	//Folder level
	var request model.Folder
	errorRQ := db.First(&request, "id = ?", DelRQ.ID).Error
	if errorRQ != nil {
		//Collection level
		var requests model.Item
		err := db.First(&requests, "id = ?", DelRQ.ID).Error
		if err != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "No request found with ID", "data": nil})
		}
		if requests.IsFolder {
			return c.JSON(fiber.Map{"status": "error", "message": "Not a request"})
		}
		var examp []model.Example
		db.Where("item_id = ?", requests.ID).Find(&examp)
		for i := range examp {
			DeleteResponseData(c, examp[i].ID)
		}
		errorF := db.Where("id = ?", requests.ID).Delete(&requests).Error
		if errorF != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot delete request", "data": nil})
		}
		//del rqdt
		errDetele := DeleteRequestData(c, requests.ID)
		if errDetele != nil {
			return errDetele
		}
		return c.JSON(requests)
	}
	var examp []model.Example
	db.Where("folder_id = ?", request.ID).Find(&examp)
	for i := range examp {
		// db.Delete(&examp[i])
		DeleteResponseData(c, examp[i].ID)
	}
	errorF := db.Where("id = ?", request.ID).Delete(&request).Error
	if errorF != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot delete request", "data": nil})
	}
	//del rqdt
	errDetele := DeleteRequestData(c, request.ID)
	if errDetele != nil {
		return errDetele
	}
	return c.JSON(request)
}

// DeleteProduct delete product
func DeleteFolder(c *fiber.Ctx) error {
	db := database.DB

	type DelFolder struct {
		ID string `json:"ID"`
	}
	var DelF DelFolder
	if err := c.BodyParser(&DelF); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	var folder model.Item
	err := db.First(&folder, "id = ?", DelF.ID).Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "No folder found with ID", "data": nil})
	}
	if !folder.IsFolder {
		return c.JSON(fiber.Map{"status": "error", "message": "Not a folder"})
	}
	var Request []model.Folder
	errorFindRq := db.Where("item_id = ?", folder.ID).Find(&Request).Error
	if errorFindRq == nil {
		var RqDt model.Request
		var examp []model.Example
		for i := range Request {
			db.Where("folder_id = ?", Request[i].ID).Find(&examp)
			for i := range examp {
				DeleteResponseData(c, examp[i].ID)
			}
			db.Where("parent_id = ?", Request[i].ID).Delete(&RqDt)
		}
	}
	errorF := db.Select("Folders").Where("id = ?", folder.ID).Delete(&folder).Error
	if errorF != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot delete folder", "data": nil})
	}
	return c.JSON(folder)
}

// DeleteProduct delete product
func DeleteCollection(c *fiber.Ctx) error {
	type DelClt struct {
		ID string `json:"ID"`
	}
	var DelC DelClt
	if err := c.BodyParser(&DelC); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	db := database.DB

	var info model.Info
	// db.First(&user, "id = ?", "1b74413f-f3b8-409f-ac47-e8c062e3472a")
	err := db.First(&info, "id = ?", DelC.ID).Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "No collection found with ID", "data": nil})
	}
	//Delete Info Collection
	er := db.Delete(&info).Error
	if er != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot delete info", "data": nil})
	}
	//Delete Collection Item
	var collection model.Collection
	erro := db.Where("info_id = ?", DelC.ID).First(&collection).Error
	if erro != nil {
		return c.JSON(info)
	}
	var items []model.Item
	errorFindItem := db.Where("collection_id = ?", collection.ID).Find(&items).Error
	if errorFindItem == nil {
		var folders []model.Folder
		var RqDt model.Request
		for i := range items {
			if !items[i].IsFolder {
				var examp []model.Example
				db.Where("item_id = ?", items[i].ID).Find(&examp)
				for i := range examp {
					DeleteResponseData(c, examp[i].ID)
				}
				db.Where("parent_id = ?", items[i].ID).Delete(&RqDt)
				continue
			}
			err := db.Where("item_id = ?", items[i].ID).Find(&folders).Error
			if err != nil {
				continue
			}
			for i := range folders {
				var examp []model.Example
				db.Where("folder_id = ?", folders[i].ID).Find(&examp)
				for i := range examp {
					DeleteResponseData(c, examp[i].ID)
				}
				db.Where("parent_id = ?", folders[i].ID).Delete(&RqDt)
			}
		}

	}
	errorI := db.Select("Items").Delete(&collection).Error
	if errorI != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot delete Collection item", "data": nil})
	}

	return c.JSON(info)
}
