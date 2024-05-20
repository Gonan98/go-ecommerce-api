package cart

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gonan98/go-ecommerce-api/helper"
	"github.com/gonan98/go-ecommerce-api/service/auth"
	"github.com/gonan98/go-ecommerce-api/types"
)

type Handler struct {
	orderStore   types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(orderStore types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{
		orderStore:   orderStore,
		productStore: productStore,
		userStore:    userStore,
	}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /cart/checkout", helper.MakeHTTPHandler(auth.WithJWTAuth(h.handleCheckout, h.userStore)))
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) error {
	userID := auth.GetUserIDFromContext(r.Context())
	var req types.CartCheckoutRequest

	if err := helper.ParseJSON(r, &req); err != nil {
		return helper.InvalidJSON()
	}

	if err := helper.Validate.Struct(req); err != nil {
		errors := err.(validator.ValidationErrors)
		return helper.InvalidRequestData(errors)
	}

	productIDs, err := getCartItemsIDs(req.Items)
	if err != nil {
		return err
	}

	products, err := h.productStore.GetByIDs(productIDs)
	if err != nil {
		return err
	}

	orderID, totalPrice, err := h.createOrder(products, req.Items, userID)
	if err != nil {
		return err
	}

	return helper.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"totalPrice": totalPrice,
		"orderId":    orderID,
	})
}

func (h *Handler) createOrder(products []types.Product, cartItems []types.CartItem, userID int) (int, float64, error) {
	productsMap := make(map[int]types.Product)
	for _, p := range products {
		productsMap[p.ID] = p
	}

	if err := checkIfCartIsInStock(cartItems, productsMap); err != nil {
		return 0, 0, err
	}

	totalPrice := calculateTotalPrice(cartItems, productsMap)

	for _, item := range cartItems {
		product := productsMap[item.ProductID]
		product.Stock -= item.Quantity
		h.productStore.Update(product)
	}

	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID: userID,
		Status: "pending",
	})

	if err != nil {
		return 0, 0, err
	}

	for _, item := range cartItems {
		h.orderStore.CreateOrderDetail(types.OrderDetail{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: productsMap[item.ProductID].UnitPrice,
		})
	}

	return orderID, totalPrice, nil
}
