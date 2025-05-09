package utilsHandl

import "frappuccino/models"

func ValidateInventory(i models.InventoryItem) string {
	if i.IngredientID == "" {
		return "Ingredient_id cannot be empty!"
	}
	if i.Name == "" {
		return "Name cannot be empty!"
	}
	if i.Quantity <= 0 {
		return "Quantity must be greater than zero!"
	}
	if i.Unit == "" {
		return "Unit cannot be empty!"
	}
	if i.Unit != "shots" && i.Unit != "ml" && i.Unit != "g" {
		return "Unit cannot be empty!"
	}
	return "OK"
}

func ValidateMenu(i models.MenuItem) string {
	if i.ID == "" {
		return "Product_id cannot be empty"
	}
	if i.Name == "" {
		return "name cannot be empty"
	}
	if i.Description == "" {
		return "Validation error: description cannot be empty"
	}
	if i.Price < 0 {
		return "Validation error: price cannot be negative"
	}
	if len(i.Ingredients) == 0 {
		return "Validation error: menu item must have at least one ingredient"
	}

	for _, ing := range i.Ingredients {
		if ing.IngredientID == "" {
			return "Validation error: ingredient_id cannot be empty"
		}
		if ing.Quantity <= 0 {
			return "Validation error: ingredient quantity must be greater than zero"
		}
	}

	if i.Type != models.TypeBeverage && i.Type != models.TypeBreakfast && i.Type != models.TypeSnack && i.Type != models.TypeDessert && i.Type != models.TypeMainCourse {
		return "Validation error: type must be either beverage or breakfast"
	}

	return "OK"
}

func ValidateOrder(o models.Order) string {
	if o.CustomerName == "" {
		return "Customer name cannot be empty!"
	}
	if len(o.Items) == 0 {
		return "Order must contain at least one item!"
	}
	for _, item := range o.Items {
		if item.ProductID == "" {
			return "Product ID cannot be empty in order items!"
		}
		if item.Quantity <= 0 {
			return "Quantity must be greater than zero in order items!"
		}
	}
	return "OK"
}
