package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvNoPisaElEntorno(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	os.WriteFile(path, []byte("# comentario\nexport BE_TEST_A=\"uno\"\nBE_TEST_B='dos'\nBE_TEST_C=tres\nmal\n"), 0600)
	t.Setenv("BE_TEST_C", "entorno")
	os.Unsetenv("BE_TEST_A")
	os.Unsetenv("BE_TEST_B")
	defer os.Unsetenv("BE_TEST_A")
	defer os.Unsetenv("BE_TEST_B")

	LoadDotEnv(path)

	if os.Getenv("BE_TEST_A") != "uno" || os.Getenv("BE_TEST_B") != "dos" {
		t.Fatalf("valores no cargados: %q %q", os.Getenv("BE_TEST_A"), os.Getenv("BE_TEST_B"))
	}
	if os.Getenv("BE_TEST_C") != "entorno" {
		t.Fatal("pisó una variable ya definida en el entorno")
	}
}

func TestLoadDotEnvSinArchivoNoFalla(t *testing.T) {
	LoadDotEnv(filepath.Join(t.TempDir(), "no-existe"))
}
