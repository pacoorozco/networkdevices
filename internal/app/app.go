package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pacoorozco/networkdevices/internal/storer"
)

// App represents a running application with all the dependant services.
type App struct {
	Router *mux.Router
	Storer storer.Storer
}

// Initialize set up the application.
func (a *App) Initialize() {
	a.Storer = storer.NewInMemory()
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run starts the application, in this case the HTTP server.
func (a *App) Run(addr string) {
	fmt.Printf("Starting server on %s...", addr)

	log.Fatal(http.ListenAndServe(addr, a.Router))
}
