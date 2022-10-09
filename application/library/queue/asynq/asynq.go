package asynq

import (
	"github.com/hibiken/asynq"
)

func New(redisOptions asynq.RedisConnOpt) (*Asynq, error) {
	a := &Asynq{}
	err := a.SetRedisOptions(redisOptions)
	return a, err
}

type Asynq struct {
	// redis
	redisOptions asynq.RedisConnOpt

	// for client
	clientInstance *asynq.Client

	// for server
	serverInstance *asynq.Server
	serverHandler  *asynq.ServeMux
	serverConfig   *asynq.Config
}

func ParseRedisURI(connURI string) (asynq.RedisConnOpt, error) {
	return asynq.ParseRedisURI(connURI)
}

func (a *Asynq) SetRedisOptions(redisOptions asynq.RedisConnOpt) error {
	if redisOptions == nil {
		redisOptions = &asynq.RedisClientOpt{
			Addr: "127.0.0.1:6379",
		}
	}
	a.redisOptions = redisOptions
	a.Close()
	if a.clientInstance != nil {
		a.clientInstance = nil
	}
	if a.serverInstance != nil {
		a.serverInstance = nil
	}
	return nil
}

func (a *Asynq) Client() *asynq.Client {
	if a.clientInstance == nil {
		a.clientInstance = asynq.NewClient(a.redisOptions)
	}
	return a.clientInstance
}

func (a *Asynq) Server() *asynq.Server {
	if a.serverInstance == nil {
		a.serverInstance = a.newServer()
	}
	return a.serverInstance
}

func (a *Asynq) Close() error {
	if a.serverInstance != nil {
		a.serverInstance.Stop()
	}
	if a.clientInstance != nil {
		return a.clientInstance.Close()
	}
	return nil
}
