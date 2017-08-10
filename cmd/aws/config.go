package aws

// Config represents configuration of the AWS proxy, including
// the compute environment, job queue, and base job definition.
type Config struct {
	Container       string
	DefaultTaskName string
	ComputeEnv      ComputeEnvConfig
	JobDef          JobDefConfig
	JobQueue        JobQueueConfig
}

// ComputeEnvConfig represents configuration of the AWS Batch
// Compute Environment.
type ComputeEnvConfig struct {
	Name             string
	MinVCPUs         int64
	MaxVCPUs         int64
	SecurityGroupIds []string
	Subnets          []string
	Tags             map[string]string
	ServiceRole      string
	InstanceRole     string
	InstanceTypes    []string
}

// JobQueueConfig represents configuration of the AWS Batch
// Job Queue.
type JobQueueConfig struct {
	Name        string
	ComputeEnvs []string
}

// JobDefConfig represents configuration of the AWS Batch
// base Job Definition.
type JobDefConfig struct {
	Name   string
	Memory int64
	VCPUs  int64
}

// DefaultConfig returns default configuration of AWS.
func DefaultConfig() Config {
	return Config{
		Container:       "docker.io/adamstruck/funnel:aws-auth",
		DefaultTaskName: "funnel task",
		ComputeEnv: ComputeEnvConfig{
			Name:         "funnel-compute-environment",
			InstanceRole: "ecsInstanceRole",
			InstanceTypes: []string{
				"optimal",
			},
			MinVCPUs: 0,
			MaxVCPUs: 256,
			Tags: map[string]string{
				"Name": "Funnel",
			},
		},
		JobQueue: JobQueueConfig{
			Name: "funnel-job-queue",
			ComputeEnvs: []string{
				"funnel-compute-environment",
			},
		},
		JobDef: JobDefConfig{
			Name:   "funnel-job-def",
			Memory: 128,
			VCPUs:  1,
		},
	}
}