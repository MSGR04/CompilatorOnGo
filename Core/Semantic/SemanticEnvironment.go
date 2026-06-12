package Semantic

type variableInfo struct {
	initialized bool
	used        bool
	valueType   ValueType
}

type functionInfo struct {
	used       bool
	paramCount int
	returnType ValueType
}

type symbolKind int

const (
	symbolMissing symbolKind = iota
	symbolVariable
	symbolFunction
)

type SemanticEnvironment struct {
	parent    *SemanticEnvironment
	vars      map[string]*variableInfo
	functions map[string]*functionInfo
}

func NewSemanticEnvironment(parent *SemanticEnvironment) *SemanticEnvironment {
	return &SemanticEnvironment{
		parent:    parent,
		vars:      make(map[string]*variableInfo),
		functions: make(map[string]*functionInfo),
	}
}

func (e *SemanticEnvironment) DefineVariable(name string, valueType ValueType, initialized bool) bool {
	if _, ok := e.vars[name]; ok {
		return false
	}
	if _, ok := e.functions[name]; ok {
		return false
	}

	e.vars[name] = &variableInfo{
		initialized: initialized,
		used:        false,
		valueType:   valueType,
	}
	return true
}

func (e *SemanticEnvironment) DefineFunction(name string, paramCount int) (*functionInfo, bool) {
	if _, ok := e.functions[name]; ok {
		return nil, false
	}
	if _, ok := e.vars[name]; ok {
		return nil, false
	}

	info := &functionInfo{
		used:       false,
		paramCount: paramCount,
		returnType: UnknownType,
	}
	e.functions[name] = info
	return info, true
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

func (e *SemanticEnvironment) UseFunction(name string) (*functionInfo, bool) {
	if info, ok := e.functions[name]; ok {
		info.used = true
		return info, true
	}
	if e.parent != nil {
		return e.parent.UseFunction(name)
	}
	return nil, false
}

func (e *SemanticEnvironment) ResolveFunction(name string) (*functionInfo, bool) {
	if info, ok := e.functions[name]; ok {
		return info, true
	}
	if e.parent != nil {
		return e.parent.ResolveFunction(name)
	}
	return nil, false
}

func (e *SemanticEnvironment) ResolveName(name string) (symbolKind, *variableInfo, *functionInfo) {
	if variable, ok := e.vars[name]; ok {
		return symbolVariable, variable, nil
	}
	if function, ok := e.functions[name]; ok {
		return symbolFunction, nil, function
	}
	if e.parent != nil {
		return e.parent.ResolveName(name)
	}
	return symbolMissing, nil, nil
}

func (e *SemanticEnvironment) CollectUnused() []string {
	unused := make([]string, 0)

	for name, info := range e.vars {
		if !info.used {
			unused = append(unused, name)
		}
	}
	for name, info := range e.functions {
		if !info.used {
			unused = append(unused, name)
		}
	}

	return unused
}
