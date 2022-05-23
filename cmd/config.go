package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyCluster = "cluster"
	keyService = "service"
	keyRegion  = "region"
	keyVerbose = "verbose"
	keyNoColor = "nocolor"
	keyTask    = "task"
	keyRule    = "rule"
)

//configure viper to manage parameter input
func initConfig(cmd *cobra.Command) {

	//config file (fargate.yml)
	viper.SetConfigName("fargate")
	viper.AddConfigPath("./")
	viper.ReadInConfig()

	//env vars
	viper.BindEnv(keyCluster, "FARGATE_CLUSTER")
	viper.BindEnv(keyService, "FARGATE_SERVICE")
	viper.BindEnv(keyVerbose, "FARGATE_VERBOSE")
	viper.BindEnv(keyNoColor, "FARGATE_NOCOLOR")
	viper.BindEnv(keyTask, "FARGATE_TASK")
	viper.BindEnv(keyRule, "FARGATE_RULE")

	//cli arg
	initPFlag(keyCluster, cmd)
	initPFlag(keyVerbose, cmd)
	initPFlag(keyRegion, cmd)
	initPFlag(keyNoColor, cmd)
}

func initPFlag(key string, cmd *cobra.Command) {
	viper.BindPFlag(key, cmd.PersistentFlags().Lookup(key))
}

//region can come from fargate.yml, AWS_REGION, AWS_DEFAULT_REGION or --region
func getRegion() string {
	result := viper.GetString(keyRegion)
	if result == "" {
		envAwsDefaultRegion := os.Getenv("AWS_DEFAULT_REGION")
		envAwsRegion := os.Getenv("AWS_REGION")

		if envAwsDefaultRegion != "" {
			result = envAwsDefaultRegion
		} else if envAwsRegion != "" {
			result = envAwsRegion
		} else {
			result = defaultRegion
		}
	}

	return result
}

//cluster can come from fargate.yml, FARGATE_CLUSTER envar, or --cluster cli arg
func getClusterName() string {
	result := viper.GetString(keyCluster)
	if result == "" {
		fmt.Println("please specify cluster using: fargate.yml, FARGATE_SERVICE envvar, or --cluster")
		os.Exit(-1)
	}
	return result
}

//service can come from fargate.yml, FARGATE_SERVICE, or --service cli arg
func getServiceName() string {
	result := viper.GetString(keyService)
	if result == "" {
		fmt.Println("please specify service using: fargate.yml, FARGATE_SERVICE envvar, or --service")
		os.Exit(-1)
	}
	return result
}

//task can come from fargate.yml, FARGATE_TASK, or --task cli arg
func getTaskName() string {
	result := viper.GetString(keyTask)
	if result == "" {
		fmt.Println("please specify task family using: fargate.yml, FARGATE_TASK envvar, or --task")
		os.Exit(-1)
	}
	return result
}

//rule can come from fargate.yml, FARGATE_RULE, or --task cli arg
func getRuleName() string {
	result := viper.GetString(keyRule)
	if result == "" {
		fmt.Println("please specify rule using: fargate.yml, FARGATE_RULE envvar, or --task")
		os.Exit(-1)
	}
	return result
}

func getVerbose() bool {
	return viper.GetBool(keyVerbose)
}

func getNoColor() bool {
	return viper.GetBool(keyNoColor)
}
