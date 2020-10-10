package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pacoorozco/networkdevices/internal/models"
	"github.com/pacoorozco/networkdevices/internal/storer"
)

func (a *App) getDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fqdn := vars["fqdn"]

	device, err := a.Storer.GetDevice(fqdn)
	if err != nil {
		switch err {
		case storer.ErrDeviceNotFound:
			respondWithError(w, http.StatusNotFound, "Device not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, device.Presenter())
}

func (a *App) getDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := a.Storer.GetAllDevices()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, devices)
}

func (a *App) createDevice(w http.ResponseWriter, r *http.Request) {
	device := models.Device{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&device); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := device.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.Storer.SetDevice(device); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, device.Presenter())
}

func (a *App) updateDevice(w http.ResponseWriter, r *http.Request) {
	device := models.Device{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&device); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := device.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.Storer.SetDevice(device); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, device.Presenter())
}

func (a *App) deleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fqdn := vars["fqdn"]

	if err := a.Storer.DeleteDevice(fqdn); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
