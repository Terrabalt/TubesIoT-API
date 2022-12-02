package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Status struct {
	IsOn bool `json:"isOn"`
	Suhu int  `json:"suhu"`
}

func (s Status) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (l *Status) Bind(r *http.Request) error {
	return nil
}

type ControlSignal struct {
	IsOn bool `json:"isOn"`
}

func (s ControlSignal) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type config struct {
	SQLiteURL string `env:"SQLDB"`
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	var conf config
	if err := envconfig.Process(context.Background(), &conf); err != nil {
		log.Fatalf("Required enviroment keys was not set up: %s", err)
	}

	db, err := StartDB(conf.SQLiteURL)
	if err != nil {
		log.Fatalf("Error starting database: %s", err)
	}
	if err := db.InitDB(context.Background()); err != nil {
		log.Fatalf("Error starting database: %s", err)
	}
	r := chi.NewRouter()

	r.Use(middleware.CleanPath)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		a := Status{
			true,
			30,
		}

		render.Status(r, http.StatusOK)
		w.Header().Set("content-type", "application/json")
		render.Render(w, r, a)
	})
	r.Post("/TurnOnLamp", func(w http.ResponseWriter, r *http.Request) {
		a := ControlSignal{
			true,
		}

		render.Status(r, http.StatusOK)
		w.Header().Set("content-type", "application/json")
		render.Render(w, r, a)
	})
	r.Post("/TurnOffLamp", func(w http.ResponseWriter, r *http.Request) {
		a := ControlSignal{
			false,
		}

		render.Status(r, http.StatusOK)
		w.Header().Set("content-type", "application/json")
		render.Render(w, r, a)
	})
	r.Post("/post0data", func(w http.ResponseWriter, r *http.Request) {
		data := &Status{}
		if err := render.Bind(r, data); err != nil {
			render.Status(r, http.StatusBadRequest)
			w.Header().Set("content-type", "application/json")
			return
		}

		render.Status(r, http.StatusNoContent)
	})
}
