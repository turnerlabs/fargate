package cmd

import (
	"github.com/turnerlabs/fargate/console"
	"github.com/turnerlabs/fargate/ecs"
	"github.com/turnerlabs/fargate/sts"

	"github.com/spf13/cobra"
	"github.com/turnerlabs/fargate/cloudwatchevents"
)

var flagEventsTargetRule string
var flagEventsTargetRevision string

//represents an events target operation
type eventsTargetOperation struct {
	Cluster  string
	Task     string
	Rule     string
	Revision string
	Region   string
}

func (o *eventsTargetOperation) validate() {
	if o.Revision == "" {
		console.IssueExit("--revision is required")
	}
}

var eventsTargetCmd = &cobra.Command{
	Use:   "target",
	Short: "Updates an event rule target to run a particular task definition revision.",
	Example: `
fargate events target --revision <revision>
fargate events target --rule <rule> --revision <revision>
	`,
	Run: func(cmd *cobra.Command, args []string) {

		operation := eventsTargetOperation{
			Cluster:  getClusterName(),
			Task:     getTaskName(),
			Rule:     getRuleName(),
			Revision: flagEventsTargetRevision,
			Region:   region,
		}

		operation.validate()
		eventsTarget(operation)
	},
}

func init() {
	eventsTargetCmd.PersistentFlags().StringVarP(&flagEventsTargetRule, "rule", "", "", `CloudWatch Events Rule`)

	eventsTargetCmd.PersistentFlags().StringVarP(&flagEventsTargetRevision, "revision", "r", "", `Task Definition Revision Number`)

	initPFlag(keyRule, eventsTargetCmd)

	eventsCmd.AddCommand(eventsTargetCmd)
}

func eventsTarget(op eventsTargetOperation) {
	events := cloudwatchevents.New(sess)
	ecs := ecs.New(sess, op.Cluster)
	sts := sts.New(sess)

	//look up account
	account := sts.GetCallerIdentity().Account

	//build full task definiton arn with revision
	tdARN := ecs.GetTaskDefinitionARN(op.Region, account, op.Task, op.Revision)

	//update target
	events.UpdateTargetRevision(op.Rule, tdARN)

	console.Info("rule %v now targeting revision %v", op.Rule, op.Revision)
}
