package controllers

import (
	"github.com/Dadard29/go-api-utils/log"
	"github.com/Dadard29/go-api-utils/log/logLevel"
	"github.com/Dadard29/go-subscription-connector/subChecker"
	"github.com/Dadard29/go-warehouse/api"
	"net/http"
)

var Sc *subChecker.SubChecker
var logger = log.NewLogger("CONTROLLER", logLevel.DEBUG)

const (
	apiName        = "warehouse"
	accessTokenKey = "X-Access-Token"
)

func checkToken(token string, w http.ResponseWriter) bool {
	msg, err := Sc.CheckToken(token, apiName)
	if err != nil {
		logger.Error(err.Error())
		api.Api.BuildErrorResponse(http.StatusUnauthorized,
			msg, w)
		return false
	}

	return true
}
