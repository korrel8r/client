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

// Start Start identifies a set of starting objects for correlation.
//
// swagger:model Start
type Start struct {

	// Class for `objects`
	Class string `json:"class,omitempty"`

	// constraint
	Constraint *RestConstraint `json:"constraint,omitempty"`

	// Objects of `class` serialized as JSON
	Objects interface{} `json:"objects,omitempty"`

	// Queries for starting objects
	Queries []string `json:"queries"`
}

// Validate validates this start
func (m *Start) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateConstraint(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Start) validateConstraint(formats strfmt.Registry) error {
	if swag.IsZero(m.Constraint) { // not required
		return nil
	}

	if m.Constraint != nil {
		if err := m.Constraint.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("constraint")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("constraint")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this start based on the context it is used
func (m *Start) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateConstraint(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Start) contextValidateConstraint(ctx context.Context, formats strfmt.Registry) error {

	if m.Constraint != nil {

		if swag.IsZero(m.Constraint) { // not required
			return nil
		}

		if err := m.Constraint.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("constraint")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("constraint")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Start) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Start) UnmarshalBinary(b []byte) error {
	var res Start
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
