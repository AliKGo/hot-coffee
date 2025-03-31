package models

type Type string

const (
	TypeBeverage   Type = "beverage"
	TypeDessert    Type = "dessert"
	TypeMainCourse Type = "main_course"
	TypeSnack      Type = "snack"
	TypeBreakfast  Type = "breakfast"
)

type MenuItem struct {
	ID          string               `json:"product_id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Price       float64              `json:"price"`
	Ingredients []MenuItemIngredient `json:"ingredients"`
	Type        Type                 `json:"type"`
}

type MenuItemIngredient struct {
	IngredientID string  `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
}
