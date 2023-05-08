package handler

import (
	"collection-format/database"
	"collection-format/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func FindCollection(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	var info model.Info
	fmt.Println(id)
	db.Where("id = ?", id).First(&info)
	var collection model.Collection
	db.Where("info_id = ?", info.ID).Preload("Items.Examples").Preload("Items.Folders.Examples").First(&collection)

	Info := model.EpInfo{
		ID:     info.ID,
		Name:   info.Name,
		Schema: "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
	}

	Collection := model.EpCollection{
		EpInfo: Info,
	}
	if len(collection.Items) == 0 {
		return c.JSON(collection)
	}

	for i := range collection.Items {
		if !collection.Items[i].IsFolder {
			Item := model.EpItemRequest{
				Name: collection.Items[i].Name,
				Protocol: map[string]interface{}{
					"disableBodyPruning": true,
				},
			}
			var Rquest model.Request
			err := db.Where("parent_id = ?", collection.Items[i].ID).Preload("Body").Preload("Headers").Preload("Url.Queries").First(&Rquest).Error
			if err != nil {
				return err
			}
			request := MountRqData(Rquest)
			Item.Request = request
			// if len(Item.Request.Method) == 0 {
			// 	Item.Request = model.EpRequest{}
			// }
			for j := range collection.Items[i].Examples {
				var ex model.Response
				var Rquest model.Request
				db.Where("example_id = ?", collection.Items[i].Examples[j].ID).First(&ex)
				db.Where("parent_id = ?", ex.ID).Preload("Body").Preload("Headers").Preload("Url.Queries").First(&Rquest)
				request := MountRqData(Rquest)
				response := model.EpResponse{
					Name:     collection.Items[i].Examples[j].Name,
					Request:  request,
					Status:   "200",
					Code:     "OK",
					Language: "json",
					Body:     ex.Body,
				}
				Item.Response = append(Item.Response, response)

			}
			if len(Item.Response) == 0 {
				Item.Response = []model.EpResponse{}
			}
			Collection.EpItems = append(Collection.EpItems, Item)
			continue
		}

		Folder := model.EpItemFolder{
			Name: collection.Items[i].Name,
		}

		for j := range collection.Items[i].Folders {
			ItemFolder := model.EpItemRequest{
				Name: collection.Items[i].Folders[j].Name,
				Protocol: map[string]interface{}{
					"disableBodyPruning": true,
				},
			}
			var Rquest model.Request
			err := db.Where("parent_id = ?", collection.Items[i].Folders[j].ID).Preload("Body").Preload("Headers").Preload("Url.Queries").First(&Rquest).Error
			if err != nil {
				return err
			}
			request := MountRqData(Rquest)
			ItemFolder.Request = request
			for k := range collection.Items[i].Folders[j].Examples {
				var ex model.Response
				var Rquest model.Request
				db.Where("example_id = ?", collection.Items[i].Folders[j].Examples[k].ID).First(&ex)
				db.Where("parent_id = ?", ex.ID).Preload("Body").Preload("Headers").Preload("Url.Queries").First(&Rquest)
				request := MountRqData(Rquest)
				response := model.EpResponse{
					Name:     collection.Items[i].Folders[j].Examples[k].Name,
					Request:  request,
					Status:   "200",
					Code:     "OK",
					Language: "json",
					Body:     ex.Body,
				}
				ItemFolder.Response = append(ItemFolder.Response, response)
			}
			if len(ItemFolder.Response) == 0 {
				ItemFolder.Response = []model.EpResponse{}
			}
			Folder.Folder = append(Folder.Folder, ItemFolder)
		}
		if len(Folder.Folder) == 0 {
			Folder.Folder = []model.EpItemRequest{}
		}
		Collection.EpItems = append(Collection.EpItems, Folder)

	}

	// return c.JSON(Collection)

	// Tạo đường dẫn đến thư mục "data"
	dirPath := "./file_export"

	// Kiểm tra thư mục có tồn tại hay không, nếu chưa thì tạo mới
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			return errors.New("error creating directory")
		}
	}
	// Tạo đường dẫn đến file mới
	filePath := dirPath + "/" + info.Name + ".json"

	// Kiểm tra file có tồn tại hay không, nếu có thì không tạo mới và trả về lỗi
	// if _, err := os.Stat(filePath); !os.IsNotExist(err) {
	// 	return errors.New("file already exists")
	// }

	// Tạo tệp mới
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Viết nội dung vào tệp
	jsonBytes, err := json.MarshalIndent(Collection, "", "	")
	if err != nil {
		return err
	}

	_, err = file.Write(jsonBytes)
	if err != nil {
		return err
	}

	fmt.Println(filePath)
	errDown := c.Download(filePath, info.Name+".json")
	if errDown != nil {
		return errDown
	}
	return c.SendFile(filePath)
}

func MountRqData(Rquest model.Request) model.EpRequest {
	Raw := Rquest.Url.Raw
	Url := *ParseUrl(Raw)

	Body := model.EpBody{
		Mode: "raw",
		Raw:  Rquest.Body.Raw,
		Options: map[string]interface{}{
			"raw": map[string]interface{}{
				"language": "json",
			},
		},
	}

	Request := model.EpRequest{
		Method:    Rquest.Method,
		EpHeaders: []model.EpHeader{},
		EpBody:    &Body,
		EpUrl:     &Url,
	}

	for _, header := range Rquest.Headers {
		HD := model.EpHeader{
			Key:         header.Key,
			Value:       header.Value,
			Description: header.Description,
			Disabled:    header.Disabled,
		}
		Request.EpHeaders = append(Request.EpHeaders, HD)
	}

	return Request
}

func ParseUrl(u string) *model.EpUrl {
	// Parse the URL
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil
	}

	// Create a new EpUrl struct
	epUrl := model.EpUrl{
		Raw:      u,
		Protocol: parsedURL.Scheme,
	}

	if parsedURL.Port() != "" {
		epUrl.Port = parsedURL.Port()
	}

	path := strings.Split(parsedURL.Path, "/")
	for _, part := range path {
		if part != "" {
			epUrl.Path = append(epUrl.Path, part)
		}
	}

	var queryArray []model.EpQuery

	// Convert the query parameters to a JSON array
	// Convert the query parameters to a JSON array
	for key, values := range parsedURL.Query() {
		for _, value := range values {
			item := model.EpQuery{
				Key:   key,
				Value: value,
			}
			queryArray = append(queryArray, item)
		}
	}

	epUrl.EpQueries = append(epUrl.EpQueries, queryArray...)
	// Split the hostname into its parts and add to the EpUrl
	hostParts := strings.Split(parsedURL.Hostname(), ".")
	for _, part := range hostParts {
		epUrl.Host = append(epUrl.Host, part)
	}

	// Return a pointer to the EpUrl struct
	return &epUrl
}
