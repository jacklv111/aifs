/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import "github.com/google/uuid"

const (
	TABLE_IMAGE_EXT = "image_exts"
)

// 更多 image 的元数据存在 ext 表中
type ImageExtDo struct {
	ID        uuid.UUID `gorm:"primaryKey;<-:create"`
	Thumbnail uuid.UUID
	// bytes
	Size   int64
	Sha256 string `gorm:"uniqueIndex:image_sha256"`
	Width  int32
	Height int32
}

func (ImageExtDo) TableName() string {
	return TABLE_IMAGE_EXT
}

func GetHashList(imageExtDoList []ImageExtDo) []string {
	var sha256List []string
	for _, data := range imageExtDoList {
		sha256List = append(sha256List, data.Sha256)
	}
	return sha256List
}
