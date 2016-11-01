// Code generated by protoc-gen-go.
// source: task_execution.proto
// DO NOT EDIT!

/*
Package ga4gh_task_exec is a generated protocol buffer package.

It is generated from these files:
	task_execution.proto

It has these top-level messages:
	TaskParameter
	DockerExecutor
	Volume
	Resources
	Task
	TaskListRequest
	TaskListResponse
	JobListRequest
	JobListResponse
	TaskId
	JobId
	JobLog
	Job
	ServiceInfoRequest
	ServiceInfo
*/
package ga4gh_task_exec

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type State int32

const (
	State_Unknown     State = 0
	State_Queued      State = 1
	State_Running     State = 2
	State_Paused      State = 3
	State_Complete    State = 4
	State_Error       State = 5
	State_SystemError State = 6
	State_Canceled    State = 7
)

var State_name = map[int32]string{
	0: "Unknown",
	1: "Queued",
	2: "Running",
	3: "Paused",
	4: "Complete",
	5: "Error",
	6: "SystemError",
	7: "Canceled",
}
var State_value = map[string]int32{
	"Unknown":     0,
	"Queued":      1,
	"Running":     2,
	"Paused":      3,
	"Complete":    4,
	"Error":       5,
	"SystemError": 6,
	"Canceled":    7,
}

func (x State) String() string {
	return proto.EnumName(State_name, int32(x))
}
func (State) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Parameters for Task
type TaskParameter struct {
	// name of the parameter
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// optional text description
	Description string `protobuf:"bytes,2,opt,name=description" json:"description,omitempty"`
	// location in long term storage, is a url specific to the implementing
	// system. For example s3://my-object-store/file1 or gs://my-bucket/file2 or
	// file:///path/to/my/file
	Location string `protobuf:"bytes,3,opt,name=location" json:"location,omitempty"`
	// path in the machine file system. Note, this MUST be a path that exists
	// within one of the defined volumes
	Path string `protobuf:"bytes,4,opt,name=path" json:"path,omitempty"`
	// Type of data, "File" or "Directory"
	// if used for an output all the files in the directory
	// will be copied to the storage location
	Class string `protobuf:"bytes,5,opt,name=class" json:"class,omitempty"`
	// if the parameter is an output, should the element be created before executing
	// the command. For example if saving the working directory as an output,
	// the directory needs to exist before the command is called. If false, it is
	// assumed that the user will create the element as a part of the operation
	Create bool `protobuf:"varint,6,opt,name=create" json:"create,omitempty"`
}

func (m *TaskParameter) Reset()                    { *m = TaskParameter{} }
func (m *TaskParameter) String() string            { return proto.CompactTextString(m) }
func (*TaskParameter) ProtoMessage()               {}
func (*TaskParameter) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// A command line to be executed and the docker container to run it
type DockerExecutor struct {
	// Docker Image name
	ImageName string `protobuf:"bytes,1,opt,name=imageName" json:"imageName,omitempty"`
	// The command to be executed
	Cmd []string `protobuf:"bytes,2,rep,name=cmd" json:"cmd,omitempty"`
	// The working directory that the command will be executed in
	Workdir string `protobuf:"bytes,3,opt,name=workdir" json:"workdir,omitempty"`
	// Path for stdout recording, blank if not storing to file
	Stdout string `protobuf:"bytes,4,opt,name=stdout" json:"stdout,omitempty"`
	// Path for stderr recording, blank if not storing to file
	Stderr string `protobuf:"bytes,5,opt,name=stderr" json:"stderr,omitempty"`
}

func (m *DockerExecutor) Reset()                    { *m = DockerExecutor{} }
func (m *DockerExecutor) String() string            { return proto.CompactTextString(m) }
func (*DockerExecutor) ProtoMessage()               {}
func (*DockerExecutor) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// Attached volume request.
type Volume struct {
	// Name of attached volume
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Minimum size
	SizeGb uint32 `protobuf:"varint,2,opt,name=sizeGb" json:"sizeGb,omitempty"`
	// Source volume, this would refer to an existing volume the execution engine
	// could identify. Leave blank if is to be a newly created volume
	// Volumes loaded from a source will be mounted as read only
	Source string `protobuf:"bytes,3,opt,name=source" json:"source,omitempty"`
	// mount point for volume inside the docker container
	MountPoint string `protobuf:"bytes,6,opt,name=mountPoint" json:"mountPoint,omitempty"`
}

