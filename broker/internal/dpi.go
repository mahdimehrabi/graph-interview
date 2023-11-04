package internal

import (
	"github.com/mahdimehrabi/graph-interview/broker/external/broker"
	infrastructures "github.com/mahdimehrabi/graph-interview/broker/internal/infrastructure"
)

// singleton dependency injection
var DPI *dpi

type dpi struct {
	BrokerSocket *broker.Socket
}

func SetupDPI(env *infrastructures.Env) {
	DPI = &dpi{
		BrokerSocket: broker.NewSocket(env.ServerPort),
	}
}
