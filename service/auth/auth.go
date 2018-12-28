package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type key int

const ctxKey key = 2

type UserContext struct {
	ID int
}

type jwtClaims struct {
	UserContext UserContext
	jwt.StandardClaims
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
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}

	// todo: do not work
	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		return &claims.UserContext, nil
	} else {
		panic(ok)
		return nil, errors.New("invalid authorization token")
	}
}

func MiddlewareLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r
		login := ctx.FormValue("login")
		password := ctx.FormValue("password")

		user, err := getUser(login)
		if err != nil {
			panic(err.Error())
		}
		if user == nil || !CheckPasswordHash(password, user.Password) {
			panic("authorization failed:1")
		}

		expires := 24 * time.Hour
		clm := &jwtClaims{
			UserContext: UserContext{ID:*user.ID},
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(expires).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, clm)

		tokenString, _ := token.SignedString(getSecretKey())

		w.Write([]byte(tokenString))
	}
}

func MiddlewareRegistration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*ctx := r

		if r.Method == http.MethodPost {
			login := ctx.FormValue("login")
			password := ctx.FormValue("password")
			name := ctx.FormValue("name")
			surname := ctx.FormValue("surname")
			patronymic := ctx.FormValue("patronymic")
			setUser()
		} else {
			io.WriteString(w, r.Method + " is not implemented")
		}*/
	}
}



//todo clarify model
type User struct {
	ID 			*int	`json:"id" db:"id"`
	Login		string	`json:"login" db:"login" validate:"required"`
	Password   	string  `json:"password,omitempty" db:"password" validate:"required"`
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func getUser(login string) (*User, error) {//todo clarify model
	//todo find user in db
	if login != "kaz_avto_trans" {
		return nil, errors.New("authorization failed:2")
	}
	id := 1
	return &User{
		&id,
		"kaz_avto_trans",
		"$2a$10$JewUnzN6b6dVrnWNvtlbdOzzticUee.MNPX6.h.2L8.8pI/nQU4sa",
	}, nil
}
/*func setUser(login string) (error) {//todo clarify model
	return nil
}*/
