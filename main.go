//
// create-refresh-token
//

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"
)

// https://console.developers.google.com/apis/credentials
// 1. create OAuth 2.0 client IDs
// 2. type: Other
// 3. enter any Name
// 4. click on: "Download JSON"
//
// default json filename created from google console credentials "Download JSON"
// client_secret_000000000000-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.apps.googleusercontent.com.json

// When running the app, use the same google account
// that created the OAuth 2.0 client IDs.

// sample output refresh token json file content
// https://github.com/firebase/firebase-admin-go/blob/master/testdata/refresh_token.json
// {
//     "type": "authorized_user",
//     "client_id": "000000000000-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.apps.googleusercontent.com.json",
//     "client_secret": "aaaaaaaaaaaaaaaaaaaaaaaa",
//     "refresh_token": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
// }

const _credFile = "./client_secret_000000000000-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.apps.googleusercontent.com.json" // input Edit 1.
const _tokenFile = "./token.json"                                                                                 // output

// Firebase
// https://firebase.google.com/docs/admin/setup#using_an_oauth_20_refresh_token
//

func main() {

	b, err := ioutil.ReadFile(_credFile)

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/firebase.database"}

	config, err := google.ConfigFromJSON(b, scopes...)

	if err != nil {
		log.Fatalf("Unable to parse client secret file \"%s\"\n to config\n error = %v", _credFile, err)
	}

	err = createRefreshToken(config, _tokenFile)
	if err != nil {
		log.Printf("%v\n", err)
		return // exit app
	}

	fmt.Println("OK! Token Creation")

	// test newly created token file
	tok, err := tokenFromFile(_tokenFile)

	if err != nil {
		log.Printf("ERROR: tokenFromFile(\"%s\") error = %v\n", _tokenFile, err)
		return // exit app
	}

	client := config.Client(context.Background(), tok)

	fmt.Printf("Test created or existing token file, *http.Client = %v\n", client)

}
