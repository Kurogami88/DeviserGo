package main

var tpLogger = `/***
	Author: Leong Kai Khee (Kurogami)
	Date: 2020

	Generated by DeviserGO
***/

package main

import (
	"log"
	"os"
)

// Log printing function
func Log(s string) {
	log.Println(s)
}

// LogDebug for debug log entry
func LogDebug(s string) {
	envLogLevel := os.Getenv("LOG_LEVEL")
	if envLogLevel == "DEBUG" {
		Log("[DEBUG] " + s)
	}
}

// LogInfo for info log entry
func LogInfo(s string) {
	envLogLevel := os.Getenv("LOG_LEVEL")
	if envLogLevel == "DEBUG" || envLogLevel == "INFO" {
		Log("[INFO] " + s)
	}
}

// LogWarning for warning log entry
func LogWarning(s string) {
	envLogLevel := os.Getenv("LOG_LEVEL")
	if envLogLevel == "DEBUG" || envLogLevel == "INFO" || envLogLevel == "WARNING" {
		Log("[WARNING] " + s)
	}
}

// LogError for error log entry
func LogError(s string) {
	envLogLevel := os.Getenv("LOG_LEVEL")
	if envLogLevel == "DEBUG" || envLogLevel == "INFO" || envLogLevel == "WARNING" || envLogLevel == "ERROR" {
		Log("[ERROR] " + s)
	}
}
`
