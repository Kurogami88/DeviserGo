# Deviser

Deviser is a database and quick basic project setup in GO aimed to effectively create your backend API in minutes.

## Requirements

- Golang 1.14.7+
- MySQL 8.0+
- MS Excel 2007+

## Database Generator

Creates or updates your database with the values in the database.xlsx.

Do not remove 'ACCOUNT', 'ACL' and 'AUDIT' from the excel sheet if you wish to use the login and audit feature.

- All field names should be lowercase
- Please re-initialize the database if there are PRIMARY KEY changes.
- All inputs are subjected to your selected database limitation

### Database Excel Field
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

### Configurations

You need to enable Login for JWT and ACL to work.

- Enabling Login/Audit/ACL/JWT without the initial database.xlsx configurations will generate codes that are incomplete
- Existing file will be replaced if the file name has conflicts

### Future Roadmaps

- Websocket
- Unit Test as part of generated code
- API document
- Support More DB
- Maintaining token cache for logout

### Contributors

-   Leong Kai Khee (Kurogami)
