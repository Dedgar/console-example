package controllers

import (
	"github.com/dedgar/console-example/datastores"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"

	"github.com/gorilla/sessions"

	"fmt"
	"net/http"
)

// GET /
func GetMain(c echo.Context) error {
	return c.Render(http.StatusOK, "main.html", datastores.PostMap)
}

// GET /login
func GetLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

// GET /dev/login
func GetDevLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "devlogin.html", nil)
}

// GET /about
func GetAbout(c echo.Context) error {
	return c.Render(http.StatusOK, "about.html", nil)
}

// GET /contact
func GetContact(c echo.Context) error {
	return c.Render(http.StatusOK, "contact.html", nil)
}

// GET /register
func GetRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

// GET /dev/register
func GetDevRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "devregister.html", nil)
}

// GET /privacy
func GetPrivacy(c echo.Context) error {
	return c.Render(http.StatusOK, "privacy.html", nil)
}

// GET /graph
func GetGraph(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sess)
	//if _, ok := sess.Values["current_user"].(string); ok {
	graphGet := map[string]int{"January": 100, "February": 200, "March": 300, "April": 400, "May": 500, "June": 600, "July": 700, "August": 800, "September": 900, "October": 1000, "November": 1100, "December": 1200}

	graphMap := map[string]interface{}{
		"catCount":  graphGet,
		"appDomain": "localhost:8080",
	}
	return c.Render(http.StatusOK, "graph_j_pie.html", graphMap)
	//}
	//return c.Redirect(http.StatusPermanentRedirect, "/login")
}

// GET /api/graph
func GetApiGraph(c echo.Context) error {
	callback := c.QueryParam("callback")
	month := []string{"January", "February", "March", "April", "May"} //, "June", "July", "August", "September", "October", "November", "December"}
	content := make(map[string]int)
	for i, item := range month {
		content[item] = (i + 1) * 300
	}
	return c.JSONP(http.StatusOK, callback, &content)
}

// GET /api/file
func GetApiFile(c echo.Context) error {
	callback := c.QueryParam("callback")

	content := "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg=="
	return c.JSONP(http.StatusOK, callback, &content)
}

// GET /trial
func GetTrial(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		fmt.Println("error getting session")
	}

	if sess.Values["Authenticated"] == "true" {
		fmt.Println("User is authenticated! Along with the following session info:")
		//fmt.Println(sess.Values)
		return c.String(http.StatusOK, "You are logged in.")
	}
	//if _, ok := sess.Values["current_user"].(string); ok {
	//loggedInUser := sess.Values["current_user"].(string)
	//return c.String(http.StatusOK, loggedInUser)
	//}
	//return c.Redirect(http.StatusPermanentRedirect, "/login")
	return c.String(http.StatusOK, "You are not logged in.")
}

// handle any error by attempting to render a custom page for it
func Custom404Handler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.Render(code, errorPage, code); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

// AuthMiddleware ensures a user has been logged in before continuing
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
			MainSession(c)
			sess, err := session.Get("session", c)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("DEBUG: authmiddleware")

			if sess.Values["Authenticated"] == "true" {
				fmt.Println("in authmiddleware if block")
				fmt.Println(sess.Values)
				return next(c)
			}
			return c.Redirect(http.StatusTemporaryRedirect, "/login/google")
		}
	}
}

// MainSession sets initial keys for a new user session
func MainSession(c echo.Context) { //error {
	sess, err := session.Get("session", c)
	if err != nil {
		fmt.Println("Error getting session info from MainSession")
	}
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 86400, // * 7,
		//HttpOnly: true,
	}
	sess.Values["Authenticated"] = "true"
	sess.Save(c.Request(), c.Response())
}
