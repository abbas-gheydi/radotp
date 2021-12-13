package monitoring

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Listen = "0.0.0.0:2112"

func Start() {
	make_inMemory_users_storage()

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", gathermetrics)

	server := http.Server{
		Addr:    Listen,
		Handler: mux,
	}

	log.Println("metrics Listen on:", Listen)
	log.Println(server.ListenAndServe())
}

func gathermetrics(w http.ResponseWriter, r *http.Request) {

	make_metrics()

	promhttp.Handler().ServeHTTP(w, r)
}
