package config

import (
	"demo/internal/util"
	"os"
	"strings"
)

const (
	EnvironmentTest  = "test"  // Для тестов
	EnvironmentLocal = "local" // Из среды разработки
	EnvironmentDev   = "dev"   // Из контейнера
	EnvironmentProd  = "prod"  // На проде
)

var (
	EnvironmentList = []string{
		EnvironmentTest,
		EnvironmentLocal,
		EnvironmentDev,
		EnvironmentProd,
	}
)

func GetEnvironment() string {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if !util.Contains(EnvironmentList, env) {
		env = EnvironmentProd
	}

	return env
}

func IsProduction() bool {
	return GetEnvironment() == EnvironmentProd
}
