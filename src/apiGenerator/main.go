package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cf, err := os.Open("deviser.conf")
	if err != nil {
		log.Fatalf("Config file error: %v\n", err)
	}
	defer cf.Close()

	// Get database generator config
	generatorConfig := map[string]string{}
	scanner := bufio.NewScanner(cf)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		generatorConfig[line[0]] = line[1]
	}

	lf, err := os.OpenFile(generatorConfig["LOG_PATH"], os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Log file error: %v\n", err)
	}
	defer lf.Close()
	log.SetOutput(lf)

	fDB, err := excelize.OpenFile(generatorConfig["EXCEL_DB_PATH"])
	if err != nil {
		log.Fatalf("Excel DB file error: %v\n", err)
	}

	// Get database excel config
	databaseConfig := map[string]string{}

	databaseConfig["addr"], err = fDB.GetCellValue("config", "B2")
	if err != nil {
		log.Fatalln("DB Config B2 Invalid")
	}
	databaseConfig["port"], err = fDB.GetCellValue("config", "B3")
	if err != nil {
		log.Fatalln("DB Config B3 Invalid")
	}
	databaseConfig["table"], err = fDB.GetCellValue("config", "B4")
	if err != nil {
		log.Fatalln("DB Config B4 Invalid")
	}
	databaseConfig["user"], err = fDB.GetCellValue("config", "B5")
	if err != nil {
		log.Fatalln("DB Config B5 Invalid")
	}
	databaseConfig["pw"], err = fDB.GetCellValue("config", "B6")
	if err != nil {
		log.Fatalln("DB Config B6 Invalid")
	}
	databaseConfig["type"], err = fDB.GetCellValue("config", "B7")
	if err != nil {
		log.Fatalln("DB Config B7 Invalid")
	}
	databaseConfig["init"], err = fDB.GetCellValue("config", "B10")
	if err != nil {
		log.Fatalln("DB Config B10 Invalid")
	}

	log.Println(databaseConfig)

	fAPI, err := excelize.OpenFile(generatorConfig["EXCEL_API_PATH"])
	if err != nil {
		log.Fatalf("Excel API file error: %v\n", err)
	}

	// Get api excel config
	apiConfig := map[string]string{}

	apiConfig["proj"], err = fAPI.GetCellValue("api", "B2")
	if apiConfig["proj"] == "" || err != nil {
		log.Fatalln("API Config B2 Invalid")
	}
	apiConfig["addr"], err = fAPI.GetCellValue("api", "B3")
	if err != nil {
		log.Fatalln("API Config B3 Invalid")
	}
	apiConfig["port"], err = fAPI.GetCellValue("api", "B4")
	if err != nil {
		log.Fatalln("API Config B4 Invalid")
	}
	apiConfig["login"], err = fAPI.GetCellValue("api", "B7")
	if err != nil {
		log.Fatalln("API Config B7 Invalid")
	}
	apiConfig["audit"], err = fAPI.GetCellValue("api", "B8")
	if err != nil {
		log.Fatalln("API Config B8 Invalid")
	}
	apiConfig["ver"], err = fAPI.GetCellValue("api", "B9")
	if err != nil {
		log.Fatalln("API Config B9 Invalid")
	}
	apiConfig["cors"], err = fAPI.GetCellValue("api", "B10")
	if err != nil {
		log.Fatalln("API Config B10 Invalid")
	}

	log.Println(apiConfig)

	// Create project folder structure
	err = os.MkdirAll("out/"+apiConfig["proj"], os.ModePerm)
	if err != nil {
		log.Fatalf("Making project folder error: %v\n", err)
	}

	// Create go mod/sum
	err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/go.mod", []byte(fmt.Sprintf(tpMod, apiConfig["proj"])), os.ModePerm)
	if err != nil {
		log.Fatalf("Making go.mod error: %v\n", err)
	}
	err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/go.sum", []byte(tpSum), os.ModePerm)
	if err != nil {
		log.Fatalf("Making go.sum error: %v\n", err)
	}
	// Logger
	err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/logger.go", []byte(tpLogger), os.ModePerm)
	if err != nil {
		log.Fatalf("Making logger.go error: %v\n", err)
	}
	// Env
	err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/.env", []byte(
		fmt.Sprintf(tpEnv, apiConfig["addr"], apiConfig["port"], databaseConfig["addr"], databaseConfig["port"], databaseConfig["table"], databaseConfig["user"], databaseConfig["pw"])), os.ModePerm)
	if err != nil {
		log.Fatalf("Making .env error: %v\n", err)
	}
	// APIAlive
	err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/apiAlive.go", []byte(apiAlive), os.ModePerm)
	if err != nil {
		log.Fatalf("Making apiAlive.go error: %v\n", err)
	}
	// Login base files
	if apiConfig["login"] == "Yes" {
		err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/auth.go", []byte(tpAuth), os.ModePerm)
		if err != nil {
			log.Fatalf("Making auth.go error: %v\n", err)
		}
		err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/acl.go", []byte(tpAcl), os.ModePerm)
		if err != nil {
			log.Fatalf("Making acl.go error: %v\n", err)
		}
		err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/apiLogin.go", []byte(apiLogin), os.ModePerm)
		if err != nil {
			log.Fatalf("Making apiLogin.go error: %v\n", err)
		}
		err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/apiLogout.go", []byte(apiLogout), os.ModePerm)
		if err != nil {
			log.Fatalf("Making apiLogout.go error: %v\n", err)
		}
		err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/apiPurgeToken.go", []byte(apiPurgeToken), os.ModePerm)
		if err != nil {
			log.Fatalf("Making apiPurgeToken.go error: %v\n", err)
		}
	}
	// Audit base files
	if apiConfig["audit"] == "Yes" {
		err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/audit.go", []byte(tpAudit), os.ModePerm)
		if err != nil {
			log.Fatalf("Making audit.go error: %v\n", err)
		}
	}
	// Main
	mainCORS1 := ""
	mainCORS2 := ""
	if apiConfig["cors"] == "Yes" {
		mainCORS1 = tpMainCORS
		mainCORS2 = tpMainCORSListener
	}
	err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/main.go", []byte(fmt.Sprintf(tpMain, apiConfig["proj"], mainCORS1, mainCORS2)), os.ModePerm)
	if err != nil {
		log.Fatalf("Making main.go error: %v\n", err)
	}

	// Prepare router
	routes := ""
	if apiConfig["login"] == "Yes" {
		routes += tpRouterPurgeToken
		routes += tpRouterLogin
		routes += tpRouterLogout
	}
	// Read all API excel configurations
	rows, err := fAPI.GetRows("api")
	if err != nil {
		log.Fatalln("Read API row error")
	}
	fieldRow := false
	for _, row := range rows {
		var rowValue [4]string
		index := 0

		// Retrieve values from excel
		for _, colCell := range row {
			if !fieldRow {
				if colCell == "Method" {
					fieldRow = true
					rowValue[index] = colCell
				}
				break
			}
			rowValue[index] = colCell
			index++
		}
		if fieldRow && rowValue[0] == "" {
			break
		}
		if rowValue[0] == "" || rowValue[0] == "Method" {
			continue
		}

		apiURL := apiConfig["ver"] + "/" + rowValue[1]
		apiTitle := strings.Title(apiConfig["ver"]) + strings.Title(rowValue[1])
		apiFunc := apiTitle
		apiMethod := rowValue[0]

		if apiConfig["audit"] == "Yes" || rowValue[2] == "Yes" || rowValue[3] == "Yes" {
			apiFunc = fmt.Sprintf(tpRouterHandler, apiFunc)
			if apiConfig["audit"] == "Yes" {
				apiFunc = fmt.Sprintf(tpRouterAudit, apiFunc)
			}
			if rowValue[2] == "Yes" {
				apiFunc = fmt.Sprintf(tpRouterAcl, apiFunc)
			}
			if rowValue[3] == "Yes" {
				apiFunc = fmt.Sprintf(tpRouterAuth, apiFunc)
			}
			routes += fmt.Sprintf(tpRouterMiddleware, apiURL, apiFunc, apiMethod)
		} else {
			routes += fmt.Sprintf(tpRouterSimple, apiURL, apiFunc, apiMethod)
		}

		// API template
		apiFilename := apiConfig["ver"] + strings.Title(rowValue[1])
		err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/"+apiFilename+".go", []byte(
			fmt.Sprintf(apiBase, apiTitle, apiTitle, apiTitle, apiTitle, apiTitle, apiTitle, apiTitle)), os.ModePerm)
		if err != nil {
			log.Fatalf("Making "+apiTitle+".go error: %v\n", err)
		}
	}
	// Router
	err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/router.go", []byte(fmt.Sprintf(tpRouter, routes)), os.ModePerm)
	if err != nil {
		log.Fatalf("Making router.go error: %v\n", err)
	}

	// Read all DB excel configurations
	for _, sheet := range fDB.GetSheetList() {
		if sheet == "config" {
			continue
		}

		rows, err := fDB.GetRows(sheet)
		if err != nil {
			log.Fatalln(sheet + " Error")
		}

		tableName, err := fDB.GetCellValue(sheet, "B2")
		if err != nil {
			log.Fatalln("Config B2 Invalid")
		}
		tableNameTitle := strings.Title(tableName)

		fieldRow := false
		fieldStruct := ""
		for _, row := range rows {
			var rowValue [8]string
			index := 0

			// Retrieve values from excel
			for _, colCell := range row {
				if !fieldRow {
					if colCell == "Field" {
						fieldRow = true
						rowValue[index] = colCell
					}
					break
				}
				rowValue[index] = colCell
				index++
			}
			if fieldRow && rowValue[0] == "" {
				break
			}
			if rowValue[0] == "" || rowValue[0] == "Field" {
				continue
			}

			fieldType := ""
			fieldJSON := fmt.Sprintf(dbTableFieldJSON, rowValue[0])
			fieldGORM := ""
			switch rowValue[1] {
			case "VARCHAR":
				fieldType = "*string"
				fieldGORM = fmt.Sprintf(dbTableFieldGormType, "varchar", "("+rowValue[2]+")")
				if rowValue[3] != "" {
					fieldGORM += " " + dbTableFieldGormDefault + "'" + rowValue[1] + "'"
				}
				if rowValue[4] == "Yes" {
					fieldGORM += " " + dbTableFieldGormKey
				}
			case "INT":
				fieldType = "*int"
				fieldGORM = fmt.Sprintf(dbTableFieldGormType, "int", "("+rowValue[2]+")")
				fieldGORM += " " + dbTableFieldGormAuto + "false"
				if rowValue[3] != "" {
					fieldGORM += " " + dbTableFieldGormDefault + rowValue[1]
				}
				if rowValue[4] == "Yes" {
					fieldGORM += " " + dbTableFieldGormKey
				}
			case "DECIMAL":
				fieldType = "*float64"
				fieldGORM = fmt.Sprintf(dbTableFieldGormType, "decimal", "("+rowValue[2]+")")
				if rowValue[3] != "" {
					fieldGORM += " " + dbTableFieldGormDefault + rowValue[1]
				}
				if rowValue[4] == "Yes" {
					fieldGORM += " " + dbTableFieldGormKey
				}
			case "BLOB":
				fieldType = "*string"
				fieldGORM = fmt.Sprintf(dbTableFieldGormType, "blob", "("+rowValue[2]+")")
				if rowValue[3] != "" {
					fieldGORM += " " + dbTableFieldGormDefault + "'" + rowValue[1] + "'"
				}
			case "BOOLEAN":
				fieldType = "*bool"
				fieldGORM = fmt.Sprintf(dbTableFieldGormType, "boolean", "")
				if rowValue[3] != "" {
					fieldGORM += " " + dbTableFieldGormDefault + rowValue[1]
				}
				if rowValue[4] == "Yes" {
					fieldGORM += " " + dbTableFieldGormKey
				}
			case "DATETIME":
				fieldType = "*time.Time"
				fieldGORM = fmt.Sprintf(dbTableFieldGormType, "boolean", "")
				if rowValue[3] == "CURRENT_TIMESTAMP()" {
					fieldGORM += " " + dbTableFieldGormDefault + dbTableFieldGormCT
				}
				if rowValue[4] == "Yes" {
					fieldGORM += " " + dbTableFieldGormKey
				}
			case "AUTO":
				fieldType = "int"
				fieldGORM = dbTableFieldGormAuto + "true"
				if rowValue[4] == "Yes" {
					fieldGORM += " " + dbTableFieldGormKey
				}
			}

			fieldStruct += fmt.Sprintf(dbTableField, strings.Title(rowValue[0]), fieldType, fieldGORM+" "+fieldJSON)
		}

		fieldStruct += dbFixedFields

		// DB template
		err = ioutil.WriteFile("out/"+apiConfig["proj"]+"/db"+tableNameTitle+".go", []byte(
			fmt.Sprintf(dbBase,
				tableNameTitle, tableName, tableNameTitle, fieldStruct, tableNameTitle, tableNameTitle,
				tableName, tableNameTitle, tableNameTitle, tableName, tableName, tableName, tableNameTitle,
				tableNameTitle, tableNameTitle, tableName, tableNameTitle, tableName, tableName, tableName,
				tableName, tableName, tableNameTitle, tableNameTitle, tableNameTitle, tableName,
				tableNameTitle, tableName, tableName, tableName, tableName, tableName, tableNameTitle,
				tableNameTitle, tableNameTitle, tableName, tableNameTitle, tableName, tableName, tableName,
				tableNameTitle, tableNameTitle, tableName, tableNameTitle, tableNameTitle, tableName,
				tableName, tableName, tableName, tableName, tableNameTitle, tableNameTitle, tableName,
				tableNameTitle, tableNameTitle, tableName, tableName, tableName, tableNameTitle,
				tableNameTitle, tableName, tableNameTitle, tableNameTitle, tableName, tableName,
				tableName, tableNameTitle, tableNameTitle, tableNameTitle)), os.ModePerm)
		if err != nil {
			log.Fatalf("Making db"+tableNameTitle+".go error: %v\n", err)
		}
	}

	log.Println("Generate API successfully")
}
