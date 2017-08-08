package aws

import (
	"fmt"
)

func deploy(conf Config, region string) error {
	cli := newBatchSvc(conf)

	a, aerr := cli.CreateComputeEnvironment(region)
	fmt.Println(a, aerr)

	b, berr := cli.CreateJobQueue(region)
	fmt.Println(b, berr)

	c, cerr := cli.CreateJobDef(region)
	fmt.Println(c, cerr)

	return nil
}
