package elbv2

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awselbv2 "github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/turnerlabs/fargate/console"
)

// Listener accepts incoming traffic on a load balancer based upon the provided routing rules.
type Listener struct {
	ARN             string
	CertificateARNs []string
	Port            int64
	Protocol        string
	Rules           []Rule
}

// String returns a friendly representation of the listener.
func (l Listener) String() string {
	return fmt.Sprintf("%s:%d", l.Protocol, l.Port)
}

// Listeners is a collection of listeners.
type Listeners []Listener

// Listeners returns a comma-separated friendly representation of the listeners.
func (l Listeners) String() string {
	var listenerStrings []string

	for _, listener := range l {
		listenerStrings = append(listenerStrings, listener.String())
	}

	return strings.Join(listenerStrings, ", ")
}

// Rule defines a routing rule defining how traffic should be routed to a listener.
type Rule struct {
	ARN            string
	IsDefault      bool
	Priority       int
	TargetGroupARN string
	Type           string
	Value          string
}

// String returns a friendly representation of a rule.
func (r Rule) String() string {
	return strings.Join([]string{r.Type, r.Value}, "=")
}

// CreateListenerParameters are the parameters required to create a new listener.
type CreateListenerParameters struct {
	CertificateARNs       []string
	DefaultTargetGroupARN string
	LoadBalancerARN       string
	Port                  int64
	Protocol              string
}

// SetCertificateARNs sets the certificate ARNs with the given ARNs.
func (input *CreateListenerParameters) SetCertificateARNs(arns []string) {
	input.CertificateARNs = arns
}

// CreateListener creates a new listener and returns the listener ARN if successfully created.
func (elbv2 SDKClient) CreateListener(p CreateListenerParameters) (string, error) {
	action := &awselbv2.Action{
		TargetGroupArn: aws.String(p.DefaultTargetGroupARN),
		Type:           aws.String(awselbv2.ActionTypeEnumForward),
	}

	i := &awselbv2.CreateListenerInput{
		Port:            aws.Int64(p.Port),
		Protocol:        aws.String(p.Protocol),
		LoadBalancerArn: aws.String(p.LoadBalancerARN),
		DefaultActions:  []*awselbv2.Action{action},
	}

	if len(p.CertificateARNs) > 0 {
		certificates := []*awselbv2.Certificate{}

		for _, certificateARN := range p.CertificateARNs {
			certificates = append(certificates,
				&awselbv2.Certificate{
					CertificateArn: aws.String(certificateARN),
				},
			)
		}

		i.SetCertificates(certificates)
	}

	resp, err := elbv2.client.CreateListener(i)

	if err != nil {
		return "", err
	}

	return aws.StringValue(resp.Listeners[0].ListenerArn), nil
}

// DescribeListeners returns all of the listeners for a given load balancer ARN.
func (elbv2 SDKClient) DescribeListeners(lbARN string) (Listeners, error) {
	var listeners []Listener

	input := &awselbv2.DescribeListenersInput{
		LoadBalancerArn: aws.String(lbARN),
	}

	err := elbv2.client.DescribeListenersPages(
		input,
		func(resp *awselbv2.DescribeListenersOutput, lastPage bool) bool {
			for _, l := range resp.Listeners {
				listener := Listener{
					ARN:      aws.StringValue(l.ListenerArn),
					Port:     aws.Int64Value(l.Port),
					Protocol: aws.StringValue(l.Protocol),
				}

				for _, certificate := range l.Certificates {
					listener.CertificateARNs = append(listener.CertificateARNs, aws.StringValue(certificate.CertificateArn))
				}

				listeners = append(listeners, listener)
			}

			return true
		},
	)

	return listeners, err
}

func (elbv2 SDKClient) ModifyLoadBalancerDefaultAction(lbARN, targetGroupARN string) {
	for _, listener := range elbv2.GetListeners(lbARN) {
		elbv2.ModifyListenerDefaultAction(listener.ARN, targetGroupARN)
	}
}

func (elbv2 SDKClient) ModifyListenerDefaultAction(listenerARN, targetGroupARN string) {
	action := &awselbv2.Action{
		TargetGroupArn: aws.String(targetGroupARN),
		Type:           aws.String(awselbv2.ActionTypeEnumForward),
	}

	elbv2.client.ModifyListener(
		&awselbv2.ModifyListenerInput{
			ListenerArn:    aws.String(listenerARN),
			DefaultActions: []*awselbv2.Action{action},
		},
	)
}

