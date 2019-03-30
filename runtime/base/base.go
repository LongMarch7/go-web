package base

import (
    "google.golang.org/grpc"
    "context"
)

type BaseGatewayManager struct {
    Handler func(rctx context.Context, conn *grpc.ClientConn) error
    Manager interface{}
}

func NewManager(manager interface{}) *BaseGatewayManager {
    return &BaseGatewayManager{
        Manager: manager,
    }
}
