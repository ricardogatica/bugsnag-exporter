package daa

import (
	"bufio"
	"fmt"
	neturl "net/url"
	"os"
	"regexp"
	"strings"

	"github.com/ricardogatica/bugsnag-exporter/pkg/common"
)

func addQueryParams(url string, filters string) string {
	params := []string{"per_page=100"} // máximo de la API: menos llamadas y menos rate limit
	if !common.MinimalReports {
		params = append(params, "full_reports=true")
	}
	if common.Since != "" {
		params = append(params, eqFilter("event.since", common.Since))
	}
	if common.ReleaseStage != "" {
		params = append(params, eqFilter("app.release_stage", common.ReleaseStage))
	}
	if filters != "" {
		if f := formatFilters(filters); f != "BAD_FILTERS" {
			params = append(params, strings.TrimLeft(f, "?&"))
		} else {
			common.ExitWithString(1, "Filters provided were not valid")
		}
	}
	return url + "?" + strings.Join(params, "&")
}

/*
Arma un filtro de igualdad en el formato de la Data Access API:
filters[campo][][type]=eq&filters[campo][][value]=valor
*/
func eqFilter(field string, value string) string {
	return "filters[" + field + "][][type]=eq&filters[" + field + "][][value]=" + neturl.QueryEscape(value)
}

func formatFilters(filters string) string {
	r := regexp.MustCompile(`\??(filters(\[[^]]*\])+=([^&]*))&?`)
	matches := r.FindAllString(filters, -1)
	common.PrintVerbose("filters count: " + fmt.Sprint(len(matches)))
	common.PrintVerbose("filters:       " + fmt.Sprint(matches))
	if matches == nil {
		return "BAD_FILTERS"
	} else {
		return filters
	}
}

/*
Waits for a user to confirm a message. If they reject the message, the
application will exit with a message.

With thanks to @r0l1: https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b
*/
func getConfirmation(confirmationMessage string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", confirmationMessage)
		response, err := reader.ReadString('\n')
		if err != nil {
			common.ExitWithError(1, err)
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			common.ExitWithString(1, "User rejected further API calls")
		}
	}
}
