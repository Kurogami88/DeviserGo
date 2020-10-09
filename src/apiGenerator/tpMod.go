package main

/* PARSE IN VALUE
1. Project Name
*/
var tpMod = `module %s

go 1.14

require (
	github.com/auth0/go-jwt-middleware v0.0.0-20200810150920-a32d7af194d1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/jinzhu/gorm v1.9.16
	github.com/joho/godotenv v1.3.0
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
)
`
