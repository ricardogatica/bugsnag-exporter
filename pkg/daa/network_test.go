package daa

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ricardogatica/bugsnag-exporter/pkg/common"
)

// Recorre 2 páginas con un 429 en medio: debe reintentar, seguir rel="next"
// aunque venga después de prev, y mandar los headers de autenticación/versión.
func TestBugsnagGetArrayPaginaYReintenta429(t *testing.T) {
	common.NoWarn = true
	PersonalAuthToken = "tok"
	calls := 0
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Header.Get("Authorization") != "token tok" || r.Header.Get("X-Version") != "2" {
			t.Errorf("headers: %v", r.Header)
		}
		w.Header().Set("X-Total-Count", "2")
		switch {
		case r.URL.Query().Get("page") == "":
			w.Header().Set("Link", fmt.Sprintf(`<%s/x?prev=1>; rel="prev", <%s/x?page=2>; rel="next"`, srv.URL, srv.URL))
			fmt.Fprint(w, `[{"id":"a"}]`)
		case calls == 2:
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
		default:
			fmt.Fprint(w, `[{"id":"b"}]`)
		}
	}))
	defer srv.Close()

	got := BugsnagGetArray(srv.URL + "/x")
	if len(got) != 2 || got[0]["id"] != "a" || got[1]["id"] != "b" {
		t.Fatalf("elementos: %v", got)
	}
	if calls != 3 {
		t.Fatalf("llamadas: %d, esperado 3 (página 1, 429, página 2)", calls)
	}
}
