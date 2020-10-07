package main

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"strings"
	"unicode"

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

	f, err := excelize.OpenFile(generatorConfig["EXCEL_DB_PATH"])
	if err != nil {
		log.Fatalf("Excel DB file error: %v\n", err)
	}

	// Get database excel config
	databaseConfig := map[string]string{}

	databaseConfig["addr"], err = f.GetCellValue("config", "B2")
	if err != nil {
		log.Fatalln("DB Config B2 Invalid")
	}
	databaseConfig["port"], err = f.GetCellValue("config", "B3")
	if err != nil {
		log.Fatalln("DB Config B3 Invalid")
	}
	databaseConfig["table"], err = f.GetCellValue("config", "B4")
	if err != nil {
		log.Fatalln("DB Config B4 Invalid")
	}
	databaseConfig["user"], err = f.GetCellValue("config", "B5")
	if err != nil {
		log.Fatalln("DB Config B5 Invalid")
	}
	databaseConfig["pw"], err = f.GetCellValue("config", "B6")
	if err != nil {
		log.Fatalln("DB Config B6 Invalid")
	}
	databaseConfig["type"], err = f.GetCellValue("config", "B7")
	if err != nil {
		log.Fatalln("DB Config B7 Invalid")
	}
	databaseConfig["init"], err = f.GetCellValue("config", "B10")
	if err != nil {
		log.Fatalln("DB Config B10 Invalid")
	}

	log.Println(databaseConfig)

	db, err := sql.Open(strings.ToLower(databaseConfig["type"]), databaseConfig["user"]+":"+
		databaseConfig["pw"]+"@tcp("+databaseConfig["addr"]+":"+
		databaseConfig["port"]+")/")
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}

	// Drop and create database
	if databaseConfig["init"] == "Yes" {
		_, err = db.Exec("DROP DATABASE IF EXISTS " + databaseConfig["table"])
		if err != nil {
			log.Fatalln(err)
		}
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + databaseConfig["table"])
	if err != nil {
		log.Fatalln(err)
	}

	// Loop through excel sheet
	for _, sheet := range f.GetSheetList() {
		if sheet == "config" {
			continue
		}

		rows, err := f.GetRows(sheet)
		if err != nil {
			log.Fatalln(sheet + " error")
		}

		tableName := ""
		preTableName, err := f.GetCellValue(sheet, "B2")
		if err != nil {
			log.Fatalln("Config B2 Invalid")
		}

		for _, char := range preTableName {
			if unicode.IsUpper(char) && unicode.IsLetter(char) {
				tableName = tableName + "_" + strings.ToLower(string(char))
			} else {
				tableName = tableName + string(char)
			}
		}

		res, _ := db.Query("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name='" + tableName + "' AND table_schema='" + databaseConfig["table"] + "'")

		// Check if table already exist
		var column string
		res.Next()
		res.Scan(&column)
		if column == "" {
			// Do creation
			_, err = db.Exec("USE " + databaseConfig["table"])
			if err != nil {
				log.Fatalln(err)
			}

			tableStatement := "CREATE TABLE IF NOT EXISTS `" + tableName + "` ( "
			keyStatement := ""
			indexStatement := ""
			fieldStart := 4
			for _, row := range rows {
				if fieldStart > 0 {
					fieldStart--
					continue
				}

				var rowValue [8]string
				index := 0

				// Retrieve values from excel
				for _, colCell := range row {
					rowValue[index] = colCell
					index++
				}
				// Stop processing if empty field found
				if rowValue[0] == "" {
					break
				}

				// Field and field type
				tableStatement += "`" + strings.ToLower(rowValue[0]) + "` "
				if rowValue[1] == "AUTO" {
					tableStatement += "INT AUTO_INCREMENT "
				} else if rowValue[1] == "BOOLEAN" || rowValue[1] == "DATETIME" {
					tableStatement += rowValue[1] + " "
				} else {
					tableStatement += rowValue[1] + "(" + rowValue[2] + ") "
				}

				// Unsigned
				if rowValue[1] == "INT" && rowValue[7] == "Yes" {
					tableStatement += "UNSIGNED "
				}

				// Nullable
				if rowValue[6] == "No" {
					tableStatement += "NOT NULL "
				}

				// Default value
				if rowValue[3] != "" {
					tableStatement += "DEFAULT "
					if rowValue[1] == "VARCHAR" {
						tableStatement += " '" + rowValue[3] + "', "
					} else {
						tableStatement += rowValue[3] + ", "
					}
				} else {
					tableStatement += ", "
				}

				// Primary key
				if rowValue[4] == "Yes" {
					if keyStatement != "" {
						keyStatement += ", "
					}
					keyStatement += "`" + strings.ToLower(rowValue[0]) + "`"
				}

				// Index
				if rowValue[5] == "Yes" && rowValue[4] == "No" {
					if indexStatement != "" {
						indexStatement += ", "
					}
					indexStatement += "`" + strings.ToLower(rowValue[0]) + "`"
				}
			}
			// CreatedAt, UpdatedAt, DeletedAt
			tableStatement += "`created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(), "
			tableStatement += "`updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(), "
			tableStatement += "`deleted_at` DATETIME DEFAULT NULL, "

			// Append index
			if indexStatement != "" {
				tableStatement += "INDEX `idx` (" + indexStatement + "), "
			}
			// Append primary key
			if keyStatement != "" {
				tableStatement += "PRIMARY KEY (" + keyStatement + ")"
			}

			tableStatement += ")"
			log.Println(tableStatement)
			_, err = db.Exec(tableStatement)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			// Do alteration
			// Prepare process map
			tableSchema := map[string]bool{}
			tableSchema[column] = false
			for res.Next() {
				res.Scan(&column)
				tableSchema[column] = false
			}

			log.Println("Table " + tableName)

			keyStatement := ""
			indexStatement := ""
			fieldRow := false
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

				_, err = db.Exec("USE " + databaseConfig["table"])
				if err != nil {
					log.Fatalln(err)
				}

				tableStatement := "ALTER TABLE " + tableName + " "
				if _, ok := tableSchema[strings.ToLower(rowValue[0])]; ok {
					tableSchema[strings.ToLower(rowValue[0])] = true
					tableStatement += "MODIFY COLUMN "
				} else {
					tableStatement += "ADD COLUMN "
				}

				// Field and field type
				tableStatement += "`" + strings.ToLower(rowValue[0]) + "` "
				if rowValue[1] == "AUTO" {
					tableStatement += "INT AUTO_INCREMENT "
				} else if rowValue[1] == "BOOLEAN" || rowValue[1] == "DATETIME" {
					tableStatement += rowValue[1] + " "
				} else {
					tableStatement += rowValue[1] + "(" + rowValue[2] + ") "
				}

				// Unsigned
				if rowValue[1] == "INT" && rowValue[7] == "Yes" {
					tableStatement += "UNSIGNED "
				}

				// Nullable
				if rowValue[6] == "No" {
					tableStatement += "NOT NULL "
				}

				// Default value
				if rowValue[3] != "" {
					tableStatement += "DEFAULT "
					if rowValue[1] == "VARCHAR" {
						tableStatement += " '" + rowValue[3] + "'"
					} else {
						tableStatement += rowValue[3]
					}
				}

				// Alter row
				log.Println(tableStatement)
				_, err = db.Exec(tableStatement)
				if err != nil {
					log.Fatalln(err)
				}

				// Primary key
				if rowValue[4] == "Yes" {
					if keyStatement != "" {
						keyStatement += ", "
					}
					keyStatement += "`" + strings.ToLower(rowValue[0]) + "`"
				}

				// Index
				if rowValue[5] == "Yes" && rowValue[4] == "No" {
					if indexStatement != "" {
						indexStatement += ", "
					}
					indexStatement += "`" + strings.ToLower(rowValue[0]) + "`"
				}
			}

			// Alter primary key and index
			tableStatement := "ALTER TABLE " + tableName + " DROP PRIMARY KEY, ADD PRIMARY KEY (" + keyStatement + ")"
			log.Println(tableStatement)
			_, err = db.Exec(tableStatement)
			if err != nil {
				log.Fatalln(err)
			}

			if indexStatement != "" {
				// Attempt to drop index but ignore drop error
				tableStatement = "ALTER TABLE " + tableName + " DROP INDEX `idx`"
				log.Println(tableStatement)
				db.Exec(tableStatement)

				tableStatement = "ALTER TABLE " + tableName + " ADD INDEX `idx` (" + indexStatement + ")"
				log.Println(tableStatement)
				_, err = db.Exec(tableStatement)
				if err != nil {
					log.Fatalln(err)
				}
			} else {
				// Attempt to drop index but ignore drop error
				tableStatement = "ALTER TABLE " + tableName + " DROP INDEX `idx`"
				log.Println(tableStatement)
				db.Exec(tableStatement)
			}

			// Drop extra columns
			for key, value := range tableSchema {
				if key == "created_at" || key == "updated_at" || key == "deleted_at" {
					continue
				}

				if !value {
					tableStatement = "ALTER TABLE " + tableName + " DROP COLUMN " + key
					log.Println(tableStatement)
					_, err = db.Exec(tableStatement)
					if err != nil {
						log.Fatalln(err)
					}
				}
			}
		}
	}

	log.Println("Generate DB successfully")
}
