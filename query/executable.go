package query

type (
	FunctionStatement struct {
		Name           string
		InputVariables []any

		IsPartOfQuery bool
		IsTVF         bool
	}

	ProcedureStatement struct {
		Name           string
		InputVariables []any
		OutputVarNames []string
	}
)

func Function(name string, args ...any) FunctionStatement {
	return FunctionStatement{
		Name:           name,
		InputVariables: args,
		IsPartOfQuery:  false,
		IsTVF:          false,
	}
}

func (f FunctionStatement) String() string {
	functionString := f.Name + "("
	for i := range f.InputVariables {
		functionString += "?"
		if i < len(f.InputVariables)-1 {
			functionString += ", "
		}
	}
	functionString += ")"
	return functionString
}

func (f FunctionStatement) Parameters() []any {
	return f.InputVariables
}

// Procedure
const procedureTemplate = `
`

func Procedure(name string) *ProcedureStatement {
	return &ProcedureStatement{
		Name: name,
	}
}

func (p *ProcedureStatement) Input(variables ...any) *ProcedureStatement {
	p.InputVariables = variables
	return p
}

func (p *ProcedureStatement) OutputNames(variables ...string) *ProcedureStatement {
	p.OutputVarNames = variables
	return p
}

func (p ProcedureStatement) Parameters() []any {
	return p.InputVariables
}
