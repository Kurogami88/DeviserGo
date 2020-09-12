package main

/* PARSE IN VALUE
1. Redis import
1. Token Storing Method
*/
var apiLogin = `
/***
	Author: Leong Kai Khee (Kurogami)
	Date: 2020

	Generated by DeviserGO
***/

package main

import (
	"encoding/json"
	"net/http"
	"strings"
	%s
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

//LoginInput is the fields expected to be received
type LoginInput struct {
	Username string ` + "`json:\"username\"`" + `
	Password string ` + "`json:\"password\"`" + `
}

//LoginOutput is the fields expected to be sent
type LoginOutput struct {
	Token string ` + "`json:\"token\"`" + ` 
}

//Login is an API to authenticate a user and return a token
func Login(w http.ResponseWriter, r *http.Request) {
	result := DeviserResponse{HTTPStatus: 200, Result: LoginOutput{}}

	var input LoginInput
	json.NewDecoder(r.Body).Decode(&input)

	account, _ := DBAccountRetrieveCondition("` + "`username`" + ` = '" + input.Username + "'")
	if len(account) != 1 {
		result = DeviserResponse{HTTPStatus: 400, Result: "Error logging in"}
		result.DoResponse(w)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(*account[0].Password), []byte(input.Password))
	if err != nil {
		result = DeviserResponse{HTTPStatus: 400, Result: "Error logging in"}
		result.DoResponse(w)
		return
	}

	uuid := strings.Replace(uuid.New().String(), "-", "", -1)

	token, err := AuthCreateToken(uuid, *account[0].Username, *account[0].Role)
	if err != nil {
		result = DeviserResponse{HTTPStatus: 400, Result: "Error creating token"}
		result.DoResponse(w)
		return
	}

	dbToken := Token{
		Uuid:     &uuid,
		Username: account[0].Username,
		Role:     account[0].Role,
	}

	%s

	result.Result = LoginOutput{Token: token}
	
	result.DoResponse(w)
	return
}
`

var apiLoginRedisImport = `	"os"
	"time"
`
var apiLoginRedis = `envAuthDuration, err := time.ParseDuration(os.Getenv("JWT_EXPIRY_MIN") + "m")
	if err != nil {
		result = DeviserResponse{HTTPStatus: 400, Result: "Error storing tokens"}
	}

	jsonToken, err := json.Marshal(dbToken)
	if err != nil {
		result = DeviserResponse{HTTPStatus: 400, Result: "Error storing token"}
		result.DoResponse(w)
		return
	}

	err = cache.Set("Token_"+uuid, jsonToken, envAuthDuration).Err()
	if err != nil {
		result = DeviserResponse{HTTPStatus: 400, Result: "Error storing token"}
		result.DoResponse(w)
		return
	}`
var apiLoginDB = `_, err = DBTokenCreate(dbToken)
	if err != nil {
		result = DeviserResponse{HTTPStatus: 400, Result: "Error storing token"}
		result.DoResponse(w)
		return
	}`
