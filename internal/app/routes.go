package app

// initializeRoutes defines the routes that the app will handle.
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/devices", a.getDevices).Methods("GET")
	a.Router.HandleFunc("/devices", a.createDevice).Methods("POST")
	a.Router.HandleFunc("/devices", a.updateDevice).Methods("PUT")
	a.Router.HandleFunc("/devices/{fqdn}", a.getDevice).Methods("GET")
	a.Router.HandleFunc("/devices/{fqdn}", a.deleteDevice).Methods("DELETE")
}
