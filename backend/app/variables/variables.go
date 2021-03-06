package variables

import (
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"context"
	"database/sql"
)

var (
	DB          *sql.DB
	Cache       *redis.Client
	CTX         = context.Background()
	Log         *logrus.Logger
	Secret      []byte
	TemplateDir string
)
