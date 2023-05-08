package handler

import (
	"collection-format/database"
	"collection-format/model"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetRequestDetail(c *fiber.Ctx) error {
	db := database.DB
	id := c.Query("id")

	var RequestDetail model.Request
	err := db.Where("parent_id = ?", id).Preload("Body").Preload("Headers").Preload("Url.Queries").First(&RequestDetail).Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot query request", "data": nil})
	}
	return c.JSON(RequestDetail)
}

func SaveRequest(c *fiber.Ctx) error {
	db := database.DB

	type SaveRQ struct {
		ID       string `json:"ID"`
		ParentID string `json:"ParentID"` //Không dùng
		Method   string `json:"Method"`
	}
	var RQ SaveRQ
	if err := c.BodyParser(&RQ); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// fmt.Println(id)
	var Rquest model.Request
	err := db.First(&Rquest, "id = ?", RQ.ID).Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cant find request"})
	}
	// Update the product with the new information
	erro := db.Table("requests").Where("ID = ?", RQ.ID).Updates(map[string]interface{}{"Method": RQ.Method}).Error
	if erro != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Save request", "data": nil})
	}

	Rquest.Method = RQ.Method

	return c.JSON(Rquest)
}

func SaveHeader(c *fiber.Ctx) error {
	db := database.DB

	type SaveHeader struct {
		ID          string `json:"ID"`
		RequestID   string `json:"RequestID"`
		Key         string `json:"Key"`
		Value       string `json:"Value"`
		Description string `json:"Description"`
		Disabled    bool   `json:"Disabled"`
	}
	var HDer SaveHeader
	if err := c.BodyParser(&HDer); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// fmt.Println(id)
	var HD model.Header
	err := db.First(&HD, "id = ?", HDer.ID).Error
	if err != nil {
		if HDer.RequestID == "" {
			return c.JSON(fiber.Map{"status": "error", "message": "No request found"})
		}
		if HDer.Key == "" && HDer.Value == "" && HDer.Description == "" {
			return c.JSON(fiber.Map{"status": "error", "message": "Fields are empty"})
		}
		NewHeader := model.Header{
			ID:          uuid.NewString(),
			RequestID:   HDer.RequestID,
			Key:         HDer.Key,
			Value:       HDer.Value,
			Description: HDer.Description,
			Disabled:    HDer.Disabled,
		}
		errors := db.Create(&NewHeader).Error
		if errors != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create Header"})
		}
		return c.JSON(NewHeader)
	}
	//delete if fields are empty
	if HDer.Key == "" && HDer.Value == "" && HDer.Description == "" {
		errorDel := db.Where("id = ?", HDer.ID).Delete(&HD).Error
		if errorDel != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot delete header empty", "data": errorDel})
		}
		return c.JSON(fiber.Map{"status": "error", "message": "Deleted empty header"})
	}
	// Update the product with the new information
	erro := db.Table("headers").Where("ID = ?", HDer.ID).Updates(map[string]interface{}{"Key": HDer.Key, "Value": HDer.Value, "Description": HDer.Description}).Error
	if erro != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Save request", "data": nil})
	}

	HD.Key = HDer.Key
	HD.Value = HDer.Value
	HD.Description = HDer.Description

	return c.JSON(HD)
}

func SaveBody(c *fiber.Ctx) error {
	db := database.DB

	type SaveBody struct {
		ID        string `json:"ID"`
		RequestID string `json:"RequestID"` // không dùng
		Raw       string `json:"Raw"`
	}
	var Body SaveBody
	if err := c.BodyParser(&Body); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review HDeryour input", "data": err})
	}

	// fmt.Println(id)
	var BD model.Body
	err := db.First(&BD, "id = ?", Body.ID).Error
	if err != nil {
		if Body.RequestID == "" {
			return c.JSON(fiber.Map{"status": "error", "message": "No RQ found with ID"})
		}
		if Body.ID != "" {
			return c.JSON(fiber.Map{"status": "error", "message": "ID found but RQ not found"})
		}
		NewBody := model.Body{
			ID:        uuid.NewString(),
			RequestID: Body.RequestID,
			Raw:       Body.Raw,
		}
		errors := db.Create(&NewBody).Error
		if errors != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create Body"})
		}
		return c.JSON(NewBody)
	}
	// Update the product with the new information
	erro := db.Table("bodies").Where("ID = ?", Body.ID).Updates(map[string]interface{}{"Raw": Body.Raw}).Error
	if erro != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cant Save body", "data": nil})
	}

	BD.Raw = Body.Raw

	return c.JSON(BD)
}

