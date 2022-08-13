package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/abbot/go-http-auth"
	"github.com/foomo/htpasswd"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/mendsley/gojwk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var mySigningKey *ecdsa.PrivateKey
var jwk []byte

func SetupConfig() {
	// defaults
	viper.SetDefault("audience", "myhome")
	viper.SetDefault("endpoint", ":3000")
	viper.SetDefault("passfile", "./passes")

	// set names and expected directories
	viper.SetConfigName("iam-conf")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.iam")

	// config via environment
	viper.SetEnvPrefix("iam")
	viper.BindEnv("passfile")
	viper.BindEnv("endpoint")
	viper.BindEnv("audience")

	// posix flags
	pflag.StringP("endpoint", "e", ":3000", "endpoint to run the IAM (default ':3000')")
	pflag.StringP("passfile", "f", "./passes", "htpasswd file to operate on")
	pflag.StringP("audience", "a", "myhome", "audience/realm that gets protected")
	viper.BindPFlag("endpoint", pflag.Lookup("endpoint"))
	viper.BindPFlag("passfile", pflag.Lookup("passfile"))
	viper.BindPFlag("audience", pflag.Lookup("audience"))

	//read in
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
	}
}

func main() {
	SetupConfig()
	StartServer()
}

/**
 * generates an ECDSA keypair
 */
func GenNewKeyPair() {
	mySigningKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pubKey, _ := gojwk.PublicKey(mySigningKey.Public())
	jwk, _ = gojwk.Marshal(pubKey)
	log.WithField("key", string(jwk)).Info("generated new Keypair")
}

func basicAuthorize(user, _ string) string {
	filepath := viper.GetString("passfile")

	if passwords, err := htpasswd.ParseHtpasswdFile(filepath); err == nil {
		if pw, ok := passwords[user]; ok {
			return pw
		} else {
			log.WithField("user", user).Warn("invalid login attempt")
		}
	} else {
		log.WithError(err).Error("could not access users backend")
	}

	// they expect an empty string... seems a bit weak to me
	return ""
}

var GetTokenHandler = auth.AuthenticatedHandlerFunc(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {

	claims := jwt.StandardClaims{
		Subject:   r.Username,
		Audience:  []string{viper.GetString("audience")},
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	/* Create the token */
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	/* Finally, write the token to the request */
	w.Write([]byte(tokenString))

	log.WithFields(log.Fields{
		"user":     r.Username,
		"clientIp": r.RemoteAddr,
	}).Info("issued new token")
})

func StartServer() {
	filepath := viper.GetString("passfile")
	endpoint := viper.GetString("endpoint")
	audience := viper.GetString("audience")

	GenNewKeyPair()

	r := mux.NewRouter()
	authenticator := auth.NewBasicAuthenticator(audience, basicAuthorize)
	r.Handle("/token", authenticator.Wrap(GetTokenHandler))
	r.HandleFunc("/key", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jwk))
	})
	http.Handle("/", r)
	log.WithFields(log.Fields{
		"endpoint":   "http://" + endpoint + "/token",
		"audience":   audience,
		"passwdfile": filepath,
	}).Info("started IAM")
	//fmt.Printf("IAM is serving tokens under http://%s/token for %s using passfile %s\n", endpoint, audience, filepath)
	panic(http.ListenAndServe(endpoint, nil))
}
