package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/nljackson2020/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
	polkaApiKey    string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	fmt.Printf("JWT_SECRET %v\n", jwtSecret)

	polkaApiKey := os.Getenv("POLKA_API_KEY")
	if polkaApiKey == "" {
		log.Fatal("POLKA_API_KEY environment variable is not set")
	}
	fmt.Printf("POLKA_API_KEY %v\n", polkaApiKey)

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
		polkaApiKey:    polkaApiKey,
	}

	router := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handlerReset)

	apiRouter.Post("/refresh", apiCfg.handlerRefresh)
	apiRouter.Post("/revoke", apiCfg.handlerRevoke)

	apiRouter.Post("/login", apiCfg.handlerLogin)

	apiRouter.Post("/users", apiCfg.handlerUsersCreate)
	apiRouter.Put("/users", apiCfg.handlerUsersUpdate)

	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsRetrieve)
	apiRouter.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	apiRouter.Delete("/chirps/{chirpID}", apiCfg.handlerChirpsDelete)

	apiRouter.Post("/polka/webhooks", apiCfg.handlerPolka)
	router.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
