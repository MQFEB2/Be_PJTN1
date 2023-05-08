package handler

import (
	"collection-format/database"
	"collection-format/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateProduct new product
func ImportCollection(c *fiber.Ctx) error {
	db := database.DB

	type Collection struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	var CollectionImport Collection
	if err := c.BodyParser(&CollectionImport); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	var collection model.Info
	erro := db.First(&collection, "id = ?", CollectionImport.ID).Error
	if erro == nil {
		NewCollectionInfo := model.Info{
			ID:   uuid.NewString(),
			Name: CollectionImport.Name,
		}
		err := db.Create(&NewCollectionInfo).Error
		if err != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create Collection", "data": nil})
		}
		return c.JSON(NewCollectionInfo)
	}

	NewCollectionInfo := model.Info{
		ID:   CollectionImport.ID,
		Name: CollectionImport.Name,
	}
	err := db.Create(&NewCollectionInfo).Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create Collection", "data": nil})
	}
	return c.JSON(NewCollectionInfo)
}

// Add folder for collection
func ImportFolder(c *fiber.Ctx) error {
	db := database.DB
	type Folder struct {
		ID   string `json:"InfoId"`
		Name string `json:"name"`
	}
	var NF Folder
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
			Name:         NF.Name,
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
		Name:         NF.Name,
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
func ImportRequest(c *fiber.Ctx) error {
	db := database.DB
	type postID struct {
		ID   string `json:"ParentId"`
		Name string `json:"name"`
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
				Name:         pID.Name,
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
			Name:         pID.Name,
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
		Name:     pID.Name,
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

func ImportExample(c *fiber.Ctx) error {
	db := database.DB

	type Example struct {
		ParentID string `json:"ParentID"`
		Name     string `json:"name"`
	}

	var Ex Example
	if err := c.BodyParser(&Ex); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	var item model.Item
	err := db.First(&item, "id = ?", Ex.ParentID).Error
	if err != nil {
		NewExample := model.Example{
			ID:       uuid.NewString(),
			Name:     Ex.Name,
			FolderID: Ex.ParentID,
		}
		errors := db.Select("ID", "Name", "FolderID").Create(&NewExample).Error
		if errors != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create"})
		}
		return c.JSON(NewExample)

	} else {
		NewExample := model.Example{
			ID:     uuid.NewString(),
			Name:   Ex.Name,
			ItemID: Ex.ParentID,
		}
		errors := db.Select("ID", "Name", "ItemID").Create(&NewExample).Error
		if errors != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create ex1"})
		}
		return c.JSON(NewExample)
	}
}

func ImportSaveRequest(c *fiber.Ctx) error {
	db := database.DB

	type SaveRQ struct {
		ParentID string `json:"ParentID"`
		Method   string `json:"Method"`
	}
	var RQ SaveRQ
	if err := c.BodyParser(&RQ); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// fmt.Println(id)
	var Rquest model.Request
	err := db.First(&Rquest, "parent_id = ?", RQ.ParentID).Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cant find request"})
	}
	// Update the product with the new information
	erro := db.Table("requests").Where("ID = ?", Rquest.ID).Updates(map[string]interface{}{"Method": RQ.Method}).Error
	if erro != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Save request", "data": nil})
	}

	Rquest.Method = RQ.Method

	return c.JSON(Rquest)
}
