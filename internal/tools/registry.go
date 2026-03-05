package tools

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {

	return &Registry{
		tools: make(map[string]Tool),
	}
}

func (r *Registry) Register(t Tool) {

	r.tools[t.Name()] = t
}

func (r *Registry) Get(name string) Tool {

	return r.tools[name]
}
