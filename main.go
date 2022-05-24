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

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	tt "github.com/kkdai/twitter"
)

var ConsumerKey string
var ConsumerSecret string
var twitterClient *tt.ServerClient

func init() {
	ConsumerKey = os.Getenv("ConsumerKey")
	ConsumerSecret = os.Getenv("ConsumerSecret")
}

const (
	//This URL need note as follow:
	// 1. Could not be localhost, change your hosts to a specific domain name
	// 2. This setting must be identical with your app setting on twitter Dev
	CallbackURL string = "http://YOURDOMAIN.com/maketoken"
)

func main() {

	if ConsumerKey == "" && ConsumerSecret == "" {
		fmt.Println("Please setup ConsumerKey and ConsumerSecret.")
		return
	}

	var port *int = flag.Int(
		"port",
		8888,
		"Port to listen on.")

	flag.Parse()

	fmt.Println("[app] Init server key=", ConsumerKey, " secret=", ConsumerSecret)
	twitterClient = tt.NewServerClient(ConsumerKey, ConsumerSecret)
	http.HandleFunc("/maketoken", GetTwitterToken)
	http.HandleFunc("/request", RedirectUserToTwitter)
	http.HandleFunc("/follow", GetFollower)
	http.HandleFunc("/followids", GetFollowerIDs)
	http.HandleFunc("/time", GetTimeLine)
	http.HandleFunc("/user", GetUserDetail)
	http.HandleFunc("/", MainProcess)

	u := fmt.Sprintf(":%d", *port)
	fmt.Printf("Listening on '%s'\n", u)
	http.ListenAndServe(u, nil)
}

func MainProcess(w http.ResponseWriter, r *http.Request) {

	if !twitterClient.HasAuth() {
		fmt.Fprintf(w, "<BODY><CENTER><A HREF='/request'><IMG SRC='https://g.twimg.com/dev/sites/default/files/images_documentation/sign-in-with-twitter-gray.png'></A></CENTER></BODY>")
		return
	} else {
		//Logon, redirect to display time line
		timelineURL := fmt.Sprintf("http://%s/time", r.Host)
		http.Redirect(w, r, timelineURL, http.StatusTemporaryRedirect)
	}
}

func RedirectUserToTwitter(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter redirect to twitter")
	fmt.Println("Token URL=", CallbackURL)
	requestUrl, _ := twitterClient.GetAuthURL(CallbackURL)

	http.Redirect(w, r, requestUrl, http.StatusTemporaryRedirect)
	fmt.Println("Leave redirtect")
}

func GetTimeLine(w http.ResponseWriter, r *http.Request) {
	timeline, bits, _ := twitterClient.QueryTimeLine(1)
	fmt.Println("TimeLine=", timeline)
	fmt.Fprintf(w, "The item is: "+string(bits))

}
func GetTwitterToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter Get twitter token")
	values := r.URL.Query()
	verificationCode := values.Get("oauth_verifier")
	tokenKey := values.Get("oauth_token")

	twitterClient.CompleteAuth(tokenKey, verificationCode)
	timelineURL := fmt.Sprintf("http://%s/time", r.Host)

	http.Redirect(w, r, timelineURL, http.StatusTemporaryRedirect)
}

func GetFollower(w http.ResponseWriter, r *http.Request) {
	followers, bits, _ := twitterClient.QueryFollower(10)
	fmt.Println("Followers=", followers)
	fmt.Fprintf(w, "The item is: "+string(bits))
}

func GetFollowerIDs(w http.ResponseWriter, r *http.Request) {
	followers, bits, _ := twitterClient.QueryFollowerIDs(10)
	fmt.Println("Follower IDs=", followers)
	fmt.Fprintf(w, "The item is: "+string(bits))
}
func GetUserDetail(w http.ResponseWriter, r *http.Request) {
	followers, bits, _ := twitterClient.QueryFollowerById(2244994945)
	fmt.Println("Follower Detail of =", followers)
	fmt.Fprintf(w, "The item is: "+string(bits))
}
