package cloudwatchevents

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/turnerlabs/fargate/console"
)

//CloudWatchEvents represents the cloudwatchevents api
type CloudWatchEvents struct {
	svc *cloudwatchevents.CloudWatchEvents
}

//New creates a new cloudwatchevents client
func New(sess *session.Session) CloudWatchEvents {
	return CloudWatchEvents{
		svc: cloudwatchevents.New(sess),
	}
}

//UpdateTargetRevision updates a cloudwatch events rule target
func (c *CloudWatchEvents) UpdateTargetRevision(rule string, revision string) {

	//fetch existing rule targets
	input := &cloudwatchevents.ListTargetsByRuleInput{
		Rule: &rule,
	}
	resp, err := c.svc.ListTargetsByRule(input)
	if err != nil {
		console.ErrorExit(err, "ListTargetsByRuleInput failed")
	}
	if len(resp.Targets) == 0 {
		console.IssueExit("no targets found for rule: ", rule)
	}

	//update the task definition revision
	resp.Targets[0].EcsParameters.TaskDefinitionArn = &revision

	//update the rule target
	putInput := &cloudwatchevents.PutTargetsInput{
		Rule:    &rule,
		Targets: resp.Targets,
	}
	putResp, err := c.svc.PutTargets(putInput)
	if err != nil {
		console.ErrorExit(err, "PutTargets failed")
	}
	if putResp != nil {
		if *putResp.FailedEntryCount != 0 && len(putResp.FailedEntries) != 0 {
			for _, entry := range putResp.FailedEntries {
				fmt.Printf("TargetId: %s; ErrorCode: %s; ErrorMessage: %s", *entry.TargetId, *entry.ErrorCode, *entry.ErrorMessage)
				fmt.Println()
			}
			console.IssueExit("PutTargets failed")
		}
	}
}