func (m *Volume) Reset()                    { *m = Volume{} }
func (m *Volume) String() string            { return proto.CompactTextString(m) }
func (*Volume) ProtoMessage()               {}
func (*Volume) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Resources struct {
	// Minimum number of CPUs
	MinimumCpuCores uint32 `protobuf:"varint,1,opt,name=minimumCpuCores" json:"minimumCpuCores,omitempty"`
	// Can schedule on resource that resource that can be preempted, like AWS Spot Instances
	Preemptible bool `protobuf:"varint,2,opt,name=preemptible" json:"preemptible,omitempty"`
	// Minimum RAM required
	MinimumRamGb uint32 `protobuf:"varint,3,opt,name=minimumRamGb" json:"minimumRamGb,omitempty"`
	// Volumes to be mounted into the docker container
	Volumes []*Volume `protobuf:"bytes,4,rep,name=volumes" json:"volumes,omitempty"`
	// optional scheduling information for systems where multiple compute zones are avalible
	Zones []string `protobuf:"bytes,5,rep,name=zones" json:"zones,omitempty"`
}

func (m *Resources) Reset()                    { *m = Resources{} }
func (m *Resources) String() string            { return proto.CompactTextString(m) }
func (*Resources) ProtoMessage()               {}
func (*Resources) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Resources) GetVolumes() []*Volume {
	if m != nil {
		return m.Volumes
	}
	return nil
}

// The description of a task to be run
type Task struct {
	// user name for task
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// parameter for execution engine to define/store group information
	ProjectId string `protobuf:"bytes,2,opt,name=projectId" json:"projectId,omitempty"`
	// free text description of task
	Description string `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
	// Files to be copied into system before tasks
	Inputs []*TaskParameter `protobuf:"bytes,4,rep,name=inputs" json:"inputs,omitempty"`
	// Files to be copied out of the system after tasks
	Outputs []*TaskParameter `protobuf:"bytes,5,rep,name=outputs" json:"outputs,omitempty"`
	// Define required system resources to run job
	Resources *Resources `protobuf:"bytes,6,opt,name=resources" json:"resources,omitempty"`
	// System defined identifier of task
	TaskId string `protobuf:"bytes,7,opt,name=taskId" json:"taskId,omitempty"`
	// An array of docker executions that will be run sequentially
	Docker []*DockerExecutor `protobuf:"bytes,8,rep,name=docker" json:"docker,omitempty"`
}

func (m *Task) Reset()                    { *m = Task{} }
func (m *Task) String() string            { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()               {}
func (*Task) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Task) GetInputs() []*TaskParameter {
	if m != nil {
		return m.Inputs
	}
	return nil
}

func (m *Task) GetOutputs() []*TaskParameter {
	if m != nil {
		return m.Outputs
	}
	return nil
}

func (m *Task) GetResources() *Resources {
	if m != nil {
		return m.Resources
	}
	return nil
}

func (m *Task) GetDocker() []*DockerExecutor {
	if m != nil {
		return m.Docker
	}
	return nil
}

type TaskListRequest struct {
	// Required. The name of the project to search for pipelines. Caller must have READ access to this project.
	ProjectId string `protobuf:"bytes,1,opt,name=projectId" json:"projectId,omitempty"`
	// Pipelines with names that match this prefix should be returned. If unspecified, all pipelines in the project, up to pageSize, will be returned.
	NamePrefix string `protobuf:"bytes,2,opt,name=namePrefix" json:"namePrefix,omitempty"`
	// Number of pipelines to return at once. Defaults to 256, and max is 2048.
	PageSize uint32 `protobuf:"varint,3,opt,name=pageSize" json:"pageSize,omitempty"`
	// Token to use to indicate where to start getting results. If unspecified, returns the first page of results.
	PageToken string `protobuf:"bytes,4,opt,name=pageToken" json:"pageToken,omitempty"`
}

func (m *TaskListRequest) Reset()                    { *m = TaskListRequest{} }
func (m *TaskListRequest) String() string            { return proto.CompactTextString(m) }
func (*TaskListRequest) ProtoMessage()               {}
func (*TaskListRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type TaskListResponse struct {
	Tasks         []*Task `protobuf:"bytes,1,rep,name=tasks" json:"tasks,omitempty"`
	NextPageToken string  `protobuf:"bytes,2,opt,name=nextPageToken" json:"nextPageToken,omitempty"`
}

func (m *TaskListResponse) Reset()                    { *m = TaskListResponse{} }
func (m *TaskListResponse) String() string            { return proto.CompactTextString(m) }
func (*TaskListResponse) ProtoMessage()               {}
func (*TaskListResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *TaskListResponse) GetTasks() []*Task {
	if m != nil {
		return m.Tasks
	}
	return nil
}

type JobListRequest struct {
	// Required. The name of the project to search for pipelines. Caller must have READ access to this project.
	ProjectId string `protobuf:"bytes,1,opt,name=projectId" json:"projectId,omitempty"`
	// Pipelines with names that match this prefix should be returned. If unspecified, all pipelines in the project, up to pageSize, will be returned.
	NamePrefix string `protobuf:"bytes,2,opt,name=namePrefix" json:"namePrefix,omitempty"`
	// Number of pipelines to return at once. Defaults to 256, and max is 2048.
	PageSize uint32 `protobuf:"varint,3,opt,name=pageSize" json:"pageSize,omitempty"`
	// Token to use to indicate where to start getting results. If unspecified, returns the first page of results.
	PageToken string `protobuf:"bytes,4,opt,name=pageToken" json:"pageToken,omitempty"`
}

func (m *JobListRequest) Reset()                    { *m = JobListRequest{} }
func (m *JobListRequest) String() string            { return proto.CompactTextString(m) }
func (*JobListRequest) ProtoMessage()               {}
func (*JobListRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type JobListResponse struct {
	Jobs          []*Job `protobuf:"bytes,1,rep,name=jobs" json:"jobs,omitempty"`
	NextPageToken string `protobuf:"bytes,2,opt,name=nextPageToken" json:"nextPageToken,omitempty"`
}

func (m *JobListResponse) Reset()                    { *m = JobListResponse{} }
func (m *JobListResponse) String() string            { return proto.CompactTextString(m) }
func (*JobListResponse) ProtoMessage()               {}
func (*JobListResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *JobListResponse) GetJobs() []*Job {
	if m != nil {
		return m.Jobs
	}
	return nil
}

// ID of a Task description
type TaskId struct {
	Value string `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
}

