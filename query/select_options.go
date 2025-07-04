package query

// _ "embed"

// //go:embed templates/select.tmpl
// var selectTemplate string

// const (
// 	selectTmpl  string = "SELECT %s"
// 	fromTmpl    string = "FROM %s"
// 	whereTmpl   string = "WHERE %s"
// 	orderbyTmpl string = "ORDER BY %s"
// 	groupbyTmpl string = "GROUP BY %s"
// 	limitTmpl   string = "LIMIT %d"
// 	offsetTmpl  string = "OFFSET %d"
// 	joinTmpl    string = "JOIN %s ON %s"
// 	onTmpl      string = "%s = %s"
// 	asTmpl      string = "AS %s"

// 	oneLineSelect string = selectTmpl + " " + fromTmpl + " " + whereTmpl + " " + orderbyTmpl + " " + groupbyTmpl + " " + limitTmpl + " " + offsetTmpl
// 	oneLineJoined string = oneLineSelect + " " + joinTmpl + " " + onTmpl + " " + asTmpl

// 	formatedSelect string = selectTmpl + "\n" + fromTmpl + "\n" + whereTmpl + "\n" + orderbyTmpl + "\n" + groupbyTmpl + "\n" + limitTmpl + "\n" + offsetTmpl
// 	formatedJoined string = formatedSelect + "\n" + joinTmpl + "\n" + onTmpl + "\n" + asTmpl
// )

// type SelectStatementFuncOpt struct {
// 	Distinct bool
// 	Fields   []string
// 	Table    string

// 	Joins []struct {
// 		Type  string
// 		Table string
// 		On    string
// 	}

// 	Where string

// 	GroupBy []string
// 	Having  string
// 	OrderBy []struct {
// 		Field string
// 		Desc  bool
// 	}

// 	Limit  *int
// 	Offset *int
// }

// // Parameters implements database.Query.
// func (s *SelectStatementFuncOpt) Parameters() []any {
// 	panic("unimplemented")
// }

// // Render implements database.Query.
// func (s *SelectStatementFuncOpt) Render() string {
// 	panic("unimplemented")
// }

// // Overrider implements database.Query.
// func (s *SelectStatementFuncOpt) Overrider(func(*database.Statement) *database.Statement) {
// 	panic("unimplemented")
// }

// func SelectWith(opts ...SelectOption) *SelectStatementFuncOpt {
// 	selectStatement := &SelectStatementFuncOpt{
// 		Fields: []string{"*"},
// 	}

// 	for _, opt := range opts {
// 		opt(selectStatement)
// 	}

// 	return selectStatement
// }

// // Formatted implements Query.
// func (s *SelectStatementFuncOpt) RenderFormatted() string {
// 	tmpl, err := template.New("select").Parse(selectTemplate)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var sb strings.Builder
// 	err = tmpl.Execute(&sb, s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return sb.String()
// }

// type FunctionStatement struct {
// 	name string
// 	args []any
// }

// // Parameters implements database.Query.
// func (f *FunctionStatement) Parameters() []any {
// 	panic("unimplemented")
// }

// func (f FunctionStatement) Render() string {
// 	return f.get()
// }

// // Overrider implements database.Query.
// func (f *FunctionStatement) Overrider(func(*database.Statement) *database.Statement) {
// 	panic("unimplemented")
// }

// func Function(name string, args ...any) FunctionStatement {
// 	return FunctionStatement{
// 		name: strings.TrimSuffix(name, "()"),
// 		args: args,
// 	}
// }

// // Formatted implements database.Query.
// func (f FunctionStatement) RenderFormatted() string {
// 	return f.get()
// }

// // OneLine implements database.Query.
// func (f FunctionStatement) OneLine() string {
// 	return f.get()
// }

// const functionStatementFormat string = "SELECT %s(%s)"

// func (f FunctionStatement) get() string {
// 	var args []string
// 	for _, arg := range f.args {
// 		args = append(args, fmt.Sprint(arg))
// 	}

// 	if len(args) == 0 {
// 		return fmt.Sprintf(functionStatementFormat, f.name, "")
// 	}
// 	// "SELECT %s(%s)"
// 	return fmt.Sprintf(functionStatementFormat, f.name, strings.Join(args, ", "))
// }

// //--- JoinedSelect

// type JoinedSelect SelectStatement

// func NewJoinedSelect() JoinedSelect {
// 	return JoinedSelect{}
// }

// // Get implements Query.
// func (j JoinedSelect) Get(formated bool) string {
// 	if formated {
// 		return j.Formatted()
// 	}
// 	return j.OneLine()
// }

// // Formatted implements Query.
// func (j JoinedSelect) Formatted() string {
// 	tmpl, err := template.New("select").Parse(formatedSelect)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var sb strings.Builder
// 	err = tmpl.Execute(&sb, j)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return sb.String()
// }

// // OneLine implements Query.
// func (j JoinedSelect) OneLine() string {
// 	return strings.ReplaceAll(j.Formatted(), "\n", "")
// }
