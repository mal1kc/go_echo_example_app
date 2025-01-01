package app

import (
	"fmt"
	"os"
	"prod_app/common"
)

type ConfigurationManager struct {
	PostgreSqlConnectionString string
}

func NewConfigurationManager() *ConfigurationManager {
	pqsqlConnectionString := getPostgreSqlConfig()
	return &ConfigurationManager{
		PostgreSqlConnectionString: pqsqlConnectionString,
	}
}

func getPostgreSqlConfig() string {
	// TODO: first try to get data fron env then try to get data from config file

	envVar := os.Getenv(string(common.PqSqlEnv))
	if len(envVar) != 0 {
		return envVar
	}
	panic(fmt.Sprintf("envVar:%s data is invalid", common.PqSqlEnv))
}
