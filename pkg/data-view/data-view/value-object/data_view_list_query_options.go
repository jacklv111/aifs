/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

// 查询列表参数
type DataViewListQueryOptions struct {
	Offset         int
	Limit          int
	DataViewIdList []string
	DataViewName   string
}

// HasNameFilter 名字的 filter 是否存在
//
//	@receiver option
//	@return bool
func (option DataViewListQueryOptions) HasNameFilter() bool {
	return option.DataViewName != ""
}

// HasDataViewList dataViewList filter 是否存在
//
//	@receiver option
//	@return bool
func (option DataViewListQueryOptions) HasDataViewList() bool {
	return option.DataViewIdList != nil
}
