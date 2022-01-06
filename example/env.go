package example

import (
	"fmt"
	"time"
)

//go:generate go run ../main.go -type=FooEnv
//go:generate go run ../main.go -type=BarEnv
//go:generate go run ../main.go -type=RabEnv -optionType=RabOption
//go:generate go run ../main.go -type=fooEnv -optionType=PrivOption -output=private_foo_env.go

type FooEnv struct {
	Home         string        `env:"HOME"`
	Port         int           `env:"PORT" envDefault:"3000"`
	Password     string        `env:"PASSWORD,unset"`
	IsProduction bool          `env:"PRODUCTION"`
	Hosts        []string      `env:"HOSTS" envSeparator:":"`
	Duration     time.Duration `env:"DURATION"`
	TempFolder   string        `env:"TEMP_FOLDER" envDefault:"${HOME}/tmp" envExpand:"true"`
}

// BarEnv is one with a private field should result in nothing being generated
type BarEnv struct {
	bar string `env:"FOO"`
}

type RabEnv struct {
	Home string `env:"HOME"`
}

// private version of Foo
type fooEnv struct {
	Hosts    []string      `env:"HOSTS" envSeparator:":"`
	Duration time.Duration `env:"DURATION"`
}

// Run make the warnings go away
func RunBar(be BarEnv) {
	_ = fooEnv{}
	fmt.Println(be.bar)
}
