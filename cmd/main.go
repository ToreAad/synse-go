package main

import (
	"fmt"
	"net/http"

	"github.com/toreaad/aad.land/apps/blog/internals/synse"
)

func main() {

	http.HandleFunc("/synse/api/seismic", synse.SeismicHandler)
	http.HandleFunc("/synse/api/psftemporal", synse.PsfTemporalHandler)
	http.HandleFunc("/synse/api/psfspatial", synse.PsfSpatialHandler)

	fmt.Println("Listening and serving on http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
