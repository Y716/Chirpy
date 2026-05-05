package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Y716/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database       database.Queries
	environment    string
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

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	env := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error getting database: %v\n", err)
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		database:       *dbQueries,
		environment:    env,
	}

	fileHandler := http.StripPrefix("/app", fileServer)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileHandler))

	mux.HandleFunc("GET /api/healthz", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Add("Content-Type", "text/plain; charset=utf-8")

		writer.WriteHeader(200)
		writer.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /admin/metrics", apiCfg.printMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerDeleteAllUsers)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Printf("Serving file from %s on port %s\n", ".", "8080")
	log.Fatal(server.ListenAndServe())
}
