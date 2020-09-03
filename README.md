# DeviserGo

DeviserGo is a database and quick basic project generator written in Golang aimed to effectively create your backend API in minutes.

## Requirements

- Golang 1.14.7+
- MySQL 8.0+
- MS Excel 2007+

## Modifying the template generator

All source code is in the /src/{project} folder

- You need to configure deviser.conf for each project individually

## How to generate your Golang Project

1. Modify `doc/database.xlsx` to set the configuration and add your own database tables
2. Modify `doc/api.xlsx` to set the configuration and add your own REST API endpoints
3. Run `dbGenerator.exe` to create/update your database with your configuration  
If you already created your database by your own means, please skip this step
4. Run `apiGenerator.exe` to generate your Golang project code  
For step 3. and 4., check that there are no errors in output.log
5. Extract your new project code in `out/{Project}` to your `$GOPATH`
6. Open your project with your favourite IDE and run `go test` and `go run {Project}`

Do you need more information about DeviserGo?  
Click on [DeviserGo](https://github.com/Kurogami88/DeviserGo/wiki) for more details  

### Features

- Login/Logout with JWT Auth Token
- Access Control for API
- Audit logging
- CRUD functions for DB tables

### Future Roadmaps

- Generating Unit Test Function
- Generating API document
- Support More Database
- Sub-directory project structure
- Websocket

### Contributors

-   Leong Kai Khee  
(Kurogami)