func SaveURL(c *fiber.Ctx) error {
	db := database.DB

	type SaveURL struct {
		ID        string `json:"ID"`
		RequestID string `json:"RequestID"`
		Raw       string `json:"Raw"`
	}
	var URL SaveURL
	if err := c.BodyParser(&URL); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// fmt.Println(id)
	var url model.Url
	err := db.First(&url, "id = ?", URL.ID).Error
	if err != nil {
		if URL.RequestID == "" {
			return c.JSON(fiber.Map{"status": "error", "message": "No RQ found with ID"})
		}
		if URL.ID != "" {
			return c.JSON(fiber.Map{"status": "error", "message": "ID found but RQ not found"})
		}
		NewURL := model.Url{
			ID:        uuid.NewString(),
			RequestID: URL.RequestID,
			Raw:       URL.Raw,
			Queries:   []model.Query{},
		}
		errors := db.Create(&NewURL).Error
		if errors != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create URL"})
		}
		return c.JSON(NewURL)
	}
	// Update the product with the new information
	erro := db.Table("urls").Where("ID = ?", URL.ID).Updates(map[string]interface{}{"Raw": URL.Raw}).Error
	if erro != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Save request", "data": nil})
	}

	url.Raw = URL.Raw

	return c.JSON(url)
}

func SaveQuery(c *fiber.Ctx) error {
	db := database.DB

	type SaveQuery struct {
		ID          string `json:"ID"`
		UrlID       string `json:"UrlID"`
		Key         string `json:"Key"`
		Value       string `json:"Value"`
		Description string `json:"Description"`
		Disabled    bool   `json:"Disabled"`
	}
	var Query SaveQuery
	if err := c.BodyParser(&Query); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// fmt.Println(id)
	var QR model.Query
	err := db.First(&QR, "id = ?", Query.ID).Error
	if err != nil {
		if Query.UrlID == "" {
			return c.JSON(fiber.Map{"status": "error", "message": "No Url found with ID"})
		}
		if Query.Key == "" && Query.Value == "" && Query.Description == "" {
			return c.JSON(fiber.Map{"status": "error", "message": "Fields are empty"})
		}
		NewQuery := model.Query{
			ID:          uuid.NewString(),
			UrlID:       Query.UrlID,
			Key:         Query.Key,
			Value:       Query.Value,
			Description: Query.Description,
		}
		errors := db.Create(&NewQuery).Error
		if errors != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create Header"})
		}
		return c.JSON(NewQuery)
		// return c.JSON(fiber.Map{"status": "error", "message": "No header found with ID", "data": nil})
	}
	//delete if fields are empty
	if Query.Key == "" && Query.Value == "" && Query.Description == "" {
		errorDel := db.Where("id = ?", Query.ID).Delete(&QR).Error
		if errorDel != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot delete query empty", "data": errorDel})
		}
		return c.JSON(fiber.Map{"status": "error", "message": "Deleted empty query"})
	}
	// Update the product with the new information
	erro := db.Table("headers").Where("ID = ?", Query.ID).Updates(map[string]interface{}{"Key": Query.Key, "Value": Query.Value, "Description": Query.Description}).Error
	if erro != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Save request", "data": nil})
	}

	QR.Key = Query.Key
	QR.Value = Query.Value
	QR.Description = Query.Description

	return c.JSON(QR)
}

func CreateExampleSidebar(c *fiber.Ctx) error {
	db := database.DB

	type Example struct {
		ParentID string `json:"ParentID"`
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
			Name:     "New Example",
			FolderID: Ex.ParentID,
		}
		errors := db.Select("ID", "Name", "FolderID").Create(&NewExample).Error
		if errors != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create"})
		}
		Response, err := CreateNewResponse(NewExample.ID, "")
		if err != nil {
			return err
		}
		Request, err := CreateNewRqDt(Response.ID)
		if err != nil {
			return err
		}
		CreateUrlMock(Request.ID, NewExample.ID)
		return c.JSON(NewExample)

	} else {
		NewExample := model.Example{
			ID:     uuid.NewString(),
			Name:   "New Example",
			ItemID: Ex.ParentID,
		}
		errors := db.Select("ID", "Name", "ItemID").Create(&NewExample).Error
		if errors != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create ex1"})
		}
		Response, err := CreateNewResponse(NewExample.ID, "")
		if err != nil {
			return err
		}
		Request, err := CreateNewRqDt(Response.ID)
		if err != nil {
			return err
		}
		CreateUrlMock(Request.ID, NewExample.ID)
		return c.JSON(NewExample)
	}
}

