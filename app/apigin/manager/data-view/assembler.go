/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package manager

import (
	"github.com/google/uuid"
	"github.com/jacklv111/aifs/app/apigin/view-object/openapi"
	datamodule "github.com/jacklv111/aifs/pkg/data-view/data-module"
	annodo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	dvdo "github.com/jacklv111/aifs/pkg/data-view/data-view/do"
	vb "github.com/jacklv111/aifs/pkg/data-view/data-view/value-object"
	"github.com/jacklv111/common-sdk/log"
)

func assembleToDataViewDetails(details vb.DataViewDetails, annoType string) openapi.DataViewDetails {
	res := openapi.DataViewDetails{
		Id:          details.ID.String(),
		Name:        details.Name,
		ViewType:    openapi.DataViewType(details.ViewType),
		Description: details.Description,
		CreateAt:    details.CreateAt,
	}
	if details.DataViewDo.ViewType == datamodule.MODEL {
		res.Progress = float32(details.Progress)
		res.CommitId = details.CommitId.String
	}
	if details.DataViewDo.ViewType == datamodule.DATASET_ZIP {
		res.ZipFormat = openapi.ZipFormat(details.ZipFormat.String)
		res.Progress = float32(details.Progress)
		res.Status = details.Status.String
		res.TrainRawDataViewId = details.TrainRawDataViewId.String
		res.TrainAnnotationViewId = details.TrainAnnotationViewId.String
		res.ValRawDataViewId = details.ValRawDataViewId.String
		res.ValAnnotationViewId = details.ValAnnotationViewId.String
		res.RawDataViewId = details.RawDataViewId.String
		res.AnnotationViewId = details.AnnotationViewId.String
		if details.AnnotationTemplateId != uuid.Nil {
			res.AnnotationTemplateId = details.AnnotationTemplateId.String()
		}
	}
	if details.DataViewDo.ViewType == datamodule.RAW_DATA {
		res.RawDataType = openapi.RawDataType(details.RawDataType)
	}
	if details.DataViewDo.ViewType == datamodule.ANNOTATION {
		res.AnnotationTemplateId = details.AnnotationTemplateId.String()
		res.AnnotationTemplateType = annoType
	}
	return res
}

func assembleToDataViewStatistics(statistics vb.DataViewStatistics) openapi.DataViewStatistics {
	res := openapi.DataViewStatistics{
		ItemCount:     statistics.ItemCount,
		LabelCount:    statistics.LabelCount,
		TotalDataSize: statistics.TotalDataSize,
	}
	for _, ld := range statistics.LabelDistribution {
		res.LabelDistribution = append(res.LabelDistribution, openapi.LabelDistribution{
			LabelId: ld.LabelId,
			Count:   ld.Count,
			Ratio:   ld.Ratio,
		})
	}
	return res
}

func assembleToDataViewListItemList(dolist []dvdo.DataViewDo, annoTypeMap map[uuid.UUID]string) []openapi.DataViewListItem {
	result := make([]openapi.DataViewListItem, 0)
	for _, data := range dolist {
		item := openapi.DataViewListItem{
			Id:       data.ID.String(),
			Name:     data.Name,
			ViewType: openapi.DataViewType(data.ViewType),
			CreateAt: data.CreateAt,
		}
		if data.ViewType == datamodule.ANNOTATION {
			item.AnnotationTemplateId = data.AnnotationTemplateId.String()
			item.AnnotationTemplateType = annoTypeMap[data.ID]
		}
		if data.ViewType == datamodule.RAW_DATA {
			item.RawDataType = openapi.RawDataType(data.RawDataType)
		}
		result = append(result, item)
	}
	return result
}

func assembleRawDataLocations(result vb.RawDataLocationResult) openapi.RawDataViewLocations {
	var resp openapi.RawDataViewLocations
	resp.DataViewId = result.DataViewId
	resp.ViewType = openapi.DataViewType(result.ViewType)
	resp.RawDataType = openapi.RawDataType(result.RawDataType)
	for _, data := range result.DataItemDoList {
		location, ok := result.LocationMap[data.ID]
		if !ok {
			log.Errorf("raw data (id, name) = (%s, %s) doesn't have data location", data.ID.String(), data.Name)
			continue
		}
		resp.DataItems = append(resp.DataItems, openapi.S3StorageInner{
			DataItemId: data.ID.String(),
			DataName:   data.Name,
			ObjectKey:  location.ObjectKey,
		})
	}
	return resp
}

