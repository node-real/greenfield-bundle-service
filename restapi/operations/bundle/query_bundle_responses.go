// Code generated by go-swagger; DO NOT EDIT.

package bundle

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/node-real/greenfield-bundle-service/models"
)

// QueryBundleOKCode is the HTTP code returned for type QueryBundleOK
const QueryBundleOKCode int = 200

/*
QueryBundleOK Successfully queried bundle

swagger:response queryBundleOK
*/
type QueryBundleOK struct {

	/*
	  In: Body
	*/
	Payload *models.QueryBundleResponse `json:"body,omitempty"`
}

// NewQueryBundleOK creates QueryBundleOK with default headers values
func NewQueryBundleOK() *QueryBundleOK {

	return &QueryBundleOK{}
}

// WithPayload adds the payload to the query bundle o k response
func (o *QueryBundleOK) WithPayload(payload *models.QueryBundleResponse) *QueryBundleOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the query bundle o k response
func (o *QueryBundleOK) SetPayload(payload *models.QueryBundleResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *QueryBundleOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// QueryBundleBadRequestCode is the HTTP code returned for type QueryBundleBadRequest
const QueryBundleBadRequestCode int = 400

/*
QueryBundleBadRequest Invalid request or file format

swagger:response queryBundleBadRequest
*/
type QueryBundleBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewQueryBundleBadRequest creates QueryBundleBadRequest with default headers values
func NewQueryBundleBadRequest() *QueryBundleBadRequest {

	return &QueryBundleBadRequest{}
}

// WithPayload adds the payload to the query bundle bad request response
func (o *QueryBundleBadRequest) WithPayload(payload *models.Error) *QueryBundleBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the query bundle bad request response
func (o *QueryBundleBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *QueryBundleBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// QueryBundleNotFoundCode is the HTTP code returned for type QueryBundleNotFound
const QueryBundleNotFoundCode int = 404

/*
QueryBundleNotFound Bundle not found

swagger:response queryBundleNotFound
*/
type QueryBundleNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewQueryBundleNotFound creates QueryBundleNotFound with default headers values
func NewQueryBundleNotFound() *QueryBundleNotFound {

	return &QueryBundleNotFound{}
}

// WithPayload adds the payload to the query bundle not found response
func (o *QueryBundleNotFound) WithPayload(payload *models.Error) *QueryBundleNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the query bundle not found response
func (o *QueryBundleNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *QueryBundleNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// QueryBundleInternalServerErrorCode is the HTTP code returned for type QueryBundleInternalServerError
const QueryBundleInternalServerErrorCode int = 500

/*
QueryBundleInternalServerError Internal server error

swagger:response queryBundleInternalServerError
*/
type QueryBundleInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewQueryBundleInternalServerError creates QueryBundleInternalServerError with default headers values
func NewQueryBundleInternalServerError() *QueryBundleInternalServerError {

	return &QueryBundleInternalServerError{}
}

// WithPayload adds the payload to the query bundle internal server error response
func (o *QueryBundleInternalServerError) WithPayload(payload *models.Error) *QueryBundleInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the query bundle internal server error response
func (o *QueryBundleInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *QueryBundleInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
