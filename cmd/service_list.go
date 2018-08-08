package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/turnerlabs/fargate/console"
	ECS "github.com/turnerlabs/fargate/ecs"
	ELBV2 "github.com/turnerlabs/fargate/elbv2"
	"github.com/spf13/cobra"
)

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List services",
	Run: func(cmd *cobra.Command, args []string) {
		listServices()
	},
}

func init() {
	serviceCmd.AddCommand(serviceListCmd)
}

func listServices() {
	var targetGroupArns []string
	var loadBalancerArns []string

	targetGroups := make(map[string]ELBV2.TargetGroup)
	loadBalancers := make(map[string]ELBV2.LoadBalancer)

	ecs := ECS.New(sess, getClusterName())
	elbv2 := ELBV2.New(sess)
	services := ecs.ListServices()

	for _, service := range services {
		if service.TargetGroupArn != "" {
			targetGroupArns = append(targetGroupArns, service.TargetGroupArn)
		}
	}

	if len(targetGroupArns) > 0 {
		for _, targetGroup := range elbv2.DescribeTargetGroups(targetGroupArns) {
			targetGroups[targetGroup.Arn] = targetGroup

			if targetGroup.LoadBalancerARN != "" {
				loadBalancerArns = append(loadBalancerArns, targetGroup.LoadBalancerARN)
			}
		}
	}

	if len(loadBalancerArns) > 0 {
		lbs, _ := elbv2.DescribeLoadBalancersByARN(loadBalancerArns)
		for _, loadBalancer := range lbs {
			loadBalancers[loadBalancer.ARN] = loadBalancer
		}
	}

	if len(services) > 0 {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 1, '\t', 0)
		fmt.Fprintln(w, "NAME\tIMAGE\tCPU\tMEMORY\tLOAD BALANCER\tDESIRED\tRUNNING\tPENDING\t")

		for _, service := range services {
			var loadBalancer string

			if service.TargetGroupArn != "" {
				tg := targetGroups[service.TargetGroupArn]
				lb := loadBalancers[tg.LoadBalancerARN]

				loadBalancer = lb.Name
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\t%d\t%d\t\n",
				service.Name,
				service.Image,
				service.Cpu,
				service.Memory,
				loadBalancer,
				service.DesiredCount,
				service.RunningCount,
				service.PendingCount,
			)
		}

		w.Flush()
	} else {
		console.Info("No services found")
	}
}
