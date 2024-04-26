package ctx

import (
	"context"
	"watchAlert/internal/cache"
	"watchAlert/internal/repo"
)

type Context struct {
	DB    repo.InterEntryRepo
	Redis cache.InterEntryCache
	Ctx   context.Context
}

var (
	DB    repo.InterEntryRepo
	Redis cache.InterEntryCache
	Ctx   context.Context
)

func NewContext(ctx context.Context, db repo.InterEntryRepo, redis cache.InterEntryCache) *Context {
	DB = db
	Redis = redis
	Ctx = ctx
	return &Context{
		DB:    db,
		Redis: redis,
		Ctx:   ctx,
	}
}

func DO() *Context {
	return &Context{
		DB:    DB,
		Redis: Redis,
		Ctx:   Ctx,
	}
}
