/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

import "github.com/google/uuid"

// 存储请求的输入参数
type DeleteParams struct {
	deleteIdList []uuid.UUID
}

func (params *DeleteParams) AddItem(id uuid.UUID) {
	params.deleteIdList = append(params.deleteIdList, id)
}

func (params *DeleteParams) GetIdList() []uuid.UUID {
	return params.deleteIdList
}

func (params *DeleteParams) IsEmpty() bool {
	return len(params.deleteIdList) == 0
}
