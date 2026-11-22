package common

import (
	"fmt"
	"os"
	"strings"
)

/*
Variable to control output verbosity.
Defaults to false
*/
var Verbose bool = false

/*
Variable to control whether a warning is given
if more than 10 API calls will be required to complete
an operation. 10 API calls would encounter rate limiting
under normal conditions.
*/
var NoWarn bool = false

/*
Variable to control whether we output as a JSON file
or as a CSV. JSON is preferable as data in Bugsnag events is
complex and does not suit CSVs well.
*/
var UseCsv bool = false

/*
Variable to store the output directory to store the files created
*/
var OutputDir string = "./"

/*
Variable to control whether we should be fetching minimal reports
*/
var MinimalReports bool = false

/*
Ventana de tiempo para filtrar eventos (p. ej. "7d", "24h" o una fecha
ISO 8601). Se traduce al filtro event.since de la API.
*/
var Since string = ""

/*
Release stage para filtrar (p. ej. "production"). Se traduce al filtro
app.release_stage de la API.
*/
var ReleaseStage string = ""

/*
Exits the application with the error code given after
printing an error message string to the console.
*/
func ExitWithString(errorCode int, errorMessage string) {
	Print("[ERROR] %s\n", errorMessage)
	os.Exit(errorCode)
}

/*
Exits the application with the error code given after
printing an error message from an error object to console.
*/
func ExitWithError(errorCode int, err error) {
	Print("[ERROR] %s\n", err.Error())
	os.Exit(errorCode)
}

/*
Exits the application with the error code given after
printing a string and message from an error object to console.
*/
func ExitWithErrorAndString(errorCode int, err error, errorMessage string) {
	Print("[ERROR] %s\n", errorMessage)
	Print("[ERROR] %s\n", err.Error())
	os.Exit(errorCode)
}

/*
Prints a message to the console, can take in arguments
after the initial message so you can treat this as a Sprintf()
*/
func Print(format string, args ...interface{}) {
	fmt.Printf("[bugsnag-exporter] "+format+"\r\n", args...)
}

/*
Prints a message to the console, in the same way as Print()
but ONLY if `Verbose` is set to true
*/
func PrintVerbose(format string, args ...interface{}) {
	if Verbose {
		fmt.Printf("[b-e-verbose-----] "+format+"\r\n", args...)
	}
}

/*
Prints a header to the console, including the application version
*/
func PrintHeader() {
	if Verbose {
		Print("####################################################")
		Print("#                                                  #")
		Print("#              bugsnag-exporter, v%s            #", PackageVersion)
		Print("#                Xander Jones, 2021                #")
		Print("# https://github.com/ricardogatica/bugsnag-exporter #")
		Print("#                                                  #")
		Print("####################################################")
	}
}

/*
Carga variables KEY=VALUE desde un archivo .env (si existe). Ignora líneas
vacías y comentarios (#), acepta "export KEY=VALUE" y comillas simples o
dobles alrededor del valor. No pisa variables ya definidas en el entorno.
*/
func LoadDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		eq := strings.Index(line, "=")
		if eq <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		value := strings.TrimSpace(line[eq+1:])
		if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') && value[len(value)-1] == value[0] {
			value = value[1 : len(value)-1]
		}
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}
