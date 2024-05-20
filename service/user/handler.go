package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gonan98/go-ecommerce-api/config"
	"github.com/gonan98/go-ecommerce-api/helper"
	"github.com/gonan98/go-ecommerce-api/service/auth"
	"github.com/gonan98/go-ecommerce-api/types"
)

type Handler struct {
	store types.UserStore
}

func NewHanlder(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /login", helper.MakeHTTPHandler(h.HandleLogin))
	router.HandleFunc("POST /register", helper.MakeHTTPHandler(h.HandleRegister))
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	var req types.LoginUserRequest
	if err := helper.ParseJSON(r, &req); err != nil {
		return helper.InvalidJSON()
	}

	if err := helper.Validate.Struct(req); err != nil {
		errors := err.(validator.ValidationErrors)
		return helper.InvalidRequestData(errors)
	}

	u, err := h.store.GetByEmail(req.Email)
	if err != nil {
		return helper.InvalidEmailOrPassword()
	}

	if !auth.ComparePasswords(u.Password, []byte(req.Password)) {
		return helper.InvalidEmailOrPassword()
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.GenerateJWT(secret, u.ID)
	if err != nil {
		return err
	}

	return helper.WriteJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) error {

	var req types.RegisterUserRequest

	if err := helper.ParseJSON(r, &req); err != nil {
		return helper.InvalidJSON()
	}

	if err := helper.Validate.Struct(req); err != nil {
		errors := err.(validator.ValidationErrors)
		return helper.InvalidRequestData(errors)
	}

	_, err := h.store.GetByEmail(req.Email)
	if err == nil {
		return helper.BadRequest(fmt.Sprintf("User with email %s already exists", req.Email))
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	err = h.store.Create(types.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		return err
	}

	return helper.WriteJSON(w, http.StatusCreated, map[string]string{"msg": "User created"})
}
