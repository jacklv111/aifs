/*
 * Aifs api
 *
 * aifs api
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type CreateDataViewRequest struct {

	// the name of the data view
	DataViewName string `json:"dataViewName,omitempty"`

	// the description of the data view
	Description string `json:"description,omitempty"`

	ViewType DataViewType `json:"viewType,omitempty"`

	RawDataType RawDataType `json:"rawDataType,omitempty"`

	ZipFormat ZipFormat `json:"zipFormat,omitempty"`

	// If it is an annotation type data view, it must have a related raw-data data view
	RelatedDataViewId string `json:"relatedDataViewId,omitempty"`

	// If it is an annotation type data view, it must have a related annotation template id. If it is a dataset-zip data view, it can have an annotation template id to indicate the annotation template of the annotation data.
	AnnotationTemplateId string `json:"annotationTemplateId,omitempty"`

	// If it is a dataset-zip type data view, it can have a raw data view id to upload raw data to the data view
	RawDataViewId string `json:"rawDataViewId,omitempty"`

	// If it is a dataset-zip type data view, it can have a annotation view id to upload annotation data to the data view
	AnnotationViewId string `json:"annotationViewId,omitempty"`
}
