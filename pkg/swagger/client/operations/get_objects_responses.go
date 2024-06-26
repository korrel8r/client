// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// GetObjectsReader is a Reader for the GetObjects structure.
type GetObjectsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetObjectsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetObjectsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetObjectsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetObjectsOK creates a GetObjectsOK with default headers values
func NewGetObjectsOK() *GetObjectsOK {
	return &GetObjectsOK{}
}

/*
GetObjectsOK describes a response with status code 200, with default header values.

OK
*/
type GetObjectsOK struct {
	Payload []interface{}
}

// IsSuccess returns true when this get objects o k response has a 2xx status code
func (o *GetObjectsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get objects o k response has a 3xx status code
func (o *GetObjectsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get objects o k response has a 4xx status code
func (o *GetObjectsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get objects o k response has a 5xx status code
func (o *GetObjectsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get objects o k response a status code equal to that given
func (o *GetObjectsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get objects o k response
func (o *GetObjectsOK) Code() int {
	return 200
}

func (o *GetObjectsOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /objects][%d] getObjectsOK %s", 200, payload)
}

func (o *GetObjectsOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /objects][%d] getObjectsOK %s", 200, payload)
}

func (o *GetObjectsOK) GetPayload() []interface{} {
	return o.Payload
}

func (o *GetObjectsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetObjectsDefault creates a GetObjectsDefault with default headers values
func NewGetObjectsDefault(code int) *GetObjectsDefault {
	return &GetObjectsDefault{
		_statusCode: code,
	}
}

/*
GetObjectsDefault describes a response with status code -1, with default header values.

GetObjectsDefault get objects default
*/
type GetObjectsDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this get objects default response has a 2xx status code
func (o *GetObjectsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get objects default response has a 3xx status code
func (o *GetObjectsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get objects default response has a 4xx status code
func (o *GetObjectsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get objects default response has a 5xx status code
func (o *GetObjectsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get objects default response a status code equal to that given
func (o *GetObjectsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get objects default response
func (o *GetObjectsDefault) Code() int {
	return o._statusCode
}

func (o *GetObjectsDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /objects][%d] GetObjects default %s", o._statusCode, payload)
}

func (o *GetObjectsDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /objects][%d] GetObjects default %s", o._statusCode, payload)
}

func (o *GetObjectsDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *GetObjectsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
