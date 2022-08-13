package main

import (
	"crypto"
	"encoding/json"
	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mendsley/gojwk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
)

func SetupConfig() {
	// defaults
	viper.SetDefault("key", "http://localhost:3000/key")

	// conf name & locations
	viper.SetConfigName("iam-conf")
	viper.AddConfigPath(".")

	// posix flags
	pflag.StringP("key", "k", "http://localhost:3000/key", "iam endpoint for signing key")
	viper.BindPFlag("key", pflag.Lookup("key"))

	// reading config
	err := viper.ReadInConfig()
	if err != nil {
		log.Warn("No config loaded - reverting to defaults")
	}
}

func main() {
	SetupConfig()
	StartServer()
}

func GetKey(uri string) (crypto.PublicKey, error) {
	log.WithField("uri", uri).Debug("acquiring key")

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jwk, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.WithField("key", string(jwk)).Debug("loaded key")
	jwkObj, err := gojwk.Unmarshal(jwk)
	if err != nil {
		return nil, err
	}
	return jwkObj.DecodePublicKey()
}

func StartServer() {
	r := mux.NewRouter()
	keyEndpoint := viper.GetString("key")
	myKey, err := GetKey(keyEndpoint)
	if err != nil {
		panic(err)
	}

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return myKey, nil
		},
		SigningMethod: jwt.SigningMethodES256,
		UserProperty:  "token",
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			log.WithFields(log.Fields{
				"ip":    r.RemoteAddr,
				"error": err,
			}).Error("authentification failure")
			jwtmiddleware.OnError(w, r, err)
		},
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
	log.WithField("endpoint", "http://:3001/").Info("started server")
	log.Fatal(http.ListenAndServe(":3001", nil))
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
	expires := time.Until(time.Unix(int64(claims["exp"].(float64)), 0))

	log.WithFields(log.Fields{
		"sub":    claims["sub"],
		"aud":    claims["aud"],
		"exp":    expires,
		"ip":     r.RemoteAddr,
		"action": "open",
	}).Warn("opening door")

	respondJson(map[string]interface{}{
		"sub":    claims["sub"],
		"aud":    claims["aud"],
		"exp":    expires,
		"ip":     r.RemoteAddr,
		"action": "open",
	}, w)
}

func CloseHandler(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "token").(*jwt.Token).Claims.(jwt.MapClaims)
	expires := time.Until(time.Unix(int64(claims["exp"].(float64)), 0))

	log.WithFields(log.Fields{
		"sub":    claims["sub"],
		"aud":    claims["aud"],
		"exp":    expires,
		"ip":     r.RemoteAddr,
		"action": "close",
	}).Warn("closing door")

	respondJson(map[string]interface{}{
		"sub":    claims["sub"],
		"aud":    claims["aud"],
		"exp":    expires,
		"ip":     r.RemoteAddr,
		"action": "close",
	}, w)
}

func RootHandler(w http.ResponseWriter, _ *http.Request) {
	links := map[string]interface{}{
		"links": []Link{
			{"action", "/open"},
			{"action", "/close"},
		},
	}
	respondJson(links, w)
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}
