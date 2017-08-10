---
title: AWS Batch

menu:
  main:
    parent: guides
    weight: 20
---

# AWS Batch

This guide covers setting up a Funnel server to submit tasks via the [AWS Batch][batch]
API. You will need to have an AWS account and the AWS Command Line Interface (CLI). 

To install the AWS CLI run:

```bash
pip install awscli
```

You will also need to [configure][aws-cli] the client by running:

```bash
aws configure
```

## Start Server

Start a sever that submits to AWS batch.

```bash
funnel aws proxy
```

## Setup

To configure AWS Batch for your account for use with Funnel run:

```bash
funnel aws delploy us-west-2
```

This will create a [ComputeEnviroment][1], [JobQueue][2], and [JobDefinition][3]
in the region you sepcifiy. 

## Authentication

The server uses HTTP Basic Authentication for all requests when the AWS Batch backend is active.

Set the `FUNNEL_KEY` and `FUNNEL_SECRET` environment variables to your aws_access_key_id  and aws_secret_access_key.

```bash
$ export FUNNEL_KEY= ASFDGSFDHGEFWFAKE
$ export FUNNEL_SECRET=E1ADFAfca/4r/YgfaGDSAD2+FAKE
$ funnel task get 110ec58a-a0f2-4ac4-8393-c866d813b8d1
```

## Submitting jobs

Specify the compute region in your task message in the [Zones][4] field. 
This field is required for this backend for scoping requests to AWS batch. 

```bash
funnel task create <task-message>.json
```

## Limitations

[aws-cli]: http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html
[batch]: https://docs.aws.amazon.com/batch/latest/APIReference/Welcome.html
[1]: https://docs.aws.amazon.com/batch/latest/APIReference/API_ComputeEnvironmentDetail.html
[2]: https://docs.aws.amazon.com/batch/latest/APIReference/API_JobQueueDetail.html
[3]: https://docs.aws.amazon.com/batch/latest/APIReference/API_JobDefinition.html
[4]: https://github.com/ga4gh/task-execution-schemas/blob/v0.2/task_execution.proto#L208
