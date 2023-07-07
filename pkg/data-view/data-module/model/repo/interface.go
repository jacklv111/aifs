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
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/model/do"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type modelRepoInterface interface {
	CreateBatch(dataItemDoList []basicdo.DataItemDo, modelExts []do.ModelExtDo) error
	DeleteBatch(idList []uuid.UUID) error
	GetByIdList(idList []uuid.UUID) (dataItemDoList []basicdo.DataItemDo, modelExts []do.ModelExtDo, err error)
}
