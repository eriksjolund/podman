// Code generated by go-swagger; DO NOT EDIT.

//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package entries

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

	"github.com/sigstore/rekor/pkg/generated/models"
)

// NewSearchLogQueryParams creates a new SearchLogQueryParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSearchLogQueryParams() *SearchLogQueryParams {
	return &SearchLogQueryParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSearchLogQueryParamsWithTimeout creates a new SearchLogQueryParams object
// with the ability to set a timeout on a request.
func NewSearchLogQueryParamsWithTimeout(timeout time.Duration) *SearchLogQueryParams {
	return &SearchLogQueryParams{
		timeout: timeout,
	}
}

// NewSearchLogQueryParamsWithContext creates a new SearchLogQueryParams object
// with the ability to set a context for a request.
func NewSearchLogQueryParamsWithContext(ctx context.Context) *SearchLogQueryParams {
	return &SearchLogQueryParams{
		Context: ctx,
	}
}

// NewSearchLogQueryParamsWithHTTPClient creates a new SearchLogQueryParams object
// with the ability to set a custom HTTPClient for a request.
func NewSearchLogQueryParamsWithHTTPClient(client *http.Client) *SearchLogQueryParams {
	return &SearchLogQueryParams{
		HTTPClient: client,
	}
}

/*
SearchLogQueryParams contains all the parameters to send to the API endpoint

	for the search log query operation.

	Typically these are written to a http.Request.
*/
type SearchLogQueryParams struct {

	// Entry.
	Entry *models.SearchLogQuery

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the search log query params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchLogQueryParams) WithDefaults() *SearchLogQueryParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the search log query params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchLogQueryParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the search log query params
func (o *SearchLogQueryParams) WithTimeout(timeout time.Duration) *SearchLogQueryParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the search log query params
func (o *SearchLogQueryParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the search log query params
func (o *SearchLogQueryParams) WithContext(ctx context.Context) *SearchLogQueryParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the search log query params
func (o *SearchLogQueryParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the search log query params
func (o *SearchLogQueryParams) WithHTTPClient(client *http.Client) *SearchLogQueryParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the search log query params
func (o *SearchLogQueryParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithEntry adds the entry to the search log query params
func (o *SearchLogQueryParams) WithEntry(entry *models.SearchLogQuery) *SearchLogQueryParams {
	o.SetEntry(entry)
	return o
}

// SetEntry adds the entry to the search log query params
func (o *SearchLogQueryParams) SetEntry(entry *models.SearchLogQuery) {
	o.Entry = entry
}

// WriteToRequest writes these params to a swagger request
func (o *SearchLogQueryParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Entry != nil {
		if err := r.SetBodyParam(o.Entry); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}