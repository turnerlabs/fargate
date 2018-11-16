package ecs

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
)

func TestSortEnvVars(t *testing.T) {
	sess := session.Must(session.NewSession())
	ecs := New(sess, "my-app-dev")

	result := ecs.SortEnvVars([]EnvVar{
		EnvVar{Key: "PORT", Value: "8080"},
		EnvVar{Key: "ENVIRONMENT", Value: "prod"},
		EnvVar{Key: "PRODUCT", Value: "my-app-prod"},
		EnvVar{Key: "ENABLE_LOGGING", Value: "false"},
		EnvVar{Key: "HEALTHCHECK", Value: "/hc"},
	})

	sorted := []EnvVar{
		EnvVar{Key: "ENABLE_LOGGING", Value: "false"},
		EnvVar{Key: "ENVIRONMENT", Value: "prod"},
		EnvVar{Key: "HEALTHCHECK", Value: "/hc"},
		EnvVar{Key: "PORT", Value: "8080"},
		EnvVar{Key: "PRODUCT", Value: "my-app-prod"},
	}

	for i, env := range result {
		expected := sorted[i].Key
		got := env.Key

		if got != expected {
			t.Errorf("Expected %s, got %s", expected, got)
		}
	}
}

func TestGetTaskFamily(t *testing.T) {
	var got, expected string

	sess := session.Must(session.NewSession())
	ecs := New(sess, "my-app-dev")

	got = ecs.GetTaskFamily("arn:aws:ecs:us-east-1:000000000000:task-definition/my-app-dev:25")
	expected = "my-app-dev"
	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}

	got = ecs.GetTaskFamily("arn:aws:ecs:us-east-1:000000000000:task-definition/app-prod:2")
	expected = "app-prod"
	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}

func TestResolveRevisionNumber_Absolute(t *testing.T) {
	sess := session.Must(session.NewSession())
	ecs := New(sess, "my-app-dev")
	taskDefinitionArn := "arn:aws:ecs:us-east-1:000000000000:task-definition/my-app-dev:25"

	if ecs.ResolveRevisionNumber(taskDefinitionArn, "12") != "12" {
		t.Error("Expected 12")
	}

	if ecs.ResolveRevisionNumber(taskDefinitionArn, "37") != "37" {
		t.Error("Expected 37")
	}
}

func TestResolveRevisionNumber_NegativeExpression(t *testing.T) {
	sess := session.Must(session.NewSession())
	ecs := New(sess, "my-app-dev")
	taskDefinitionArn := "arn:aws:ecs:us-east-1:000000000000:task-definition/my-app-dev:50"

	if ecs.ResolveRevisionNumber(taskDefinitionArn, "-1") != "49" {
		t.Error("Expected 49")
	}

	if ecs.ResolveRevisionNumber(taskDefinitionArn, "-10") != "40" {
		t.Error("Expected 40")
	}
}
func TestResolveRevisionNumber_PositiveExpression(t *testing.T) {
	sess := session.Must(session.NewSession())
	ecs := New(sess, "my-app-dev")
	taskDefinitionArn := "arn:aws:ecs:us-east-1:000000000000:task-definition/my-app-dev:20"

	if ecs.ResolveRevisionNumber(taskDefinitionArn, "+1") != "21" {
		t.Error("Expected 21")
	}

	if ecs.ResolveRevisionNumber(taskDefinitionArn, "+33") != "53" {
		t.Error("Expected 53")
	}
}
func TestResolveRevisionNumber_NoInput(t *testing.T) {
	sess := session.Must(session.NewSession())
	ecs := New(sess, "my-app-dev")

	if ecs.ResolveRevisionNumber("arn:aws:ecs:us-east-1:000000000000:task-definition/my-app-dev:5", "") != "5" {
		t.Error("Expected 5")
	}

	if ecs.ResolveRevisionNumber("arn:aws:ecs:us-east-1:000000000000:task-definition/my-app-dev:12", "") != "12" {
		t.Error("Expected 12")
	}

}
func TestResolveRevisionNumber_InvalidInput(t *testing.T) {
	sess := session.Must(session.NewSession())
	ecs := New(sess, "my-app-dev")
	taskDefinitionArn := "arn:aws:ecs:us-east-1:000000000000:task-definition/my-app-dev:2"

	if ecs.ResolveRevisionNumber(taskDefinitionArn, "q") != "" {
		t.Error("Expected empty string")
	}

	if ecs.ResolveRevisionNumber(taskDefinitionArn, "-10") != "" {
		t.Error("Expected empty string")
	}
}