func CreateExample(c *fiber.Ctx) error {
	db := database.DB

	type Example struct {
		ParentID string `json:"ParentID"`
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
			Name:     "New Example",
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
			Name:   "New Example",
			ItemID: Ex.ParentID,
		}
		errors := db.Select("ID", "Name", "ItemID").Create(&NewExample).Error
		if errors != nil {
			return c.JSON(fiber.Map{"status": "error", "message": "Cannot Create ex1"})
		}
		return c.JSON(NewExample)
	}
}

func CreateNewResponse(ExampleID string, Body string) (model.Response, error) {
	db := database.DB
	NewResponse := model.Response{
		ID:        uuid.NewString(),
		ExampleID: ExampleID,
		Body:      Body,
	}
	erro := db.Create(&NewResponse).Error
	if erro != nil {
		return model.Response{}, erro
	}
	return NewResponse, nil
}

func CreateUrlMock(RequestID string, MockID string) (model.Url, error) {
	db := database.DB
	NewUrl := model.Url{
		ID:        uuid.NewString(),
		RequestID: RequestID,
		Raw:       "http://127.0.0.1:8000/rp/body/" + MockID,
	}
	if errors := db.Create(&NewUrl).Error; errors != nil {
		return model.Url{}, errors
	}
	return NewUrl, nil
}

func CreateNewRqDt(ParentID string) (model.Request, error) {
	db := database.DB
	NewRequestDt := model.Request{
		ID:       uuid.NewString(),
		ParentID: ParentID,
		Headers:  []model.Header{},
	}
	if errors := db.Create(&NewRequestDt).Error; errors != nil {
		return model.Request{}, errors
	}
	return NewRequestDt, nil
}

func CreateRequestForResponse(c *fiber.Ctx) error {
	type NewExample struct {
		ExampleID string `json:"ParentID"`
		Body      string `json:"Body"`
	}
	var NE NewExample
	if err := c.BodyParser(&NE); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	Response, err := CreateNewResponse(NE.ExampleID, NE.Body)
	if err != nil {
		return err
	}
	Request, err := CreateNewRqDt(Response.ID)
	if err != nil {
		return err
	}
	return c.JSON(Request)
}

func GetResponse(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	fmt.Println(id)
	var ex model.Response
	err := db.Where("example_id = ?", id).First(&ex).Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot query request", "data": nil})
	}

	return c.JSON(ex)
}

func SaveResponse(c *fiber.Ctx) error {
	db := database.DB

	type Response struct {
		ExampleID string `json:"ExampleID"`
		Body      string `json:"Body"`
	}
	var Rp Response
	if err := c.BodyParser(&Rp); err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	var Rps model.Response
	err := db.First(&Rps, "example_id = ?", Rp.ExampleID).Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cant find request"})

	}
	// Update the product with the new information
	erro := db.Table("responses").Where("example_id = ?", Rp.ExampleID).Updates(map[string]interface{}{"Body": Rp.Body}).Error
	if erro != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Save request", "data": nil})
	}

	Rps.Body = Rp.Body

	return c.JSON(Rps)
}

func GetResponseDetail(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")

	var response model.Response
	errRp := db.Where("example_id = ?", id).First(&response).Error
	if errRp != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot query response", "data": nil})
	}

	var request model.Request
	err := db.Where("parent_id = ?", response.ID).First(&request).Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Cannot query request", "data": nil})
	}

	method := request.Method

	// Kiểm tra nếu phương thức không khớp với yêu cầu hiện tại
	if method != c.Method() {
		// Trả về mã trạng thái 405 Method Not Allowed và thông báo lỗi
		return c.Status(fiber.StatusMethodNotAllowed).SendString("Method Not Allowed")
	}

	var data interface{}
	errCv := json.Unmarshal([]byte(response.Body), &data)
	if errCv != nil {
		return c.JSON(response.Body)
	}

	return c.JSON(data)
}
