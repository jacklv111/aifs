/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"github.com/google/uuid"
	annodo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type annotationRepoInterface interface {
	// CreateBatch 批量插入 annotation 数据
	//  @param dataItemDoList
	//  @param annoDoList
	//  @param rawDataLabelList
	//  @return error
	CreateBatch(dataItemDoList []basicdo.DataItemDo, annoDoList []annodo.AnnotationDo, rawDataLabelList []annodo.RawDataLabelDo) error

	// DeleteBatch
	//  @param idList
	//  @return error
	DeleteBatch(idList []uuid.UUID) error

	// GetByIdList
	//  @param idList
	//  @return dataItemDoList
	//  @return annoDoList
	//  @return err
	GetByIdList(idList []uuid.UUID) (dataItemDoList []basicdo.DataItemDo, annoDoList []annodo.AnnotationDo, err error)

	// GetRawDataLabelByIdList
	//  @param idList
	//  @return rawDataLabelList
	//  @return err
	GetRawDataLabelByIdList(idList []uuid.UUID) (rawDataLabelList []annodo.RawDataLabelDo, err error)
}
