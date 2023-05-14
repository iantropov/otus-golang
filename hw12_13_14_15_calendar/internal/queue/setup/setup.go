package setup

import (
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/queue"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/queue/rabbit"
)

func Setup(config config.QueueConf) (queue.Connection, error) {
	return rabbit.NewConnection(config)
}
