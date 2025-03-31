package dal

import (
	"encoding/json"
	"errors"
	"frappuccino/models"
	"frappuccino/tools"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type OrderRepoImpl struct {
	inventoryFilePath string
}

func NewOrderRepoImpl() *OrderRepoImpl {
	return &OrderRepoImpl{inventoryFilePath: filepath.Join(*tools.Dir, "/json/orders.json")}
}

func (repo *OrderRepoImpl) ReadOrderOfDal() (map[string]models.Order, string, int) {
	file, err := os.Open(repo.inventoryFilePath)
	if err != nil {
		return nil, "Repository: Server error", http.StatusInternalServerError
	}
	defer file.Close()

	var orderMap map[string]models.Order
	if err := json.NewDecoder(file).Decode(&orderMap); err != nil {
		if errors.Is(err, io.EOF) {
			return make(map[string]models.Order), "Success (empty file)", http.StatusOK
		}
		return nil, "Repository: Server error", http.StatusInternalServerError
	}

	return orderMap, "Success", http.StatusOK
}

func (repo *OrderRepoImpl) UpdateOrderOfDal(item models.Order) (string, int) {
	items, msg, code := repo.ReadOrderOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	if _, exists := items[item.ID]; !exists {
		return "Repository: Not found", http.StatusNotFound
	}

	items[item.ID] = item

	if err := repo.write(items); err != "Success" {
		return "Repository: " + err, http.StatusInternalServerError
	}
	return "Success", http.StatusOK
}

func (repo *OrderRepoImpl) DeleteOrderOfDal(item models.Order) (string, int) {
	items, msg, code := repo.ReadOrderOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	if _, exists := items[item.ID]; !exists {
		return "Repository: Not found", http.StatusNotFound
	}

	delete(items, item.ID)

	if err := repo.write(items); err != "Success" {
		return "Repository: " + err, http.StatusInternalServerError
	}

	return "Success", http.StatusNoContent
}

func (repo *OrderRepoImpl) AddOrderOfDal(item models.Order) (string, int) {
	items, msg, code := repo.ReadOrderOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	if _, exists := items[item.ID]; exists {
		return "Repository: There is already an order for this ID", http.StatusBadRequest
	}

	items[item.ID] = item
	if err := repo.write(items); err != "Success" {
		return "Repository: " + err, http.StatusInternalServerError
	}
	return "Success", http.StatusOK
}

func (repo *OrderRepoImpl) write(items map[string]models.Order) string {
	file, err := os.Create(repo.inventoryFilePath)
	if err != nil {
		return "Repository: Failed to create inventory file - " + err.Error()
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(items); err != nil {
		return "Repository: Failed to encode inventory data: " + err.Error()
	}

	return "Success"
}
