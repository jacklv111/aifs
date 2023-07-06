/*
 * Created on Wed Jul 05 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */
package apigin

import "net/http"

const (
	INVALID_PARAMS = "0001"
	NOT_FOUND      = "0002"
	// used for test, default error
	UNDEFINED_ERROR = "1000"
)

var appErrorCodeHttpStatusMap = map[string]int{
	INVALID_PARAMS: http.StatusBadRequest,
	NOT_FOUND:      http.StatusNotFound,

	UNDEFINED_ERROR: http.StatusInternalServerError,
}

var appErrorCodeMsgTemplateMap = map[string]string{
	INVALID_PARAMS: "invalid params %s.",
	NOT_FOUND:      "%s not found.",

	UNDEFINED_ERROR: "undefined error %s.",
}

func GetHttpStatusByAppErrorCode(appErrorCode string) int {
	if httpStatus, ok := appErrorCodeHttpStatusMap[appErrorCode]; ok {
		return httpStatus
	}
	return http.StatusInternalServerError
}

func GetMsgTemplateByAppErrorCode(appErrorCode string) string {
	if msgTemplate, ok := appErrorCodeMsgTemplateMap[appErrorCode]; ok {
		return msgTemplate
	}
	return "server internal error."
}
