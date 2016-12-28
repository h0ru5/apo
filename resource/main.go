package main

import (
	"encoding/json"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func SetupConfig() {
	// defaults
	viper.SetDefault("key", "secret key please change this")

	// conf name & locations
	viper.SetConfigName("iam-conf")
	viper.AddConfigPath(".")

	// posix flags
	pflag.StringP("key", "k", "secret", "key for HS256 signature")
	viper.BindPFlag("key", pflag.Lookup("key"))

	// reading config
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
	}
}

func main() {
	SetupConfig()
	StartServer()
}

func StartServer() {
	r := mux.NewRouter()
	myKey := []byte(viper.GetString("key"))

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return myKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		UserProperty:  "token",
	})

	r.Handle("/open", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(OpenHandler)),
	))

	http.Handle("/", r)
	fmt.Println("serving secured endpoint under http://localhost:3001/open")
	http.ListenAndServe(":3001", nil)
}

func respondJsonText(text string, w http.ResponseWriter) {
	resp := map[string]interface{}{"text": text}
	respondJson(resp, w)
}

func respondJson(response map[string]interface{}, w http.ResponseWriter) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func OpenHandler(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "token").(*jwt.Token).Claims.(jwt.MapClaims)
	expires := time.Since(time.Unix(int64(claims["exp"].(float64)), 0))

	msg := fmt.Sprintf("Authenticated as: %v for %v, expires in %v", claims["sub"], claims["aud"], expires)

	fmt.Println(msg)

	// do your action

	respondJsonText(msg, w)
}