func assembleModelDataLocations(result vb.ModelLocationResult) openapi.ModelDataViewLocations {
	var resp openapi.ModelDataViewLocations
	resp.DataViewId = result.DataViewId
	resp.ViewType = openapi.DataViewType(result.ViewType)
	for _, data := range result.DataItemDoList {
		location, ok := result.LocationMap[data.ID]
		if !ok {
			log.Errorf("model data (id, name) = (%s, %s) doesn't have data location", data.ID.String(), data.Name)
			continue
		}
		resp.DataItems = append(resp.DataItems, openapi.S3StorageInner{
			DataItemId: data.ID.String(),
			DataName:   data.Name,
			ObjectKey:  location.ObjectKey,
		})
	}
	return resp
}

func assembleDatasetZipLocations(result vb.DatasetZipLocationResult) openapi.DatasetZipLocation {
	var resp openapi.DatasetZipLocation
	if len(result.DataItemDoList) == 0 {
		return resp
	}
	resp.DataViewId = result.DataViewId
	resp.ViewType = openapi.DataViewType(result.ViewType)
	data := result.DataItemDoList[0]
	location, ok := result.LocationMap[data.ID]
	if !ok {
		log.Errorf("dataset zip data (id, name) = (%s, %s) doesn't have data location", data.ID.String(), data.Name)
	}

	resp.DataItems = append(resp.DataItems, openapi.S3StorageInner{
		DataItemId: data.ID.String(),
		DataName:   data.Name,
		ObjectKey:  location.ObjectKey,
	})
	return resp
}

func assembleArtifactLocations(result vb.ArtifactLocationResult) openapi.ArtifactLocations {
	var resp openapi.ArtifactLocations
	if len(result.DataItemDoList) == 0 {
		return resp
	}
	resp.DataViewId = result.DataViewId
	resp.ViewType = openapi.DataViewType(result.ViewType)
	for _, data := range result.DataItemDoList {
		location, ok := result.LocationMap[data.ID]
		if !ok {
			log.Errorf("dataset zip data (id, name) = (%s, %s) doesn't have data location", data.ID.String(), data.Name)
		}

		resp.DataItems = append(resp.DataItems, openapi.S3StorageInner{
			DataItemId: data.ID.String(),
			DataName:   data.Name,
			ObjectKey:  location.ObjectKey,
		})
	}

	return resp
}

func assembleAnnotationLocations(result vb.AnnotationLocationResult) openapi.AnnotationViewLocations {
	var resp openapi.AnnotationViewLocations
	resp.DataViewId = result.DataViewId
	resp.ViewType = openapi.DataViewType(result.ViewType)
	resp.AnnotationTemplateId = result.AnnotationTemplateId
	for _, data := range result.DataItemDoList {
		location, ok := result.LocationMap[data.ID]
		if !ok {
			log.Errorf("annotation data (id, name) = (%s, %s) doesn't have data location", data.ID.String(), data.Name)
			continue
		}
		anno, ok := result.AnnoDoMap[data.ID]
		if !ok {
			log.Errorf("annotation (id, name) = (%s, %s) doesn't have annotation data", data.ID.String(), data.Name)
			continue
		}

		resp.DataItems = append(resp.DataItems, openapi.S3StorageInner{
			DataItemId: anno.DataItemId.String(),
			DataName:   data.Name,
			ObjectKey:  location.ObjectKey,
		})
	}
	return resp
}

func assembleAnnotationData(result vb.AnnotationData) openapi.AnnotationViewData {
	var resp openapi.AnnotationViewData
	resp.DataViewId = result.DataViewId
	resp.ViewType = openapi.DataViewType(result.ViewType)
	resp.AnnotationTemplateId = result.AnnotationTemplateId
	for _, data := range result.DataItemDoList {
		anno, ok := result.AnnoDoMap[data.ID]
		if !ok {
			log.Errorf("annotation (id, name) = (%s, %s) doesn't have annotation data", data.ID.String(), data.Name)
			continue
		}

		resp.DataItems = append(resp.DataItems, openapi.AnnotationDataInner{
			DataItemId: anno.DataItemId.String(),
			Labels:     annodo.GetLabelIdStrList(result.RawDataLabelMap[data.ID]),
			TextData:   anno.TextData.String,
		})
	}
	return resp
}

func assembleRawDataHashList(result []basicdo.IdHash) []openapi.RawDataHashListInner {
	resp := make([]openapi.RawDataHashListInner, 0)
	for _, data := range result {
		resp = append(resp, openapi.RawDataHashListInner{RawDataId: data.ID.String(), Sha256: data.Sha256})
	}
	return resp
}
