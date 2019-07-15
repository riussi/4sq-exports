// Copyright © 2019 Juha Ristolainen <juha.ristolainen@iki.fi>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
)

var authoriseCmd = &cobra.Command{
	Use:   "authorise",
	Short: "authorise the app",
	Run: func(cmd *cobra.Command, args []string) {
		startAuthoriseFlow()
	},
}

func init() {
	RootCmd.AddCommand(authoriseCmd)
}

const (
	callbackURI  = "http://localhost:12345/4sq"
	clientID     = "ZYV105YQOQLKGR0MSSFBQK0EHZF5WJ5MLM4ETPBXCQZ1A5F4"
	clientSecret = "43Z0AOAXU1KWEZ4MOHHYSVBM0K0P5A3PPVNPZB4K1GUO0M2Y"
)

func startAuthoriseFlow() {
	//4. curl https://api.foursquare.com/v2/users/self/checkins?oauth_token=ACCESS_TOKEN&v=YYYYMMDD

	// 1. Setup a channel to listen to the authorisation code from localhost callback
	authCodeCallbackChannel := make(chan string)
	// 2. Start HTTP server to get the callback
	go listenToAuthCodeCallbackOnLocalhost(authCodeCallbackChannel)

	// 3. Open browser for the user to login, grant permissions and redirect to our localhost
	var getAuthCodeURI = fmt.Sprintf("https://foursquare.com/oauth2/authenticate?client_id=%s&response_type=code&redirect_uri=%s", clientID, url.QueryEscape(callbackURI))
	openInBrowser(getAuthCodeURI)

	// 4. Wait for the callback to be called with the auth code
	authCode := <-authCodeCallbackChannel
	fmt.Println("Your Foursquare authentication detais. Please make a note of them. You will need them for other commands.")
	fmt.Printf("- Authorisation code: %v\n", authCode)

	// 5. Get access token
	var accessTokenURI = fmt.Sprintf("https://foursquare.com/oauth2/access_token?client_id=%s&client_secret=%s&grant_type=authorization_code&redirect_uri=%s&code=%s", clientID, clientSecret, url.QueryEscape(callbackURI), authCode)

	accessToken, err := getAccessToken(accessTokenURI)
	if err != nil {
		fmt.Printf("Error getting access token: %s\n", err)
	} else {
		fmt.Printf("- Access token: %v\n", accessToken)
	}
}

func listenToAuthCodeCallbackOnLocalhost(code chan string) {
	callbackHandler := getLocalhostHandlerFunc(code)
	http.HandleFunc("/4sq", callbackHandler)
	log.Fatal(http.ListenAndServe(":12345", nil))
}

func getLocalhostHandlerFunc(code chan string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		keys, ok := req.URL.Query()["code"]
		authCode := "ERROR"
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'code' is missing")
		} else {
			authCode = keys[0]
		}

		fmt.Fprintf(w, "Received authorisation code %v\n", authCode)
		code <- authCode
	}
}

func openInBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func getAccessToken(uri string) (string, error) {
	resp, err := http.Get(uri)
	defer resp.Body.Close()

	if err != nil {
		return "error", err
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return "error", err
	}
	response, err3 := unmarshalResponse(body)
	if err3 != nil {
		return "error", err
	}

	return response.AccessToken, nil
}

func unmarshalResponse(data []byte) (Response, error) {
	var r Response
	err := json.Unmarshal(data, &r)
	return r, err
}

// Response is the struct returned for access code
type Response struct {
	AccessToken string `json:"access_token"`
}
