/*
 * Aifs api
 *
 * aifs api
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type S3StorageInner struct {

	DataItemId string `json:"dataItemId,omitempty"`

	DataName string `json:"dataName,omitempty"`

	// the object key of the data in s3 storage
	ObjectKey string `json:"objectKey,omitempty"`
}
