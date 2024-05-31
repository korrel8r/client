// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Graph Graph resulting from a correlation search.
//
// swagger:model Graph
type Graph struct {

	// edges
	Edges []*Edge `json:"edges"`

	// nodes
	Nodes []*Node `json:"nodes"`
}

// Validate validates this graph
func (m *Graph) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEdges(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNodes(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Graph) validateEdges(formats strfmt.Registry) error {
	if swag.IsZero(m.Edges) { // not required
		return nil
	}

	for i := 0; i < len(m.Edges); i++ {
		if swag.IsZero(m.Edges[i]) { // not required
			continue
		}

		if m.Edges[i] != nil {
			if err := m.Edges[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("edges" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("edges" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Graph) validateNodes(formats strfmt.Registry) error {
	if swag.IsZero(m.Nodes) { // not required
		return nil
	}

	for i := 0; i < len(m.Nodes); i++ {
		if swag.IsZero(m.Nodes[i]) { // not required
			continue
		}

		if m.Nodes[i] != nil {
			if err := m.Nodes[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("nodes" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("nodes" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this graph based on the context it is used
func (m *Graph) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateEdges(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNodes(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Graph) contextValidateEdges(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Edges); i++ {

		if m.Edges[i] != nil {

			if swag.IsZero(m.Edges[i]) { // not required
				return nil
			}

			if err := m.Edges[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("edges" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("edges" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Graph) contextValidateNodes(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Nodes); i++ {

		if m.Nodes[i] != nil {

			if swag.IsZero(m.Nodes[i]) { // not required
				return nil
			}

			if err := m.Nodes[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("nodes" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("nodes" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Graph) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Graph) UnmarshalBinary(b []byte) error {
	var res Graph
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
