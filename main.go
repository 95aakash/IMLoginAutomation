package main

import (
	"encoding/json"
	"reflect"
	// "time"
	"crypto/rand"
    "encoding/base64"
	"fmt"
	"io/ioutil"
    // "os"
	"net/http"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"database/sql"
   _ "github.com/lib/pq"
   "strings"
   "log"
)

func init() {
	fmt.Println("Initialising..")
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "testdb"
	  )
	  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	  "password=%s dbname=%s sslmode=disable",
	  host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	  }
	  defer db.Close()
	
	  err = db.Ping()
	  if err != nil {
		panic(err)
	  }
	
	  fmt.Println("Successfully connected!")
}
var cid = "1045104154777-qpo2aust2h8sl8rtjpvqnnlt6pr4mfpa.apps.googleusercontent.com"
var csecret = "xDIGP6Wi1YOdDO_4FxZEk830"
var conf = &oauth2.Config{
	// ClientID:  os.Getenv("client_id"),
	ClientID: cid,
	ClientSecret: csecret,
	
	RedirectURL: "http://127.0.0.1:8000/callback",
	Scopes: []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",

},
Endpoint: google.Endpoint,
}

func main() {
	fmt.Println("inside main confid ",conf.ClientID)
	if conf.ClientID == ""{
		fmt.Println("client id is null")
	}else {
		fmt.Println("id not null")
		}

	fmt.Println("inside main- going to start server")
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	// http.HandleFunc("/wrongUser", handleWrongUser)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
// func handleWrongUser(w http.ResponseWriter, r *http.Request){
// // 	var htmlstring string
// // 	htmlstring = `<html>
// // <body>
// // 	<h1>Not an IM User<h1>
// // </body>
// // </html>`

// // fmt.Fprintf(w, htmlstring)
// time.Sleep(3 * time.Second)
// fmt.Fprintf(w,"Now redirecting to login page")
// fmt.Println("Now redirecting to login")
// http.Redirect(w,r,"/login",http.StatusTemporaryRedirect)


// }


func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlstring string
	htmlstring = `<html>
	<body>
		<a href="/login">Google Log In</a>
	</body>
	</html>`
	fmt.Fprintf(w, htmlstring)
	}

var state = randToken()
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
		// w.Write([]byte("Not an Indiamart User, redirecting to login"))
		// http.Redirect(w,r,"http://127.0.0.1:8000/",http.StatusTemporaryRedirect)
		// 
		var htmlstring string
		htmlstring = `<html>
		<body>
			<h1>Not an IM User<h1>
			<a href="/login">Go to Login again</a>
		</body>
		</html>`
		fmt.Fprintf(w, htmlstring)
		return
	}

	if (err != nil){
	
		fmt.Println(err.Error())
		w.Write([]byte("Error in getting response from google"))
		http.Redirect(w, r, "http://127.0.0.1:8000/", http.StatusTemporaryRedirect)
		
		return
	}
	
	

	// fmt.Fprintf(w, "Welcome IM User")
	// http.Redirect(w,r, "http://127.0.0.1:8000/",http.StatusTemporaryRedirect)  //redirecting an im user

	var htmlstring string
	htmlstring = `<html>
	<body>
		<h1>Welcome IM User<h1>
		<a href="https://github.com/golang/go/issues/14115">Click for redirection</a>
	</body>
	</html>`
	fmt.Fprintf(w, htmlstring)
	//database store
	
	// seen :=false
	

	// sqlStatement := `SELECT email FROM users WHERE email=$1;`
	// var email string


	// row := db.QueryRow(sqlStatement, //emailid )
	// switch err := row.Scan(&id, &email); err {
	// case sql.ErrNoRows:
	// fmt.Println("No rows were returned!")
	// case nil:
	// fmt.Println(id, email)
	// default:
	// panic(err)
	// }

	  

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

