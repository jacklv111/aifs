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
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type BasicDataRepoInterface interface {

	// GetNameAndType 获取 data item 的 name 和 type
	//  @param idList
	//  @return []do.DataItemDo
	//  @return error
	GetNameAndType(idList []uuid.UUID) ([]basicdo.DataItemDo, error)

	// CreateBatch 批量创建
	//  @param dataItemDoList
	//  @return error
	CreateBatch(dataItemDoList []basicdo.DataItemDo) error

	// DeleteBatch 批量删除
	//  @param idList
	//  @return error
	DeleteBatch(idList []uuid.UUID) error

	// GetByIdList 通过 id 获取数据
	//  @param idList
	//  @return dataItemDoList
	//  @return err
	GetByIdList(idList []uuid.UUID) (dataItemDoList []basicdo.DataItemDo, err error)

	// Create 单个创建数据
	//  @param dataItemDo
	//  @return error
	Create(dataItemDo basicdo.DataItemDo) error
}
