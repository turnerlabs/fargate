package ecs

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awsecs "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turnerlabs/fargate/console"
)

const logStreamPrefix = "fargate"

var taskDefinitionCache = make(map[string]*awsecs.TaskDefinition)

type CreateTaskDefinitionInput struct {
	Cpu              string
	EnvVars          []EnvVar
	ExecutionRoleArn string
	Image            string
	Memory           string
	Name             string
	Port             int64
	LogGroupName     string
	LogRegion        string
	TaskRole         string
	Type             string
}

type EnvVar struct {
	Key   string
	Value string
}

type envSorter []EnvVar

func (a envSorter) Len() int {
	return len(a)
}
func (a envSorter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a envSorter) Less(i, j int) bool {
	return a[i].Key < a[j].Key
}

func (ecs *ECS) CreateTaskDefinition(input *CreateTaskDefinitionInput) string {
	console.Debug("Creating ECS task definition")

	logConfiguration := &awsecs.LogConfiguration{
		LogDriver: aws.String(awsecs.LogDriverAwslogs),
		Options: map[string]*string{
			"awslogs-region":        aws.String(input.LogRegion),
			"awslogs-group":         aws.String(input.LogGroupName),
			"awslogs-stream-prefix": aws.String(logStreamPrefix),
		},
	}

	containerDefinition := &awsecs.ContainerDefinition{
		Environment:      input.Environment(),
		Essential:        aws.Bool(true),
		Image:            aws.String(input.Image),
		LogConfiguration: logConfiguration,
		Name:             aws.String(input.Name),
	}

	if input.Port != 0 {
		containerDefinition.SetPortMappings(
			[]*awsecs.PortMapping{
				&awsecs.PortMapping{
					ContainerPort: aws.Int64(int64(input.Port)),
				},
			},
		)
	}

	resp, err := ecs.svc.RegisterTaskDefinition(
		&awsecs.RegisterTaskDefinitionInput{
			ContainerDefinitions:    []*awsecs.ContainerDefinition{containerDefinition},
			Cpu:                     aws.String(input.Cpu),
			ExecutionRoleArn:        aws.String(input.ExecutionRoleArn),
			Family:                  aws.String(fmt.Sprintf("%s_%s", input.Type, input.Name)),
			Memory:                  aws.String(input.Memory),
			NetworkMode:             aws.String(awsecs.NetworkModeAwsvpc),
			RequiresCompatibilities: aws.StringSlice([]string{awsecs.CompatibilityFargate}),
			TaskRoleArn:             aws.String(input.TaskRole),
		},
	)

	if err != nil {
		console.ErrorExit(err, "Couldn't register ECS task definition")
	}

	td := resp.TaskDefinition

	console.Debug("Created ECS task definition [%s:%d]", aws.StringValue(td.Family), aws.Int64Value(td.Revision))

	return aws.StringValue(td.TaskDefinitionArn)
}

func (input *CreateTaskDefinitionInput) Environment() []*awsecs.KeyValuePair {
	return convertEnvVars(input.EnvVars)
}

func convertEnvVars(envvars []EnvVar) []*awsecs.KeyValuePair {
	var environment []*awsecs.KeyValuePair

	for _, envVar := range envvars {
		environment = append(environment,
			&awsecs.KeyValuePair{
				Name:  aws.String(envVar.Key),
				Value: aws.String(envVar.Value),
			},
		)
	}

	return environment

}

func (ecs *ECS) DescribeTaskDefinition(taskDefinitionArn string) *awsecs.TaskDefinition {
	if taskDefinitionCache[taskDefinitionArn] != nil {
		return taskDefinitionCache[taskDefinitionArn]
	}

	resp, err := ecs.svc.DescribeTaskDefinition(
		&awsecs.DescribeTaskDefinitionInput{
			TaskDefinition: aws.String(taskDefinitionArn),
		},
	)

	if err != nil {
		console.ErrorExit(err, "Could not describe ECS task definition")
	}

	taskDefinitionCache[taskDefinitionArn] = resp.TaskDefinition

	return taskDefinitionCache[taskDefinitionArn]
}

