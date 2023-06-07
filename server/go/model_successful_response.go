/*
 * Rate Limit API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package server

type SuccessfulResponse struct {
	BookName string `json:"book-name,omitempty"`

	LineNumber int32 `json:"line-number,omitempty"`

	Text string `json:"text,omitempty"`
}

// AssertSuccessfulResponseRequired checks if the required fields are not zero-ed
func AssertSuccessfulResponseRequired(obj SuccessfulResponse) error {
	return nil
}

// AssertRecurseSuccessfulResponseRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of SuccessfulResponse (e.g. [][]SuccessfulResponse), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseSuccessfulResponseRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aSuccessfulResponse, ok := obj.(SuccessfulResponse)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertSuccessfulResponseRequired(aSuccessfulResponse)
	})
}
