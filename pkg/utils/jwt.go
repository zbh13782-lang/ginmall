package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSerect = []byte("zbh")

type Claims struct {
	Username  string `json:"username"`
	ID        uint   `json:"id"`
	Authority int    `json:"authority"`
	jwt.StandardClaims
}

type EmailClaims struct {
	UserID        uint   `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	OperationType uint   `json:"operation_type"`
	jwt.StandardClaims
}

func GenerateToken(id uint, username string, authority int) (token string, err error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Hour * 24)
	claims := Claims{
		Username:  username,
		ID:        id,
		Authority: authority,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "mall",
		},
	}

	tokenclaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenclaim.SignedString(jwtSerect)
	return
}

func ParseToken(token string) (*Claims, error) {
	tokenclaim, er := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSerect, nil
	})
	if tokenclaim != nil {
		if claims, ok := tokenclaim.Claims.(*Claims); ok && tokenclaim.Valid {
			return claims, nil
		}
	}
	return nil, er
}

func GenerateEmailToken(id ,operator uint ,email string,password string)(string,error){
	nowTime:=time.Now()
	expireTime:=nowTime.Add(time.Minute*30)
	claims:=EmailClaims{
		UserID: id,
		Email: email,
		Password: password,
		OperationType: operator,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer: "email",
		},
	}
	tokenClaims :=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	token,err:=tokenClaims.SignedString(jwtSerect)
	return token , err
}

//ParseEmailToken 验证邮箱验证token
func ParseEmailToken(token string) (*EmailClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSerect, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*EmailClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}