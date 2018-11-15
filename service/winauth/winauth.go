package winauth

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	//"io/ioutil"
	//"log"

	//"time"

	//"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	//"gopkg.in/jcmturner/gokrb5.v6/client"
	//"gopkg.in/jcmturner/gokrb5.v6/config"

	//"github.com/davecgh/go-spew/spew"
	//httpntlm "github.com/vadimi/go-http-ntlm"
)

const (
	Realm = "NEFIS.LOCAL"
	// todo перенести в config
	KDC = "192.168.2.10:88"
	kRB5CONF = `
[libdefaults]
  default_realm = NEFIS.LOCAL
[realms]
 NEFIS.LOCAL = {
  kdc = nefis.local
  admin_server = nefis.local
  default_domain = nefis.local
 }
[domain_realm]
 .nefis.local = NEFIS.LOCAL
 nefis.local = NEFIS.LOCAL
 `
)

type key int

const ctxKey key = 4

type UserContext struct {
	ID int
}

func withContext(req *http.Request, userCtx *UserContext) *http.Request {
	ctx := context.WithValue(req.Context(), ctxKey, userCtx)
	req = req.WithContext(ctx)
	return req
}

func GetContext(ctx context.Context) *UserContext {
	val := ctx.Value(ctxKey)
	if val == nil {
		return nil
	}
	return val.(*UserContext)
}

func getSecretKey() []byte {
	return []byte("secret")
}

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				data, err := ejectClaims(bearerToken[1])
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
				r = withContext(r, data)
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid template of authorization token"))
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing authorization token"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ejectClaims(tokenString string) (*UserContext, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return getSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user_id := int(claims["user_id"].(float64))
		return &UserContext{user_id}, nil
	} else {
		return nil, errors.New("invalid authorization token")
	}
}

func MiddlewareLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		fmt.Println(header)
		authReq := strings.Split(header, " ")

		if len(authReq) == 2 || authReq[0] == "Negotiate" {
			spew.Dump(r)
			w.WriteHeader(200)
		} else {
			w.Header().Set("WWW-Authenticate", "Negotiate")
			w.WriteHeader(http.StatusUnauthorized)
		}

		return

		/*ctx := r
		username := ctx.FormValue("username")
		password := ctx.FormValue("password")

		cl := client.NewClientWithPassword(username, Realm, password)

		// Load the client krb5 config
		conf, err := config.NewConfigFromString(kRB5CONF)
		if err != nil {
			panic(err)
		}

		conf.Realms[0].KDC = []string{KDC}

		// Apply the config to the client
		cl.WithConfig(conf)
		cl.GoKrb5Conf.DisablePAFXFast = true

		// Log in the client
		err = cl.Login()
		if err != nil {
			panic(err)
		}

		spew.Dump(cl)

		// Создаем новый токен
		claims := jwt.MapClaims{
			//todo откуда брать user_id
			"user_id": 1,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Подписываем токен нашим секретным ключем
		tokenString, _ := token.SignedString(getSecretKey())

		w.Write([]byte(tokenString))*/
	}
}