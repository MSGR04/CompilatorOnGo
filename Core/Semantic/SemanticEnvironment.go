package Semantic

type variableInfo struct {
	initialized bool
	used        bool
}

type SemanticEnvironment struct {
	parent *SemanticEnvironment
	vars   map[string]*variableInfo
}

func NewSemanticEnvironment(parent *SemanticEnvironment) *SemanticEnvironment {
	return &SemanticEnvironment{
		parent: parent,
		vars:   make(map[string]*variableInfo),
	}
}

func (e *SemanticEnvironment) DefineVariable(name string, initialized bool) bool {
	if _, ok := e.vars[name]; ok {
		return false
	}
	e.vars[name] = &variableInfo{initialized: initialized, used: false}
	return true
}

func (e *SemanticEnvironment) UseVariable(name string) (defined, initialized bool) {
	if info, ok := e.vars[name]; ok {
		info.used = true
		return true, info.initialized
	}
	if e.parent != nil {
		return e.parent.UseVariable(name)
	}
	return false, false
}

func (e *SemanticEnvironment) AssignVariable(name string) bool {
	if info, ok := e.vars[name]; ok {
		info.initialized = true
		return true
	}
	if e.parent != nil {
		return e.parent.AssignVariable(name)
	}
	return false
}

func (e *SemanticEnvironment) CollectUnused() []string {
	var unused []string
	for name, info := range e.vars {
		if !info.used {
			unused = append(unused, name)
		}
	}
	return unused
}
