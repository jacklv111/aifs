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
	"github.com/jacklv111/aifs/pkg/store/do"
	"github.com/jacklv111/common-sdk/collection/mapset"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type locationRepoInterface interface {
	// 创建位置记录
	Create([]do.LocationDo) error

	// FindExistedUkey 返回存在记录的数据的 ukey set
	//  @param ukList
	//  @return mapset.Set[do.LocationUkey]
	//  @return error
	FindExistedUkey(ukList []do.LocationUkey) (mapset.Set[do.LocationUkey], error)

	// FindByIdList
	//  @param idList
	//  @return result key: location ukey, value location do
	//  @return err
	FindByIdList(idList []uuid.UUID) (result map[do.LocationUkey]do.LocationDo, err error)

	// DeleteByUk
	//  @param doList
	DeleteByUk(doList []do.LocationUkey) error
}
