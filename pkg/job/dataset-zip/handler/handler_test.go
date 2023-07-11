/*
 * Created on Mon Jul 10 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package handler

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/google/uuid"
	aifsclientgo "github.com/jacklv111/aifs-client-go"
	"github.com/jacklv111/aifs/app/apigin/view-object/openapi"
	annotempconst "github.com/jacklv111/aifs/pkg/annotation-template-type"
	rawdataconst "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/constant"
	"github.com/jacklv111/common-sdk/annotation"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/utils"
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	AIFS_HOST = "http://localhost:8080"
)

func getAifsClient() *aifsclientgo.APIClient {
	clientConfig := aifsclientgo.NewConfiguration()
	clientConfig.Servers = aifsclientgo.ServerConfigurations{
		{
			URL:         AIFS_HOST,
			Description: "No description provided",
		},
	}
	return aifsclientgo.NewAPIClient(clientConfig)
}

func getDataViewDetails(zipfmt *aifsclientgo.ZipFormat) *aifsclientgo.DataViewDetails {
	id := uuid.New().String()
	viewType := aifsclientgo.DATASET_ZIP
	viewName := "test"
	return &aifsclientgo.DataViewDetails{
		Id:        &id,
		ViewType:  &viewType,
		ZipFormat: zipfmt,
		Name:      &viewName,
	}
}

func createIds() (annoTempId, trainRawDataViewId, trainAnnoViewId, valRawDataViewId, valAnnoViewId string) {
	return uuid.New().String(), uuid.New().String(), uuid.New().String(), uuid.New().String(), uuid.New().String()
}

func TestCocoFormatHandler(t *testing.T) {
	log.ValidateAndApply(log.LogConfig)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	Convey("Test Coco Format Handler", t, func() {
		Convey("execute success", func() {
			annoTempId, trainRawDataViewId, trainAnnoViewId, valRawDataViewId, valAnnoViewId := createIds()

			zipView := getDataViewDetails(aifsclientgo.COCO.Ptr())
			var annoTempDetails openapi.AnnotationTemplateDetails
			annoTempDetails.Id = annoTempId

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/annotation-templates", AIFS_HOST),
				func(req *http.Request) (*http.Response, error) {
					var body openapi.CreateAnnotationTemplateRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.Type, ShouldEqual, annotempconst.COCO_TYPE)
					So(len(body.Labels), ShouldEqual, 80)

					annoTempDetails.Type = annotempconst.COCO_TYPE
					annoTempDetails.Labels = body.Labels
					annoTempDetails.Name = body.Name
					for idx := range annoTempDetails.Labels {
						annoTempDetails.Labels[idx].Id = uuid.New().String()
					}
					return httpmock.NewJsonResponse(201, map[string]any{
						"annotationTemplateId": annoTempDetails.Id,
					})
				},
			)

			httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/data-views/%s/dataset-zip", AIFS_HOST, *zipView.Id),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempDetails.Id)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.TrainRawDataViewId, ShouldEqual, *zipView.TrainRawDataViewId)
					So(body.TrainAnnotationViewId, ShouldEqual, *zipView.TrainAnnotationViewId)
					So(body.ValRawDataViewId, ShouldEqual, *zipView.ValRawDataViewId)
					So(body.ValAnnotationViewId, ShouldEqual, *zipView.ValAnnotationViewId)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					// update progress
					return httpmock.NewJsonResponse(200, map[string]any{})
				}),
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/annotation-templates/%s/details", AIFS_HOST, annoTempDetails.Id),
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(200, annoTempDetails)
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views", AIFS_HOST),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					zipView.TrainRawDataViewId = &trainRawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.IMAGE)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": trainRawDataViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.TrainAnnotationViewId = &trainAnnoViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempDetails.Id)
					So(body.RelatedDataViewId, ShouldEqual, *zipView.TrainRawDataViewId)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": trainAnnoViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.ValRawDataViewId = &valRawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.IMAGE)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": valRawDataViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.ValAnnotationViewId = &valAnnoViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempDetails.Id)
					So(body.RelatedDataViewId, ShouldEqual, *zipView.ValRawDataViewId)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": valAnnoViewId,
					})
				}),
			)
			// train data check
			rawDataId := uuid.New().String()

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, trainRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File["files"]))
					So("000000002587.jpg", ShouldEqual, req.MultipartForm.File["files"][0].Filename)

					file, err := req.MultipartForm.File["files"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256FromFile("test-data/coco-data/train/raw-data/000000002587.jpg")
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/data-views/%s/raw-data-hash-list?limit=1000&offset=0", AIFS_HOST, trainRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					hash, err := utils.GetFileSha256FromFile("test-data/coco-data/train/raw-data/000000002587.jpg")
					So(err, ShouldBeNil)
					return httpmock.NewJsonResponse(200, []any{
						map[string]string{
							"rawDataId": rawDataId,
							"sha256":    hash,
						},
					})
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/annotations", AIFS_HOST, trainAnnoViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File))
					_, ok := req.MultipartForm.File[rawDataId]
					So(ok, ShouldBeTrue)

					file, err := req.MultipartForm.File[rawDataId][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()

					var anno annotation.CocoAnnoFormat
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					err = json.Unmarshal(fileBytes, &anno)
					sort.Slice(anno.AnnoData, func(i, j int) bool {
						return anno.AnnoData[i].Id < anno.AnnoData[j].Id
					})
					So(err, ShouldBeNil)

					So(anno.RawDataId, ShouldEqual, rawDataId)
					So(anno.AnnotationTemplateId, ShouldEqual, annoTempId)
					So(len(anno.AnnoData), ShouldEqual, 2)
					So(anno.AnnoData[0].Id, ShouldEqual, 1042327)
					So(anno.AnnoData[1].Id, ShouldEqual, 1078530)
					// add more checks here if necessary
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			// val data check
			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, valRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File["files"]))
					So("000000032081.jpg", ShouldEqual, req.MultipartForm.File["files"][0].Filename)

					file, err := req.MultipartForm.File["files"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256FromFile("test-data/coco-data/val/raw-data/000000032081.jpg")
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/data-views/%s/raw-data-hash-list?limit=1000&offset=0", AIFS_HOST, valRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					hash, err := utils.GetFileSha256FromFile("test-data/coco-data/val/raw-data/000000032081.jpg")
					So(err, ShouldBeNil)
					return httpmock.NewJsonResponse(200, []any{
						map[string]string{
							"rawDataId": rawDataId,
							"sha256":    hash,
						},
					})
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/annotations", AIFS_HOST, valAnnoViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File))
					_, ok := req.MultipartForm.File[rawDataId]
					So(ok, ShouldBeTrue)

					file, err := req.MultipartForm.File[rawDataId][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()

					var anno annotation.CocoAnnoFormat
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					err = json.Unmarshal(fileBytes, &anno)
					sort.Slice(anno.AnnoData, func(i, j int) bool {
						return anno.AnnoData[i].Id < anno.AnnoData[j].Id
					})
					So(err, ShouldBeNil)

					So(anno.RawDataId, ShouldEqual, rawDataId)
					So(anno.AnnotationTemplateId, ShouldEqual, annoTempId)
					So(len(anno.AnnoData), ShouldEqual, 2)
					So(anno.AnnoData[0].Id, ShouldEqual, 469488)
					So(anno.AnnoData[1].Id, ShouldEqual, 657735)
					// add more checks here if necessary
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)
			client := getAifsClient()
			handler, err := BuildHandler(client, zipView, "test-data/coco-data")
			So(err, ShouldBeNil)
			err = handler.Exec()
			So(err, ShouldBeNil)
		})
	})
}

func TestSegmentationMasksFormatHandler(t *testing.T) {
	log.ValidateAndApply(log.LogConfig)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	Convey("Test Segmentation Masks Format Handler", t, func() {
		Convey("execute success", func() {
			annoTempId, trainRawDataViewId, trainAnnoViewId, valRawDataViewId, valAnnoViewId := createIds()

			zipView := getDataViewDetails(aifsclientgo.IMAGE_SEGMENTATION_MASKS.Ptr())
			var annoTempDetails openapi.AnnotationTemplateDetails
			annoTempDetails.Id = annoTempId

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/annotation-templates", AIFS_HOST),
				func(req *http.Request) (*http.Response, error) {
					var body openapi.CreateAnnotationTemplateRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.Type, ShouldEqual, annotempconst.SEGMENTATION_MASKS)
					So(len(body.Labels), ShouldEqual, 2)

					annoTempDetails.Type = annotempconst.SEGMENTATION_MASKS
					annoTempDetails.Labels = body.Labels
					annoTempDetails.Name = body.Name
					for idx := range annoTempDetails.Labels {
						annoTempDetails.Labels[idx].Id = uuid.New().String()
					}
					return httpmock.NewJsonResponse(201, map[string]any{
						"annotationTemplateId": annoTempDetails.Id,
					})
				},
			)

			httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/data-views/%s/dataset-zip", AIFS_HOST, *zipView.Id),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempDetails.Id)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.TrainRawDataViewId, ShouldEqual, *zipView.TrainRawDataViewId)
					So(body.TrainAnnotationViewId, ShouldEqual, *zipView.TrainAnnotationViewId)
					So(body.ValRawDataViewId, ShouldEqual, *zipView.ValRawDataViewId)
					So(body.ValAnnotationViewId, ShouldEqual, *zipView.ValAnnotationViewId)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					// update progress
					return httpmock.NewJsonResponse(200, map[string]any{})
				}),
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/annotation-templates/%s/details", AIFS_HOST, annoTempDetails.Id),
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(200, annoTempDetails)
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views", AIFS_HOST),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					zipView.TrainRawDataViewId = &trainRawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.IMAGE)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": trainRawDataViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.TrainAnnotationViewId = &trainAnnoViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempDetails.Id)
					So(body.RelatedDataViewId, ShouldEqual, *zipView.TrainRawDataViewId)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": trainAnnoViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.ValRawDataViewId = &valRawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.IMAGE)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": valRawDataViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.ValAnnotationViewId = &valAnnoViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempDetails.Id)
					So(body.RelatedDataViewId, ShouldEqual, *zipView.ValRawDataViewId)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": valAnnoViewId,
					})
				}),
			)
			// train data check
			rawDataId := uuid.New().String()

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, trainRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File["files"]))
					So("87.jpg", ShouldEqual, req.MultipartForm.File["files"][0].Filename)

					file, err := req.MultipartForm.File["files"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256FromFile("test-data/segmentation-masks/train/raw-data/87.jpg")
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/data-views/%s/raw-data-hash-list?limit=1000&offset=0", AIFS_HOST, trainRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					hash, err := utils.GetFileSha256FromFile("test-data/segmentation-masks/train/raw-data/87.jpg")
					So(err, ShouldBeNil)
					return httpmock.NewJsonResponse(200, []any{
						map[string]string{
							"rawDataId": rawDataId,
							"sha256":    hash,
						},
					})
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/annotations", AIFS_HOST, trainAnnoViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File))
					_, ok := req.MultipartForm.File[rawDataId]
					So(ok, ShouldBeTrue)
					So("87.png", ShouldEqual, req.MultipartForm.File[rawDataId][0].Filename)

					file, err := req.MultipartForm.File[rawDataId][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256FromFile("test-data/segmentation-masks/train/annotation/87.png")
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)
					// add more checks here if necessary
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			// val data check
			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, valRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File["files"]))
					So("632.jpg", ShouldEqual, req.MultipartForm.File["files"][0].Filename)

					file, err := req.MultipartForm.File["files"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256FromFile("test-data/segmentation-masks/val/raw-data/632.jpg")
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/data-views/%s/raw-data-hash-list?limit=1000&offset=0", AIFS_HOST, valRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					hash, err := utils.GetFileSha256FromFile("test-data/segmentation-masks/val/raw-data/632.jpg")
					So(err, ShouldBeNil)
					return httpmock.NewJsonResponse(200, []any{
						map[string]string{
							"rawDataId": rawDataId,
							"sha256":    hash,
						},
					})
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/annotations", AIFS_HOST, valAnnoViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File))
					_, ok := req.MultipartForm.File[rawDataId]
					So(ok, ShouldBeTrue)
					So("632.png", ShouldEqual, req.MultipartForm.File[rawDataId][0].Filename)

					file, err := req.MultipartForm.File[rawDataId][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256FromFile("test-data/segmentation-masks/val/annotation/632.png")
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)
					// add more checks here if necessary
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)
			client := getAifsClient()
			handler, err := BuildHandler(client, zipView, "test-data/segmentation-masks")
			So(err, ShouldBeNil)
			err = handler.Exec()
			So(err, ShouldBeNil)
		})
	})
}

func TestRgbdFormatHandler(t *testing.T) {
	log.ValidateAndApply(log.LogConfig)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	Convey("Test Rgbd Format Handler", t, func() {
		Convey("execute success", func() {
			annoTempId, trainRawDataViewId, trainAnnoViewId, valRawDataViewId, valAnnoViewId := createIds()

			zipView := getDataViewDetails(aifsclientgo.RGBD_BOUNDING_BOX_2D_AND_3D.Ptr())
			var annoTempDetails openapi.AnnotationTemplateDetails
			annoTempDetails.Id = annoTempId
			client := getAifsClient()
			handler, err := BuildHandler(client, zipView, "test-data/rgbd")
			So(err, ShouldBeNil)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/annotation-templates", AIFS_HOST),
				func(req *http.Request) (*http.Response, error) {
					var body openapi.CreateAnnotationTemplateRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.Type, ShouldEqual, annotempconst.RGBD)
					So(len(body.Labels), ShouldEqual, 4)

					annoTempDetails.Type = annotempconst.RGBD
					annoTempDetails.Labels = body.Labels
					annoTempDetails.Name = body.Name
					for idx := range annoTempDetails.Labels {
						annoTempDetails.Labels[idx].Id = uuid.New().String()
					}
					return httpmock.NewJsonResponse(201, map[string]any{
						"annotationTemplateId": annoTempDetails.Id,
					})
				},
			)

			httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/data-views/%s/dataset-zip", AIFS_HOST, *zipView.Id),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempDetails.Id)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.TrainRawDataViewId, ShouldEqual, *zipView.TrainRawDataViewId)
					So(body.TrainAnnotationViewId, ShouldEqual, *zipView.TrainAnnotationViewId)
					So(body.ValRawDataViewId, ShouldEqual, *zipView.ValRawDataViewId)
					So(body.ValAnnotationViewId, ShouldEqual, *zipView.ValAnnotationViewId)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					// update progress
					return httpmock.NewJsonResponse(200, map[string]any{})
				}),
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/annotation-templates/%s/details", AIFS_HOST, annoTempDetails.Id),
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(200, annoTempDetails)
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views", AIFS_HOST),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					zipView.TrainRawDataViewId = &trainRawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.RGBD)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": trainRawDataViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.TrainAnnotationViewId = &trainAnnoViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempDetails.Id)
					So(body.RelatedDataViewId, ShouldEqual, *zipView.TrainRawDataViewId)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": trainAnnoViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.ValRawDataViewId = &valRawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.RGBD)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": valRawDataViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.ValAnnotationViewId = &valAnnoViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempDetails.Id)
					So(body.RelatedDataViewId, ShouldEqual, *zipView.ValRawDataViewId)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": valAnnoViewId,
					})
				}),
			)
			// train data check
			rawDataId := uuid.New().String()
			var hash string
			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, trainRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File["files"]))
					So("005051.bin", ShouldEqual, req.MultipartForm.File["files"][0].Filename)

					file, err := req.MultipartForm.File["files"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					binBytes, meta, err := handler.(*RgbdHandler).GetBinAndMeta("005051.bin", "test-data/rgbd/train/image/005051.jpg", "test-data/rgbd/train/depth/005051.png", "test-data/rgbd/train/calib/005051.txt")
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256Bytes(binBytes)
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)
					hash = hashActual

					// check meta
					fileMeta, err := req.MultipartForm.File["fileMeta"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err = io.ReadAll(fileMeta)
					So(err, ShouldBeNil)
					So(strings.Trim(string(fileBytes), " \n"), ShouldEqual, meta.String())

					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/data-views/%s/raw-data-hash-list?limit=1000&offset=0", AIFS_HOST, trainRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(200, []any{
						map[string]string{
							"rawDataId": rawDataId,
							"sha256":    hash,
						},
					})
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/annotations", AIFS_HOST, trainAnnoViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File))
					_, ok := req.MultipartForm.File[rawDataId]
					So(ok, ShouldBeTrue)
					So("005051.json", ShouldEqual, req.MultipartForm.File[rawDataId][0].Filename)

					file, err := req.MultipartForm.File[rawDataId][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)

					var anno annotation.RgbdAnnotation
					err = json.Unmarshal(fileBytes, &anno)
					So(err, ShouldBeNil)
					So(anno.RawDataId, ShouldEqual, rawDataId)
					So(anno.AnnotationTemplateId, ShouldEqual, annoTempId)
					So(len(anno.BoundingBoxList), ShouldEqual, 7)
					// add more checks here if necessary

					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			// val data check

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, valRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File["files"]))
					So("000001.bin", ShouldEqual, req.MultipartForm.File["files"][0].Filename)

					file, err := req.MultipartForm.File["files"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					binBytes, meta, err := handler.(*RgbdHandler).GetBinAndMeta("000001.bin", "test-data/rgbd/val/image/000001.jpg", "test-data/rgbd/val/depth/000001.png", "test-data/rgbd/val/calib/000001.txt")
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256Bytes(binBytes)
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)
					hash = hashActual

					// check meta
					fileMeta, err := req.MultipartForm.File["fileMeta"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err = io.ReadAll(fileMeta)
					So(err, ShouldBeNil)
					So(strings.Trim(string(fileBytes), " \n"), ShouldEqual, meta.String())

					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/data-views/%s/raw-data-hash-list?limit=1000&offset=0", AIFS_HOST, valRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(200, []any{
						map[string]string{
							"rawDataId": rawDataId,
							"sha256":    hash,
						},
					})
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/annotations", AIFS_HOST, valAnnoViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File))
					_, ok := req.MultipartForm.File[rawDataId]
					So(ok, ShouldBeTrue)
					So("000001.json", ShouldEqual, req.MultipartForm.File[rawDataId][0].Filename)

					file, err := req.MultipartForm.File[rawDataId][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)

					var anno annotation.RgbdAnnotation
					err = json.Unmarshal(fileBytes, &anno)
					So(err, ShouldBeNil)
					So(anno.RawDataId, ShouldEqual, rawDataId)
					So(anno.AnnotationTemplateId, ShouldEqual, annoTempId)
					So(len(anno.BoundingBoxList), ShouldEqual, 3)
					// add more checks here if necessary

					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			err = handler.Exec()
			So(err, ShouldBeNil)
		})
	})
}

func TestClassificationFormatHandler(t *testing.T) {
	log.ValidateAndApply(log.LogConfig)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	Convey("Test Classification Format Handler", t, func() {
		Convey("execute success", func() {
			annoTempId, trainRawDataViewId, trainAnnoViewId, valRawDataViewId, valAnnoViewId := createIds()

			zipView := getDataViewDetails(aifsclientgo.IMAGE_CLASSIFICATION.Ptr())
			var annoTempDetails openapi.AnnotationTemplateDetails
			annoTempDetails.Id = annoTempId
			client := getAifsClient()
			handler, err := BuildHandler(client, zipView, "test-data/classification")
			So(err, ShouldBeNil)
			labelIdNameMap := make(map[string]string)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/annotation-templates", AIFS_HOST),
				func(req *http.Request) (*http.Response, error) {
					var body openapi.CreateAnnotationTemplateRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.Type, ShouldEqual, annotempconst.CATEGORY)
					So(len(body.Labels), ShouldEqual, 2)

					annoTempDetails.Type = annotempconst.CATEGORY
					annoTempDetails.Labels = body.Labels
					annoTempDetails.Name = body.Name
					for idx := range annoTempDetails.Labels {
						annoTempDetails.Labels[idx].Id = uuid.New().String()
					}
					for _, label := range annoTempDetails.Labels {
						labelIdNameMap[label.Id] = label.Name
					}
					return httpmock.NewJsonResponse(201, map[string]any{
						"annotationTemplateId": annoTempId,
					})
				},
			)

			httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/data-views/%s/dataset-zip", AIFS_HOST, *zipView.Id),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempId)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.TrainRawDataViewId, ShouldEqual, *zipView.TrainRawDataViewId)
					So(body.TrainAnnotationViewId, ShouldEqual, *zipView.TrainAnnotationViewId)
					So(body.ValRawDataViewId, ShouldEqual, *zipView.ValRawDataViewId)
					So(body.ValAnnotationViewId, ShouldEqual, *zipView.ValAnnotationViewId)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					// update progress
					return httpmock.NewJsonResponse(200, map[string]any{})
				}),
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/annotation-templates/%s/details", AIFS_HOST, annoTempId),
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(200, annoTempDetails)
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views", AIFS_HOST),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					zipView.TrainRawDataViewId = &trainRawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.IMAGE)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": trainRawDataViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.TrainAnnotationViewId = &trainAnnoViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempId)
					So(body.RelatedDataViewId, ShouldEqual, *zipView.TrainRawDataViewId)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": trainAnnoViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.ValRawDataViewId = &valRawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.IMAGE)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": valRawDataViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.ValAnnotationViewId = &valAnnoViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempId)
					So(body.RelatedDataViewId, ShouldEqual, *zipView.ValRawDataViewId)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": valAnnoViewId,
					})
				}),
			)
			// train data check
			type IdHash struct {
				Id        string
				Hash      string
				LabelName string
			}

			rawDataIdHashMap := make(map[string]IdHash, 0)
			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, trainRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(5, ShouldEqual, len(req.MultipartForm.File["files"]))
					set := mapset.NewSet[string]()
					for _, form := range req.MultipartForm.File["files"] {
						file, err := form.Open()
						So(err, ShouldBeNil)
						defer file.Close()
						fileBytes, err := io.ReadAll(file)
						So(err, ShouldBeNil)
						hashActual, err := utils.GetFileSha256Bytes(fileBytes)
						So(err, ShouldBeNil)
						set.Add(hashActual)
					}
					fileList, err := utils.ReadAllFiles("test-data/classification/train")
					So(err, ShouldBeNil)
					for _, filePath := range fileList {
						hashExpected, err := utils.GetFileSha256FromFile(filePath)
						So(err, ShouldBeNil)
						So(set.Contains(hashExpected), ShouldBeTrue)
						idhash := IdHash{Id: uuid.New().String(), Hash: hashExpected, LabelName: handler.(*ClassificationHandler).GetLabelName(filePath)}
						rawDataIdHashMap[idhash.Id] = idhash
					}

					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/data-views/%s/raw-data-hash-list?limit=1000&offset=0", AIFS_HOST, trainRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					body := make([]map[string]string, 0)
					for _, data := range rawDataIdHashMap {
						body = append(body, map[string]string{
							"rawDataId": data.Id,
							"sha256":    data.Hash,
						})
					}

					return httpmock.NewJsonResponse(200, body)
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/annotations", AIFS_HOST, trainAnnoViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					_, ok := req.MultipartForm.File["fileMeta"]
					So(ok, ShouldBeTrue)
					fileMeta, err := req.MultipartForm.File["fileMeta"][0].Open()
					So(err, ShouldBeNil)
					scanner := bufio.NewScanner(fileMeta)

					for i := 0; i < 2; i++ {
						So(scanner.Scan(), ShouldBeTrue)
						line := scanner.Text()
						stList := strings.Split(line, " ")
						So(labelIdNameMap[stList[1]], ShouldEqual, rawDataIdHashMap[stList[0]].LabelName)
					}

					// add more checks here if necessary

					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			// val data check
			rawDataIdHashMap = make(map[string]IdHash, 0)
			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, valRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(5, ShouldEqual, len(req.MultipartForm.File["files"]))
					set := mapset.NewSet[string]()
					for _, form := range req.MultipartForm.File["files"] {
						file, err := form.Open()
						So(err, ShouldBeNil)
						defer file.Close()
						fileBytes, err := io.ReadAll(file)
						So(err, ShouldBeNil)
						hashActual, err := utils.GetFileSha256Bytes(fileBytes)
						So(err, ShouldBeNil)
						set.Add(hashActual)
					}
					fileList, err := utils.ReadAllFiles("test-data/classification/val")
					So(err, ShouldBeNil)
					for _, filePath := range fileList {
						hashExpected, err := utils.GetFileSha256FromFile(filePath)
						So(err, ShouldBeNil)
						So(set.Contains(hashExpected), ShouldBeTrue)
						idhash := IdHash{Id: uuid.New().String(), Hash: hashExpected, LabelName: handler.(*ClassificationHandler).GetLabelName(filePath)}
						rawDataIdHashMap[idhash.Id] = idhash
					}

					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/data-views/%s/raw-data-hash-list?limit=1000&offset=0", AIFS_HOST, valRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					body := make([]map[string]string, 0)
					for _, data := range rawDataIdHashMap {
						body = append(body, map[string]string{
							"rawDataId": data.Id,
							"sha256":    data.Hash,
						})
					}

					return httpmock.NewJsonResponse(200, body)
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/annotations", AIFS_HOST, valAnnoViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					_, ok := req.MultipartForm.File["fileMeta"]
					So(ok, ShouldBeTrue)
					fileMeta, err := req.MultipartForm.File["fileMeta"][0].Open()
					So(err, ShouldBeNil)
					scanner := bufio.NewScanner(fileMeta)

					for i := 0; i < 2; i++ {
						So(scanner.Scan(), ShouldBeTrue)
						line := scanner.Text()
						stList := strings.Split(line, " ")
						So(labelIdNameMap[stList[1]], ShouldEqual, rawDataIdHashMap[stList[0]].LabelName)
					}

					// add more checks here if necessary

					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)
			err = handler.Exec()
			So(err, ShouldBeNil)
		})
	})
}

func TestOcrFormatHandler(t *testing.T) {
	log.ValidateAndApply(log.LogConfig)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	Convey("Test Ocr Format Handler", t, func() {
		Convey("execute success", func() {
			annoTempId, trainRawDataViewId, trainAnnoViewId, valRawDataViewId, valAnnoViewId := createIds()

			zipView := getDataViewDetails(aifsclientgo.OCR.Ptr())
			var annoTempDetails openapi.AnnotationTemplateDetails
			annoTempDetails.Id = annoTempId
			client := getAifsClient()
			handler, err := BuildHandler(client, zipView, "test-data/ocr")
			So(err, ShouldBeNil)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/annotation-templates", AIFS_HOST),
				func(req *http.Request) (*http.Response, error) {
					var body openapi.CreateAnnotationTemplateRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.Type, ShouldEqual, annotempconst.OCR)
					So(len(body.WordList), ShouldEqual, 11)
					// check word list details if necessary
					annoTempDetails.Type = annotempconst.OCR
					annoTempDetails.WordList = body.WordList
					annoTempDetails.Name = body.Name
					for idx := range annoTempDetails.Labels {
						annoTempDetails.Labels[idx].Id = uuid.New().String()
					}

					return httpmock.NewJsonResponse(201, map[string]any{
						"annotationTemplateId": annoTempId,
					})
				},
			)

			httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/data-views/%s/dataset-zip", AIFS_HOST, *zipView.Id),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempId)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.TrainRawDataViewId, ShouldEqual, *zipView.TrainRawDataViewId)
					So(body.TrainAnnotationViewId, ShouldEqual, *zipView.TrainAnnotationViewId)
					So(body.ValRawDataViewId, ShouldEqual, *zipView.ValRawDataViewId)
					So(body.ValAnnotationViewId, ShouldEqual, *zipView.ValAnnotationViewId)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					// update progress
					return httpmock.NewJsonResponse(200, map[string]any{})
				}),
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/annotation-templates/%s/details", AIFS_HOST, annoTempId),
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(200, annoTempDetails)
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views", AIFS_HOST),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					zipView.TrainRawDataViewId = &trainRawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.IMAGE)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": trainRawDataViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.TrainAnnotationViewId = &trainAnnoViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempId)
					So(body.RelatedDataViewId, ShouldEqual, *zipView.TrainRawDataViewId)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": trainAnnoViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.ValRawDataViewId = &valRawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.IMAGE)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": valRawDataViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.ValAnnotationViewId = &valAnnoViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempId)
					So(body.RelatedDataViewId, ShouldEqual, *zipView.ValRawDataViewId)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": valAnnoViewId,
					})
				}),
			)
			// train data check
			rawDataId := uuid.New().String()
			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, trainRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File["files"]))
					So("0_0_0_3_27_32_32_31_33.jpg", ShouldEqual, req.MultipartForm.File["files"][0].Filename)

					file, err := req.MultipartForm.File["files"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256FromFile("test-data/ocr/train/images/0_0_0_3_27_32_32_31_33.jpg")
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/data-views/%s/raw-data-hash-list?limit=1000&offset=0", AIFS_HOST, trainRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					hash, err := utils.GetFileSha256FromFile("test-data/ocr/train/images/0_0_0_3_27_32_32_31_33.jpg")
					So(err, ShouldBeNil)
					return httpmock.NewJsonResponse(200, []any{
						map[string]string{
							"rawDataId": rawDataId,
							"sha256":    hash,
						},
					})
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/annotations", AIFS_HOST, trainAnnoViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					_, ok := req.MultipartForm.File["fileMeta"]
					So(ok, ShouldBeTrue)
					fileMeta, err := req.MultipartForm.File["fileMeta"][0].Open()
					So(err, ShouldBeNil)
					scanner := bufio.NewScanner(fileMeta)

					for i := 0; i < 1; i++ {
						So(scanner.Scan(), ShouldBeTrue)
						line := scanner.Text()
						stList := strings.Split(line, " ")
						So(stList[0], ShouldEqual, rawDataId)
						So(stList[1], ShouldEqual, "皖AD38879")
					}
					So(scanner.Scan(), ShouldBeFalse)

					// add more checks here if necessary

					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			// val data check
			rawDataId = uuid.New().String()
			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, valRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File["files"]))
					So("1_0_0_3_28_29_32_32_33.jpg", ShouldEqual, req.MultipartForm.File["files"][0].Filename)

					file, err := req.MultipartForm.File["files"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256FromFile("test-data/ocr/val/images/1_0_0_3_28_29_32_32_33.jpg")
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/data-views/%s/raw-data-hash-list?limit=1000&offset=0", AIFS_HOST, valRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					hash, err := utils.GetFileSha256FromFile("test-data/ocr/val/images/1_0_0_3_28_29_32_32_33.jpg")
					So(err, ShouldBeNil)
					return httpmock.NewJsonResponse(200, []any{
						map[string]string{
							"rawDataId": rawDataId,
							"sha256":    hash,
						},
					})
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/annotations", AIFS_HOST, valAnnoViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					_, ok := req.MultipartForm.File["fileMeta"]
					So(ok, ShouldBeTrue)
					fileMeta, err := req.MultipartForm.File["fileMeta"][0].Open()
					So(err, ShouldBeNil)
					scanner := bufio.NewScanner(fileMeta)

					for i := 0; i < 1; i++ {
						So(scanner.Scan(), ShouldBeTrue)
						line := scanner.Text()
						stList := strings.Split(line, " ")
						So(stList[0], ShouldEqual, rawDataId)
						So(stList[1], ShouldEqual, "皖AD45889")
					}
					So(scanner.Scan(), ShouldBeFalse)

					// add more checks here if necessary

					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			err = handler.Exec()
			So(err, ShouldBeNil)
		})
	})
}

func TestRawDataImagesHandler(t *testing.T) {
	log.ValidateAndApply(log.LogConfig)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	Convey("Test Raw Data Images Handler", t, func() {
		Convey("execute success", func() {
			rawDataViewId := uuid.New().String()

			zipView := getDataViewDetails(aifsclientgo.RAW_DATA_IMAGES.Ptr())

			httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/data-views/%s/dataset-zip", AIFS_HOST, *zipView.Id),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataViewId, ShouldEqual, *zipView.RawDataViewId)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					// update progress
					return httpmock.NewJsonResponse(200, map[string]any{})
				}),
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views", AIFS_HOST),
				func(req *http.Request) (*http.Response, error) {
					zipView.RawDataViewId = &rawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.IMAGE)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": rawDataViewId,
					})
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, rawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File["files"]))
					So("000000002587.jpg", ShouldEqual, req.MultipartForm.File["files"][0].Filename)

					file, err := req.MultipartForm.File["files"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256FromFile("test-data/coco-data/train/raw-data/000000002587.jpg")
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			client := getAifsClient()
			handler, err := BuildHandler(client, zipView, "test-data/raw-data-images")
			So(err, ShouldBeNil)
			err = handler.Exec()
			So(err, ShouldBeNil)
		})
	})
}

func TestSamHandler(t *testing.T) {
	log.ValidateAndApply(log.LogConfig)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	Convey("Test Coco Format Handler", t, func() {
		Convey("execute success", func() {
			annoTempId := uuid.New().String()
			rawDataViewId := uuid.New().String()
			annoViewId := uuid.New().String()

			zipView := getDataViewDetails(aifsclientgo.SAM.Ptr())
			zipView.AnnotationTemplateId = &annoTempId
			var annoTempDetails openapi.AnnotationTemplateDetails
			annoTempDetails.Id = annoTempId

			httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/data-views/%s/dataset-zip", AIFS_HOST, *zipView.Id),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataViewId, ShouldEqual, *zipView.RawDataViewId)
					So(body.AnnotationViewId, ShouldEqual, *zipView.AnnotationViewId)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					// update progress
					return httpmock.NewJsonResponse(200, map[string]any{})
				}),
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/annotation-templates/%s/details", AIFS_HOST, annoTempDetails.Id),
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(200, annoTempDetails)
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views", AIFS_HOST),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					zipView.RawDataViewId = &rawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.IMAGE)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": rawDataViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.AnnotationViewId = &annoViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempDetails.Id)
					So(body.RelatedDataViewId, ShouldEqual, *zipView.RawDataViewId)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": annoViewId,
					})
				}),
			)

			rawDataId := uuid.New().String()

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, rawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File["files"]))
					So("sa_1.jpg", ShouldEqual, req.MultipartForm.File["files"][0].Filename)

					file, err := req.MultipartForm.File["files"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256FromFile("test-data/sam/sa00000/sa_1.jpg")
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/data-views/%s/raw-data-hash-list?limit=1000&offset=0", AIFS_HOST, rawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					hash, err := utils.GetFileSha256FromFile("test-data/sam/sa00000/sa_1.jpg")
					So(err, ShouldBeNil)
					return httpmock.NewJsonResponse(200, []any{
						map[string]string{
							"rawDataId": rawDataId,
							"sha256":    hash,
						},
					})
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/annotations", AIFS_HOST, annoViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File))
					_, ok := req.MultipartForm.File[rawDataId]
					So(ok, ShouldBeTrue)

					file, err := req.MultipartForm.File[rawDataId][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()

					var anno annotation.CocoAnnoFormat
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					err = json.Unmarshal(fileBytes, &anno)
					sort.Slice(anno.AnnoData, func(i, j int) bool {
						return anno.AnnoData[i].Id < anno.AnnoData[j].Id
					})
					So(err, ShouldBeNil)

					So(anno.RawDataId, ShouldEqual, rawDataId)
					So(anno.AnnotationTemplateId, ShouldEqual, annoTempId)
					So(len(anno.AnnoData), ShouldEqual, 65)
					So(anno.AnnoData[0].Id, ShouldEqual, 523353737)

					// add more checks here if necessary
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			client := getAifsClient()
			handler, err := BuildHandler(client, zipView, "test-data/sam")
			So(err, ShouldBeNil)
			err = handler.Exec()
			So(err, ShouldBeNil)
		})
	})
}

func TestPoints3DFormatHandler(t *testing.T) {
	log.ValidateAndApply(log.LogConfig)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	Convey("Test points 3d Format Handler", t, func() {
		Convey("execute success", func() {
			annoTempId, trainRawDataViewId, trainAnnoViewId, valRawDataViewId, valAnnoViewId := createIds()

			zipView := getDataViewDetails(aifsclientgo.POINTS_3D_ZIP.Ptr())
			var annoTempDetails openapi.AnnotationTemplateDetails
			annoTempDetails.Id = annoTempId

			client := getAifsClient()
			handler, err := BuildHandler(client, zipView, "test-data/points-3d")
			So(err, ShouldBeNil)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/annotation-templates", AIFS_HOST),
				func(req *http.Request) (*http.Response, error) {
					var body openapi.CreateAnnotationTemplateRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.Type, ShouldEqual, annotempconst.POINTS_3D)
					So(len(body.Labels), ShouldEqual, 13)

					annoTempDetails.Type = annotempconst.POINTS_3D
					annoTempDetails.Labels = body.Labels
					annoTempDetails.Name = body.Name
					for idx := range annoTempDetails.Labels {
						annoTempDetails.Labels[idx].Id = uuid.New().String()
					}
					return httpmock.NewJsonResponse(201, map[string]any{
						"annotationTemplateId": annoTempDetails.Id,
					})
				},
			)

			httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/data-views/%s/dataset-zip", AIFS_HOST, *zipView.Id),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempDetails.Id)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					var body openapi.UpdateDatasetZipRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.TrainRawDataViewId, ShouldEqual, *zipView.TrainRawDataViewId)
					So(body.TrainAnnotationViewId, ShouldEqual, *zipView.TrainAnnotationViewId)
					So(body.ValRawDataViewId, ShouldEqual, *zipView.ValRawDataViewId)
					So(body.ValAnnotationViewId, ShouldEqual, *zipView.ValAnnotationViewId)
					return httpmock.NewJsonResponse(200, map[string]any{})
				}).Then(func(req *http.Request) (*http.Response, error) {
					// update progress
					return httpmock.NewJsonResponse(200, map[string]any{})
				}),
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/annotation-templates/%s/details", AIFS_HOST, annoTempDetails.Id),
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(200, annoTempDetails)
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views", AIFS_HOST),
				httpmock.Responder(func(req *http.Request) (*http.Response, error) {
					zipView.TrainRawDataViewId = &trainRawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.POINTS_3D)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": trainRawDataViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.TrainAnnotationViewId = &trainAnnoViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempDetails.Id)
					So(body.RelatedDataViewId, ShouldEqual, *zipView.TrainRawDataViewId)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": trainAnnoViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.ValRawDataViewId = &valRawDataViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.RawDataType, ShouldEqual, rawdataconst.POINTS_3D)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": valRawDataViewId,
					})
				}).Then(func(req *http.Request) (*http.Response, error) {
					zipView.ValAnnotationViewId = &valAnnoViewId

					var body openapi.CreateDataViewRequest
					content, err := io.ReadAll(req.Body)
					defer req.Body.Close()
					So(err, ShouldBeNil)
					err = json.Unmarshal(content, &body)
					So(err, ShouldBeNil)
					So(body.AnnotationTemplateId, ShouldEqual, annoTempDetails.Id)
					So(body.RelatedDataViewId, ShouldEqual, *zipView.ValRawDataViewId)

					return httpmock.NewJsonResponse(201, map[string]any{
						"dataViewId": valAnnoViewId,
					})
				}),
			)
			// train data check
			rawDataId := uuid.New().String()
			var binHash string

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, trainRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File["files"]))
					So("room_1.bin", ShouldEqual, req.MultipartForm.File["files"][0].Filename)

					file, err := req.MultipartForm.File["files"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)

					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					binBytes, meta, err := handler.(*Points3DHandler).GetBinAndMeta("room_1.bin", "test-data/points-3d/train/raw-data/pos/room_1", "test-data/points-3d/train/raw-data/rgb/room_1")
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256Bytes(binBytes)
					binHash = hashExpected
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)

					So(meta.Sha256, ShouldEqual, hashExpected)
					So(meta.FileName, ShouldEqual, "room_1.bin")
					So(meta.Size, ShouldEqual, 3)
					So(meta.Xmin, ShouldAlmostEqual, 1.604)
					So(meta.Xmax, ShouldAlmostEqual, 1.633)
					So(meta.Ymin, ShouldAlmostEqual, 1.8390000000000004)
					So(meta.Ymax, ShouldAlmostEqual, 1.9670000000000005)
					So(meta.Zmin, ShouldAlmostEqual, 0.017000000000000015)
					So(meta.Zmax, ShouldAlmostEqual, 0.019000000000000017)
					So(meta.Rmean, ShouldAlmostEqual, 118.33333333333333)
					So(meta.Gmean, ShouldAlmostEqual, 130)
					So(meta.Bmean, ShouldAlmostEqual, 126)
					So(meta.Rstd, ShouldAlmostEqual, 0.9428090415823364)
					So(meta.Gstd, ShouldAlmostEqual, 1.4142135623730951)
					So(meta.Bstd, ShouldAlmostEqual, 1.4142135623730951)
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/data-views/%s/raw-data-hash-list?limit=1000&offset=0", AIFS_HOST, trainRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					So(err, ShouldBeNil)
					return httpmock.NewJsonResponse(200, []any{
						map[string]string{
							"rawDataId": rawDataId,
							"sha256":    binHash,
						},
					})
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/annotations", AIFS_HOST, trainAnnoViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File))
					_, ok := req.MultipartForm.File[rawDataId]
					So(ok, ShouldBeTrue)
					So("room_1", ShouldEqual, req.MultipartForm.File[rawDataId][0].Filename)

					file, err := req.MultipartForm.File[rawDataId][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					reader, err := zlib.NewReader(file)
					So(err, ShouldBeNil)

					actualAnno, err := io.ReadAll(reader)
					So(err, ShouldBeNil)

					labelNameIdMap := make(map[string]string, 0)
					for _, label := range annoTempDetails.Labels {
						labelNameIdMap[label.Name] = label.Id
					}

					annoBytes := bytes.Buffer{}
					annoBytes.WriteString(labelNameIdMap["chair"] + "\n")
					annoBytes.WriteString(labelNameIdMap["chair"] + "\n")
					annoBytes.WriteString(labelNameIdMap["ceiling"] + "\n")

					So(actualAnno, ShouldResemble, annoBytes.Bytes())
					// add more checks here if necessary
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			// val data check
			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/raw-data", AIFS_HOST, valRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File["files"]))
					So("room_2.bin", ShouldEqual, req.MultipartForm.File["files"][0].Filename)

					file, err := req.MultipartForm.File["files"][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()
					fileBytes, err := io.ReadAll(file)
					So(err, ShouldBeNil)
					hashActual, err := utils.GetFileSha256Bytes(fileBytes)
					So(err, ShouldBeNil)
					binBytes, meta, err := handler.(*Points3DHandler).GetBinAndMeta("room_2.bin", "test-data/points-3d/val/raw-data/pos/room_2", "test-data/points-3d/val/raw-data/rgb/room_2")
					So(err, ShouldBeNil)
					hashExpected, err := utils.GetFileSha256Bytes(binBytes)
					binHash = hashExpected
					So(err, ShouldBeNil)
					So(hashActual, ShouldEqual, hashExpected)

					So(meta.Sha256, ShouldEqual, hashExpected)
					So(meta.FileName, ShouldEqual, "room_2.bin")
					So(meta.Size, ShouldEqual, 3)
					So(meta.Xmin, ShouldAlmostEqual, 1.604)
					So(meta.Xmax, ShouldAlmostEqual, 1.633)
					So(meta.Ymin, ShouldAlmostEqual, 1.8390000000000004)
					So(meta.Ymax, ShouldAlmostEqual, 1.9670000000000005)
					So(meta.Zmin, ShouldAlmostEqual, 0.017000000000000015)
					So(meta.Zmax, ShouldAlmostEqual, 0.019000000000000017)
					So(meta.Rmean, ShouldAlmostEqual, 118.33333333333333)
					So(meta.Gmean, ShouldAlmostEqual, 130)
					So(meta.Bmean, ShouldAlmostEqual, 126)
					So(meta.Rstd, ShouldAlmostEqual, 0.9428090415823364)
					So(meta.Gstd, ShouldAlmostEqual, 1.4142135623730951)
					So(meta.Bstd, ShouldAlmostEqual, 1.4142135623730951)
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/data-views/%s/raw-data-hash-list?limit=1000&offset=0", AIFS_HOST, valRawDataViewId),
				func(req *http.Request) (*http.Response, error) {
					So(err, ShouldBeNil)
					return httpmock.NewJsonResponse(200, []any{
						map[string]string{
							"rawDataId": rawDataId,
							"sha256":    binHash,
						},
					})
				},
			)

			httpmock.RegisterResponder("POST", fmt.Sprintf("%s/data-views/%s/annotations", AIFS_HOST, valAnnoViewId),
				func(req *http.Request) (*http.Response, error) {
					req.ParseMultipartForm(32 << 20)
					So(1, ShouldEqual, len(req.MultipartForm.File))
					_, ok := req.MultipartForm.File[rawDataId]
					So(ok, ShouldBeTrue)
					So("room_2", ShouldEqual, req.MultipartForm.File[rawDataId][0].Filename)

					file, err := req.MultipartForm.File[rawDataId][0].Open()
					So(err, ShouldBeNil)
					defer file.Close()

					reader, err := zlib.NewReader(file)
					So(err, ShouldBeNil)

					actualAnno, err := io.ReadAll(reader)
					So(err, ShouldBeNil)

					labelNameIdMap := make(map[string]string, 0)
					for _, label := range annoTempDetails.Labels {
						labelNameIdMap[label.Name] = label.Id
					}

					annoBytes := bytes.Buffer{}
					annoBytes.WriteString(labelNameIdMap["ceiling"] + "\n")
					annoBytes.WriteString(labelNameIdMap["chair"] + "\n")
					annoBytes.WriteString(labelNameIdMap["chair"] + "\n")

					So(actualAnno, ShouldResemble, annoBytes.Bytes())
					// add more checks here if necessary
					return httpmock.NewJsonResponse(200, map[string]any{})
				},
			)

			err = handler.Exec()
			So(err, ShouldBeNil)
		})
	})
}
