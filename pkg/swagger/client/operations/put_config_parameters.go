// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewPutConfigParams creates a new PutConfigParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPutConfigParams() *PutConfigParams {
	return &PutConfigParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPutConfigParamsWithTimeout creates a new PutConfigParams object
// with the ability to set a timeout on a request.
func NewPutConfigParamsWithTimeout(timeout time.Duration) *PutConfigParams {
	return &PutConfigParams{
		timeout: timeout,
	}
}

// NewPutConfigParamsWithContext creates a new PutConfigParams object
// with the ability to set a context for a request.
func NewPutConfigParamsWithContext(ctx context.Context) *PutConfigParams {
	return &PutConfigParams{
		Context: ctx,
	}
}

// NewPutConfigParamsWithHTTPClient creates a new PutConfigParams object
// with the ability to set a custom HTTPClient for a request.
func NewPutConfigParamsWithHTTPClient(client *http.Client) *PutConfigParams {
	return &PutConfigParams{
		HTTPClient: client,
	}
}

/*
PutConfigParams contains all the parameters to send to the API endpoint

	for the put config operation.

	Typically these are written to a http.Request.
*/
type PutConfigParams struct {

	/* Verbose.

	   verbose setting for logging
	*/
	Verbose *int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the put config params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutConfigParams) WithDefaults() *PutConfigParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the put config params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutConfigParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the put config params
func (o *PutConfigParams) WithTimeout(timeout time.Duration) *PutConfigParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the put config params
func (o *PutConfigParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the put config params
func (o *PutConfigParams) WithContext(ctx context.Context) *PutConfigParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the put config params
func (o *PutConfigParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the put config params
func (o *PutConfigParams) WithHTTPClient(client *http.Client) *PutConfigParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the put config params
func (o *PutConfigParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithVerbose adds the verbose to the put config params
func (o *PutConfigParams) WithVerbose(verbose *int64) *PutConfigParams {
	o.SetVerbose(verbose)
	return o
}

// SetVerbose adds the verbose to the put config params
func (o *PutConfigParams) SetVerbose(verbose *int64) {
	o.Verbose = verbose
}

// WriteToRequest writes these params to a swagger request
func (o *PutConfigParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Verbose != nil {

		// query param verbose
		var qrVerbose int64

		if o.Verbose != nil {
			qrVerbose = *o.Verbose
		}
		qVerbose := swag.FormatInt64(qrVerbose)
		if qVerbose != "" {

			if err := r.SetQueryParam("verbose", qVerbose); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
