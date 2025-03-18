package dal

import (
	"encoding/json"
	"frappuccino/models"
	"frappuccino/tools"
	"net/http"
	"os"
	"path/filepath"
)

type OrderRepoImpl struct {
	inventoryFilePath string
}

func NewOrderRepoImpl() *OrderRepoImpl {
	return &OrderRepoImpl{inventoryFilePath: filepath.Join(*tools.Dir, "/orders.json")}
}

func (repo OrderRepoImpl) ReadOrderOfDal() (map[string]models.Order, string, int) {
	file, err := os.Open(repo.inventoryFilePath)
	if err != nil {
		return nil, "Repository: Server Error", http.StatusInternalServerError
	}
	defer file.Close()

	var orders []models.Order
	if err := json.NewDecoder(file).Decode(&orders); err != nil {
		return nil, "Repository: Server Error", http.StatusInternalServerError
	}

	OrderMap := make(map[string]models.Order)
	for _, item := range orders {
		OrderMap[item.ID] = item // Используем ID как ключ (замени, если надо другое)
	}
	return OrderMap, "Success", http.StatusOK
}

func (repo OrderRepoImpl) UpdateOrderOfDal(item models.Order) (string, int) {
	items, msg, code := repo.ReadOrderOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	if _, exists := items[item.ID]; !exists {
		return "Repository: Not Found", http.StatusNotFound
	}

	items[item.ID] = item

	if err := repo.write(items); err != "Success" {
		return "Repository: " + err, http.StatusInternalServerError
	}
	return "Success", http.StatusOK
}

func (repo OrderRepoImpl) write(items map[string]models.Order) string {
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

func (repo OrderRepoImpl) DeleteOrderOfDal(item models.Order) (string, int) {
	items, msg, code := repo.ReadOrderOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	if _, exists := items[item.ID]; !exists {
		return "Repository: Not Found", http.StatusNotFound
	}

	delete(items, item.ID)

	if err := repo.write(items); err != "Success" {
		return "Repository: " + err, http.StatusInternalServerError
	}

	return "Success", http.StatusNoContent
}

func (repo OrderRepoImpl) AddOrderOfDal(item models.Order) (string, int) {
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
