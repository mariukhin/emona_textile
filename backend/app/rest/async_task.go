package api

import (
	"amifactory.team/sequel/coton-app-backend/app/logger"
	"context"
	"time"
)

const (
	sendEmailTaskTimeout = time.Second * 15
)

func NewAsyncTaskContext(ctx context.Context) context.Context {
	log := logger.GetLogger(ctx)
	return logger.WithLogger(context.Background(), log)
}
