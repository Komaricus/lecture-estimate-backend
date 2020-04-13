package main

import (
	"awesomeProject3/api"
	"awesomeProject3/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"os"
	"time"
)

const timeOut = 60 * time.Second

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	if err := models.ConnectToDB(); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	//r.Use(middleware.RequestID)
	//r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	//r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(timeOut))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		MaxAge:         300,
	})
	r.Use(c.Handler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/groups", api.GetGroups)
		r.Get("/lessons", api.GetLessons)
		r.Get("/students", api.GetStudents)
		r.Get("/answers", api.GetAnswers)
		r.Post("/answers", api.PostAnswers)
		r.Post("/answer", api.PostAnswer)
	})

	log.Println("Server is running:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
