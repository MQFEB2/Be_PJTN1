package router

import (
	"collection-format/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {

	//Product
	clt := app.Group("/collection")

	//Get the product list
	clt.Get("/list.json", handler.GetAllCollection)

	//Create new product
	clt.Post("/new.json", handler.CreateCollection)

	//Create new folder
	clt.Post("/newfolder.json", handler.AddFolder)

	//Create new Request
	clt.Post("/newrequest.json", handler.AddRequest)

	//Get detail
	clt.Get("/detailcollection.json/:id", handler.GetDetailCollection)

	//Update Collection Name
	clt.Patch("/name.json", handler.UpdateCollection)

	//Update Folder name
	clt.Patch("/foldername.json", handler.UpdateFolder)

	//Update request name
	clt.Patch("/requestname.json", handler.UpdateRequest)

	//Update example name
	clt.Patch("/examplename.json", handler.UpdateExample)

	//Delete
	clt.Delete("/delete.json", handler.DeleteCollection)

	//Delete Folder
	clt.Delete("/delfolder.json", handler.DeleteFolder)

	//Delete Request
	clt.Delete("/delrequest.json", handler.DeleteRequest)

	clt.Delete("/delexamp.json", handler.DeleteExamp)

	//export
	clt.Get("/export.json/:id", handler.FindCollection)

	rq := app.Group("/request")

	rq.Get("/data.json", handler.GetRequestDetail)

	rq.Post("/save.json", handler.SaveRequest)

	rq.Post("/saveheader.json", handler.SaveHeader)

	rq.Post("/savebody.json", handler.SaveBody)

	rq.Post("/saveurl.json", handler.SaveURL)

	rq.Post("/savequery.json", handler.SaveQuery)

	rq.Post("/createexample.json", handler.CreateExample)

	rq.Post("/createexamplesidebar.json", handler.CreateExampleSidebar)

	rq.Get("/responsedata.json/:id", handler.GetResponse)

	rq.Post("/createrq.json", handler.CreateRequestForResponse)

	rq.Post("/saveresponse.json", handler.SaveResponse)

	rp := app.Group("/rp")

	rp.All("/body/:id", handler.GetResponseDetail)

	imp := app.Group("/import")

	imp.Post("/collection", handler.ImportCollection)

	imp.Post("/folder", handler.ImportFolder)

	imp.Post("/request", handler.ImportRequest)

	imp.Post("/response", handler.ImportExample)

	imp.Post("/saverequest", handler.ImportSaveRequest)
}
