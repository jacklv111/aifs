/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"database/sql"
	"io"

	"github.com/google/uuid"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	rawdatatype "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/constant"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/points-3d/do"
	p3dvb "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/points-3d/value-object"
)

func BuildWithReader(reader io.ReadSeeker, fileName string, meta p3dvb.Points3DMeta) basicbo.DataInterface {
	rawDataId := uuid.New()
	return &Points3DBo{
		DataBaseImpl: basicbo.DataBaseImpl{
			DataItemDo: basicdo.DataItemDo{ID: rawDataId, Type: rawdatatype.POINTS_3D, Name: fileName},
			ReadSeeker: reader,
		},
		Points3DExtDo: do.Points3DExtDo{
			ID:     rawDataId,
			Sha256: meta.Sha256,
			Size:   meta.Size,
			Xmin:   meta.Xmin,
			Xmax:   meta.Xmax,
			Ymin:   meta.Ymin,
			Ymax:   meta.Ymax,
			Zmin:   meta.Zmin,
			Zmax:   meta.Zmax,
			Rmean:  sql.NullFloat64{Float64: meta.Rmean, Valid: meta.Rmean >= 0},
			Gmean:  sql.NullFloat64{Float64: meta.Gmean, Valid: meta.Gmean >= 0},
			Bmean:  sql.NullFloat64{Float64: meta.Bmean, Valid: meta.Bmean >= 0},
			Rstd:   sql.NullFloat64{Float64: meta.Rstd, Valid: meta.Rstd >= 0},
			Gstd:   sql.NullFloat64{Float64: meta.Gstd, Valid: meta.Gstd >= 0},
			Bstd:   sql.NullFloat64{Float64: meta.Bstd, Valid: meta.Bstd >= 0},
		},
	}
}
