/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import (
	"database/sql"

	"github.com/google/uuid"
)

const (
	TABLE_POINTS_3D_EXT = "points_3d_exts"
)

type Points3DExtDo struct {
	ID     uuid.UUID `gorm:"primaryKey;<-:create"`
	Sha256 string    `gorm:"uniqueIndex:points3d_sha256"`
	Size   int64
	Xmin   float64
	Xmax   float64
	Ymin   float64
	Ymax   float64
	Zmin   float64
	Zmax   float64
	Rmean  sql.NullFloat64
	Gmean  sql.NullFloat64
	Bmean  sql.NullFloat64
	Rstd   sql.NullFloat64
	Gstd   sql.NullFloat64
	Bstd   sql.NullFloat64
}

func (Points3DExtDo) TableName() string {
	return TABLE_POINTS_3D_EXT
}

func GetHashList(points3DExtDoList []Points3DExtDo) []string {
	var sha256List []string
	for _, data := range points3DExtDoList {
		sha256List = append(sha256List, data.Sha256)
	}
	return sha256List
}
