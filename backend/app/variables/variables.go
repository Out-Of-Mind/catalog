package variables

import (
	"github.com/go-redis/redis/v8"
	
	"database/sql"
	"context"
)

var (
	DB *sql.DB
	Cache *redis.Client
	CTX = context.Background()
	TemplateDir = "/catalog/frontend/templates/"
)