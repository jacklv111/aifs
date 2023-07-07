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
	TABLE_RGBD_EXT = "rgbd_exts"
)

type RgbdExtDo struct {
	ID     uuid.UUID `gorm:"primaryKey;<-:create"`
	Sha256 string    `gorm:"uniqueIndex:rgbd_sha256"`
	// bytes
	ImageSize   int64
	ImageWidth  int32
	ImageHeight int32
	// bytes
	DepthSize   int64
	DepthWidth  int32
	DepthHeight int32
}

func (RgbdExtDo) TableName() string {
	return TABLE_RGBD_EXT
}

func GetHashList(rgbdExtDoList []RgbdExtDo) []string {
	var sha256List []string
	for _, data := range rgbdExtDoList {
		sha256List = append(sha256List, data.Sha256)
	}
	return sha256List
}
