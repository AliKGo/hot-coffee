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

func (repo MenuRepoImpl) ReadMenuOfDal() ([]models.MenuItem, string, int) {
	file, err := os.Open(repo.inventoryFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.MenuItem{}, "Repository: File Not Found", http.StatusNotFound
		}
		log.Printf("File opening error: %v", err)
		return []models.MenuItem{}, "Repository: Server Error", http.StatusInternalServerError
	}
	defer file.Close()

	var items []models.MenuItem
	if err := json.NewDecoder(file).Decode(&items); err != nil {
		log.Printf("JSON decoding error: %v", err)
		return []models.MenuItem{}, "Repository: Invalid JSON Format", http.StatusBadRequest
	}

	return items, "Success", http.StatusOK
}

func (repo MenuRepoImpl) AddMenuOfDal(item models.MenuItem) (string, int) {
	items, msg, code := repo.ReadMenuOfDal()
	if code != http.StatusOK {
		return msg, code
	}
	items = append(items, item)

	err := repo.write(items)
	if err != "Success" {
		return err, http.StatusInternalServerError
	}
	return "Success", http.StatusOK
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

func (repo MenuRepoImpl) DeleteMenuOfDal() (string, int) {
	items, msg, code := repo.ReadMenuOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	for i, item := range items {
		if item.ID == items[i].ID {
			items = append(items[:i], items[i+1:]...)
			err := repo.write(items)
			if err != "Success" {
				return err, http.StatusInternalServerError
			}
		}
	}

	return "Success", http.StatusNoContent
}

func (repo MenuRepoImpl) write(items []models.MenuItem) string {
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
