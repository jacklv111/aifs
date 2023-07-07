/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

import (
	"github.com/google/uuid"
	basicDo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	valueobject "github.com/jacklv111/aifs/pkg/store/value-object"
)

type RawDataLocationResult struct {
	DataViewId string

	ViewType string

	RawDataType string

	DataItemDoList []basicDo.DataItemDo

	LocationMap map[uuid.UUID]valueobject.LocationResult
}
