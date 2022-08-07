package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/drive/v3"
)

var currentSyncDirectory = getDirectory()
var driveFiles = []DriveFile{}
var startPageToken string

func listenChanges() {
	log.Println("Start listen changes")

	for {
		changes, err := server.Changes.List(startPageToken).Do()

		if err != nil {
			log.Fatalf("Unable to retrieve changes: %v", err)
		}

		for _, i := range changes.Changes {
			if i.File == nil {
				continue
			}

			log.Println("File change: ", i.FileId)
			changeDriveSyncStatus(i.File, false)

			if !i.Removed || !i.File.Trashed {
				saveFile(i.File)
			} else {
				removeFile(i.File)
			}
		}

		startPageToken = changes.NewStartPageToken
	}
}

func startSync() {
	if len(directory) <= 0 {
		return
	}

	if server == nil {
		log.Fatalf("Unable to retrieve Drive client")
	}

	log.Println("Start sync")

	r, err := server.Files.List().Q("trashed=false").Do()

	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	driveFiles = getStorage().Files

	for _, i := range r.Files {
		if checkIfExistsInStorageAndIfSynced(i) {
			continue
		}

		log.Println("Sync file: " + i.Name)
		changeDriveSyncStatus(i, false)

		if !checkIfFileExists(i) || !checkIfFileIsCompleted(i) {
			saveFile(i)
		} else {
			log.Println("Already downloaded: ", i.Name)
			changeDriveSyncStatus(i, true)
		}
	}

	log.Println("Sync finished")
	listenChanges()
}

func checkIfExistsInStorageAndIfSynced(driveFile *drive.File) bool {
	storage := getStorage()
	synced := false

	for _, file := range storage.Files {
		if file.Id == driveFile.Id {
			synced = file.Synced
			changeDriveSyncStatus(driveFile, file.Synced)
		}
	}

	return synced
}

func changeDriveSyncStatus(driveFile *drive.File, status bool) {
	founded := false

	for i, file := range driveFiles {
		if file.Id == driveFile.Id {
			driveFiles[i].Synced = status
			founded = true
		}
	}

	if !founded {
		driveFiles = append(driveFiles, DriveFile{Name: driveFile.Name, Id: driveFile.Id, Synced: false})
	}

	storage := getStorage()
	storage.Files = driveFiles

	writeInStorage(storage)
}

func checkIfFileExists(file *drive.File) bool {
	filePath := getDriveFilePath(file.Id)
	_, err := os.Stat(filePath)

	return !os.IsNotExist(err)
}

func checkIfFileIsCompleted(file *drive.File) bool {
	neededFile, err := server.Files.Get(file.Id).Fields("size").Do()

	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	currentFileSize := getFileSize(getDriveFilePath(file.Id))

	return currentFileSize == neededFile.Size
}

func getFileSize(filePath string) int64 {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Unable to get file info: %v", err)
	}

	return fileInfo.Size()
}

func saveFile(file *drive.File) {
	if file.MimeType == "application/vnd.google-apps.folder" {
		createFolder(file)
		return
	}

	log.Println("Save file: ", file.Name)
	downloadFile(file)
}

func removeFile(file *drive.File) {
	filePath := getDriveFilePath(file.Id)
	err := os.Remove(filePath)

	if err != nil {
		log.Fatalf("Unable to remove file: %v", err)
	}
}

func createFolder(file *drive.File) {
	filePath := getDriveFilePath(file.Id)
	os.Mkdir(filePath, 0755)
	changeDriveSyncStatus(file, true)
}

func downloadFile(file *drive.File) {
	r, err := server.Files.Get(file.Id).Download()

	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	defer r.Body.Close()
	saveFileInDirectory(file, r)
}

func getDriveFilePath(id string) string {
	getFile, err := server.Files.Get(id).Fields("parents", "name", "id").Do()

	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	parentFolders := getFile.Parents
	var folders = ""

	for _, i := range parentFolders {
		folder, err := server.Files.Get(i).Fields("name").Do()

		if err != nil {
			log.Fatalf("Unable to retrieve files: %v", err)
		}

		if folder.Name == "My Drive" {
			continue
		}

		folders = folders + folder.Name + "/"
	}

	return currentSyncDirectory + "/" + folders + getFile.Name
}

func saveFileInDirectory(file *drive.File, resp *http.Response) {
	f, err := os.Create(getDriveFilePath(file.Id))
	if err != nil {
		log.Fatalf("Unable to create file: %v", err)
	}
	defer f.Close()

	io.Copy(f, resp.Body)
	changeDriveSyncStatus(file, true)
}
