package tools

import (
	"fmt"
	"os"
)

func CheckJsonFiles() error {
	var dir = *Dir
	// Создание папки, если её нет
	jsonDir := dir + "/json"
	if err := os.MkdirAll(jsonDir, os.ModePerm); err != nil {
		return fmt.Errorf("ошибка при создании папки: %w", err)
	}

	// Список файлов для создания
	files := []string{"inventory.json", "orders.json", "menu_items.json"}

	// Проверка и создание файлов
	for _, file := range files {
		filePath := jsonDir + "/" + file
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("ошибка при создании файла %s: %w", filePath, err)
			}
			defer f.Close()
		}
	}

	return nil
}
