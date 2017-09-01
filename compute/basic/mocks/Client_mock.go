package mocks

import context "golang.org/x/net/context"
import grpc "google.golang.org/grpc"
import mock "github.com/stretchr/testify/mock"
import scheduler "github.com/ohsu-comp-bio/funnel/proto/scheduler"
import tasklogger "github.com/ohsu-comp-bio/funnel/proto/tasklogger"
import tes "github.com/ohsu-comp-bio/funnel/proto/tes"

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Client) Close() {
	_m.Called()
}

// GetNode provides a mock function with given fields: ctx, in, opts
func (_m *Client) GetNode(ctx context.Context, in *scheduler.GetNodeRequest, opts ...grpc.CallOption) (*scheduler.Node, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *scheduler.Node
	if rf, ok := ret.Get(0).(func(context.Context, *scheduler.GetNodeRequest, ...grpc.CallOption) *scheduler.Node); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*scheduler.Node)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *scheduler.GetNodeRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListNodes provides a mock function with given fields: ctx, in, opts
func (_m *Client) ListNodes(ctx context.Context, in *scheduler.ListNodesRequest, opts ...grpc.CallOption) (*scheduler.ListNodesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *scheduler.ListNodesResponse
	if rf, ok := ret.Get(0).(func(context.Context, *scheduler.ListNodesRequest, ...grpc.CallOption) *scheduler.ListNodesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*scheduler.ListNodesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *scheduler.ListNodesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueueTask provides a mock function with given fields: ctx, in, opts
func (_m *Client) QueueTask(ctx context.Context, in *tes.Task, opts ...grpc.CallOption) (*scheduler.QueueTaskResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *scheduler.QueueTaskResponse
	if rf, ok := ret.Get(0).(func(context.Context, *tes.Task, ...grpc.CallOption) *scheduler.QueueTaskResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*scheduler.QueueTaskResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *tes.Task, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateExecutorLogs provides a mock function with given fields: ctx, in, opts
func (_m *Client) UpdateExecutorLogs(ctx context.Context, in *tasklogger.UpdateExecutorLogsRequest, opts ...grpc.CallOption) (*tasklogger.UpdateExecutorLogsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *tasklogger.UpdateExecutorLogsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *tasklogger.UpdateExecutorLogsRequest, ...grpc.CallOption) *tasklogger.UpdateExecutorLogsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tasklogger.UpdateExecutorLogsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *tasklogger.UpdateExecutorLogsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateNode provides a mock function with given fields: ctx, in, opts
func (_m *Client) UpdateNode(ctx context.Context, in *scheduler.Node, opts ...grpc.CallOption) (*scheduler.UpdateNodeResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *scheduler.UpdateNodeResponse
	if rf, ok := ret.Get(0).(func(context.Context, *scheduler.Node, ...grpc.CallOption) *scheduler.UpdateNodeResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*scheduler.UpdateNodeResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *scheduler.Node, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTaskLogs provides a mock function with given fields: ctx, in, opts
func (_m *Client) UpdateTaskLogs(ctx context.Context, in *tasklogger.UpdateTaskLogsRequest, opts ...grpc.CallOption) (*tasklogger.UpdateTaskLogsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *tasklogger.UpdateTaskLogsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *tasklogger.UpdateTaskLogsRequest, ...grpc.CallOption) *tasklogger.UpdateTaskLogsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tasklogger.UpdateTaskLogsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *tasklogger.UpdateTaskLogsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTaskState provides a mock function with given fields: ctx, in, opts
func (_m *Client) UpdateTaskState(ctx context.Context, in *tasklogger.UpdateTaskStateRequest, opts ...grpc.CallOption) (*tasklogger.UpdateTaskStateResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *tasklogger.UpdateTaskStateResponse
	if rf, ok := ret.Get(0).(func(context.Context, *tasklogger.UpdateTaskStateRequest, ...grpc.CallOption) *tasklogger.UpdateTaskStateResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tasklogger.UpdateTaskStateResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *tasklogger.UpdateTaskStateRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}