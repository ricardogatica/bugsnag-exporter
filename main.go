package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ricardogatica/bugsnag-exporter/pkg/common"
	"github.com/ricardogatica/bugsnag-exporter/pkg/daa"
)

func main() {
	// Sin bugsnag.Configure: la versión original reportaba sesiones/errores de
	// esta herramienta al proyecto Bugsnag de su autor (API key hardcodeada).

	// Carga BUGSNAG_TOKEN (y otras variables) desde ./.env si existe. Las
	// variables ya definidas en el entorno tienen prioridad.
	common.LoadDotEnv(".env")

	token := flag.String("token", "", "[String] Your Bugsnag personal auth token. Prefer the BUGSNAG_TOKEN env var (or .env): a flag stays in your shell history")
	getProjectIds := flag.Bool("show-project-ids", false, "[Flag] List the project IDs accessible with your token (also the default when --project-id is missing)")
	projectId := flag.String("project-id", "", "[String] The Project ID you wish to download from")
	errorId := flag.String("error-id", "", "[String] An error ID to download. If provided, downloads all events within filters for this error ID")
	events := flag.Bool("events", false, "[Flag] Download events rather than error groups when this flag is enabled. Requires --project-id (and optionally --error-id)")
	//affectedUsers := flag.Bool("affected-users", false, "[Flag] Download a list of users affected by a specific error. Requires --project-id and --error-id")
	outputDir := flag.String("output-dir", "./data", "[String] Directory to store the downloaded file")
	filters := flag.String("filters", "", "[String] A string array of filters to apply (URL format)")
	//rateLimit := flag.Int("rate-limit", 0, "[Int] Set the number of calls to make per minute. Defaults to 0, no rate limit")
	minimalReports := flag.Bool("minimal", false, "[Flag] Download minimal event reports only for smaller file sizes")
	useCsv := flag.Bool("csv", false, "[Flag] Output data to file as CSV. Default false, noramally outputs as JSON")
	noWarn := flag.Bool("no-warn", false, "[Flag] Don't warn me if this call will take more than 5 calls to the API")
	verbose := flag.Bool("verbose", false, "[Flag] Set the output to be verbose for debugging purposes.")
	since := flag.String("since", "", "[String] Only events since this time window or date, e.g. 7d, 24h, 2026-11-01T00:00:00Z (filter event.since)")
	releaseStage := flag.String("release-stage", "", "[String] Only this release stage, e.g. production (filter app.release_stage)")
	flag.Parse()

	common.Verbose = *verbose
	common.NoWarn = *noWarn
	common.UseCsv = *useCsv
	common.OutputDir = *outputDir
	common.MinimalReports = *minimalReports
	common.Since = *since
	common.ReleaseStage = *releaseStage

	// Preferir la variable de entorno: un --token queda en el historial del shell.
	if *token == "" {
		*token = os.Getenv("BUGSNAG_TOKEN")
	}

	if *token == "" {
		common.ExitWithString(1, "Missing token. Use BUGSNAG_TOKEN or the --token flag")
	} else {
		daa.PersonalAuthToken = *token
		if *getProjectIds || *projectId == "" {
			common.Print("Getting your project IDs with provided token")
			orgs := daa.GetUsersOrganizations()
			for _, org := range orgs {
				common.Print("Organization: " + fmt.Sprint(org["name"]) + " [" + fmt.Sprint(org["id"]) + "]")
				projects := daa.GetOrganizationsProjects(org["id"].(string))
				for _, proj := range projects {
					common.Print("  > " + fmt.Sprint(proj["name"]) + " [" + fmt.Sprint(proj["id"]) + "]")
				}
			}
		} else {
			if *projectId == "" {
				common.ExitWithString(2, "Missing Project ID. Please supply a Project ID with --project-id flag")
			} else {
				if *errorId == "" {
					if *events {
						common.Print("Downloading all events for projectId within filters")
						daa.GetProjectEvents(*projectId, *filters)
					} else {
						common.Print("Downloading all errors from projectId within filters")
						daa.GetProjectErrors(*projectId, *filters)
					}
				} else {
					if *events {
						common.Print("Downloading all events for projectId & errorId within filters")
						daa.GetErrorEvents(*projectId, *errorId, *filters)
					} else {
						common.Print("Downloading error from projectId & errorId within filters")
						daa.GetError(*projectId, *errorId, *filters)
					}
				}
			}
		}
	}
}
