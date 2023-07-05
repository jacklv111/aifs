/*
 * Aifs api
 *
 * aifs api
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type UploadAnnotationToDataViewRequest1 struct {

	// the folder path of the resource
	ResourcePath string `json:"resourcePath,omitempty"`

	AnnotationTemplateId string `json:"annotationTemplateId,omitempty"`

	Format UploadAnnotationFormat `json:"format,omitempty"`
}
