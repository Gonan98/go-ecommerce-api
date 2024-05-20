package product

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gonan98/go-ecommerce-api/helper"
	"github.com/gonan98/go-ecommerce-api/service/auth"
	"github.com/gonan98/go-ecommerce-api/types"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHanlder(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /products", helper.MakeHTTPHandler(h.handleGetProducts))
	router.HandleFunc("GET /products/{id}", helper.MakeHTTPHandler(h.handleGetProductByID))

	// admin route
	router.HandleFunc("POST /products", helper.MakeHTTPHandler(auth.WithJWTAuth(h.handleCreateProduct, h.userStore)))
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) error {
	products, err := h.store.GetAll()
	if err != nil {
		return err
	}

	return helper.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handleGetProductByID(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return helper.BadRequest("productId must be an integer")
	}

	product, err := h.store.GetByID(id)
	if err != nil {
		return helper.EntityNotFound("product", id)
	}

	return helper.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) error {
	var req types.CreateProductRequest
	if err := helper.ParseJSON(r, &req); err != nil {
		return helper.InvalidJSON()
	}

	if err := helper.Validate.Struct(req); err != nil {
		errors := err.(validator.ValidationErrors)
		return helper.InvalidRequestData(errors)
	}

	err := h.store.Create(types.Product{
		Name:        req.Name,
		Description: req.Description,
		Brand:       req.Brand,
		UnitPrice:   req.UnitPrice,
		Stock:       req.Stock,
	})

	if err != nil {
		return err
	}

	return helper.WriteJSON(w, http.StatusCreated, map[string]string{
		"msg": "Product created successfully",
	})
}
