// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.2
// - protoc             v4.25.2
// source: api/gateway/v1/command.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationCommandServiceSendCommand = "/api.gateway.v1.CommandService/SendCommand"

type CommandServiceHTTPServer interface {
	// SendCommand 发送命令
	SendCommand(context.Context, *Command) (*CommandReply, error)
}

func RegisterCommandServiceHTTPServer(s *http.Server, srv CommandServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/command", _CommandService_SendCommand0_HTTP_Handler(srv))
}

func _CommandService_SendCommand0_HTTP_Handler(srv CommandServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in Command
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCommandServiceSendCommand)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SendCommand(ctx, req.(*Command))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CommandReply)
		return ctx.Result(200, reply)
	}
}

type CommandServiceHTTPClient interface {
	SendCommand(ctx context.Context, req *Command, opts ...http.CallOption) (rsp *CommandReply, err error)
}

type CommandServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewCommandServiceHTTPClient(client *http.Client) CommandServiceHTTPClient {
	return &CommandServiceHTTPClientImpl{client}
}

func (c *CommandServiceHTTPClientImpl) SendCommand(ctx context.Context, in *Command, opts ...http.CallOption) (*CommandReply, error) {
	var out CommandReply
	pattern := "/v1/command"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationCommandServiceSendCommand))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
