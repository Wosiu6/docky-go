package domain

import "time"

type ContainerType string

const (
	ContainerTypeGeneric       ContainerType = "generic"
	ContainerTypePostgreSQL    ContainerType = "postgresql"
	ContainerTypeMinecraft     ContainerType = "minecraft"
	ContainerTypePortainer     ContainerType = "portainer"
	ContainerTypeTraefik       ContainerType = "traefik"
	ContainerTypeImmich        ContainerType = "immich"
	ContainerTypeOwnCloud      ContainerType = "owncloud"
	ContainerTypeNginx         ContainerType = "nginx"
	ContainerTypeRedis         ContainerType = "redis"
	ContainerTypeMySQL         ContainerType = "mysql"
	ContainerTypeMongoDB       ContainerType = "mongodb"
	ContainerTypeGrafana       ContainerType = "grafana"
	ContainerTypePrometheus    ContainerType = "prometheus"
	ContainerTypeNextcloud     ContainerType = "nextcloud"
	ContainerTypeMinio         ContainerType = "minio"
	ContainerTypeMariaDB       ContainerType = "mariadb"
	ContainerTypeRabbitMQ      ContainerType = "rabbitmq"
	ContainerTypeElasticsearch ContainerType = "elasticsearch"
	ContainerTypeKibana        ContainerType = "kibana"
	ContainerTypeJenkins       ContainerType = "jenkins"
	ContainerTypeWordPress     ContainerType = "wordpress"
	ContainerTypeVaultwarden   ContainerType = "vaultwarden"
	ContainerTypeMosquitto     ContainerType = "mosquitto"
	ContainerTypePlex          ContainerType = "plex"
	ContainerTypeJellyfin      ContainerType = "jellyfin"
	ContainerTypeHomeAssistant ContainerType = "homeassistant"
	ContainerTypeSonarr        ContainerType = "sonarr"
	ContainerTypeRadarr        ContainerType = "radarr"
)

type Stats struct {
	CPUTotal   uint64
	SystemCPU  uint64
	OnlineCPUs uint64
	Collected  time.Time
}

type DetailProvider interface {
	DetailFields() map[string]string
}

type Container struct {
	ID         string
	Names      []string
	Image      string
	Status     string
	CPUPercent float64
	MemoryMB   uint64
	Type       ContainerType
	Details    DetailProvider
}
