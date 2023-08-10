package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	echojwt "github.com/labstack/echo-jwt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
)

var (
	t *jwt.Token
)

func main() {

	key := []byte("secret")
	t = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "my-auth-server",
			"sub": "john",
			"foo": 2,
		})
	fmt.Println(t.SignedString(key))
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    key,
		SigningMethod: echojwt.AlgorithmHS256,
		ErrorHandler: func(c echo.Context, err error) error {
			fmt.Println(err.Error())
			return nil
		},
	}))

	e.GET("/", func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token) // by default token is stored under `user` key
		if !ok {
			return errors.New("JWT token missing or invalid")
		}
		claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
		if !ok {
			return errors.New("failed to cast claims as jwt.MapClaims")
		}
		return c.JSON(http.StatusOK, claims)
	})

	if err := e.Start(":18888"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
