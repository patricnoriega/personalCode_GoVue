package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-yaml/yaml"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Handler
func main() {
	// setup echo
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/apiobj/:filename", ObjHandler())

	// start webserver
	http.Handle("/", e)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// convert yaml to json
func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

// ObjectHandler endpoint handler
func ObjHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// name := c.QueryParam("name")
		name := c.Param("filename")
		file := "./objects/" + name + ".yaml"
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error()})

		}
		var body interface{}
		if err := yaml.Unmarshal([]byte(data), &body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error()})
		}

		body = convert(body)

		return c.JSON(http.StatusOK, body)
	}
}
