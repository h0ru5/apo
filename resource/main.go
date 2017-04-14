package main

import (
	"crypto"
	"encoding/json"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mendsley/gojwk"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
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

func GetKey(uri string) (crypto.PublicKey, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jwk, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("using key: ", string(jwk))
	jwkObj, err := gojwk.Unmarshal(jwk)
	if err != nil {
		return nil, err
	}
	return jwkObj.DecodePublicKey()
}

func StartServer() {
	r := mux.NewRouter()
	//myKey := []byte(viper.GetString("key"))
	myKey, _ := GetKey("http://localhost:3000/key")

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return myKey, nil
		},
		SigningMethod: jwt.SigningMethodES256,
		UserProperty:  "token",
	})

	r.Handle("/open", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(OpenHandler)),
	))

	r.Handle("/close", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(CloseHandler)),
	))

	r.Handle("/", negroni.New(
		negroni.Wrap(http.HandlerFunc(RootHandler)),
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

	msg = msg + " - opening"

	respondJsonText(msg, w)
}

func CloseHandler(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "token").(*jwt.Token).Claims.(jwt.MapClaims)
	expires := time.Since(time.Unix(int64(claims["exp"].(float64)), 0))

	msg := fmt.Sprintf("Authenticated as: %v for %v, expires in %v", claims["sub"], claims["aud"], expires)

	fmt.Println(msg)

	msg = msg + " - closing"

	respondJsonText(msg, w)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	links := map[string]interface{}{
		"links": []Link{
			{"action", "open", "POST"},
			{"action", "close", "POST"},
		},
	}
	respondJson(links, w)
}

type Link struct {
	Rel    string
	Href   string
	Method string
}
