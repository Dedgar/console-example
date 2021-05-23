package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	sitePath = os.Getenv("SITE_PATH")
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// GET /home
func getHome(c echo.Context) error {
	return c.Render(http.StatusOK, "home.html", nil)
}

func main() {
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

	e := echo.New()
	e.Static("/", sitePath+"/static")
	e.Renderer = t

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", getHome)
	e.GET("/home", getHome)
	e.Logger.Info(e.Start(":8080"))
}
