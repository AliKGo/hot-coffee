package tools

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var (
	Dir  = flag.String("dir", "./data", "dirr")
	Port = flag.String("port", "8080", "port")
	Help = flag.Bool("help", false, "help")
)

func ParseFlag() {
	if err := validateDir(); err != nil {
		log.Fatal(err)
	}
	if err := validatePort(); err != nil {
		log.Fatal(err)
	}
}

func HelpFunck() {
	fmt.Println(
		`Coffee Shop Management System

Usage:
hot-coffee [--port <N>] [--dir <S>] 
hot-coffee --help

Options:
--help       Show this screen.
--port N     Port number.
--dir S      Path to the data directory.`)
	flag.PrintDefaults()
}

func validatePort() error {
	port, err := strconv.Atoi(*Port)
	if err != nil {
		return fmt.Errorf("port should be number")
	}

	if port < 1024 || port > 49151 {
		return fmt.Errorf("invalid port, must be between 1024 and 49151")
	}

	return nil
}

func SplittingThePath(str string) []string {
	return strings.Split(str, "/")
}

func validateDir() error {
	pathElements := SplittingThePath(*Dir)
	for i := 0; i < len(pathElements); i++ {
		if pathElements[i] == "tools" || pathElements[i] == "app" || pathElements[i] == "frappuccino" || pathElements[i] == "cmd" || pathElements[i] == "internal" || pathElements[i] == "models" {
			return fmt.Errorf("forbidden dir")
		}
	}
	return nil
}
