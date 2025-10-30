package strategies

import "testing"

var sampleImages = map[string]string{
	"postgres":                     "postgresql",
	"redis:7":                      "redis",
	"mysql:8":                      "mysql",
	"mongo:6":                      "mongodb",
	"grafana/grafana:latest":       "grafana",
	"prom/prometheus":              "prometheus",
	"nextcloud:stable":             "nextcloud",
	"minio/minio:latest":           "minio",
	"mariadb:10":                   "mariadb",
	"rabbitmq:3-management":        "rabbitmq",
	"elasticsearch:8":              "elasticsearch",
	"kibana:8":                     "kibana",
	"jenkins/jenkins:lts":          "jenkins",
	"wordpress:php8":               "wordpress",
	"vaultwarden/server:latest":    "vaultwarden",
	"eclipse-mosquitto:latest":     "mosquitto",
	"plexinc/pms-docker":           "plex",
	"jellyfin/jellyfin":            "jellyfin",
	"homeassistant/home-assistant": "homeassistant",
	"sonarr:latest":                "sonarr",
	"radarr:latest":                "radarr",
	"traefik:v2":                   "traefik",
	"owncloud/server":              "owncloud",
	"immich-server:latest":         "immich",
}

func TestStrategyMatchSamples(t *testing.T) {
	entries := Registry()
	for img, expectedType := range sampleImages {
		matched := false
		for _, e := range entries {
			if e.Strategy.Match(img) {
				if string(e.Type) != expectedType {
					t.Errorf("image %s matched wrong type %s want %s", img, e.Type, expectedType)
				}
				matched = true
				break
			}
		}
		if !matched {
			if expectedType != "" {
				t.Errorf("image %s expected type %s but no strategy matched", img, expectedType)
			}
		}
	}
}
