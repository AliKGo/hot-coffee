package dal

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"frappuccino/models"
)

type MenuRepoImpl struct {
	inventoryFilePath string
}

func MenuFilePath(baseDir string) InventoryRepoImpl {
	return InventoryRepoImpl{inventoryFilePath: filepath.Join(baseDir, "/inventory.json")}
}

func (repo MenuRepoImpl) ReadMenuOfDal() (map[string]models.MenuItem, string, int) {
	file, err := os.Open(repo.inventoryFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]models.MenuItem{}, "Repository: File Not Found", http.StatusNotFound
		}
		log.Printf("File opening error: %v", err)
		return map[string]models.MenuItem{}, "Repository: Server Error", http.StatusInternalServerError
	}
	defer file.Close()

	var items map[string]models.MenuItem
	if err := json.NewDecoder(file).Decode(&items); err != nil {
		log.Printf("JSON decoding error: %v", err)
		return map[string]models.MenuItem{}, "Repository: Invalid JSON Format", http.StatusBadRequest
	}

	return items, "Success", http.StatusOK
}

func (repo MenuRepoImpl) AddMenuOfDal(item models.MenuItem) (string, int) {
	items, msg, code := repo.ReadMenuOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	// Проверяем, есть ли уже элемент с таким ID
	if _, exists := items[item.ID]; exists {
		return "Repository: Menu item with this ID already exists", http.StatusBadRequest
	}

	// Добавляем новый элемент в мапу
	items[item.ID] = item

	// Записываем обновленный список в файл
	err := repo.write(items)
	if err != "Success" {
		return "Repository: " + err, http.StatusInternalServerError
	}

	return "Success", http.StatusCreated
}

func (repo MenuRepoImpl) UpdateMenuOfDal(item models.MenuItem) (string, int) {
	items, msg, code := repo.ReadMenuOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	for i, itemMenu := range items {
		if itemMenu.ID == item.ID {
			items[i] = item
		}
	}

	err := repo.write(items)
	if err != "Success" {
		return err, http.StatusInternalServerError
	}
	return "Success", http.StatusOK
}

func (repo MenuRepoImpl) DeleteMenuOfDal(id string) (string, int) {
	items, msg, code := repo.ReadMenuOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	// Проверяем, есть ли элемент с таким ID
	if _, exists := items[id]; !exists {
		return "Repository: Menu item not found", http.StatusNotFound
	}

	// Удаляем элемент из мапы
	delete(items, id)

	// Записываем обновленные данные
	err := repo.write(items)
	if err != "Success" {
		return "Repository: " + err, http.StatusInternalServerError
	}

	return "Success", http.StatusNoContent
}

func (repo MenuRepoImpl) write(items map[string]models.MenuItem) string {
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
