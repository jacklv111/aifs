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
	"github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	vb "github.com/jacklv111/aifs/pkg/store/value-object"
)

type AnnotationLocationResult struct {
	DataViewId string

	ViewType string

	AnnotationTemplateId string

	DataItemDoList []basicdo.DataItemDo

	// key: annotation id
	AnnoDoMap map[uuid.UUID]do.AnnotationDo

	LocationMap map[uuid.UUID]vb.LocationResult
}
