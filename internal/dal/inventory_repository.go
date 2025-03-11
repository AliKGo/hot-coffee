package dal

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"frappuccino/models"
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
		return []models.InventoryItem{}, "Repository: Server Error", http.StatusInternalServerError
	}

	defer file.Close()
	var items []models.InventoryItem // Simply decode and read
	if err := json.NewDecoder(file).Decode(&items); err != nil {
		return []models.InventoryItem{}, "Repository: Server Error", http.StatusInternalServerError
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

	if err != "Success" {
		return "Repository: " + err, http.StatusInternalServerError
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
			if err := repo.write(items); err != "Success" {
				return "Repository: " + err, http.StatusInternalServerError
			}
			return "Success", http.StatusOK
		}
	}
	return "Repository: Not Found", http.StatusNotFound
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

			if err != "Success" {
				return "Repository: " + err, http.StatusInternalServerError
			}
			return "Success", http.StatusOK
		}
	}

	return "Repository: Not Found", http.StatusNotFound
}

func (repo InventoryRepoImpl) write(items []models.InventoryItem) string {
	file, err := os.Create(repo.inventoryFilePath)
	if err != nil {
		return "failed to create inventory file:" + err.Error()
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(items); err != nil {
		return "failed to encode inventory data: " + err.Error()
	}

	return "Success"
}
