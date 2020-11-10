// Code generated by go-swagger; DO NOT EDIT.

package release_manager

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	"openpitrix.io/openpitrix/test/models"
)

// NewRollbackReleaseParams creates a new RollbackReleaseParams object
// with the default values initialized.
func NewRollbackReleaseParams() *RollbackReleaseParams {
	var ()
	return &RollbackReleaseParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewRollbackReleaseParamsWithTimeout creates a new RollbackReleaseParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewRollbackReleaseParamsWithTimeout(timeout time.Duration) *RollbackReleaseParams {
	var ()
	return &RollbackReleaseParams{

		timeout: timeout,
	}
}

// NewRollbackReleaseParamsWithContext creates a new RollbackReleaseParams object
// with the default values initialized, and the ability to set a context for a request
func NewRollbackReleaseParamsWithContext(ctx context.Context) *RollbackReleaseParams {
	var ()
	return &RollbackReleaseParams{

		Context: ctx,
	}
}

// NewRollbackReleaseParamsWithHTTPClient creates a new RollbackReleaseParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewRollbackReleaseParamsWithHTTPClient(client *http.Client) *RollbackReleaseParams {
	var ()
	return &RollbackReleaseParams{
		HTTPClient: client,
	}
}

/*RollbackReleaseParams contains all the parameters to send to the API endpoint
for the rollback release operation typically these are written to a http.Request
*/
type RollbackReleaseParams struct {

	/*Body*/
	Body *models.OpenpitrixRollbackReleaseRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the rollback release params
func (o *RollbackReleaseParams) WithTimeout(timeout time.Duration) *RollbackReleaseParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the rollback release params
func (o *RollbackReleaseParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the rollback release params
func (o *RollbackReleaseParams) WithContext(ctx context.Context) *RollbackReleaseParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the rollback release params
func (o *RollbackReleaseParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the rollback release params
func (o *RollbackReleaseParams) WithHTTPClient(client *http.Client) *RollbackReleaseParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the rollback release params
func (o *RollbackReleaseParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the rollback release params
func (o *RollbackReleaseParams) WithBody(body *models.OpenpitrixRollbackReleaseRequest) *RollbackReleaseParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the rollback release params
func (o *RollbackReleaseParams) SetBody(body *models.OpenpitrixRollbackReleaseRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *RollbackReleaseParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}