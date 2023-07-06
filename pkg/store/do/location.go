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
	TABLE_LOCATION = "locations"
)

type LocationDo struct {
	DataItemId  uuid.UUID `gorm:"uniqueIndex:loc_ukey,priority:1;<-:create"`
	BucketName  string
	ObjectKey   string
	Environment string `gorm:"uniqueIndex:loc_ukey,priority:3"`
	Name        string `gorm:"uniqueIndex:loc_ukey,priority:2"`
}

func (LocationDo) TableName() string {
	return TABLE_LOCATION
}
