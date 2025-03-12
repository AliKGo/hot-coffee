package tools

import "os"

func CheckDir(path string) bool { // Если существует дирректория то true
	info, err := os.Stat(path)
	if nil != err {
		return false
	}
	return info.IsDir()
}

func CheckJsonFils() {
	if _, err := os.Stat(*Dir + "/inventory.json"); err != nil {
		os.Create(*Dir + "/inventory.json")
	}
	if _, err := os.Stat(*Dir + "/orders.json"); err != nil {
		os.Create(*Dir + "/orders.json")
	}
	if _, err := os.Stat(*Dir + "/menu_items.json"); err != nil {
		os.Create(*Dir + "/menu_items.json")
	}
}

func CheckINT() {
}
