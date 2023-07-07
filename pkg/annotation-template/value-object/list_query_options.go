/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

// 查询列表参数
type ListQueryOptions struct {
	Offset             int
	Limit              int
	AnnoTemplateIdList []string
}
