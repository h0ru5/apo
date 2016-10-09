package main

import (
	"encoding/json"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

//128 random chars
var mySigningKey = []byte("M6qOdsDc_xSDfg9esYxjA5MaARJAMPl1btKk_924lVNVh9Kw9MREuulNDq_7eT4e")

func main() {

	StartServer()

}

func StartServer() {
	r := mux.NewRouter()

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		UserProperty:  "token",
	})

	r.Handle("/open", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(SecuredPingHandler)),
	))

	http.Handle("/", r)
	fmt.Println("serving secured endpoints under http://localhost:3001/open")
	http.ListenAndServe(":3001", nil)
}

type Response struct {
	Text string `json:"text"`
}

func respondJson(text string, w http.ResponseWriter) {
	response := Response{text}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func SecuredPingHandler(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "token").(*jwt.Token).Claims.(jwt.MapClaims)
	msg := fmt.Sprintf("You are authenticated as: %v for %v, expires in %v", claims["sub"], claims["aud"],
		time.Since(time.Unix(int64(claims["exp"].(float64)), 0)))
	respondJson(msg, w)
}
