package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/pkg/browser"
)

var config *oauth2.Config
var tokenFile = "token.json"
var server *drive.Service
var client *http.Client
var ctx context.Context

// Retrieve a token, saves the token, then returns the generated client.
func getClient() {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.

	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		createTokenFromWeb()
		return
	}
	client = config.Client(context.Background(), tok)
	runAuth()
}

// Request a token from the web, then returns the retrieved token.
func createTokenFromWeb() {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	browser.OpenURL(authURL)
}

func callbackTokenFromWeb(authCode string) {
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	saveToken(tokenFile, tok)
	client = config.Client(context.Background(), tok)
	runAuth()
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func startAuth() {
	ctx = context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err = google.ConfigFromJSON(b, drive.DriveFileScope, drive.DriveMetadataScope, drive.DriveScope, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	getClient()
}

func runAuth() {
	var err error

	server, err = drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	startPageTokenResponse, err := server.Changes.GetStartPageToken().Do()

	if err != nil {
		log.Fatalf("Unable to retrieve start page token: %v", err)
	}

	startPageToken = startPageTokenResponse.StartPageToken

	startSync()
}
