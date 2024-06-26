/*
 * Aifs api
 *
 * aifs api
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type MoveDataViewItemsRequest struct {

	// the id of the source data view
	SrcDataViewId string `json:"srcDataViewId,omitempty"`

	// the id of the destination data view
	DstDataViewId string `json:"dstDataViewId,omitempty"`
}
