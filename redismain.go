package main
import (
	"encoding/json"
	"reflect"
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"fmt"
	"io/ioutil"
    "os"
	"net/http"
	"github.com/joho/godotenv"
    "log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/nu7hatch/gouuid"
	"github.com/gomodule/redigo/redis"
   "strings"
//    "time"
  
)
var err = godotenv.Load()
var CID = os.Getenv("CLIENT_ID")
var SECRET = os.Getenv("CLIENT_SECRET")
var id, _ = uuid.NewV4()
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
var cache redis.Conn
func initCache() {
	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}
	// Assign the connection to the package level `cache` variable
	cache = conn
}


func main() {
	initCache()
	fmt.Println("inside main confid ",conf.ClientID)

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	http.HandleFunc("/dummy", serveFiles)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/service", receiveAjax)
	log.Fatal(http.ListenAndServe(":8000", nil))

}
func receiveAjax(w http.ResponseWriter, r *http.Request) {
// to check cookie session
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sessionToken := c.Value

	// get user from our cache, where we set the session token
	response, err := cache.Do("GET", sessionToken)
	if err != nil {
		// error fetching from cache, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error fetching from cache"))
		return
	}
	if response == nil {
		// If the session token is not present in cache, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	
// to receive ajax data
var ajax_post_data string
	if r.Method == "POST" {
		ajax_post_data = r.FormValue("ajax_post_data")
		fmt.Println("Receive ajax post data string ", ajax_post_data)
		fmt.Println("type of ajax post data ",reflect.TypeOf(ajax_post_data))
	
	}
// unmarshalling links received from ajax
var link map[string]interface{}  // map for recieving link
	var errJson error
	errJson = json.Unmarshal([]byte(ajax_post_data), &link)   // unmarshalling json
		if errJson != nil {
			fmt.Println(errJson)
		} else {
			
			fmt.Println("link recieved is ",link["linktosend"])
			
			}
	

	// fmt.Println("path in server file is ",r.URL.Path)
	
    // p := "." + r.URL.Path
    // if p == "./dummy" {
    //     p = "./static/links.html"
	// }
	// http.ServeFile(w, r, p)
   
}
func serveFiles(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// We then get the name of the user from our cache, where we set the session token
	response, err := cache.Do("GET", sessionToken)
	if err != nil {
		// If there is an error fetching from cache, return an internal server error status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response == nil {
		// If the session token is not present in cache, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// if c != id.String(){
	// 	w.WriteHeader(http.StatusUnauthorized)
	// }
	// Finally, return the welcome message to the user
	// fmt.Println("Welcome ", string(response))


	fmt.Println("path in server file is ",r.URL.Path)
	
    p := "." + r.URL.Path
    if p == "./dummy" {
        p = "./links.html"
	}
	http.ServeFile(w, r, p)
   
}
func handleHome(w http.ResponseWriter, r *http.Request) {
	
	var tmpl = template.Must(template.New("templ").ParseFiles("index.html"))
	
	if err = tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}	
	}

var state = randToken()
// func DelKey(key string) (error) {
// 	return cache.Del(key).Err()
//  }
func handleLogout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
	if err != nil{
		w.WriteHeader(http.StatusUnauthorized)
	} 
	fmt.Println("inside logout")
	// deleted_key:=DelKey("session_token")
	// fmt.Println("in logout deleted key is ",deleted_key)
	t := &http.Cookie{
		Name:"session_token",
		
		MaxAge: -1,
	
		HttpOnly: true,
	}
	
	http.SetCookie(w, t)
	return
	 }
	

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside handleGoogleLogin")

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
	//user authentication
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if (err != nil){
	
		fmt.Println(err.Error())
		w.Write([]byte("Error in getting response from google"))
		http.Redirect(w, r, "http://127.0.0.1:8000/", http.StatusTemporaryRedirect)
		
		return
	}
	//for checking indiamart id
	var result map[string]interface{}
	var errJson error
	errJson = json.Unmarshal(content, &result)   // to parse content byte I have used map
		if errJson != nil {
			fmt.Println(errJson)
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
	//after this im id will go on and create a session
	fmt.Println(w, "Welcome IM User")

	// Create a new random session token
	sessionToken := id.String()
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	_, err = cache.Do("SETEX", sessionToken, "1200", emailString)
	if err != nil {
		// If there is an error in setting the cache, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		MaxAge:  1200,//time.Now().Add(time.Hour * 1),
	})
	  
	http.Redirect(w,r, "/dummy",http.StatusTemporaryRedirect)
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
