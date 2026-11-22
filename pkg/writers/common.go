package writers

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ricardogatica/bugsnag-exporter/pkg/common"
)

func CreateNewOutputFile(projectId string, filename string) *os.File {
	// TODO: Warn user if the file is going to overwrite an old file
	var outputDirpath string = filepath.Join(common.OutputDir, projectId)
	common.PrintVerbose("Setting up directory '" + outputDirpath + "' to download data to")
	// TODO: Check MkdirAll is not destructive where not wanted!
	err := os.MkdirAll(outputDirpath, os.ModePerm)
	if err != nil {
		common.ExitWithError(1, err)
	}
	var handle *os.File
	var outputFilepath string
	// Sufijo con fecha y hora: cada descarga queda como un archivo nuevo y no
	// pisa la anterior, así se pueden comparar exportaciones en el tiempo.
	filename += "-" + time.Now().Format("20060102-150405")
	if common.UseCsv {
		outputFilepath = filepath.Join(common.OutputDir, projectId, filename+".csv")
	} else {
		outputFilepath = filepath.Join(common.OutputDir, projectId, filename+".json")
	}
	common.PrintVerbose("Creating file " + outputFilepath)
	handle, err = os.Create(outputFilepath)
	if err != nil {
		common.ExitWithError(1, err)
	} else {
		return handle
	}
	return nil
}

func WriteObjectToFile(handle *os.File, element map[string]interface{}) {
	common.PrintVerbose("Writing object to output file '" + handle.Name() + "'")
	if common.UseCsv {
		writeToCsvFile(handle, []map[string]interface{}{element})
	} else {
		writeObjectToJsonFile(handle, element)
	}
}

func WriteArrayToFile(handle *os.File, elements []map[string]interface{}) {
	common.PrintVerbose("Writing array to output file '" + handle.Name() + "'")
	if common.UseCsv {
		writeToCsvFile(handle, elements)
	} else {
		writeArrayToJsonFile(handle, elements)
	}
}

func CloseOutputFile(handle *os.File) {
	common.PrintVerbose("Closing output file '" + handle.Name() + "'")
	handle.Close()
}
