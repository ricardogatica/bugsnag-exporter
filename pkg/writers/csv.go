package writers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/ricardogatica/bugsnag-exporter/pkg/common"
)

/*
Escribe los elementos como CSV: una fila por elemento y una columna por cada
llave de primer nivel (la unión de todas, ordenadas). Los valores anidados
(objetos y arreglos) se guardan como JSON dentro de la celda.
*/
func writeToCsvFile(handle *os.File, elements []map[string]interface{}) {
	keySet := map[string]bool{}
	for _, element := range elements {
		for k := range element {
			keySet[k] = true
		}
	}
	headers := make([]string, 0, len(keySet))
	for k := range keySet {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	w := csv.NewWriter(handle)
	if err := w.Write(headers); err != nil {
		common.ExitWithError(1, err)
	}
	for _, element := range elements {
		row := make([]string, len(headers))
		for i, k := range headers {
			row[i] = csvCell(element[k])
		}
		if err := w.Write(row); err != nil {
			common.ExitWithError(1, err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		common.ExitWithError(1, err)
	}
}

func csvCell(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case map[string]interface{}, []interface{}:
		b, err := json.Marshal(t)
		if err != nil {
			return fmt.Sprint(t)
		}
		return string(b)
	default:
		return fmt.Sprint(t)
	}
}
