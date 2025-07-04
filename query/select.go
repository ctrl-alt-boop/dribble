package query

// //go:embed templates/b_select.tmpl
// var selectBuilderTemplate string

// func  ToSql() (queryString string, params []any, err error) {
// 	render, params, err := s.ToSqlFormatted(dialect)
// 	if err != nil {
// 		return "", nil, err
// 	}
// 	replacer := strings.NewReplacer("\n", " ", "\t", "", "\r", "", "\r\n", " ")
// 	return strings.TrimSpace(replacer.Replace(render)), params, err
// }

// func ToSqlFormatted() (queryString string, params []any, err error) {
// 	// selectTemplate := dialect.SelectTemplate()
// 	tmpl, err := template.New("select").Parse(selectBuilderTemplate)
// 	if err != nil {
// 		return "", nil, err
// 	}
// 	var sb strings.Builder
// 	err = tmpl.Execute(&sb, s)
// 	if err != nil {
// 		return "", nil, err
// 	}
// 	return strings.TrimSpace(sb.String()), s.params, nil
// }