func (elbv2 SDKClient) AddRule(lbARN, targetGroupARN string, rule Rule) {
	console.Debug("Adding ELB listener rule [%s=%s]", rule.Type, rule.Value)

	listeners := elbv2.GetListeners(lbARN)

	for _, listener := range listeners {
		elbv2.AddRuleToListener(listener.ARN, targetGroupARN, rule)
	}
}

func (elbv2 SDKClient) AddRuleToListener(listenerARN, targetGroupARN string, rule Rule) {
	var ruleType string

	if rule.Type == "HOST" {
		ruleType = "host-header"
	} else {
		ruleType = "path-pattern"
	}

	ruleCondition := &awselbv2.RuleCondition{
		Field:  aws.String(ruleType),
		Values: aws.StringSlice([]string{rule.Value}),
	}
	highestPriority := elbv2.GetHighestPriorityFromListener(listenerARN)
	priority := highestPriority + 10
	action := &awselbv2.Action{
		TargetGroupArn: aws.String(targetGroupARN),
		Type:           aws.String(awselbv2.ActionTypeEnumForward),
	}

	elbv2.client.CreateRule(
		&awselbv2.CreateRuleInput{
			Priority:    aws.Int64(priority),
			ListenerArn: aws.String(listenerARN),
			Actions:     []*awselbv2.Action{action},
			Conditions:  []*awselbv2.RuleCondition{ruleCondition},
		},
	)
}

func (elbv2 SDKClient) DescribeRules(listenerARN string) []Rule {
	var rules []Rule

	resp, err := elbv2.client.DescribeRules(
		&awselbv2.DescribeRulesInput{
			ListenerArn: aws.String(listenerARN),
		},
	)

	if err != nil {
		console.ErrorExit(err, "Could not describe ELB rules")
	}

	for _, r := range resp.Rules {
		for _, c := range r.Conditions {
			var field string

			switch aws.StringValue(c.Field) {
			case "host-header":
				field = "HOST"
			case "path-pattern":
				field = "PATH"
			}

			for _, v := range c.Values {
				priority, _ := strconv.Atoi(aws.StringValue(r.Priority))

				rule := Rule{
					ARN:            aws.StringValue(r.RuleArn),
					Priority:       priority,
					TargetGroupARN: aws.StringValue(r.Actions[0].TargetGroupArn),
					Type:           field,
					Value:          aws.StringValue(v),
				}

				rules = append(rules, rule)
			}
		}

		if aws.BoolValue(r.IsDefault) == true {
			rule := Rule{
				TargetGroupARN: aws.StringValue(r.Actions[0].TargetGroupArn),
				Type:           "DEFAULT",
				IsDefault:      true,
			}

			rules = append(rules, rule)
		}
	}

	return rules
}

func (elbv2 SDKClient) GetHighestPriorityFromListener(listenerARN string) int64 {
	var priorities []int

	resp, err := elbv2.client.DescribeRules(
		&awselbv2.DescribeRulesInput{
			ListenerArn: aws.String(listenerARN),
		},
	)

	if err != nil {
		console.ErrorExit(err, "Could not retrieve ELB listener rules")
	}

	for _, rule := range resp.Rules {
		priority, _ := strconv.Atoi(aws.StringValue(rule.Priority))
		priorities = append(priorities, priority)
	}

	sort.Ints(priorities)

	return int64(priorities[len(priorities)-1])
}

func (elbv2 SDKClient) GetListeners(lbARN string) []Listener {
	var listeners []Listener

	input := &awselbv2.DescribeListenersInput{
		LoadBalancerArn: aws.String(lbARN),
	}

	err := elbv2.client.DescribeListenersPages(
		input,
		func(resp *awselbv2.DescribeListenersOutput, lastPage bool) bool {
			for _, l := range resp.Listeners {
				listener := Listener{
					ARN:      aws.StringValue(l.ListenerArn),
					Port:     aws.Int64Value(l.Port),
					Protocol: aws.StringValue(l.Protocol),
				}

				for _, certificate := range l.Certificates {
					listener.CertificateARNs = append(
						listener.CertificateARNs,
						aws.StringValue(certificate.CertificateArn),
					)
				}

				listeners = append(listeners, listener)
			}

			return true
		},
	)

	if err != nil {
		console.ErrorExit(err, "Could not retrieve ELB listeners")
	}

	return listeners
}

func (elbv2 SDKClient) DeleteRule(ruleARN string) {
	_, err := elbv2.client.DeleteRule(
		&awselbv2.DeleteRuleInput{
			RuleArn: aws.String(ruleARN),
		},
	)

	if err != nil {
		console.ErrorExit(err, "Could not delete ELB rule")
	}
}
