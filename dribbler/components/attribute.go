package components

type (
	// Identifiable is used when a widget or component supports identification
	Identifiable interface {
		ID() int
		SetID(int)
	}
	Identification struct {
		id int
	}

	Selection interface {
		CursorX() int
		CursorY() int
		Cursor() (int, int)

		SetCursor(x, y int)
		MoveCursor(dX, dY int)

		MoveCursorUp(y ...int)
		MoveCursorDown(y ...int)
		MoveCursorLeft(x ...int)
		MoveCursorRight(x ...int)

		GetSelected() any
	}
)

func NewIdentification(id int) Identification {
	return Identification{
		id: id,
	}
}

func (i Identification) ID() int {
	return i.id
}

func (i *Identification) SetID(id int) {
	i.id = id
}

type (
	// ShouldAlwaysUpdate is implemented by widgets or components that shall always be updated
	ShouldAlwaysUpdate interface {
		ShouldAlwaysUpdate()
	}
	AlwaysUpdate struct{}
)

func (a AlwaysUpdate) ShouldAlwaysUpdate() {}

type (
	// Named is implemented by widgets or components that wants their name to be visible to other widgets or components
	Named interface {
		SetName(string)
		Name() string
	}
	Name struct {
		name string
	}
)

func NewName(name string) Name {
	return Name{
		name: name,
	}
}

func (n Name) Name() string {
	return n.name
}

func (n *Name) SetName(name string) {
	n.name = name
}
