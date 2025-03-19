package dal

import (
	"encoding/json"
	"errors"
	"frappuccino/tools"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"frappuccino/models"
)

type MenuRepoImpl struct {
	inventoryFilePath string
}

func MenuFilePath() MenuRepoImpl {
	return MenuRepoImpl{inventoryFilePath: filepath.Join(*tools.Dir, "/json/menu_items.json")}
}

func (repo MenuRepoImpl) ReadMenuOfDal() (map[string]models.MenuItem, string, int) {
	file, err := os.Open(repo.inventoryFilePath)
	if err != nil {
		return nil, "failed to open menu file: " + err.Error(), http.StatusInternalServerError
	}
	defer file.Close()

	var menuMap map[string]models.MenuItem
	if err := json.NewDecoder(file).Decode(&menuMap); err != nil {
		if errors.Is(err, io.EOF) {
			return make(map[string]models.MenuItem), "Success (empty file)", http.StatusOK
		}
		return nil, "failed to decode menu file: " + err.Error(), http.StatusInternalServerError
	}

	return menuMap, "Success", http.StatusOK
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
