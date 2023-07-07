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
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	vb "github.com/jacklv111/aifs/pkg/store/value-object"
)

type DatasetZipLocationResult struct {
	DataViewId string

	ViewType string

	DataItemDoList []basicdo.DataItemDo

	LocationMap map[uuid.UUID]vb.LocationResult
}
