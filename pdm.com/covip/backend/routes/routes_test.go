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

/*
func (r *mux.Route) AssertGetSumRoute(sumpath string, summethod string) {
   if path, _ := r.GetPathTemplate(); path != sumpath {
	r.t.Errorf("expected path %+v to equal %+v", string(path), sumpath)
   }
   if method, _ := r.GetMethods(); method != summethod {
	r.t.Errorf("expected path %+v to equal %+v", string(method), summethod)
   }
}
*/
func TestGetSummaryEndpoint(t *testing.T) {
	oldGetSumFunc := getSumFunc
	getSumFunc = func() []services.CovidData {
		return []services.CovidData{}
	}

	mw := &MockResponseWriter{t: t}
	getSummaryEndpoint(mw, nil)
	//nullvar []byte := nil
	mw.Assert(200, "[]")
	getSumFunc = oldGetSumFunc
}

func TestRoutes(t *testing.T) {
	router := Routes()
	if router == nil {
		t.Errorf("expected router to be not nil")
	}
	/*	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			pathTemplate, err := route.GetPathTemplate()
			if err == nil {
				fmt.Println("ROUTE:", pathTemplate)
			}
			methods, err := route.GetMethods()
			if err == nil {
				fmt.Println("Methods:", methods)
			}

			handlerfunc := route.GetHandler()
			if handlerfunc != GetSummaryEndpoint {
				t.Errorf("expected handler function %+v to equal %+v", handlerfunc, GetSummaryEndpoint)
			}

			return nil
		})

		if err != nil {
			fmt.Println(err)
		}

	*/

}
