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
	TABLE_MODEL_EXTS = "model_exts"
)

type ModelExtDo struct {
	ID     uuid.UUID `gorm:"primaryKey;<-:create"`
	Sha256 string
	// byte
	Size int
}

func (ModelExtDo) TableName() string {
	return TABLE_MODEL_EXTS
}
