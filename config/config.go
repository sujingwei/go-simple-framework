package config

import (
	"github.com/sujingwei/go-simple-framework/configuration"
	webframework "github.com/sujingwei/go-simple-framework/web-framework"
)

type AppConfigure struct {
	configuration.App `yaml: "app"`
	Web               webframework.WebConfig `yaml: "web`
}
