package daa

import (
	"strings"
	"testing"

	"github.com/ricardogatica/bugsnag-exporter/pkg/common"
)

func TestParseHeaderNextPrefiereNextAunqueVengaDespuesDePrev(t *testing.T) {
	h := []string{`<https://api.bugsnag.com/projects/p/errors?per_page=30>; rel="prev", <https://api.bugsnag.com/projects/p/errors?offset=abc&per_page=30>; rel="next"`}
	got := parseHeaderNext(h)
	if got.rel != "next" || got.url != "https://api.bugsnag.com/projects/p/errors?offset=abc&per_page=30" {
		t.Fatalf("link inesperado: %+v", got)
	}
}

func TestParseHeaderNextSinNextNoHacePanic(t *testing.T) {
	if got := parseHeaderNext([]string{`<https://x>; rel="prev"`}); got.url != "" {
		t.Fatalf("no debía haber next: %+v", got)
	}
	if got := parseHeaderNext([]string{"basura"}); got.url != "" {
		t.Fatalf("no debía haber next: %+v", got)
	}
	if got := parseHeaderNext(nil); got.url != "" {
		t.Fatalf("no debía haber next: %+v", got)
	}
}

func TestAddQueryParamsArmaFiltros(t *testing.T) {
	common.MinimalReports, common.Since, common.ReleaseStage = true, "7d", "production"
	defer func() { common.MinimalReports, common.Since, common.ReleaseStage = false, "", "" }()

	got := addQueryParams("https://api.bugsnag.com/projects/p/errors", "")
	for _, want := range []string{
		"per_page=100",
		"filters[event.since][][type]=eq&filters[event.since][][value]=7d",
		"filters[app.release_stage][][type]=eq&filters[app.release_stage][][value]=production",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("falta %q en %s", want, got)
		}
	}
	if strings.Contains(got, "full_reports") || strings.Contains(got, "?&") {
		t.Errorf("query mal armada: %s", got)
	}
}

func TestRetryDelaySinHeaderEspera60(t *testing.T) {
	if retryDelay(-1) != 60 || retryDelay(0) != 60 || retryDelay(5) != 5 {
		t.Fatal("retryDelay inesperado")
	}
}