func (m *TaskId) Reset()                    { *m = TaskId{} }
func (m *TaskId) String() string            { return proto.CompactTextString(m) }
func (*TaskId) ProtoMessage()               {}
func (*TaskId) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

// ID of an instance of a Task
type JobId struct {
	Value string `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
}

func (m *JobId) Reset()                    { *m = JobId{} }
func (m *JobId) String() string            { return proto.CompactTextString(m) }
func (*JobId) ProtoMessage()               {}
func (*JobId) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type JobLog struct {
	// The command line that was run
	Cmd []string `protobuf:"bytes,1,rep,name=cmd" json:"cmd,omitempty"`
	// When the command was executed
	StartTime string `protobuf:"bytes,2,opt,name=startTime" json:"startTime,omitempty"`
	// When the command completed
	EndTime string `protobuf:"bytes,3,opt,name=endTime" json:"endTime,omitempty"`
	// Sample of stdout (not guaranteed to be entire log)
	Stdout string `protobuf:"bytes,4,opt,name=stdout" json:"stdout,omitempty"`
	// Sample of stderr (not guaranteed to be entire log)
	Stderr string `protobuf:"bytes,5,opt,name=stderr" json:"stderr,omitempty"`
	// Exit code of the program
	ExitCode int32 `protobuf:"varint,6,opt,name=exitCode" json:"exitCode,omitempty"`
}

func (m *JobLog) Reset()                    { *m = JobLog{} }
func (m *JobLog) String() string            { return proto.CompactTextString(m) }
func (*JobLog) ProtoMessage()               {}
func (*JobLog) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

// The description of the running instance of a task
type Job struct {
	JobId    string            `protobuf:"bytes,1,opt,name=jobId" json:"jobId,omitempty"`
	Metadata map[string]string `protobuf:"bytes,2,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Task     *Task             `protobuf:"bytes,3,opt,name=task" json:"task,omitempty"`
	State    State             `protobuf:"varint,4,opt,name=state,enum=ga4gh_task_exec.State" json:"state,omitempty"`
	Logs     []*JobLog         `protobuf:"bytes,5,rep,name=logs" json:"logs,omitempty"`
}

func (m *Job) Reset()                    { *m = Job{} }
func (m *Job) String() string            { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()               {}
func (*Job) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *Job) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Job) GetTask() *Task {
	if m != nil {
		return m.Task
	}
	return nil
}

func (m *Job) GetLogs() []*JobLog {
	if m != nil {
		return m.Logs
	}
	return nil
}

