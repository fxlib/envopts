package example

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

//go:generate go run ../. -type=FooEnv
//go:generate go run ../. -type=BarEnv
//go:generate go run ../. -type=RabEnv -optionType=RabOption
//go:generate go run ../. -type=fooEnv -optionType=PrivOption -output=private_foo_env.go

type FooEnv struct {
	Home         string        `env:"HOME"`
	Port         int           `env:"PORT" envDefault:"3000"`
	Password     string        `env:"PASSWORD,unset"`
	IsProduction bool          `env:"PRODUCTION"`
	Hosts        []string      `env:"HOSTS" envSeparator:":"`
	Duration     time.Duration `env:"DURATION"`
	TempFolder   string        `env:"TEMP_FOLDER" envDefault:"${HOME}/tmp" envExpand:"true"`

	// test external package and one tag for multiple fields
	Foo, Dar []env.Options `env:"FOO"`
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
