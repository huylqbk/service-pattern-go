package main

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/huylqbk/service-pattern-go/controllers"
	"github.com/huylqbk/service-pattern-go/infrastructures"
	"github.com/huylqbk/service-pattern-go/repositories"
	"github.com/huylqbk/service-pattern-go/services"
)

func main() {
	http.ListenAndServe(":8080", ChiRouter().InitRouter())
}

type IChiRouter interface {
	InitRouter() *chi.Mux
}

type router struct{}

func (router *router) InitRouter() *chi.Mux {

	playerController := ServiceContainer().InjectPlayerController()

	r := chi.NewRouter()
	r.HandleFunc("/getScore/{player1}/vs/{player2}", playerController.GetPlayerScore)

	return r
}

var (
	m          *router
	routerOnce sync.Once
)

func ChiRouter() IChiRouter {
	if m == nil {
		routerOnce.Do(func() {
			m = &router{}
		})
	}
	return m
}

type IServiceContainer interface {
	InjectPlayerController() controllers.PlayerController
}

type kernel struct{}

func (k *kernel) InjectPlayerController() controllers.PlayerController {

	sqlConn, _ := sql.Open("sqlite3", "/var/tmp/tennis.db")
	sqliteHandler := &infrastructures.SQLiteHandler{}
	sqliteHandler.Conn = sqlConn

	playerRepository := &repositories.PlayerRepository{sqliteHandler}
	playerService := &services.PlayerService{&repositories.PlayerRepositoryWithCircuitBreaker{playerRepository}}
	playerController := controllers.PlayerController{playerService}

	return playerController
}

var (
	k             *kernel
	containerOnce sync.Once
)

func ServiceContainer() IServiceContainer {
	if k == nil {
		containerOnce.Do(func() {
			k = &kernel{}
		})
	}
	return k
}
