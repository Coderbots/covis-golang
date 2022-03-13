package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"pdm.com/covip/backend/model"
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

func jsonError(response http.ResponseWriter, status int, msg string) {
	http.Error(response, msg, status)
}

func jsonHandleError(response http.ResponseWriter, err error) {
	var apiErr model.APIError
	if errors.As(err, &apiErr) {
		status, msg := apiErr.APIError()
		jsonError(response, status, msg)
	} else {
		jsonError(response, http.StatusInternalServerError, "internal error")
	}
}

var getSumFunc = services.GetSummary

func getSummaryEndpoint(response http.ResponseWriter, request *http.Request) {
	codata, errSum := getSumFunc()
	if errSum != nil {
		fmt.Println("Error on retrieving Summary:", errSum)
		jsonHandleError(response, errSum)
		return
	}
	jsonResponse, err := json.Marshal(codata)
	if err != nil {
		fmt.Println("Json response not obtained", err)
		jsonHandleError(response, err)
		return
	}
	response.Write(jsonResponse)
}

var getCCases = services.GetCountryCases

func getCountryCasesEndpoint(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	name := params["name"]
	countrycount, errCCases := getCCases(name)
	if errCCases != nil {
		fmt.Println("Error on retrieving Country Cases:", errCCases)
		jsonHandleError(response, errCCases)
		return
	}
	jsonResponse, err := json.Marshal(countrycount)
	if err != nil {
		fmt.Println("Json response not obtained", err)
		jsonHandleError(response, err)
		return
	}
	response.Write(jsonResponse)
}

func Routes() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/summary", basicAuth(getSummaryEndpoint)).Methods("GET")

	router.HandleFunc("/countryData/{name}", getCountryCasesEndpoint).Methods("GET")

	return router
}
