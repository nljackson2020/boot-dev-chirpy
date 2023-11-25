package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type apiconfig struct {
	fileserverHits int
	jwtSecret      string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	apiCfg := apiconfig{
		fileserverHits: 0,
		jwtSecret:      jwtSecret,
	}

	r := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.reset)
	apiRouter.Post("/chirps", handlerChirpsPostData)
	apiRouter.Get("/chirps/{id}", handlerChirpsGetData)
	apiRouter.Post("/users", handlerCreateUser)
	apiRouter.Post("/login", handlerLogin)
	apiRouter.Put("/users")
	r.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)
	r.Mount("/admin", adminRouter)

	corsR := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsR,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
