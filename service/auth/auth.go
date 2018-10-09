package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
	"context"
)

type key int

const ctxKey key = 2

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
		/*ctx := r
		username := ctx.FormValue("username")
		password := ctx.FormValue("password")*/

		//todo find user in db
		/*user, err := __User(username)
		if err != nil {
			panic("User not finded")
		}
		if user == nil || !CheckPasswordHash(password, user.Password) {
			panic("authorization failed")
		}
		user.Password = ""
		expires := 24 * time.Hour
		clm := &jwtClaims{
			User: *user,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(expires).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, clm)

		sign, err := token.SignedString([]byte("secret"))
		if err != nil {
			panic("token signed")
		}
		fmt.Println(sign)*/

		// Создаем новый токен
		claims := jwt.MapClaims{
			"user_id": 1,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Подписываем токен нашим секретным ключем
		tokenString, _ := token.SignedString(getSecretKey())

		w.Write([]byte(tokenString))
	}
}

/*
type jwtClaims struct {
	User User
	jwt.StandardClaims
}

type User struct {
	Username   string            `json:"username" db:"username" validate:"required"`
	Password   string            `json:"password,omitempty" db:"password" validate:"required"`
}
func __User(username string) (*User, error) {
	return &User{
		"kaz_avto_trans",
		"$2a$10$JewUnzN6b6dVrnWNvtlbdOzzticUee.MNPX6.h.2L8.8pI/nQU4sa",
	}, nil
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}*/
