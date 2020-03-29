package main

import (
	"github.com/Dadard29/go-api-utils/API"
	"github.com/Dadard29/go-api-utils/service"
	"github.com/Dadard29/go-subscription-connector/subChecker"
	"github.com/Dadard29/go-warehouse/api"
	"github.com/Dadard29/go-warehouse/controllers"
	"net/http"
)

var routes = service.RouteMapping{
	"/upload": service.Route{
		Description:   "manage uploads",
		MethodMapping: service.MethodMapping{
			http.MethodPost: controllers.FileUpload,
			http.MethodDelete: controllers.FileDelete,
		},
	},
	"/upload/list": service.Route{
		Description:   "manage the list of files",
		MethodMapping: service.MethodMapping{
			http.MethodGet: controllers.FileGetList,
		},
	},
	"/download": service.Route{
		Description:   "manage download",
		MethodMapping: service.MethodMapping{
			http.MethodGet: controllers.DownloadGet,
		},
	},
}

// ENV:
// - VERSION: ...
// - CORS_ORIGIN: ... (from dockerfile)

// - HOST_SUB: host where to check the sub token
func main() {
	api.Api = API.NewAPI(
		"warehouse", "config/config.json", routes, true)

	controllers.Sc = subChecker.NewSubChecker(api.Api.Config.GetEnv("HOST_SUB"))

	api.Api.Service.Start()
	api.Api.Service.Stop()
}
