package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) printMetrics(w http.ResponseWriter, req *http.Request) {
	hits := c.fileserverHits.Load()

	template := fmt.Sprintf("<html>\n\t<body>\n\t\t<h1>Welcome, Chirpy Admin</h1>\n\t\t<p>Chirpy has been visited %d times!</p>\n\t</body>\n</html>", hits)
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(template))
}

func (c *apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request) {
	c.fileserverHits.Store(0)

	w.WriteHeader(200)
}

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	fileHandler := http.StripPrefix("/app", fileServer)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileHandler))

	mux.HandleFunc("GET /api/healthz", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Add("Content-Type", "text/plain; charset=utf-8")

		writer.WriteHeader(200)
		writer.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /admin/metrics", apiCfg.printMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetrics)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Printf("Serving file from %s on port %s\n", ".", "8080")
	log.Fatal(server.ListenAndServe())
}
