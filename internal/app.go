package internal

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     DatabaseAccessLayer
}

func (a *App) Initialize() {
	a.DB = NewInMemory()
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/devices", a.getDevices).Methods("GET")
	a.Router.HandleFunc("/devices", a.createDevice).Methods("POST")
	a.Router.HandleFunc("/devices", a.updateDevice).Methods("PUT")
	a.Router.HandleFunc("/devices/{fqdn}", a.getDevice).Methods("GET")
	a.Router.HandleFunc("/devices/{fqdn}", a.deleteDevice).Methods("DELETE")
}

func (a *App) getDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fqdn := vars["fqdn"]

	d := Device{FQDN: fqdn}
	if err := d.getDevice(a.DB); err != nil {
		switch err {
		case ErrNotFound:
			respondWithError(w, http.StatusNotFound, "Device not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, d)
}

func (a *App) getDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := getDevices(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, devices)
}

func (a *App) createDevice(w http.ResponseWriter, r *http.Request) {
	var d Device
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&d); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := d.createDevice(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, d)
}

func (a *App) updateDevice(w http.ResponseWriter, r *http.Request) {
	var d Device
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&d); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := d.updateDevice(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, d)
}

func (a *App) deleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fqdn := vars["fqdn"]

	d := Device{FQDN: fqdn}
	if err := d.deleteDevice(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
