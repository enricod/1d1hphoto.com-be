package model

import jwt "github.com/dgrijalva/jwt-go"

var MyKey = "1d1hphoto"

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}