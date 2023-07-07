/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import "github.com/google/uuid"

type IdHash struct {
	ID     uuid.UUID
	Sha256 string
}
