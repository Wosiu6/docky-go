package strategies

import "github.com/wosiu6/docky-go/internal/domain"

type StrategyEntry struct {
	Type     domain.ContainerType
	Strategy ContainerStrategy
}

func Registry() []StrategyEntry {
	return []StrategyEntry{
		{Type: domain.ContainerTypePostgreSQL, Strategy: &PostgreSqlStrategy{}},
		{Type: domain.ContainerTypeMinecraft, Strategy: &MinecraftStrategy{}},
		{Type: domain.ContainerTypePortainer, Strategy: &PortainerStrategy{}},
		{Type: domain.ContainerTypeTraefik, Strategy: &TraefikStrategy{}},
		{Type: domain.ContainerTypeImmich, Strategy: &ImmichStrategy{}},
		{Type: domain.ContainerTypeOwnCloud, Strategy: &OwnCloudStrategy{}},
		{Type: domain.ContainerTypeNginx, Strategy: &NginxStrategy{}},
		{Type: domain.ContainerTypeRedis, Strategy: &RedisStrategy{}},
		{Type: domain.ContainerTypeMySQL, Strategy: &MySQLStrategy{}},
		{Type: domain.ContainerTypeMongoDB, Strategy: &MongoDBStrategy{}},
		{Type: domain.ContainerTypeGrafana, Strategy: &GrafanaStrategy{}},
		{Type: domain.ContainerTypePrometheus, Strategy: &PrometheusStrategy{}},
		{Type: domain.ContainerTypeNextcloud, Strategy: &NextcloudStrategy{}},
		{Type: domain.ContainerTypeMinio, Strategy: &MinioStrategy{}},
		{Type: domain.ContainerTypeMariaDB, Strategy: &MariaDBStrategy{}},
		{Type: domain.ContainerTypeRabbitMQ, Strategy: &RabbitMQStrategy{}},
		{Type: domain.ContainerTypeElasticsearch, Strategy: &ElasticsearchStrategy{}},
		{Type: domain.ContainerTypeKibana, Strategy: &KibanaStrategy{}},
		{Type: domain.ContainerTypeJenkins, Strategy: &JenkinsStrategy{}},
		{Type: domain.ContainerTypeWordPress, Strategy: &WordPressStrategy{}},
		{Type: domain.ContainerTypeVaultwarden, Strategy: &VaultwardenStrategy{}},
		{Type: domain.ContainerTypeMosquitto, Strategy: &MosquittoStrategy{}},
		{Type: domain.ContainerTypePlex, Strategy: &PlexStrategy{}},
		{Type: domain.ContainerTypeJellyfin, Strategy: &JellyfinStrategy{}},
		{Type: domain.ContainerTypeHomeAssistant, Strategy: &HomeAssistantStrategy{}},
		{Type: domain.ContainerTypeSonarr, Strategy: &SonarrStrategy{}},
		{Type: domain.ContainerTypeRadarr, Strategy: &RadarrStrategy{}},
	}
}
