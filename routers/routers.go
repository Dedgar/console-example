package routers

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"

	"github.com/gorilla/sessions"

	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dedgar/console-example/auth"
	"github.com/dedgar/console-example/controllers"
	"github.com/dedgar/console-example/datastores"
)

var (
	// Routers supplies an instance of echo to be used in the main function.
	Routers *echo.Echo
	// sitePath is the actual run path of the code. Defaults to "."
	sitePath = os.Getenv("SITE_PATH")
)

// Template contains a pointer to a template.Template
type Template struct {
	templates *template.Template
}

// Render executes stored templates that were found in sitePath+/tmpl
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func init() {
	if sitePath == "" {
		sitePath = "."
	}
	t := &Template{
		templates: func() *template.Template {
			tmpl := template.New("")
			if err := filepath.Walk(sitePath+"/tmpl", func(path string, info os.FileInfo, err error) error {
				if strings.HasSuffix(path, ".html") {
					_, err = tmpl.ParseFiles(path)
					if err != nil {
						log.Println(err)
					}
				}
				return err
			}); err != nil {
				panic(err)
			}
			return tmpl
		}(),
	}

	Routers = echo.New()
	Routers.Pre(middleware.RemoveTrailingSlash())
	Routers.Static("/", sitePath+"/static")
	Routers.Renderer = t

	Routers.Use(middleware.Logger())
	Routers.Use(middleware.Recover())
	Routers.Use(middleware.CORS())
	Routers.Use(session.Middleware(sessions.NewCookieStore([]byte(datastores.CookieSecret))))

	// AuthMiddleware requires users be logged in with a particular email
	Routers.GET("/", controllers.GetMain)
	Routers.POST("/", controllers.GetMain)
	Routers.GET("/takedowns", controllers.GetGraph, controllers.AuthMiddleware())
	Routers.GET("/api/takedowns", controllers.GetApiGraph)
	Routers.GET("/login/google", auth.HandleGoogleLogin)
	Routers.GET("/oauth/callback", auth.HandleGoogleCallback)

	Routers.GET("/", controllers.GetMain)
	Routers.POST("/", controllers.GetMain)
	Routers.GET("/watch/:show/:season/:episode", controllers.GetShow)
	//Routers.GET("/watch/:show/:season/:episode/", controllers.GetShow)
	Routers.GET("/about", controllers.GetAbout)
	Routers.GET("/about-us", controllers.GetAbout)
	Routers.GET("/register", controllers.GetRegister)
	Routers.POST("/register", auth.PostRegister)
	Routers.GET("/dev/register", controllers.GetDevRegister)
	Routers.GET("/login", controllers.GetLogin)
	Routers.POST("/login", auth.PostLogin)
	Routers.GET("/dev/login", controllers.GetDevLogin)
	Routers.GET("/graph", controllers.GetGraph) //, controllers.AuthMiddleware())
	Routers.GET("/api/graph", controllers.GetApiGraph)
	Routers.GET("/contact", controllers.GetContact)
	Routers.GET("/contact-us", controllers.GetContact)
	Routers.GET("/privacy-policy", controllers.GetPrivacy)
	Routers.GET("/privacy", controllers.GetPrivacy)
	Routers.POST("/post-contact", controllers.PostContact)
	Routers.File("/robots.txt", sitePath+"/static/public/robots.txt")
	Routers.File("/sitemap.xml", sitePath+"/static/public/sitemap.xml")

	// Authenticated group
	authGroup := Routers.Group("/console")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &auth.JWTCustomClaims{},
		SigningKey: []byte("secret"),
	}

	authGroup.Use(middleware.JWTWithConfig(config))

	// Use the function on the base route of the group
	authGroup.GET("", controllers.restricted)
}
