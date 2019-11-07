package main

import (
	"encoding/json"
	"reflect"
	"crypto/rand"
	"encoding/base64"
	
	"fmt"
	"io/ioutil"
    "os"
	"net/http"
	"github.com/joho/godotenv"
    "log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

   "strings"
  
)

func init() {
	fmt.Println("Initialising..")
	
}
var err = godotenv.Load()
var CID = os.Getenv("CLIENT_ID")
var SECRET = os.Getenv("CLIENT_SECRET")


  var conf = &oauth2.Config{
	
	ClientID: CID,
	ClientSecret: SECRET,
	RedirectURL: "http://127.0.0.1:8000/callback",
	Scopes: []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func main() {
	fmt.Println("inside main confid ",conf.ClientID)
	

	fmt.Println("inside main- going to start server")
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	log.Fatal(http.ListenAndServe(":8000", nil))
}


func handleMain(w http.ResponseWriter, r *http.Request) {
	// fs := http.FileServer(http.Dir("static"))
  	// http.Handle("/", fs)
	var htmlstring string
	htmlstring = `<html>
	<body>
		<a href="/login"><button>Login with Google!</button></a>
	</body>
	</html>`
	fmt.Fprintf(w, htmlstring)
	}

var state = randToken()
func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside handleGoogleLogin")
	cookie := http.Cookie{Name: "oauthstate", Value: state}//,Expires: expiration}
    http.SetCookie(w, &cookie)
	url := getLoginURL(state)
	fmt.Println("Google Url is", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)  // redirecting to google login service
	}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func getLoginURL(state string) string {
    return conf.AuthCodeURL(state)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))

	var result map[string]interface{}
	var errJson error
	errJson = json.Unmarshal(content, &result)   // to parse content byte I have used map
		if errJson != nil {
			fmt.Println(err)
		} else {
			
			fmt.Println("email in response is ",result["email"])
			fmt.Println("type of result email ",reflect.TypeOf(result["email"]))
			}
			
	var emailString string
	emailString = result["email"].(string)   // type assertion i.(type)

	fmt.Println("type of emailString ",reflect.TypeOf(emailString), emailString)

	if !(strings.Contains(emailString,"@indiamart.com")){
		fmt.Println("Not an indiamart id")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
	
	
		return
	}

	if (err != nil){
	
		fmt.Println(err.Error())
		w.Write([]byte("Error in getting response from google"))
		http.Redirect(w, r, "http://127.0.0.1:8000/", http.StatusTemporaryRedirect)
		
		return
	}
	
	

	// fmt.Fprintf(w, "Welcome IM User")
	http.Redirect(w,r, "https://github.com/95aakash/IMLoginAutomation",http.StatusTemporaryRedirect)  //redirecting an im user
	  

}

func getUserInfo(oauthState string, code string) ([]byte, error) {
	if oauthState != state {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body) //give []byte
	// fmt.Println("Inside getuserinfo response body type ",reflect.TypeOf(response.Body))
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil  
}

