// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.2
// - protoc             v4.25.2
// source: api/gateway/v1/plan.proto

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

const OperationPlanCreatePlan = "/api.gateway.v1.Plan/CreatePlan"
const OperationPlanDeletePlan = "/api.gateway.v1.Plan/DeletePlan"
const OperationPlanExecPlan = "/api.gateway.v1.Plan/ExecPlan"
const OperationPlanGetPlan = "/api.gateway.v1.Plan/GetPlan"
const OperationPlanListPlan = "/api.gateway.v1.Plan/ListPlan"
const OperationPlanSwitchPlan = "/api.gateway.v1.Plan/SwitchPlan"
const OperationPlanUpdatePlan = "/api.gateway.v1.Plan/UpdatePlan"

type PlanHTTPServer interface {
	// CreatePlan 创建计划
	CreatePlan(context.Context, *CreatePlanRequest) (*CreatePlanReply, error)
	// DeletePlan 删除计划
	DeletePlan(context.Context, *DeletePlanRequest) (*DeletePlanReply, error)
	ExecPlan(context.Context, *ExecPlanRequest) (*ExecPlanReply, error)
	// GetPlan 获取计划
	GetPlan(context.Context, *GetPlanRequest) (*GetPlanReply, error)
	// ListPlan 获取计划列表
	ListPlan(context.Context, *ListPlanRequest) (*ListPlanReply, error)
	// SwitchPlan 切换计划开关
	SwitchPlan(context.Context, *SwitchPlanRequest) (*SwitchPlanReply, error)
	// UpdatePlan 更新计划
	UpdatePlan(context.Context, *UpdatePlanRequest) (*UpdatePlanReply, error)
}

func RegisterPlanHTTPServer(s *http.Server, srv PlanHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/plan/create", _Plan_CreatePlan0_HTTP_Handler(srv))
	r.POST("/v1/plan/update", _Plan_UpdatePlan0_HTTP_Handler(srv))
	r.POST("/v1/plan/delete", _Plan_DeletePlan0_HTTP_Handler(srv))
	r.GET("/v1/plan/get", _Plan_GetPlan0_HTTP_Handler(srv))
	r.GET("/v1/plan/list", _Plan_ListPlan0_HTTP_Handler(srv))
	r.POST("/v1/plan/switch", _Plan_SwitchPlan0_HTTP_Handler(srv))
	r.POST("/v1/plan/exec", _Plan_ExecPlan0_HTTP_Handler(srv))
}

func _Plan_CreatePlan0_HTTP_Handler(srv PlanHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreatePlanRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlanCreatePlan)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreatePlan(ctx, req.(*CreatePlanRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreatePlanReply)
		return ctx.Result(200, reply)
	}
}

func _Plan_UpdatePlan0_HTTP_Handler(srv PlanHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdatePlanRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlanUpdatePlan)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdatePlan(ctx, req.(*UpdatePlanRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdatePlanReply)
		return ctx.Result(200, reply)
	}
}

func _Plan_DeletePlan0_HTTP_Handler(srv PlanHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeletePlanRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlanDeletePlan)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeletePlan(ctx, req.(*DeletePlanRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeletePlanReply)
		return ctx.Result(200, reply)
	}
}

func _Plan_GetPlan0_HTTP_Handler(srv PlanHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetPlanRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlanGetPlan)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetPlan(ctx, req.(*GetPlanRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetPlanReply)
		return ctx.Result(200, reply)
	}
}

func _Plan_ListPlan0_HTTP_Handler(srv PlanHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListPlanRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlanListPlan)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListPlan(ctx, req.(*ListPlanRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListPlanReply)
		return ctx.Result(200, reply)
	}
}

func _Plan_SwitchPlan0_HTTP_Handler(srv PlanHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in SwitchPlanRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlanSwitchPlan)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SwitchPlan(ctx, req.(*SwitchPlanRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SwitchPlanReply)
		return ctx.Result(200, reply)
	}
}

func _Plan_ExecPlan0_HTTP_Handler(srv PlanHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ExecPlanRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlanExecPlan)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ExecPlan(ctx, req.(*ExecPlanRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ExecPlanReply)
		return ctx.Result(200, reply)
	}
}

type PlanHTTPClient interface {
	CreatePlan(ctx context.Context, req *CreatePlanRequest, opts ...http.CallOption) (rsp *CreatePlanReply, err error)
	DeletePlan(ctx context.Context, req *DeletePlanRequest, opts ...http.CallOption) (rsp *DeletePlanReply, err error)
	ExecPlan(ctx context.Context, req *ExecPlanRequest, opts ...http.CallOption) (rsp *ExecPlanReply, err error)
	GetPlan(ctx context.Context, req *GetPlanRequest, opts ...http.CallOption) (rsp *GetPlanReply, err error)
	ListPlan(ctx context.Context, req *ListPlanRequest, opts ...http.CallOption) (rsp *ListPlanReply, err error)
	SwitchPlan(ctx context.Context, req *SwitchPlanRequest, opts ...http.CallOption) (rsp *SwitchPlanReply, err error)
	UpdatePlan(ctx context.Context, req *UpdatePlanRequest, opts ...http.CallOption) (rsp *UpdatePlanReply, err error)
}

type PlanHTTPClientImpl struct {
	cc *http.Client
}

func NewPlanHTTPClient(client *http.Client) PlanHTTPClient {
	return &PlanHTTPClientImpl{client}
}

func (c *PlanHTTPClientImpl) CreatePlan(ctx context.Context, in *CreatePlanRequest, opts ...http.CallOption) (*CreatePlanReply, error) {
	var out CreatePlanReply
	pattern := "/v1/plan/create"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPlanCreatePlan))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlanHTTPClientImpl) DeletePlan(ctx context.Context, in *DeletePlanRequest, opts ...http.CallOption) (*DeletePlanReply, error) {
	var out DeletePlanReply
	pattern := "/v1/plan/delete"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPlanDeletePlan))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlanHTTPClientImpl) ExecPlan(ctx context.Context, in *ExecPlanRequest, opts ...http.CallOption) (*ExecPlanReply, error) {
	var out ExecPlanReply
	pattern := "/v1/plan/exec"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPlanExecPlan))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlanHTTPClientImpl) GetPlan(ctx context.Context, in *GetPlanRequest, opts ...http.CallOption) (*GetPlanReply, error) {
	var out GetPlanReply
	pattern := "/v1/plan/get"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationPlanGetPlan))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlanHTTPClientImpl) ListPlan(ctx context.Context, in *ListPlanRequest, opts ...http.CallOption) (*ListPlanReply, error) {
	var out ListPlanReply
	pattern := "/v1/plan/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationPlanListPlan))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlanHTTPClientImpl) SwitchPlan(ctx context.Context, in *SwitchPlanRequest, opts ...http.CallOption) (*SwitchPlanReply, error) {
	var out SwitchPlanReply
	pattern := "/v1/plan/switch"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPlanSwitchPlan))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlanHTTPClientImpl) UpdatePlan(ctx context.Context, in *UpdatePlanRequest, opts ...http.CallOption) (*UpdatePlanReply, error) {
	var out UpdatePlanReply
	pattern := "/v1/plan/update"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPlanUpdatePlan))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
