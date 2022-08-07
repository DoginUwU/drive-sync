package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type DriveFile struct {
	Name   string
	Id     string
	Synced bool
}

type Storage struct {
	Directory string
	Files     []DriveFile
}

func getStorage() Storage {
	checkIfStorageExists()

	file, _ := os.Open("storage.json")
	defer file.Close()

	decoder := json.NewDecoder(file)

	storage := Storage{}

	err := decoder.Decode(&storage)

	if err != nil {
		fmt.Println("error:", err)
	}

	return storage
}

func checkIfStorageExists() {
	_, err := os.Stat("storage.json")

	if os.IsNotExist(err) {
		createStorage()
	}
}

func createStorage() {
	file, _ := os.Create("storage.json")
	defer file.Close()

	storage := Storage{
		Directory: "",
		Files:     []DriveFile{},
	}

	encoder := json.NewEncoder(file)
	encoder.Encode(storage)
}

func writeInStorage(storage Storage) {
	file, _ := json.MarshalIndent(storage, "", " ")

	os.WriteFile("storage.json", file, 0644)
}
