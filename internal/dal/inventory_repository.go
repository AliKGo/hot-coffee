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

func (repo InventoryRepoImpl) ReadInventoryOfDal() (map[string]models.InventoryItem, string, int) {
	file, err := os.Open(repo.inventoryFilePath)
	if err != nil {
		return nil, "Repository: Server Error", http.StatusInternalServerError
	}
	defer file.Close()

	var items []models.InventoryItem
	if err := json.NewDecoder(file).Decode(&items); err != nil {
		return nil, "Repository: Server Error", http.StatusInternalServerError
	}

	// Преобразуем срез в map
	inventoryMap := make(map[string]models.InventoryItem)
	for _, item := range items {
		inventoryMap[item.IngredientID] = item // Используем ID как ключ (замени, если надо другое)
	}

	return inventoryMap, "Success", http.StatusOK
}

func (repo InventoryRepoImpl) AddInventoryOfDal(item models.InventoryItem) (string, int) {
	items, msg, code := repo.ReadInventoryOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	// Добавляем новый элемент в мапу
	items[item.IngredientID] = item

	// Записываем обновлённую мапу в файл
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

	if _, exists := items[itemUpdate.IngredientID]; !exists {
		return "Repository: Not Found", http.StatusNotFound
	}

	items[itemUpdate.IngredientID] = itemUpdate

	if err := repo.write(items); err != "Success" {
		return "Repository: " + err, http.StatusInternalServerError
	}

	return "Success", http.StatusOK
}

func (repo InventoryRepoImpl) DeleteInventoryOfDal(id string) (string, int) {
	items, msg, code := repo.ReadInventoryOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	if _, exists := items[id]; !exists {
		return "Repository: Not Found", http.StatusNotFound
	}

	delete(items, id)

	// Записываем обновлённую мапу в файл
	err := repo.write(items)
	if err != "Success" {
		return "Repository: " + err, http.StatusInternalServerError
	}

	return "Success", http.StatusOK
}

func (repo InventoryRepoImpl) write(items map[string]models.InventoryItem) string {
	file, err := os.Create(repo.inventoryFilePath)
	if err != nil {
		return "failed to create inventory file: " + err.Error()
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(items); err != nil {
		return "failed to encode inventory data: " + err.Error()
	}

	return "Success"
}
