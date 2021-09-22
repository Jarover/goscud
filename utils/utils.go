package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Jarover/goscud/models"
)

func FloatToString(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}

func GetDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Fail to get local dir")
		return ""
	}

	return dir
}

// SaveState Фиксация состояния во внешнем файле
// json
func SaveState(path string, id uint) {
	// Get JSON bytes for slice.
	b, _ := json.Marshal(models.Config{RemoteID: id})

	// Write entire JSON file.
	ioutil.WriteFile(path, b, 0644)
}

// ReadState Чтение состояния
func ReadState(path string) (uint, error) {
	file, _ := ioutil.ReadFile(path)
	data := models.Config{}
	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	log.Printf("Read state ID: %d", data.RemoteID)
	return data.RemoteID, nil
}

func GetBaseFile() string {
	filename := os.Args[0] // get command line first parameter
	return strings.Split(filepath.Base(filename), ".")[0]
}
