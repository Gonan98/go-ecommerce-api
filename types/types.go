package types

import "time"

type CRUD[T any] interface {
	Create(t *T) error
	GetAll() ([]T, error)
	GetById() (*T, error)
	Update(t *T, id uint) error
	Delete(id uint) error
}

type UserStore interface {
	//CRUDStore[User]
	Create(user User) error
	GetById(id int) (*User, error)
	GetByEmail(email string) (*User, error)
}

type ProductStore interface {
	GetAll() ([]Product, error)
	GetByID(int) (*Product, error)
	GetByIDs([]int) ([]Product, error)
	Create(Product) error
	Update(Product) error
}

type OrderStore interface {
	CreateOrder(Order) (int, error)
	CreateOrderDetail(OrderDetail) error
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Brand       string  `json:"brand"`
	UnitPrice   float64 `json:"unitPrice"`
	Stock       int     `json:"stock"`
}

type Order struct {
	ID        int       `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UserID    int       `json:"userID"`
}

type OrderDetail struct {
	OrderID   int     `json:"orderID"`
	ProductID int     `json:"productID"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
	Discount  float32 `json:"discount"`
}

type RegisterUserRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=16"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Brand       string  `json:"brand" validate:"required"`
	UnitPrice   float64 `json:"unitPrice" validate:"required,gte=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
}

type CartItem struct {
	ProductID int `json:"productId" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,gt=0"`
}

type CartCheckoutRequest struct {
	Items []CartItem `json:"items" validate:"required"`
}
