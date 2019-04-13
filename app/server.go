package app

import (
    "github.com/LongMarch7/go-web/runtime/consul-server"
    "github.com/LongMarch7/go-web/view"
)

type MicroServer struct{
    Server server.IServer
    Template view.View
}
var ServerObj MicroServer
func  NewMicroServer(opts ...server.SOption) *MicroServer{
    ServerObj.Server = server.NewServer(opts...)
    ServerObj.Template = view.Template
    return &ServerObj
}

func (s *MicroServer)AddPlugin(plugin interface{}){
    ServerObj.Server.AddPlugin(plugin)
}
func (s *MicroServer)Run(){
    ServerObj.Server.Run()
}