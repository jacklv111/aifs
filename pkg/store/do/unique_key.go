/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import (
	"github.com/google/uuid"
)

type LocationUkey struct {
	DataItemId  uuid.UUID
	Name        string
	Environment string
}

func GetTupleList(list []LocationUkey) (res [][]interface{}) {
	for _, data := range list {
		var param []interface{}
		param = append(param, data.DataItemId)
		param = append(param, data.Name)
		param = append(param, data.Environment)

		res = append(res, param)
	}
	return
}
