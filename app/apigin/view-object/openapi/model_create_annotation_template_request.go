/*
 * Aifs api
 *
 * aifs api
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type CreateAnnotationTemplateRequest struct {

	// name of the annotation template
	Name string `json:"name"`

	// the type of the annotation template
	Type string `json:"type"`

	// the description of the annotation template
	Description string `json:"description,omitempty"`

	Labels []Label `json:"labels,omitempty"`

	WordList []string `json:"wordList,omitempty"`
}
