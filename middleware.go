package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"os"
	"strings"
	"task/helpers"
)

var Collection = helpers.MongoConnection()

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		jwtSecret := os.Getenv("JWT_SECRET")
		currentEndPoint := req.RequestURI
		headers := make(map[string]interface{})
		requestHeaders := req.Header
		for key, value := range requestHeaders {
			headers[key] = value[0]
		}
		currentEndPointBits := strings.Split(currentEndPoint, "/")
		if len(currentEndPointBits) > 1 {
			if headers["Authorization"] != nil {
				authHeaders := headers["Authorization"].(string)
				bits := strings.Split(authHeaders, " ")
				if len(bits) == 2 && bits[1] != "" {
					tkn := bits[1]
					tokenStr := tkn
					token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, errors.New("error while getting data")
						}
						return []byte(jwtSecret), nil
					})
					if token != nil {
						preservedUser := token.Claims.(jwt.MapClaims)
						var userContext map[string]interface{}
						bytesArray, err := json.Marshal(preservedUser)
						err = json.Unmarshal(bytesArray, &userContext)
						if err == nil {
							if preservedUser["user"] != nil {
								user := preservedUser["user"].(map[string]interface{})
								checkUserLoggedOut, err := Collection["token"].Find(bson.M{"user_id": bson.ObjectIdHex(user["_id"].(string))}).Count()
								if checkUserLoggedOut != 0 && err == nil {
									ctx := context.WithValue(req.Context(), "user_id", user["_id"])
									req = req.WithContext(ctx)
									next.ServeHTTP(res, req)
									return
								}
							}
						}
					}
				}
			}
		}
		helpers.ReplyBack(res, 401, "unauthorized")
		return
	})
}
