/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"github.com/google/uuid"
	valueobject "github.com/jacklv111/aifs/pkg/annotation-template/value-object"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type AnnotationTemplateBoInterface interface {
	// 创建一个 annotation template
	Create() (uuid.UUID, error)
	// 获取 annotation template 详情。
	GetDetails() (*valueobject.AnnotationTemplateDetails, error)
	// 删除 annotation template
	Delete() error
	// 更新 annotation template
	Update() error

	// Sync 从 db 同步数据
	//  @return err db 异常或者数据不存在
	Sync() (err error)

	// HasLabel annotation template 是否存在某个 id 的 label
	//  @param id label id list
	//  @return bool true: 存在；false: 不存在
	HasLabel(id []uuid.UUID) bool

	// GetType get annotation template type
	//  @return string
	GetType() string

	// GetId
	//  @return uuid.UUID
	GetId() uuid.UUID

	// GetLabelIdByColor
	//  @return uuid.UUID
	GetLabelIdByColor(int32) uuid.UUID

	Copy() (AnnotationTemplateBoInterface, error)
}