// Blank request message for service request
type ServiceInfoRequest struct {
}

func (m *ServiceInfoRequest) Reset()                    { *m = ServiceInfoRequest{} }
func (m *ServiceInfoRequest) String() string            { return proto.CompactTextString(m) }
func (*ServiceInfoRequest) ProtoMessage()               {}
func (*ServiceInfoRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

// Information about Task Execution Service
// May include information related (but not limited to)
// resource availability and storage system information
type ServiceInfo struct {
	// System specific key/value pairs
	// Example for a shared file system based storage system:
	// storageType=sharedFile, baseDir=/path/to/shared/directory
	StorageConfig map[string]string `protobuf:"bytes,1,rep,name=storageConfig" json:"storageConfig,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *ServiceInfo) Reset()                    { *m = ServiceInfo{} }
func (m *ServiceInfo) String() string            { return proto.CompactTextString(m) }
func (*ServiceInfo) ProtoMessage()               {}
func (*ServiceInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *ServiceInfo) GetStorageConfig() map[string]string {
	if m != nil {
		return m.StorageConfig
	}
	return nil
}

func init() {
	proto.RegisterType((*TaskParameter)(nil), "ga4gh_task_exec.TaskParameter")
	proto.RegisterType((*DockerExecutor)(nil), "ga4gh_task_exec.DockerExecutor")
	proto.RegisterType((*Volume)(nil), "ga4gh_task_exec.Volume")
	proto.RegisterType((*Resources)(nil), "ga4gh_task_exec.Resources")
	proto.RegisterType((*Task)(nil), "ga4gh_task_exec.Task")
	proto.RegisterType((*TaskListRequest)(nil), "ga4gh_task_exec.TaskListRequest")
	proto.RegisterType((*TaskListResponse)(nil), "ga4gh_task_exec.TaskListResponse")
	proto.RegisterType((*JobListRequest)(nil), "ga4gh_task_exec.JobListRequest")
	proto.RegisterType((*JobListResponse)(nil), "ga4gh_task_exec.JobListResponse")
	proto.RegisterType((*TaskId)(nil), "ga4gh_task_exec.TaskId")
	proto.RegisterType((*JobId)(nil), "ga4gh_task_exec.JobId")
	proto.RegisterType((*JobLog)(nil), "ga4gh_task_exec.JobLog")
	proto.RegisterType((*Job)(nil), "ga4gh_task_exec.Job")
	proto.RegisterType((*ServiceInfoRequest)(nil), "ga4gh_task_exec.ServiceInfoRequest")
	proto.RegisterType((*ServiceInfo)(nil), "ga4gh_task_exec.ServiceInfo")
	proto.RegisterEnum("ga4gh_task_exec.State", State_name, State_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for TaskService service

type TaskServiceClient interface {
	// Get Service Info
	GetServiceInfo(ctx context.Context, in *ServiceInfoRequest, opts ...grpc.CallOption) (*ServiceInfo, error)
	// Run a task
	RunTask(ctx context.Context, in *Task, opts ...grpc.CallOption) (*JobId, error)
	// List the TaskOps
	ListJobs(ctx context.Context, in *JobListRequest, opts ...grpc.CallOption) (*JobListResponse, error)
	// Get info about a running task
	GetJob(ctx context.Context, in *JobId, opts ...grpc.CallOption) (*Job, error)
	// Cancel a running task
	CancelJob(ctx context.Context, in *JobId, opts ...grpc.CallOption) (*JobId, error)
}

type taskServiceClient struct {
	cc *grpc.ClientConn
}

func NewTaskServiceClient(cc *grpc.ClientConn) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) GetServiceInfo(ctx context.Context, in *ServiceInfoRequest, opts ...grpc.CallOption) (*ServiceInfo, error) {
	out := new(ServiceInfo)
	err := grpc.Invoke(ctx, "/ga4gh_task_exec.TaskService/GetServiceInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) RunTask(ctx context.Context, in *Task, opts ...grpc.CallOption) (*JobId, error) {
	out := new(JobId)
	err := grpc.Invoke(ctx, "/ga4gh_task_exec.TaskService/RunTask", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) ListJobs(ctx context.Context, in *JobListRequest, opts ...grpc.CallOption) (*JobListResponse, error) {
	out := new(JobListResponse)
	err := grpc.Invoke(ctx, "/ga4gh_task_exec.TaskService/ListJobs", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) GetJob(ctx context.Context, in *JobId, opts ...grpc.CallOption) (*Job, error) {
	out := new(Job)
	err := grpc.Invoke(ctx, "/ga4gh_task_exec.TaskService/GetJob", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) CancelJob(ctx context.Context, in *JobId, opts ...grpc.CallOption) (*JobId, error) {
	out := new(JobId)
	err := grpc.Invoke(ctx, "/ga4gh_task_exec.TaskService/CancelJob", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for TaskService service

type TaskServiceServer interface {
	// Get Service Info
	GetServiceInfo(context.Context, *ServiceInfoRequest) (*ServiceInfo, error)
	// Run a task
	RunTask(context.Context, *Task) (*JobId, error)
	// List the TaskOps
	ListJobs(context.Context, *JobListRequest) (*JobListResponse, error)
	// Get info about a running task
	GetJob(context.Context, *JobId) (*Job, error)
	// Cancel a running task
	CancelJob(context.Context, *JobId) (*JobId, error)
}

func RegisterTaskServiceServer(s *grpc.Server, srv TaskServiceServer) {
	s.RegisterService(&_TaskService_serviceDesc, srv)
}

func _TaskService_GetServiceInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).GetServiceInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ga4gh_task_exec.TaskService/GetServiceInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).GetServiceInfo(ctx, req.(*ServiceInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_RunTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).RunTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ga4gh_task_exec.TaskService/RunTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).RunTask(ctx, req.(*Task))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_ListJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).ListJobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ga4gh_task_exec.TaskService/ListJobs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).ListJobs(ctx, req.(*JobListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_GetJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).GetJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ga4gh_task_exec.TaskService/GetJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).GetJob(ctx, req.(*JobId))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_CancelJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).CancelJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ga4gh_task_exec.TaskService/CancelJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).CancelJob(ctx, req.(*JobId))
	}
	return interceptor(ctx, in, info, handler)
}

var _TaskService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ga4gh_task_exec.TaskService",
	HandlerType: (*TaskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetServiceInfo",
			Handler:    _TaskService_GetServiceInfo_Handler,
		},
		{
			MethodName: "RunTask",
			Handler:    _TaskService_RunTask_Handler,
		},
		{
			MethodName: "ListJobs",
			Handler:    _TaskService_ListJobs_Handler,
		},
		{
			MethodName: "GetJob",
			Handler:    _TaskService_GetJob_Handler,
		},
		{
			MethodName: "CancelJob",
			Handler:    _TaskService_CancelJob_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("task_execution.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1113 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xcc, 0x56, 0xcd, 0x6e, 0x23, 0x45,
	0x10, 0x66, 0xfc, 0xef, 0xf2, 0x3a, 0xb1, 0x9a, 0x90, 0x1d, 0x59, 0x21, 0x6b, 0x0d, 0x1c, 0x4c,
	0x16, 0x62, 0x6d, 0x40, 0x10, 0x2d, 0x12, 0x42, 0x0a, 0x51, 0x94, 0x68, 0x17, 0x99, 0x4e, 0xc2,
	0x35, 0x6a, 0x7b, 0x2a, 0xb3, 0x13, 0x7b, 0xba, 0x67, 0x7b, 0x7a, 0xb2, 0x49, 0x10, 0x17, 0x4e,
	0xc0, 0x95, 0xfb, 0x5e, 0x90, 0x78, 0x12, 0xde, 0x00, 0x21, 0x5e, 0x80, 0x07, 0x41, 0xdd, 0xd3,
	0xe3, 0x9f, 0xd8, 0x89, 0x76, 0x6f, 0xdc, 0xa6, 0xaa, 0xbf, 0xaa, 0xfe, 0xea, 0xab, 0xea, 0xb2,
	0x61, 0x4d, 0xb1, 0x64, 0x74, 0x86, 0x57, 0x38, 0x4c, 0x55, 0x28, 0xf8, 0x76, 0x2c, 0x85, 0x12,
	0x64, 0x35, 0x60, 0x9f, 0x05, 0x2f, 0xce, 0x26, 0x67, 0xed, 0x8d, 0x40, 0x88, 0x60, 0x8c, 0x3d,
	0x16, 0x87, 0x3d, 0xc6, 0xb9, 0x50, 0x4c, 0xa3, 0x93, 0x0c, 0xee, 0xfd, 0xee, 0x40, 0xf3, 0x84,
	0x25, 0xa3, 0x3e, 0x93, 0x2c, 0x42, 0x85, 0x92, 0x10, 0x28, 0x71, 0x16, 0xa1, 0xeb, 0x74, 0x9c,
	0x6e, 0x9d, 0x9a, 0x6f, 0xd2, 0x81, 0x86, 0x8f, 0xc9, 0x50, 0x86, 0xb1, 0x8e, 0x75, 0x0b, 0xe6,
	0x68, 0xd6, 0x45, 0xda, 0x50, 0x1b, 0x8b, 0xa1, 0x49, 0xed, 0x16, 0xcd, 0xf1, 0xc4, 0xd6, 0x19,
	0x63, 0xa6, 0x5e, 0xb8, 0xa5, 0x2c, 0xa3, 0xfe, 0x26, 0x6b, 0x50, 0x1e, 0x8e, 0x59, 0x92, 0xb8,
	0x65, 0xe3, 0xcc, 0x0c, 0xb2, 0x0e, 0x95, 0xa1, 0x44, 0xa6, 0xd0, 0xad, 0x74, 0x9c, 0x6e, 0x8d,
	0x5a, 0xcb, 0xfb, 0xd5, 0x81, 0x95, 0x6f, 0xc4, 0x70, 0x84, 0x72, 0xdf, 0x94, 0x2b, 0x24, 0xd9,
	0x80, 0x7a, 0x18, 0xb1, 0x00, 0xbf, 0x9d, 0x72, 0x9d, 0x3a, 0x48, 0x0b, 0x8a, 0xc3, 0xc8, 0x77,
	0x0b, 0x9d, 0x62, 0xb7, 0x4e, 0xf5, 0x27, 0x71, 0xa1, 0xfa, 0x4a, 0xc8, 0x91, 0x1f, 0x4a, 0xcb,
	0x2f, 0x37, 0xf5, 0xa5, 0x89, 0xf2, 0x45, 0xaa, 0x2c, 0x41, 0x6b, 0x59, 0x3f, 0x4a, 0x69, 0x39,
	0x5a, 0xcb, 0x1b, 0x43, 0xe5, 0x7b, 0x31, 0x4e, 0x23, 0x5c, 0x2a, 0x95, 0x8e, 0x0a, 0x6f, 0xf0,
	0x60, 0x60, 0x54, 0x6a, 0x52, 0x6b, 0x19, 0xbf, 0x48, 0xe5, 0x10, 0xed, 0xf5, 0xd6, 0x22, 0x9b,
	0x00, 0x91, 0x48, 0xb9, 0xea, 0x8b, 0x90, 0x2b, 0x53, 0x76, 0x9d, 0xce, 0x78, 0xbc, 0x3f, 0x1d,
	0xa8, 0x53, 0xcc, 0xc0, 0x09, 0xe9, 0xc2, 0x6a, 0x14, 0xf2, 0x30, 0x4a, 0xa3, 0xbd, 0x38, 0xdd,
	0x13, 0x12, 0x13, 0x73, 0x79, 0x93, 0xde, 0x76, 0xeb, 0x96, 0xc5, 0x12, 0x31, 0x8a, 0x55, 0x38,
	0x18, 0xa3, 0x21, 0x53, 0xa3, 0xb3, 0x2e, 0xe2, 0xc1, 0x03, 0x1b, 0x44, 0x59, 0x74, 0x30, 0x30,
	0xbc, 0x9a, 0x74, 0xce, 0x47, 0x9e, 0x40, 0xf5, 0xd2, 0xd4, 0x9a, 0xb8, 0xa5, 0x4e, 0xb1, 0xdb,
	0xd8, 0x79, 0xb8, 0x7d, 0x6b, 0xbe, 0xb6, 0x33, 0x2d, 0x68, 0x8e, 0xd3, 0x9d, 0xbd, 0x11, 0x1c,
	0x75, 0x67, 0xb5, 0xf8, 0x99, 0xe1, 0xfd, 0x53, 0x80, 0x92, 0x9e, 0xb3, 0xa5, 0x9a, 0x6d, 0x40,
	0x3d, 0x96, 0xe2, 0x02, 0x87, 0xea, 0xd0, 0xb7, 0xc3, 0x35, 0x75, 0xdc, 0x1e, 0xbe, 0xe2, 0xe2,
	0xf0, 0x7d, 0x0e, 0x95, 0x90, 0xc7, 0xa9, 0xca, 0x49, 0x6e, 0x2e, 0x90, 0x9c, 0x1b, 0x71, 0x6a,
	0xd1, 0x64, 0x17, 0xaa, 0x22, 0x55, 0x26, 0xb0, 0xfc, 0x46, 0x81, 0x39, 0x9c, 0xec, 0x42, 0x5d,
	0xe6, 0x4d, 0x31, 0x4d, 0x6b, 0xec, 0xb4, 0x17, 0x62, 0x27, 0x6d, 0xa3, 0x53, 0xb0, 0x9e, 0x03,
	0x8d, 0x38, 0xf4, 0xdd, 0x6a, 0x36, 0x07, 0x99, 0x45, 0xbe, 0x80, 0x8a, 0x6f, 0x26, 0xdc, 0xad,
	0x19, 0x2a, 0x8f, 0x16, 0xd2, 0xcd, 0x3f, 0x00, 0x6a, 0xe1, 0xde, 0x2f, 0x0e, 0xac, 0x6a, 0x96,
	0xcf, 0xc2, 0x44, 0x51, 0x7c, 0x99, 0x62, 0xa2, 0xe6, 0x05, 0x75, 0x6e, 0x0b, 0xba, 0x09, 0xa0,
	0x65, 0xef, 0x4b, 0x3c, 0x0f, 0xaf, 0xac, 0xde, 0x33, 0x1e, 0xfd, 0x96, 0x63, 0x16, 0xe0, 0x71,
	0x78, 0x83, 0x76, 0x28, 0x26, 0xb6, 0xc9, 0xcc, 0x02, 0x3c, 0x11, 0x23, 0xe4, 0xf6, 0xbd, 0x4c,
	0x1d, 0x1e, 0x42, 0x6b, 0x4a, 0x25, 0x89, 0x05, 0x4f, 0x90, 0x3c, 0x86, 0xb2, 0xae, 0x41, 0x0f,
	0xaa, 0xae, 0xeb, 0xbd, 0xa5, 0x12, 0xd3, 0x0c, 0x43, 0x3e, 0x84, 0x26, 0xc7, 0x2b, 0xd5, 0x9f,
	0x5c, 0x91, 0xb1, 0x9b, 0x77, 0x7a, 0x3f, 0x3b, 0xb0, 0x72, 0x24, 0x06, 0xff, 0x87, 0x8a, 0x19,
	0xac, 0x4e, 0x98, 0xd8, 0x82, 0xbb, 0x50, 0xba, 0x10, 0x83, 0xbc, 0xde, 0xb5, 0x85, 0x7a, 0x8f,
	0xc4, 0x80, 0x1a, 0xc4, 0x1b, 0x56, 0xbb, 0x09, 0x95, 0x93, 0x6c, 0x46, 0xd6, 0xa0, 0x7c, 0xc9,
	0xc6, 0x69, 0xfe, 0x78, 0x32, 0xc3, 0x7b, 0x1f, 0xca, 0x47, 0x62, 0x70, 0xe7, 0xf1, 0x6b, 0x07,
	0x2a, 0x9a, 0xa2, 0x08, 0xf2, 0xad, 0xe8, 0x4c, 0xb7, 0xe2, 0x06, 0xd4, 0x13, 0xc5, 0xa4, 0x3a,
	0x09, 0x23, 0xcc, 0x5f, 0xde, 0xc4, 0xa1, 0x77, 0x26, 0x72, 0xdf, 0x9c, 0xd9, 0x9d, 0x69, 0xcd,
	0xb7, 0xdd, 0x99, 0x5a, 0x60, 0xbc, 0x0a, 0xd5, 0x9e, 0xf0, 0xb3, 0xd5, 0x5e, 0xa6, 0x13, 0xdb,
	0x7b, 0x5d, 0x80, 0xe2, 0x91, 0x18, 0x68, 0xfa, 0x17, 0xba, 0x8e, 0x9c, 0xbe, 0x31, 0xc8, 0x57,
	0x50, 0x8b, 0x50, 0x31, 0x9f, 0x29, 0x66, 0xd6, 0x79, 0x63, 0xc7, 0x5b, 0xa6, 0xe8, 0xf6, 0x73,
	0x0b, 0xda, 0xe7, 0x4a, 0x5e, 0xd3, 0x49, 0x0c, 0xf9, 0x08, 0x4a, 0x1a, 0x68, 0x0a, 0xb8, 0x73,
	0xfa, 0x0c, 0x84, 0x7c, 0x0c, 0xe5, 0x44, 0xe9, 0x1f, 0x1f, 0x5d, 0xd3, 0xca, 0xce, 0xfa, 0x02,
	0xf6, 0x58, 0x9f, 0xd2, 0x0c, 0x44, 0x1e, 0x43, 0x69, 0x2c, 0x82, 0x7c, 0x73, 0x3c, 0x5c, 0x46,
	0xea, 0x99, 0x08, 0xa8, 0x01, 0xb5, 0xbf, 0x84, 0xe6, 0x1c, 0x41, 0xdd, 0x8a, 0x11, 0x5e, 0xdb,
	0x52, 0xf5, 0xe7, 0xb4, 0x7b, 0x85, 0x99, 0xee, 0x3d, 0x2d, 0xec, 0x3a, 0xde, 0x1a, 0x90, 0x63,
	0x94, 0x97, 0xe1, 0x10, 0x0f, 0xf9, 0xb9, 0xb0, 0x13, 0xef, 0xfd, 0xe1, 0x40, 0x63, 0xc6, 0x4d,
	0x4e, 0xa1, 0x99, 0x28, 0x21, 0x59, 0x80, 0x7b, 0x82, 0x9f, 0x87, 0x81, 0x9d, 0xbf, 0xde, 0x62,
	0x15, 0xd3, 0xa0, 0xed, 0xe3, 0xd9, 0x88, 0x4c, 0xba, 0xf9, 0x2c, 0xed, 0xaf, 0x81, 0x2c, 0x82,
	0xde, 0x86, 0xfe, 0xd6, 0x4b, 0x28, 0x1b, 0xe1, 0x48, 0x03, 0xaa, 0xa7, 0x7c, 0xc4, 0xc5, 0x2b,
	0xde, 0x7a, 0x87, 0x00, 0x54, 0xbe, 0x4b, 0x31, 0x45, 0xbf, 0xe5, 0xe8, 0x03, 0x9a, 0x72, 0x1e,
	0xf2, 0xa0, 0x55, 0xd0, 0x07, 0x7d, 0x96, 0x26, 0xe8, 0xb7, 0x8a, 0xe4, 0x01, 0xd4, 0xf6, 0x44,
	0x14, 0x8f, 0x51, 0x61, 0xab, 0x44, 0xea, 0x50, 0xde, 0x97, 0x52, 0xc8, 0x56, 0x99, 0xac, 0x42,
	0xe3, 0xf8, 0x3a, 0x51, 0x18, 0x65, 0x8e, 0x8a, 0x41, 0x32, 0x3e, 0xc4, 0x31, 0xfa, 0xad, 0xea,
	0xce, 0xdf, 0x45, 0x68, 0xe8, 0xc6, 0xda, 0x52, 0x49, 0x04, 0x2b, 0x07, 0xa8, 0x66, 0xd5, 0xfa,
	0xe0, 0x3e, 0x59, 0xac, 0xc4, 0xed, 0x8d, 0xfb, 0x40, 0x9e, 0xfb, 0xd3, 0x5f, 0xff, 0xfe, 0x56,
	0x20, 0xa4, 0xd5, 0xbb, 0x7c, 0xd2, 0xd3, 0xef, 0xf9, 0x93, 0xc4, 0x5e, 0xf7, 0xdc, 0xd4, 0x63,
	0x7e, 0xee, 0x96, 0x0f, 0x5c, 0x7b, 0x7d, 0xd9, 0xb8, 0x1c, 0xfa, 0xde, 0xbb, 0x26, 0x67, 0xd3,
	0xab, 0xe5, 0x39, 0x9f, 0x3a, 0x5b, 0xe4, 0x0c, 0x6a, 0x7a, 0xc1, 0x1c, 0xe9, 0x95, 0xf1, 0x68,
	0xe9, 0x9c, 0x4d, 0x17, 0x61, 0xbb, 0x73, 0x37, 0x20, 0xdb, 0x4f, 0x5e, 0xcb, 0xdc, 0x01, 0x64,
	0x72, 0x07, 0xe9, 0x43, 0xe5, 0x00, 0x75, 0x7e, 0x72, 0x07, 0xaf, 0xf6, 0xd2, 0x2d, 0xb6, 0xa8,
	0x40, 0xef, 0x07, 0xd3, 0xf6, 0x1f, 0xc9, 0x29, 0xd4, 0xb3, 0x76, 0xdc, 0x97, 0xf4, 0x2e, 0x11,
	0x6c, 0xda, 0xad, 0x85, 0xb4, 0x83, 0x8a, 0xf9, 0xd3, 0xfa, 0xe9, 0x7f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0xf1, 0xba, 0x1b, 0x48, 0xfb, 0x0a, 0x00, 0x00,
}
