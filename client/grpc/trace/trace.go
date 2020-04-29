package trace

import "github.com/micro/go-plugins/client/grpc/internal/tracer"

func init() {
	tracer.TurnOn()
}
