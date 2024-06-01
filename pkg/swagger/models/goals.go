// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Goals Starting point for a goals search.
//
// swagger:model Goals
type Goals struct {

	// Goal classes for correlation.
	// Example: ["domain:class"]
	Goals []string `json:"goals"`

	// start
	Start *Start `json:"start,omitempty"`
}

// Validate validates this goals
func (m *Goals) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateStart(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Goals) validateStart(formats strfmt.Registry) error {
	if swag.IsZero(m.Start) { // not required
		return nil
	}

	if m.Start != nil {
		if err := m.Start.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("start")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("start")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this goals based on the context it is used
func (m *Goals) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateStart(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Goals) contextValidateStart(ctx context.Context, formats strfmt.Registry) error {

	if m.Start != nil {

		if swag.IsZero(m.Start) { // not required
			return nil
		}

		if err := m.Start.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("start")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("start")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Goals) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Goals) UnmarshalBinary(b []byte) error {
	var res Goals
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}