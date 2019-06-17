package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func createRefreshToken(config *oauth2.Config, tokenFile string) error {

	tok, err := tokenFromFile(tokenFile)
	var map1 map[string]string
	if err != nil {
		tok, map1 = getTokenFromWeb(config)
		return saveToken(tokenFile, tok, map1)
	}

	fmt.Printf("Token file \"%s\" already exists.\n", tokenFile)

	return nil
}

// Request a token from the web, then returns the retrieved token, "client_id" and "client_secret"
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, map[string]string) {
	fmt.Printf("config.ClientID %s\n", config.ClientID)
	fmt.Printf("config.ClientSecret %s\n", config.ClientSecret)

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok, map[string]string{"client_id": config.ClientID, "client_secret": config.ClientSecret}
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

// Saves a refresh token to a file path.
func saveToken(path string, token *oauth2.Token, map1 map[string]string) error {

	type authorizedUserToken struct {
		Type         string `json:"type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RefreshToken string `json:"refresh_token"`
	}

	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.New("Unable to cache oauth token")
	}
	defer f.Close()

	au := authorizedUserToken{
		Type:         "authorized_user",
		ClientID:     map1["client_id"],
		ClientSecret: map1["client_secret"],
		RefreshToken: token.RefreshToken}

	return json.NewEncoder(f).Encode(au) // authorized_user

	//
	// return json.NewEncoder(f).Encode(token) // access_token
}
