package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dedgar/console-example/datastores"
	"github.com/dedgar/console-example/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	defaultCost, _    = strconv.Atoi(os.Getenv("DEFAULT_COST"))
	oauthStateString  = "random" // TODO randomize
	googleOauthConfig = &oauth2.Config{
		ClientID:     datastores.OAuthID,
		ClientSecret: datastores.OAuthKey,
		//RedirectURL:  "https://www.tacofreeze.com/oauth/callback",
		RedirectURL: "http://127.0.0.1:8080/oauth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
)

// JWTLogin handles the request and page
// POST /jwt/login
func JWTLogin(c echo.Context) error {
	userName := c.FormValue("username")
	password := c.FormValue("password")

	if !userFound(c.FormValue("username")) {
		fmt.Println("user not found")
		return c.String(http.StatusOK, "Username not found!")
	}

	if !compareLogin(userName, password) {
		// A bad username or password gets a 401
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &models.JWTCustomClaims{
		name:  userName,
		admin: true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// HandleGoogleCallback listens on
// GET /oauth/callback
func HandleGoogleCallback(c echo.Context) error {
	state := c.QueryParam("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	code := c.QueryParam("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("error getting response")
		fmt.Println(err)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading response")
		fmt.Println(err)
	}

	var gUser models.GoogleUser
	err = json.Unmarshal(contents, &gUser)
	if err != nil {
		fmt.Println("Error Unmarshaling google user json: ", err)
	}

	if ok := datastores.AuthMap[gUser.Email]; ok {
		sess, _ := session.Get("session", c)
		sess.Values["Authenticated"] = "true"
		sess.Values["Google_logged_in"] = gUser.Email
		sess.Save(c.Request(), c.Response())

		return c.Render(http.StatusOK, "console.html", nil)
	}
	return c.String(http.StatusOK, string(contents)+`https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=`+token.AccessToken)
}

// HandleGoogleLogin handles the request and page
// GET /login/google
func HandleGoogleLogin(c echo.Context) error {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// PostLogin handles the request and page
// POST /login
func PostLogin(c echo.Context) error {
	fmt.Println("In PostLogin")
	if !userFound(c.FormValue("username")) {
		fmt.Println("user not found")
		return c.String(http.StatusOK, "Username not found!")
	}

	fmt.Println("comparing login")
	if compareLogin(c.FormValue("username"), c.FormValue("password")) {
		sess, _ := session.Get("session", c)
		sess.Values["User"] = c.FormValue("username")
		sess.Values["Authenticated"] = "true"
		sess.Save(c.Request(), c.Response())

		return c.Redirect(http.StatusPermanentRedirect, "/")
	}
	fmt.Println("postlogin - nothing worked")

	return c.Render(http.StatusUnauthorized, "404.html", "401 not authenticated")
}

// PostRegister handles the request and page
// POST /register
func PostRegister(c echo.Context) error {
	fmt.Println("Postregister")
	TextBody := c.FormValue("login") + "\n" + c.FormValue("password")
	fmt.Println("Textbody is: ", TextBody)

	if userFound(c.FormValue("username")) || emailFound(c.FormValue("email")) {
		return c.String(http.StatusOK, "Email address or username already taken, try again!")
	}

	createUser(c.FormValue("email"), c.FormValue("username"), c.FormValue("password"))

	return c.Redirect(http.StatusPermanentRedirect, "/login")
}

// hashPass returns a bcrypt hash string from the provided password string
func hashPass(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	return string(bytes), err
}

func createUser(eName, uName, pWord string) {
	fmt.Println("createUser:", eName, uName, pWord)
	hashedPW, err := hashPass(pWord)

	if err != nil {
		log.Fatal(err)
	}

	newUser := models.User{Email: eName, UName: uName, Password: hashedPW}
	datastores.DB.Create(&newUser)
}

func compareLogin(uName, pWord string) bool {
	var user models.User
	var foundU models.User

	datastores.DB.Where(&models.User{UName: uName}).First(&user).Scan(&foundU)

	if foundU.UName == "" {
		fmt.Println("Invalid username or password!")
		return false
	}

	hashedPW := foundU.Password

	// Implements subtle.ConstantTimeCompare
	err := bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(pWord))

	if err != nil {
		fmt.Println("Invalid username or password!")
		fmt.Println(err)
		return false
	}

	fmt.Println("Found login combo matched!")
	return true
}

func userFound(uName string) bool {
	var user models.User
	var foundU models.User

	fmt.Println("in userFound")
	fmt.Println("uname is", uName)

	datastores.DB.Where(&models.User{UName: uName}).First(&user).Scan(&foundU)

	if foundU.UName != "" {
		fmt.Println("Username found.")
		return true
	}

	fmt.Println("Username not found.")
	return false
}

func emailFound(eName string) bool {
	var user models.User
	var foundE models.User

	datastores.DB.Where(&models.User{Email: eName}).First(&user).Scan(&foundE)

	if foundE.Email != "" {
		fmt.Printf("%s already taken!\n", foundE.Email)
		return true
	}

	fmt.Printf("%s not taken!\n", foundE.Email)
	return false
}
