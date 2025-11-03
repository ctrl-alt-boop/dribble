package data

type Text struct {
	Base
	string
}

func TextItem(s string) Text {
	return Text{
		string: s,
	}
}

func (t Text) Render() string {
	return t.string
}
