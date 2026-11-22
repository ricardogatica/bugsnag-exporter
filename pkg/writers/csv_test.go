package writers

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteToCsvFileUneColumnasYSerializaAnidados(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.csv")
	f, _ := os.Create(path)
	writeToCsvFile(f, []map[string]interface{}{
		{"id": "1", "meta": map[string]interface{}{"a": 1.0}},
		{"id": "2", "count": 3.0},
	})
	f.Close()

	f, _ = os.Open(path)
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	want := [][]string{{"count", "id", "meta"}, {"", "1", `{"a":1}`}, {"3", "2", ""}}
	if len(rows) != len(want) {
		t.Fatalf("filas: %v", rows)
	}
	for i := range want {
		for j := range want[i] {
			if rows[i][j] != want[i][j] {
				t.Fatalf("fila %d col %d: %q, esperado %q", i, j, rows[i][j], want[i][j])
			}
		}
	}
}
