package example

// FooService behaves differently based on configuration options
type FooService struct {
	cfg FooEnv
}

// NewFooService inits the FooService while taking options
func NewFooService(opts ...Option) (s *FooService) {
	s = &FooService{
		// set default values, overwritten by any explicitely configured options
		cfg: ApplyOptions(opts...),
	}
	return
}
