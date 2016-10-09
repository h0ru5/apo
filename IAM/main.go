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

//128 random chars
var mySigningKey = []byte("M6qOdsDc_xSDfg9esYxjA5MaARJAMPl1btKk_924lVNVh9Kw9MREuulNDq_7eT4e")
var filepath = "./passes"

func main() {
	viper.SetDefault("endpoint", ":3000")
	viper.SetDefault("passfile", "./passes")
	viper.SetConfigName("iam-conf")

	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.iam")

	viper.SetEnvPrefix("iam")
	viper.BindEnv("passfile")
	viper.BindEnv("endpoint")

	pflag.StringP("endpoint", "e", ":3000", "endpoint to run the IAM (default ':3000')")
	pflag.StringP("passfile", "f", "./passes", "htpasswd file to operate on")

	viper.BindPFlag("port", pflag.Lookup("port"))
	viper.BindPFlag("passfile", pflag.Lookup("passfile"))

	StartServer()
}

func Secret(user, _ string) string {
	//TODO config this
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
	filepath = viper.GetString("passfile")
	endpoint := viper.GetString("endpoint")

	r := mux.NewRouter()
	authenticator := auth.NewBasicAuthenticator("Wohnung", Secret)
	r.Handle("/token", authenticator.Wrap(GetTokenHandler))
	http.Handle("/", r)
	fmt.Println("IAM serving tokens under http://", endpoint, "/token using passfile ", filepath)
	http.ListenAndServe(endpoint, nil)
}
