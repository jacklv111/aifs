/*
 * Aifs api
 *
 * aifs api
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type AnnotationViewData struct {

	DataViewId string `json:"dataViewId,omitempty"`

	AnnotationTemplateId string `json:"annotationTemplateId,omitempty"`

	ViewType DataViewType `json:"viewType,omitempty"`

	DataItems []AnnotationDataInner `json:"dataItems,omitempty"`
}
