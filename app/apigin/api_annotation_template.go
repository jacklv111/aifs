/*
 * Aifs api
 *
 * aifs api
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apigin

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	manager "github.com/jacklv111/aifs/app/apigin/manager/annotation-template"
	"github.com/jacklv111/aifs/app/apigin/view-object/openapi"
	annotationtemplate "github.com/jacklv111/aifs/pkg/annotation-template"
	"github.com/jacklv111/common-sdk/errors"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/utils"
)

// CopyAnnotationTemplate - Copy an annotation template
func CopyAnnotationTemplate(c *gin.Context) {
	annoTempId := uuid.MustParse(c.Param(ANNOTATION_TEMPLATE_ID))
	resp, err := manager.AnnotationTemplateMgr.CopyAnnotationTemplate(annoTempId)
	if err != nil {
		if err == annotationtemplate.ErrAnnotationTemplateNotFound {
			log.Errorf("annotation template %s not found", annoTempId)
			c.Error(errors.NewAppErr(NOT_FOUND, err, annoTempId.String()))
			return
		}
		log.Errorf("Error occurred when copying annotation template %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// CreateAnnotationTemplate - Create an annotation template
func CreateAnnotationTemplate(c *gin.Context) {
	var req openapi.CreateAnnotationTemplateRequest
	err := c.BindJSON(&req)
	if err != nil {
		log.Errorf("Error occurred when binding json %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
	}
	result, err := manager.AnnotationTemplateMgr.Create(req)
	if err != nil {
		log.Errorf("Error occurred when creating annotation template %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, result)
}

// DeleteAnnotationTemplate - Delete an annotation template
func DeleteAnnotationTemplate(c *gin.Context) {
	annoTempId := uuid.MustParse(c.Param(ANNOTATION_TEMPLATE_ID))
	err := manager.AnnotationTemplateMgr.DeleteById(annoTempId)
	if err != nil {
		if err == annotationtemplate.ErrAnnotationTemplateNotFound {
			log.Errorf("annotation template %s not found", annoTempId)
			c.Error(errors.NewAppErr(NOT_FOUND, err, annoTempId.String()))
			return
		}
		log.Errorf("Error occurred when manager getting details by id, error", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}
	c.Status(http.StatusOK)
}

// GetAnnoTemplateDetails - Get annotation template details
func GetAnnoTemplateDetails(c *gin.Context) {
	annoTempId := uuid.MustParse(c.Param(ANNOTATION_TEMPLATE_ID))
	details, err := manager.AnnotationTemplateMgr.GetDetailsById(annoTempId)
	if err != nil {
		log.Errorf("Error occurred when manager getting details by id, error", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}
	c.JSON(http.StatusOK, details)
}

// GetAnnoTemplateList - Get annotation template list
func GetAnnoTemplateList(c *gin.Context) {
	offset, err := utils.ParseInt(c.Query(OFFSET_STR), 0, math.MaxInt, 0)
	if err != nil {
		log.Errorf("Error occurred when parsing offset type %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
		return
	}
	limit, err := utils.ParseInt(c.Query(LIMIT_STR), LIMIT_MIN, LIMIT_MAX, 10)
	if err != nil {
		log.Errorf("Error occurred when parsing limit type %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
		return
	}

	annoTempIdListStr, ok := c.GetQuery(ANNOTATION_TEMPLATE_ID_LIST)
	annoTempIdList := utils.ParseListStr(annoTempIdListStr, ok, ",")

	resultList, err := manager.AnnotationTemplateMgr.GetList(offset, limit, annoTempIdList)

	if err != nil {
		log.Errorf("Error occurred when getting annotation template list %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}

	c.JSON(http.StatusOK, resultList)
}

// UpdateAnnotationTemplate - Update an annotation template
func UpdateAnnotationTemplate(c *gin.Context) {
	var req openapi.UpdateAnnotationTemplateRequest
	err := c.BindJSON(&req)
	if err != nil {
		log.Errorf("Error occurred when binding json %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
	}
	err = manager.AnnotationTemplateMgr.Update(req)
	if err != nil {
		if err == annotationtemplate.ErrAnnotationTemplateNotFound {
			log.Errorf("annotation template %s not found", req.Id)
			c.Error(errors.NewAppErr(NOT_FOUND, err, req.Id))
			return
		}
		log.Errorf("Error occurred when updating annotation template %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
