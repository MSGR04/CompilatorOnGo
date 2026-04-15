package Semantic

type variableInfo struct {
	initialized bool
	used        bool
	valueType   ValueType
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

func (e *SemanticEnvironment) DefineVariable(name string, valueType ValueType, initialized bool) bool {
	if _, ok := e.vars[name]; ok {
		return false
	}

	e.vars[name] = &variableInfo{
		initialized: initialized,
		used:        false,
		valueType:   valueType,
	}

	return true
}

func (e *SemanticEnvironment) UseVariable(name string) (*variableInfo, bool) {
	if info, ok := e.vars[name]; ok {
		info.used = true
		return info, true
	}

	if e.parent != nil {
		return e.parent.UseVariable(name)
	}

	return nil, false
}

func (e *SemanticEnvironment) ResolveVariable(name string) (*variableInfo, bool) {
	if info, ok := e.vars[name]; ok {
		return info, true
	}

	if e.parent != nil {
		return e.parent.ResolveVariable(name)
	}

	return nil, false
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
