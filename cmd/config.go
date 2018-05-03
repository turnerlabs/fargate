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
	keyVerbose = "verbose"
	keyNoColor = "nocolor"
)

//configure viper to manage parameter input
func initConfig(cmd *cobra.Command) {

	//config file
	viper.SetConfigName("fargate")
	viper.AddConfigPath("./")
	viper.ReadInConfig()

	//env vars
	viper.BindEnv(keyCluster, "FARGATE_CLUSTER")
	viper.BindEnv(keyService, "FARGATE_SERVICE")
	viper.BindEnv(keyVerbose, "FARGATE_VERBOSE")
	viper.BindEnv(keyNoColor, "FARGATE_NOCOLOR")

	//cli arg
	initPFlag(keyCluster, cmd)
	initPFlag(keyVerbose, cmd)
	initPFlag(keyNoColor, cmd)
}

func initPFlag(key string, cmd *cobra.Command) {
	viper.BindPFlag(key, cmd.PersistentFlags().Lookup(key))
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

func getVerbose() bool {
	return viper.GetBool(keyVerbose)
}

func getNoColor() bool {
	return viper.GetBool(keyNoColor)
}
