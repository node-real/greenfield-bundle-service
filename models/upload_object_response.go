// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// UploadObjectResponse upload object response
//
// swagger:model UploadObjectResponse
type UploadObjectResponse struct {

	// The name of the bundle where the file has been uploaded
	BundleName string `json:"bundleName"`
}

// Validate validates this upload object response
func (m *UploadObjectResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this upload object response based on context it is used
func (m *UploadObjectResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *UploadObjectResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UploadObjectResponse) UnmarshalBinary(b []byte) error {
	var res UploadObjectResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
