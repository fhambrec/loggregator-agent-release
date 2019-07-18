package main

import (
	"code.cloudfoundry.org/go-loggregator/metrics"
	"log"
	"os"

	"code.cloudfoundry.org/loggregator-agent/cmd/metric-scraper/app"
)

func main() {
	log := log.New(os.Stderr, "", log.LstdFlags)
	log.Printf("starting Metrics Scraper...")
	defer log.Printf("closing Metrics Scraper...")

	cfg := app.LoadConfig(log)

	dt := map[string]string{
		"metrics_version": "2.0",
		"origin": "loggregator_metrics_scraper",
		"source_id": "metrics_scraper",
	}

	metricClient := metrics.NewRegistry(
		log,
		metrics.WithDefaultTags(dt),
		metrics.WithServer(cfg.DebugPort),
	)

	app.NewMetricScraper(cfg, log, metricClient).Run()
}
