package workflow

type Step interface {
	ID() string
	Name() string
	Run(ctx Context) (map[string]any, error)
	WithID(id string) Step
	WithInput(key string, outputKey string) Step
}

type ConditionalStep interface {
	ShouldRun(ctx Context) (bool, error)
}

type step struct {
	id     string
	name   string
	inputs map[string]string
	run    func(ctx Context) (map[string]any, error)
}

func (s step) ID() string {
	return s.id
}

func (s step) WithID(id string) Step {
	s.id = id
	return s
}

func (s step) WithInput(key string, outputKey string) Step {
	m := make(map[string]string)
	for k, v := range s.inputs {
		m[k] = v
	}
	m[key] = outputKey
	s.inputs = m
	return s
}

func (s step) Name() string {
	return s.name
}

func (s step) Run(ctx Context) (map[string]any, error) {
	return s.run(ctx)
}

func StepFunc(id string, name string, run func(ctx Context) (map[string]any, error)) Step {
	return step{
		id:   id,
		name: name,
		run:  run,
	}
}

type conditionalStep struct {
	Step
	shouldRun func(ctx Context) (bool, error)
}

func (s conditionalStep) ShouldRun(ctx Context) (bool, error) {
	return s.shouldRun(ctx)
}
