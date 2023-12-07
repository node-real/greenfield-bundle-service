// Code generated by go-swagger; DO NOT EDIT.

package bundle

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/node-real/greenfield-bundle-service/models"
)

// ViewBundleObjectOKCode is the HTTP code returned for type ViewBundleObjectOK
const ViewBundleObjectOKCode int = 200

/*
ViewBundleObjectOK Successfully retrieved file

swagger:response viewBundleObjectOK
*/
type ViewBundleObjectOK struct {

	/*
	  In: Body
	*/
	Payload io.ReadCloser `json:"body,omitempty"`
}

// NewViewBundleObjectOK creates ViewBundleObjectOK with default headers values
func NewViewBundleObjectOK() *ViewBundleObjectOK {

	return &ViewBundleObjectOK{}
}

// WithPayload adds the payload to the view bundle object o k response
func (o *ViewBundleObjectOK) WithPayload(payload io.ReadCloser) *ViewBundleObjectOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the view bundle object o k response
func (o *ViewBundleObjectOK) SetPayload(payload io.ReadCloser) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ViewBundleObjectOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ViewBundleObjectNotFoundCode is the HTTP code returned for type ViewBundleObjectNotFound
const ViewBundleObjectNotFoundCode int = 404

/*
ViewBundleObjectNotFound Bundle or object not found

swagger:response viewBundleObjectNotFound
*/
type ViewBundleObjectNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewViewBundleObjectNotFound creates ViewBundleObjectNotFound with default headers values
func NewViewBundleObjectNotFound() *ViewBundleObjectNotFound {

	return &ViewBundleObjectNotFound{}
}

// WithPayload adds the payload to the view bundle object not found response
func (o *ViewBundleObjectNotFound) WithPayload(payload *models.Error) *ViewBundleObjectNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the view bundle object not found response
func (o *ViewBundleObjectNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ViewBundleObjectNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ViewBundleObjectInternalServerErrorCode is the HTTP code returned for type ViewBundleObjectInternalServerError
const ViewBundleObjectInternalServerErrorCode int = 500

/*
ViewBundleObjectInternalServerError Internal server error

swagger:response viewBundleObjectInternalServerError
*/
type ViewBundleObjectInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewViewBundleObjectInternalServerError creates ViewBundleObjectInternalServerError with default headers values
func NewViewBundleObjectInternalServerError() *ViewBundleObjectInternalServerError {

	return &ViewBundleObjectInternalServerError{}
}

// WithPayload adds the payload to the view bundle object internal server error response
func (o *ViewBundleObjectInternalServerError) WithPayload(payload *models.Error) *ViewBundleObjectInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the view bundle object internal server error response
func (o *ViewBundleObjectInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ViewBundleObjectInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