//UpdateTaskDefinitionImage registers a new task definition with the updated image
func (ecs *ECS) UpdateTaskDefinitionImage(taskDefinitionArn, image string) string {
	taskDefinition := ecs.DescribeTaskDefinition(taskDefinitionArn)
	taskDefinition.ContainerDefinitions[0].Image = aws.String(image)
	return ecs.registerTaskDefinition(taskDefinition)
}

// UpdateTaskDefinitionImageAndReplaceEnvVars creates a new, updated task definition
// based on the specified image and env vars.
// Note that any existing envvars are replaced by the new ones
func (ecs *ECS) UpdateTaskDefinitionImageAndEnvVars(taskDefinitionArnOrFamily string, image string, environmentVariables []EnvVar, replaceEnvVars bool) string {

	//fetch task definition details (for specific or latest active)
	taskDefinition := ecs.DescribeTaskDefinition(taskDefinitionArnOrFamily)

	//which container are we updating?
	container := taskDefinition.ContainerDefinitions[0]

	//update image if specified
	if image != "" {
		container.Image = aws.String(image)
	}

	//convert envvars to aws input format
	if len(environmentVariables) > 0 {
		envvars := convertEnvVars(environmentVariables)

		//is this a replace or add operation?
		if replaceEnvVars {
			container.Environment = envvars

		} else {
			for _, e := range envvars {
				container.Environment = append(container.Environment, e)
			}
		}
	}

	return ecs.registerTaskDefinition(taskDefinition)
}

//registers a new task definition based on a task definition struct
func (ecs *ECS) registerTaskDefinition(taskDefinition *awsecs.TaskDefinition) string {

	//register a new task definition
	resp, err := ecs.svc.RegisterTaskDefinition(
		&awsecs.RegisterTaskDefinitionInput{
			ContainerDefinitions:    taskDefinition.ContainerDefinitions,
			Cpu:                     taskDefinition.Cpu,
			ExecutionRoleArn:        taskDefinition.ExecutionRoleArn,
			Family:                  taskDefinition.Family,
			Memory:                  taskDefinition.Memory,
			NetworkMode:             taskDefinition.NetworkMode,
			RequiresCompatibilities: taskDefinition.RequiresCompatibilities,
			TaskRoleArn:             taskDefinition.TaskRoleArn,
		},
	)	

	if err != nil {
		console.ErrorExit(err, "Could not register ECS task definition")
	}

	return aws.StringValue(resp.TaskDefinition.TaskDefinitionArn)	
}

//AddEnvVarsToTaskDefinition registers a new task definition with the envvars appended
func (ecs *ECS) AddEnvVarsToTaskDefinition(taskDefinitionArn string, envVars []EnvVar) string {
	taskDefinition := ecs.DescribeTaskDefinition(taskDefinitionArn)

	for _, envVar := range envVars {
		keyValuePair := &awsecs.KeyValuePair{
			Name:  aws.String(envVar.Key),
			Value: aws.String(envVar.Value),
		}

		taskDefinition.ContainerDefinitions[0].Environment = append(
			taskDefinition.ContainerDefinitions[0].Environment,
			keyValuePair,
		)
	}

	return ecs.registerTaskDefinition(taskDefinition)
}

//RemoveEnvVarsFromTaskDefinition registers a new task definition with the specified keys removed
func (ecs *ECS) RemoveEnvVarsFromTaskDefinition(taskDefinitionArn string, keys []string) string {
	var newEnvironment []*awsecs.KeyValuePair

	//look up task definition
	taskDefinition := ecs.DescribeTaskDefinition(taskDefinitionArn)
	environment := taskDefinition.ContainerDefinitions[0].Environment

	//iterate existing envvars
	for _, keyValuePair := range environment {

		//is this key a match to remove?
		match := false
		for _, key := range keys {
			if aws.StringValue(keyValuePair.Name) == key {
				match = true
				break
			}
		}

		//add this envvar since it wasn't a match to remove
		if !match {
			newEnvironment = append(newEnvironment, keyValuePair)
		}
	}

	taskDefinition.ContainerDefinitions[0].Environment = newEnvironment

	return ecs.registerTaskDefinition(taskDefinition)
}

