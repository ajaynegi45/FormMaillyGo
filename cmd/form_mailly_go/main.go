package main

import (
	"Form-Mailly-Go/internal"
	"encoding/json"
	"net/http"
)

func main() {
	mux := http.NewServeMux() // Router

	// Routes
	mux.HandleFunc("GET /api/health", checkHealth)
	mux.HandleFunc("POST /api/contact", handleForm)

	// Server
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
		return
	}

}

func checkHealth(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	_, err := res.Write([]byte("Everything is OK"))
	if err != nil {
		panic(err)
		return
	}
}

func handleForm(response http.ResponseWriter, request *http.Request) {

	var formData internal.ContactForm
	err := json.NewDecoder(request.Body).Decode(&formData)
	if err != nil {
		http.Error(response, "", http.StatusInternalServerError)
		return
	}

	var isValid bool = internal.ValidateFormData(&formData)
	if isValid {
		internal.SendEmail(&formData, response)
	} else {
		http.Error(response, "Form data field is not valid", http.StatusBadRequest)
	}
}
