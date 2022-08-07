package main

var directory = getStorage().Directory

func getDirectory() string {
	return directory
}

func setDirectory(newDirectory string) {
	directory = newDirectory

	storage := getStorage()

	storage.Directory = directory

	currentSyncDirectory = directory

	writeInStorage(storage)
}
