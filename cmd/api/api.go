package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gonan98/go-ecommerce-api/service/cart"
	"github.com/gonan98/go-ecommerce-api/service/order"
	"github.com/gonan98/go-ecommerce-api/service/product"
	"github.com/gonan98/go-ecommerce-api/service/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	subrouter := http.NewServeMux()
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", subrouter))

	userStore := user.NewStore(s.db)
	userHandler := user.NewHanlder(userStore)
	userHandler.RegisterRoutes(subrouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHanlder(productStore, userStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(s.db)
	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Printf("Server running on port %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}