//GetEnvVarsFromTaskDefinition retrieves envvars from an existing task definition
func (ecs *ECS) GetEnvVarsFromTaskDefinition(taskDefinitionArn string) []EnvVar {
	var envVars []EnvVar

	taskDefinition := ecs.DescribeTaskDefinition(taskDefinitionArn)

	for _, keyValuePair := range taskDefinition.ContainerDefinitions[0].Environment {
		envVars = append(envVars,
			EnvVar{
				Key:   aws.StringValue(keyValuePair.Name),
				Value: aws.StringValue(keyValuePair.Value),
			},
		)
	}

	return envVars
}

//UpdateTaskDefinitionCpuAndMemory registers a new task definition with the cpu/memory
func (ecs *ECS) UpdateTaskDefinitionCpuAndMemory(taskDefinitionArn, cpu, memory string) string {
	taskDefinition := ecs.DescribeTaskDefinition(taskDefinitionArn)

	if cpu != "" {
		taskDefinition.Cpu = aws.String(cpu)
	}

	if memory != "" {
		taskDefinition.Memory = aws.String(memory)
	}

	return ecs.registerTaskDefinition(taskDefinition)
}

//GetRevisionNumber returns the revision number from a task definition
func (ecs *ECS) GetRevisionNumber(taskDefinitionArn string) string {
	contents := strings.Split(taskDefinitionArn, ":")
	return contents[len(contents)-1]
}

//GetTaskDefinitionARN builds an ARN
func (ecs *ECS) GetTaskDefinitionARN(region string, account string, family string, revisionNumber string) string {
	return fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/%s:%s", region, account, family, revisionNumber)
}

//GetTaskFamily returns the task family from a task definition ARN
func (ecs *ECS) GetTaskFamily(taskDefinitionArn string) string {
	contents := strings.Split(taskDefinitionArn, ":")
	return strings.TrimPrefix(contents[len(contents)-2], "task-definition/")
}

//GetCpuAndMemoryFromTaskDefinition returns the cpu/memory from a task definition
func (ecs *ECS) GetCpuAndMemoryFromTaskDefinition(taskDefinitionArn string) (string, string) {
	taskDefinition := ecs.DescribeTaskDefinition(taskDefinitionArn)

	return aws.StringValue(taskDefinition.Cpu), aws.StringValue(taskDefinition.Memory)
}

//ResolveRevisionNumber returns a task defintion revision number by absolute value or expression
func (ecs *ECS) ResolveRevisionNumber(taskDefinitionArn string, revisionExpression string) string {
	currentRevision := ecs.GetRevisionNumber(taskDefinitionArn)
	currentRevisionNumber, err := strconv.ParseInt(currentRevision, 10, 64)

	if err != nil {
		return ""
	}

	if revisionExpression == "" {
		return currentRevision
	}

	var nextRevisionNumber int64

	// if not a delta assume absolute
	if revisionExpression[0] != '+' && revisionExpression[0] != '-' {
		if _, err := strconv.ParseInt(revisionExpression, 10, 64); err != nil {
			return ""
		}

		return revisionExpression
	}

	if s, err := strconv.ParseInt(revisionExpression[1:len(revisionExpression)], 10, 64); err == nil {
		if revisionExpression[0] == '+' {
			nextRevisionNumber = currentRevisionNumber + s
		} else if revisionExpression[0] == '-' {
			nextRevisionNumber = currentRevisionNumber - s
		}
	}

	if nextRevisionNumber <= 0 {
		return ""
	}

	result := strconv.FormatInt(nextRevisionNumber, 10)

	return result
}

// SortEnvVars sorts a slice of EnvVar's by Key
func (ecs *ECS) SortEnvVars(envVars []EnvVar) []EnvVar {
	sort.Sort(envSorter(envVars))
	return envVars
}
