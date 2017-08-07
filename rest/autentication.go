package rest

import (
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	"github.com/enricod/1h1dphoto.com-be/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"fmt"
	"github.com/enricod/1h1dphoto.com-be/db"
	"gopkg.in/gomail.v2"
)

var Tokens = make(map[string]db.User)

/**
 * riceve in post il tipo UserRegisterReq. Se utente non esiste già sul database, lo creo.
 * genero codice alfanumerico che poi sarà spedito via email.
 * Restituisce un tipo UserRegisterRes
 */
func UserRegister(res http.ResponseWriter, req *http.Request) {

	var userRegisterReq model.UserRegisterReq
	if req.Body == nil {
		http.Error(res, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&userRegisterReq)
	if err != nil {
		http.Error(res, err.Error(), 400)
		return
	}
	//fmt.Println(userRegisterReq.Username)



	// crea record per sessione
	validationCode  := model.RandStringBytes(5)
	userToken 		:= model.RandStringBytes(32)
	head := model.ResHead{Success:true}
	body := model.UserRegisterResBody{
		UserToken: userToken }
	userRegisterRes := model.UserRegisterRes{Head:head, Body:body}



	fmt.Println("validationCode: ",  validationCode );

	// crea utente nel database se necessario
	if user, err := db.UserFindByEmail(userRegisterReq.Email); err != nil {
		user := db.User{Username:userRegisterReq.Username,Email: userRegisterReq.Email}
		db.SalvaUser( &user )
		// crea record in USER_APP_TOKEN
		db.SalvaAppToken(user.ID, userToken)
	} else {
		// crea record in USER_APP_TOKEN
		db.SalvaAppToken(user.ID, userToken)
	}

	// spedisci via email il codice di validazione
	m := gomail.NewMessage()
	m.SetHeader("From", "1h1dphoto@gmail.com")
	m.SetHeader("To", userRegisterReq.Email)
	m.SetHeader("Subject", "validation code")
	m.SetBody("text/html", "Your validation code <b> " + validationCode + "</b>")


	d := gomail.NewDialer("smtp.gmail.com", 465, "1h1dphoto@gmail.com", "1h1dphotos")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}


	res.WriteHeader(http.StatusOK)
	err2 := json.NewEncoder(res).Encode(userRegisterRes)
	if err2 != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func UserCodeValidation(res http.ResponseWriter, req *http.Request) {
	var userCodeValidationReq model.UserCodeValidationReq
	if req.Body == nil {
		http.Error(res, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&userCodeValidationReq)
	if err != nil {
		http.Error(res, err.Error(), 400)
		return
	}

	// aggiorna USER_APP_TOKEN con informazione che utente ha validato email
	// CHECK_VALID => true
}

func Logout(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	sToken := vars["token"]
	delete(Tokens, sToken)
	res.WriteHeader(http.StatusOK)
}

func Sessions(res http.ResponseWriter, req *http.Request) {
	s := reflect.ValueOf(Tokens).MapKeys()
	var str string
	for _, element := range s {
		str += element.String() + ","
	}
	err := json.NewEncoder(res).Encode(model.Response{Data: str})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func createToken(username string, host string) *jwt.Token {
	expireToken := time.Now().Add(time.Hour * 1).Unix()
	claim := model.Claims{
		username,
		time.Now().Round(time.Millisecond).UnixNano(),
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    host,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token
}
