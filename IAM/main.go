package main

import (
	"fmt"
	"github.com/abbot/go-http-auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/foomo/htpasswd"
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func SetupConfig() {
	// defaults
	viper.SetDefault("endpoint", ":3000")
	viper.SetDefault("passfile", "./passes")
	viper.SetDefault("key", "secret key please change this")

	// set names and expected directories
	viper.SetConfigName("iam-conf")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.iam")

	// config via environment
	viper.SetEnvPrefix("iam")
	viper.BindEnv("passfile")
	viper.BindEnv("endpoint")

	// posix flags
	pflag.StringP("endpoint", "e", ":3000", "endpoint to run the IAM (default ':3000')")
	pflag.StringP("passfile", "f", "./passes", "htpasswd file to operate on")
	pflag.StringP("key", "k", "secret", "key for HS256 signature")
	viper.BindPFlag("endpoint", pflag.Lookup("endpoint"))
	viper.BindPFlag("passfile", pflag.Lookup("passfile"))
	viper.BindPFlag("key", pflag.Lookup("key"))

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

func Secret(user, _ string) string {
	filepath := viper.GetString("passfile")

	if passwords, err := htpasswd.ParseHtpasswdFile(filepath); err == nil {
		if pw, ok := passwords[user]; ok {
			return pw
		} else {
			fmt.Printf("user %s not found\n", user)
		}
	} else {
		fmt.Println("Error getting users: ", err)
	}

	// they expect an empty string... seems a bit weak to me
	return ""
}

var GetTokenHandler = auth.AuthenticatedHandlerFunc(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	mySigningKey := []byte(viper.GetString("key"))

	claims := jwt.StandardClaims{
		Subject:   r.Username,
		Audience:  "Wohnung",
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	/* Create the token */
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	/* Finally, write the token to the browser window */
	w.Write([]byte(tokenString))
})

func StartServer() {
	filepath := viper.GetString("passfile")
	endpoint := viper.GetString("endpoint")

	r := mux.NewRouter()
	authenticator := auth.NewBasicAuthenticator("Wohnung", Secret)
	r.Handle("/token", authenticator.Wrap(GetTokenHandler))
	http.Handle("/", r)
	fmt.Println("IAM serving tokens under http://", endpoint, "/token using passfile ", filepath)
	http.ListenAndServe(endpoint, nil)
}
