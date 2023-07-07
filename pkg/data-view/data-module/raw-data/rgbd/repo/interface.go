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
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/rgbd/do"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type RgbdRawDataRepoInterface interface {
	// CreateBatch 批量插入 rgbd 相关数据。如果该数据已经存在，则将其 id 替换成存在的 item id。判断数据是否相同使用 hash。
	//
	//  @param dataItemDoList
	//  @param rgbdExtDoList
	//  @return error
	CreateBatch(dataItemDoList []basicdo.DataItemDo, rgbdExtDoList []do.RgbdExtDo) error

	// FindExistedByHash 找出 hash code 在列表中存在的 data item
	//
	//  @param []string hash list
	//  @return map[string]uuid.UUID key: hash code; value: data item id
	//  @return error
	FindExistedByHash([]string) (map[string]uuid.UUID, error)

	// DeleteBatch 批量删除
	//  @param idList
	//  @return error
	DeleteBatch(idList []uuid.UUID) error

	// GetHashList
	//  @param dataItemIdList
	//  @return res
	//  @return err
	GetHashList(dataItemIdList []uuid.UUID) (res []basicdo.IdHash, err error)
}
