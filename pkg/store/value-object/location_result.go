/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

import "github.com/google/uuid"

type LocationResult struct {
	ID        uuid.UUID
	ObjectKey string
}
