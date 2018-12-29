package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

var DefaultConfig = Config{
	Host:     "localhost",
	Port:     8080,
}

type key int

const ctxKey key = 2

type User struct {
	ID 			*int	`json:"id" db:"id"`
	Login		string	`json:"login" db:"login" validate:"required"`
	Password   	string  `json:"password,omitempty" db:"password" validate:"required"`
	Name   		*string `json:"name,omitempty" db:"name" validate:"required"`
	Surname   	*string `json:"surname,omitempty" db:"surname" validate:"required"`
	Patronymic  *string `json:"patronymic,omitempty" db:"patronymic" validate:"required"`
}

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

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		return &claims.UserContext, nil
	} else {
		panic(ok)
		return nil, errors.New("invalid authorization token")
	}
}

func MiddlewareLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := r.FormValue("login")
		password := r.FormValue("password")

		user, err := getUser(login)
		if err != nil {
			panic(err.Error())
		}
		if user == nil || !CheckPasswordHash(password, user.Password) {
			panic("authorization failed:1")
		}

		w.Write([]byte(getJWT(user)))
	}
}

func MiddlewareRegistration(v *viper.Viper) http.HandlerFunc {
	cfg := DefaultConfig
	if err := v.UnmarshalKey("app", &cfg); err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			decoder := json.NewDecoder(r.Body)
			var user User

			err := decoder.Decode(&user)
			if err != nil {
				panic(err)
			}

			err = validate(user)
			if err != nil {
				panic(err)
			}

			err = setUser(&user)
			if err != nil {
				panic(err)
			}

			err = sendEmail(user, cfg)
			if err != nil {
				panic(err)
			}

			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, r.Method + " is not implemented")
		}
	}
}

func MiddlewareConfirm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			token := r.URL.Query().Get("token")

			userContext, err := ejectClaims(token)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			err = confirmUser(userContext)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, r.Method + " is not implemented")
		}
	}
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func sendEmail(user User, cfg Config) (error) {
	url := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
	token := getJWT(&user)
	email := user.Login
	message := fmt.Sprintf("Здравствуйте!\n" +
		"Вы получили это сообщение, так как Вы прошли регистрацию в домене ncsd.ru.\n\n" +
		"Для начала работы перейдите по ссылке: %s/confirm?token=%s\n" +
		"Обратите внимание, что домен должен быть равен ncsd.ru\n\n" +
		"Если Вы не проходили никакую регистрацию, то не обращайте внимание на сообщение.\n\n" +
		"Сообщение сгенерировано автоматически, отвечать не нужно", url, token)
	subject := fmt.Sprintf("Подтверждение регистрации в домене ncsd.ru")
	letter := map[string]string{"email": email, "subject": subject, "message": message}
	jsonLetter, _ := json.Marshal(letter)
	resp, err := http.Post("http://email.ms.stage.ncsd.ru", "application/json", bytes.NewBuffer(jsonLetter))

	if err != nil {
		return err
	}
	if (resp.StatusCode != http.StatusOK) {
		return errors.New("Не удалось отправить письма на почту для подтверждения регистрации")
	}

	return nil
}

func getJWT(user *User) (string) {
	expires := 24 * time.Hour
	clm := &jwtClaims{
		UserContext: UserContext{ID:*user.ID},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expires).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clm)

	tokenString, _ := token.SignedString(getSecretKey())

	return tokenString
}






//todo find user in db
func getUser(login string) (*User, error) {
	if login != "kaz_avto_trans" {
		return nil, errors.New("authorization failed:2")
	}
	id := 1
	return &User{
		&id,
		"kaz_avto_trans",
		"$2a$10$JewUnzN6b6dVrnWNvtlbdOzzticUee.MNPX6.h.2L8.8pI/nQU4sa",
		nil,
		nil,
		nil,
	}, nil
}

//todo insert user in db
func setUser(user *User) (error) {
	id := 1
	user.ID = &id
	return nil
}

//todo реализовать заглушки
func confirmUser(userContext *UserContext) (error) {
	fmt.Println(userContext)
	return nil
}

//todo implemente
func validate(user User) (error) {
	return nil
}
