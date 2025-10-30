package strategies

import (
	"context"

	"github.com/wosiu6/docky-go/internal/model"
)

type ContainerStrategy interface {
	Match(image string) bool
	Extract(ctx context.Context, id string, raw map[string]interface{}, base model.BaseContainerInfo, client interface{}) interface{}
}
