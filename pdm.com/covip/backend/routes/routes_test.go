package routes

import (
	"net/http"
	"testing"
	//"github.com/gorilla/mux"
	//"fmt"
	"pdm.com/covip/backend/services"
)

type MockResponseWriter struct {
	t       *testing.T
	headers http.Header
	body    []byte
	status  int
}

func (m *MockResponseWriter) Header() http.Header {
	return m.headers
}

func (m *MockResponseWriter) Write(body []byte) (int, error) {
	m.body = body
	m.status = 200
	return len(body), nil
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.status = statusCode
}

func (r *MockResponseWriter) Assert(status int, body string) {
	if r.status != status {
		r.t.Errorf("expected status %+v to equal %+v", r.status, status)
	}
	if string(r.body) != body {
		r.t.Errorf("expected body %+v to equal %+v", string(r.body), body)
	}
}

func TestGetSummaryEndpoint(t *testing.T) {
	// Setting up mock function.
	oldGetSumFunc := getSumFunc
	getSumFunc = func() ([]services.CovidData, error) {
		return []services.CovidData{}, nil
	}

	mw := &MockResponseWriter{t: t}
	getSummaryEndpoint(mw, nil)

	// Check if successful call was made.
	mw.Assert(200, "[]")
	getSumFunc = oldGetSumFunc
}

func TestGetCountryCasesEndpoint(t *testing.T) {
	oldGetCCases := getCCases
	getCCases = func(name string) ([]services.CovidData, error) {
		return []services.CovidData{}, nil
	}

	r, _ := http.NewRequest("GET", "countryData/test", nil)
	mw := &MockResponseWriter{t: t}
	getCountryCasesEndpoint(mw, r)

	mw.Assert(200, "[]")
	getCCases = oldGetCCases
}

func TestRoutes(t *testing.T) {
	router := Routes()
	if router == nil {
		t.Errorf("expected router to be not nil")
	}
}
