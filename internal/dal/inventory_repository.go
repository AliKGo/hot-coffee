package dal

import (
	"encoding/json"
	"fmt"
	"frappuccino/models"
	"net/http"
	"os"
	"path/filepath"
)

type InventoryRepoImpl struct {
	inventoryFilePath string
}

func InventoryFilePath(baseDir string) InventoryRepoImpl {
	return InventoryRepoImpl{inventoryFilePath: filepath.Join(baseDir, "/inventory.json")}
}

func (repo InventoryRepoImpl) ReadInventoryOfDal() ([]models.InventoryItem, string, int) {
	file, err := os.Open(repo.inventoryFilePath)

	if err != nil {
		return []models.InventoryItem{}, "Server Error", http.StatusInternalServerError
	}

	defer file.Close()
	var items []models.InventoryItem // Simply decode and read
	if err := json.NewDecoder(file).Decode(&items); err != nil {
		return []models.InventoryItem{}, "Server Error", http.StatusInternalServerError
	}

	return items, "Success", http.StatusOK
}

func (repo InventoryRepoImpl) AddInventoryOfDal(item models.InventoryItem) (string, int) {
	items, msg, code := repo.ReadInventoryOfDal()
	if code != http.StatusOK {
		return msg, code
	}
	items = append(items, item)

	err := repo.write(items)

	if err != nil {
		return err.Error(), http.StatusInternalServerError
	}
	return "Success", http.StatusCreated
}

func (repo InventoryRepoImpl) UpdateInventoryOfDal(itemUpdate models.InventoryItem) (string, int) {
	items, msg, code := repo.ReadInventoryOfDal()
	if code != http.StatusOK {
		return msg, code
	}
	for i, item := range items {
		if item.IngredientID == itemUpdate.IngredientID {
			items[i].IngredientID = itemUpdate.IngredientID
			if err := repo.write(items); err != nil {
				return err.Error(), http.StatusInternalServerError
			}
			return "Success", http.StatusOK
		}
	}
	return "Not Found", http.StatusNotFound
}

func (repo InventoryRepoImpl) DeleteInventoryOfDal(id string) (string, int) {
	items, msg, code := repo.ReadInventoryOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	for i, item := range items {
		if item.IngredientID == id {
			items = append(items[:i], items[i+1:]...)

			err := repo.write(items)

			if err != nil {
				return err.Error(), http.StatusInternalServerError
			}
			return "Success", http.StatusOK
		}
	}

	return "Not Found", http.StatusNotFound
}

func (repo InventoryRepoImpl) write(items []models.InventoryItem) error {
	file, err := os.Create(repo.inventoryFilePath)
	if err != nil {
		return fmt.Errorf("failed to create inventory file: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(items); err != nil {
		return fmt.Errorf("failed to encode inventory data: %w", err)
	}

	return nil
}
