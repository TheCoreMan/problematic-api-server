/*
 * Rate Limit API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package server

type CooldownResponse struct {
	// Cooldown in seconds
	Cooldown int32 `json:"cooldown,omitempty"`
}

// AssertCooldownResponseRequired checks if the required fields are not zero-ed
func AssertCooldownResponseRequired(obj CooldownResponse) error {
	return nil
}

// AssertRecurseCooldownResponseRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of CooldownResponse (e.g. [][]CooldownResponse), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseCooldownResponseRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aCooldownResponse, ok := obj.(CooldownResponse)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertCooldownResponseRequired(aCooldownResponse)
	})
}
