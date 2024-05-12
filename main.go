package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	registry := ServiceRegistry{}

	registry.Init()

	dockerClient, err := NewDockerClient()
	if err != nil {
		panic(err)
	}

	registrar := Registrar{SRegistry: &registry, DockerClient: dockerClient}

	if err := registrar.Init(); err != nil {
		panic(err)
	}

	go registrar.Observe()

	app := Application{SRegistry: &registry}

	r := chi.NewRouter()

	r.HandleFunc("/reverse-proxy", app.Handle)

	log.Fatalln(http.ListenAndServe(":3000", nil))
}
