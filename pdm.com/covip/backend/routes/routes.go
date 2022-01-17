package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"pdm.com/covip/backend/services"
)

func verifyUserPass(username, password string) bool {
	return true //FIX ME!Returning true for all users
}

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if ok && verifyUserPass(user, pass) {
			next.ServeHTTP(w, r)
			return
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}

var getSumFunc = services.GetSummary

func getSummaryEndpoint(response http.ResponseWriter, request *http.Request) {
	//fmt.Fprintf(w, "You get to see the secret\n")
	codata := getSumFunc()
	//fmt.Println("In GetSummary Endpoint", codata)
	jsonResponse, err := json.Marshal(codata)
	if err != nil {
		fmt.Println("Json response not obtained", err)
	}
	//fmt.Println("Json response is:", jsonResponse)
	response.Write(jsonResponse)
}

var getCCases = services.GetCountryCases

func getCountryCasesEndpoint(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	name := params["name"]
	countrycount := getCCases(name)
	jsonResponse, err := json.Marshal(countrycount)
	if err != nil {
		fmt.Println("Json response not obtained", err)
	}
	response.Write(jsonResponse)
}

func Routes() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/summary", basicAuth(getSummaryEndpoint)).Methods("GET")

	router.HandleFunc("/countryData/{name}", getCountryCasesEndpoint).Methods("GET")

	return router
}
