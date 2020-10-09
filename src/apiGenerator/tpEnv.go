package main

/* PARSE IN VALUE
1. Server IP
2. Server Port
3. DB IP
4. DB Port
5. DB Name
6. DB User
7. DB PW
*/
var tpEnv = `SERVER_IP=%s
SERVER_PORT=%s
DB_IP=%s
DB_PORT=%s
DB_NAME=%s
DB_USER=%s
DB_PW=%s
LOG_LEVEL=DEBUG
JWT_KEY=use-your-own-key
JWT_EXPIRY_MIN=60
REDIS_IP=localhost
REDIS_PORT=6379
REDIS_PW=
REDIS_DB=0
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=from@gmail.com
SMTP_PW=`
