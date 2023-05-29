package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	env "raptor/config"
	"raptor/logger"
	"raptor/pkg/helpers"
	"raptor/pkg/templates"
)

type CustomClaims struct {
	jwt.StandardClaims
	Role string `json:"role"`
}

// RegisterAuthentication only add to protected router groups
func RegisterAuthentication() func(next echo.HandlerFunc) echo.HandlerFunc {

	auth := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Read public key file
			keyBytes, err := ioutil.ReadFile(env.GetString("security.public_key_path"))
			if err != nil {
				logger.Info(err)
			}

			// Parse public key
			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
			if err != nil {
				logger.Info(err)
			}
			tokenString, err := helpers.GetJwtStringFromRequest(c.Request())
			if err != nil {
				templates.Unauthorized(err)
				return c.JSON(http.StatusUnauthorized, templates.Unauthorized(err))
			}
			// Retrieve JWT from context

			_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return publicKey, nil
			})
			if err != nil {
				templates.Unauthorized(err)
				return c.JSON(http.StatusUnauthorized, templates.Unauthorized(err))
			}
			// Access claims
			//claims := token.Claims.(jwt.MapClaims)
			//fmt.Println("Subject:", claims["sub"])
			//fmt.Println("Name:", claims["name"])
			//fmt.Println("Issued at:", claims["iat"])

			return next(c)
		}
	}

	return auth
}
