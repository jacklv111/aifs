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
	vb "github.com/jacklv111/aifs/pkg/store/value-object"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type StoreMgrInterface interface {

	// Save 根据真 store 管辖范围的实集群环境的特点，选择合适的存储介质进行数据的存储，将数据保存到 store 模块管辖范围内的存储区域，并留下存储的元信息
	//
	//  @param params
	//  @return error
	Upload(params vb.UploadParams) error

	// Download 将数据下载到某个环境
	//
	//	@receiver mgr
	//	@param params
	//	@return error
	Download(params vb.DownloadParams) error

	// Delete 删除 store 中的数据
	//
	//  @param params
	//  @return error
	Delete(params vb.DeleteParams) error

	// GetByIdList
	//
	//  @param idList
	//  @return map[uuid.UUID]valueobject.LocationResult
	//  @return error
	GetByIdList(idList []uuid.UUID) (map[uuid.UUID]vb.LocationResult, error)

	// GetUrlListByIdList
	//
	//  @param idList
	//  @return map[uuid.UUID]string
	//  @return error
	GetUrlListByIdList(idList []uuid.UUID) (map[uuid.UUID]string, error)
}
