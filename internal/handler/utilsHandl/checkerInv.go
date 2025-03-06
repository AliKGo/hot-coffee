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
