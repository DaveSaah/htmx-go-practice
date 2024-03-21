package main

import (
	"htmx-practice/actions"
	"io"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Templates struct for parsed html templates
type Templates struct {
	templates *template.Template
}

// Render renders html templates
func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// newTemplate creates a new template
func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type LoginError struct {
	Msg    string
	Email  string
	Passwd string
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Static("assets"))
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	e.GET("/login", func(c echo.Context) error {
		return c.Render(200, "login", nil)
	})

	e.POST("/login", func(c echo.Context) error {
		email := c.FormValue("email")
		passwd := c.FormValue("passwd")
		person, err := actions.Login(email, passwd)
		if err != nil {
			return c.Render(200, "login", &LoginError{
				Msg:    err.Error(),
				Email:  email,
				Passwd: passwd,
			})
		}

		return c.Render(200, "home", person)
	})

	// start server with logger
	e.Logger.Fatal(e.Start(":42069"))
}
