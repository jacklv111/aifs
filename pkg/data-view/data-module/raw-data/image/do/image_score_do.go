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
	TABLE_IMAGE_SCORE = "image_scores"
)

// 图片评分
type ImageScoreDo struct {
	ID uuid.UUID `gorm:"primaryKey;<-:create"`
	// 光照评分/%
	Light float32
	// 密集度评分/%
	Dense float32
	// 遮挡评分/%
	Shelter float32
	// 目标大小评分/%
	Size float32
}

func (ImageScoreDo) TableName() string {
	return TABLE_IMAGE_SCORE
}
