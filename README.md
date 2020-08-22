# Deviser

Deviser is a database and quick basic project setup in GO aimed to effectively create your backend API in minutes.

## Requirements

Golang 1.14.7+
MySQL 8.0+
MS Excel 2007+

## Database Generator

Creates or updates your database with the values in the database.xlsx

### Configurations

Do not remove 'ACCOUNT', 'ACL' and 'AUDIT' from the excel sheet if you wish to use the login and audit feature.

Note!
- All table names should be plural
- All field names should be lowercase
- Please re-initialize the database if there are PRIMARY KEY changes.

*Please note of your database limitation*
```
Length - X or X,Y
16 for INT(16)/VARCHAR(16)
64,2 for DECIMAL(64,2)
```
```
Default - X or FUNC()
XXXX or CURRENT_TIMESTAMP()
```

## Backend API Generator

Generates the baseline GO project with all the basic setup of the database tables in database.xlsx and API services in api.xlsx.
The output will be generated in the *out* folder in the executable's directory.
Existing file will be replaced if the file name has conflicts.

### Configurations

You need to enable Login for JWT and ACL to work.

Note!
- Enabling Login/Audit/ACL/JWT without the initial database.xlsx configurations will cause errors

### Contributors

-   Leong Kai Khee (Kurogami)
