// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.10.1 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// Pod defines model for pod.
type Pod struct {
	// the uptime age of the pod as a string
	Age string `json:"age"`

	// the name assigned to this pod by the cluster
	Name string `json:"name"`

	// the number of times this pod has been restarted since it was first deployed
	Restarts int `json:"restarts"`
}

// indicates the number of records per response
type RecordLimit int

// indicates the number of records a response is offset by, given the sorting of the result set
type RecordOffset int

// indicates the total number of records that can be returned in a query
type TotalRecords int

// NamespacePath defines model for namespacePath.
type NamespacePath string

// Greeting defines model for Greeting.
type Greeting struct {
	// The full response given the input name.
	Greeting string `json:"greeting"`
}

// ListPodsResponse defines model for ListPodsResponse.
type ListPodsResponse struct {
	// indicates the number of records per response
	Limit RecordLimit `json:"limit"`

	// indicates the number of records a response is offset by, given the sorting of the result set
	Offset RecordOffset `json:"offset"`
	Pods   []Pod        `json:"pods"`

	// indicates the total number of records that can be returned in a query
	Total TotalRecords `json:"total"`
}

// ListPodsParams defines parameters for ListPods.
type ListPodsParams struct {
	// allows callers to specify the number of records to return per result set. Default is 10
	Limit *RecordLimit `json:"limit,omitempty"`

	// allows callers to specify how many records to start from, given the offset. Default is 0.
	Offset *RecordOffset `json:"offset,omitempty"`

	// allows callers to sort the result set. Default is name. Current order is ascending.
	Sort *ListPodsParamsSort `json:"sort,omitempty"`
}

// ListPodsParamsSort defines parameters for ListPods.
type ListPodsParamsSort string

// SayHelloParams defines parameters for SayHello.
type SayHelloParams struct {
	// The name for which to say Hello to. 'World' is the default if not provided or is empty.
	Name *string `json:"name,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List the pods in this namespace.
	// (GET /cluster/{namespace}/pod)
	ListPods(w http.ResponseWriter, r *http.Request, namespace NamespacePath, params ListPodsParams)
	// Says hello to whomever you provide in the query parameter.
	// (GET /hello)
	SayHello(w http.ResponseWriter, r *http.Request, params SayHelloParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// ListPods operation middleware
func (siw *ServerInterfaceWrapper) ListPods(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "namespace" -------------
	var namespace NamespacePath

	err = runtime.BindStyledParameter("simple", false, "namespace", chi.URLParam(r, "namespace"), &namespace)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "namespace", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params ListPodsParams

	// ------------- Optional query parameter "limit" -------------
	if paramValue := r.URL.Query().Get("limit"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	// ------------- Optional query parameter "offset" -------------
	if paramValue := r.URL.Query().Get("offset"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "offset", r.URL.Query(), &params.Offset)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "offset", Err: err})
		return
	}

	// ------------- Optional query parameter "sort" -------------
	if paramValue := r.URL.Query().Get("sort"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "sort", r.URL.Query(), &params.Sort)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sort", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListPods(w, r, namespace, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// SayHello operation middleware
func (siw *ServerInterfaceWrapper) SayHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params SayHelloParams

	// ------------- Optional query parameter "name" -------------
	if paramValue := r.URL.Query().Get("name"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "name", r.URL.Query(), &params.Name)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SayHello(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/cluster/{namespace}/pod", wrapper.ListPods)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/hello", wrapper.SayHello)
	})

	return r
}
