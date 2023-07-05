/*
 * Aifs api
 *
 * aifs api
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type LabelDistribution struct {

	// the label id
	LabelId string `json:"labelId,omitempty"`

	// the number of the label in the data view.
	Count int32 `json:"count,omitempty"`

	// the ratio of the label in the data view.
	Ratio float32 `json:"ratio,omitempty"`
}
