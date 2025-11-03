package components

// import (
// 	"container/ring"
// 	"reflect"

// 	"github.com/charmbracelet/bubbles/v2/key"
// 	tea "github.com/charmbracelet/bubbletea/v2"
// 	"github.com/charmbracelet/lipgloss/v2"
// 	"github.com/ctrl-alt-boop/dribbler/config"
// 	"github.com/ctrl-alt-boop/dribbler/internal/panel"

// 	"github.com/ctrl-alt-boop/dribbler/logging"
// )

// var logger = logging.GlobalLogger()

// var _ Model = (*ContentArea)(nil)

// type ContentArea struct {
// 	ID   int
// 	name string

// 	Style        lipgloss.Style
// 	PanelManager panel.Manager

// 	Children []tea.Model

// 	idChildren map[int]tea.Model
// 	nextID     int

// 	// FocusedChild int
// 	focus *ring.Ring

// 	msgHandlers map[reflect.Type]func(msg tea.Msg) (Model, tea.Cmd)
// }

// // New creates a new ContentArea widget
// func New(name string, manager panel.Manager, children ...tea.Model) ContentArea {
// 	idChildren := map[int]tea.Model{}
// 	for i, child := range children {
// 		idChildren[i] = child
// 	}
// 	nextID := len(children)

// 	return ContentArea{
// 		name:         name,
// 		PanelManager: manager,
// 		Children:     children,
// 		idChildren:   map[int]tea.Model{},
// 		nextID:       nextID,
// 		Style:        lipgloss.NewStyle(),

// 		msgHandlers: map[reflect.Type]func(msg tea.Msg) (Model, tea.Cmd){},
// 	}
// }

// func (a *ContentArea) Focused() int {
// 	return a.focus.Value.(int)
// }

// // AddMsgHandler [temporary maybe method for widgets or components that need some super special handling for message type]
// func (a *ContentArea) AddMsgHandler(msg any, handler func(msg tea.Msg) (Model, tea.Cmd)) {
// 	a.msgHandlers[reflect.TypeOf(msg)] = handler
// }

// // Name implements widget.Named
// //
// // Is used by other widgets or components
// func (a *ContentArea) Name() string {
// 	return a.name
// }

// // AddChild .
// func (a *ContentArea) AddChild(child tea.Model) {
// 	a.Children = append(a.Children, child)
// 	newID := a.nextID
// 	a.idChildren[newID] = child
// 	a.nextID++
// }

// // SetStyle is used to render all children
// func (a *ContentArea) SetStyle(style lipgloss.Style) {
// 	a.Style = style
// }

// // Init implements tea.Model.
// func (a ContentArea) Init() tea.Cmd {
// 	cmds := []tea.Cmd{}
// 	for _, child := range a.Children {
// 		cmds = append(cmds, child.Init())
// 	}
// 	return tea.Batch(cmds...)
// }

// // Update implements tea.Model.
// func (a ContentArea) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	logger.Infof("%s got %T: %v", a.Name(), msg, msg)
// 	var cmds []tea.Cmd
// 	updated := a
// 	handler, ok := updated.msgHandlers[reflect.TypeOf(msg)]
// 	if ok {
// 		updated, cmd := handler(msg)
// 		if cmd != nil {
// 			cmds = append(cmds, cmd)
// 		}
// 		return updated, tea.Batch(cmds...)
// 	}

// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:

// 		updated.PanelManager.SetSize(msg.Width-updated.Style.GetHorizontalFrameSize(), msg.Height-updated.Style.GetVerticalFrameSize())

// 		updatedChildren := updated.PanelManager.Layout(updated.Children)
// 		updated.Children = updatedChildren
// 		return updated, nil

// 	case tea.KeyMsg:
// 		var cmd tea.Cmd
// 		switch {
// 		case len(updated.Children) == 0:
// 			return updated, nil
// 		// case updated.FocusedChild == -1:
// 		// 	// noop
// 		case key.Matches(msg, config.Keys.CycleViewNext):
// 			cmd := updated.cycleView(1)
// 			if cmd != nil {
// 				cmds = append(cmds, cmd)
// 			}
// 		case key.Matches(msg, config.Keys.CycleViewPrev):
// 			cmd := updated.cycleView(-1)
// 			if cmd != nil {
// 				cmds = append(cmds, cmd)
// 			}
// 		default:
// 			if updated.focus.Len() == 0 {
// 				break
// 			}
// 			updated.Children[updated.focus.Value.(int)], cmd = updated.Children[updated.focus.Value.(int)].Update(msg)
// 		}
// 		if cmd != nil {
// 			cmds = append(cmds, cmd)
// 		}

// 	// case FocusMsg:
// 	// 	updated.FocusedChild = msg.Index[0]
// 	// 	if updated.FocusedChild != -1 && updated.FocusedChild >= len(updated.Children) {
// 	// 		if len(msg.Index) > 1 {
// 	// 			msg.Index = msg.Index[1:]
// 	// 			focusedChild, cmd := updated.Children[updated.FocusedChild].Update(msg)
// 	// 			updated.Children[updated.FocusedChild] = focusedChild
// 	// 			if cmd != nil {
// 	// 				cmds = append(cmds, cmd)
// 	// 			}
// 	// 		}
// 	// 	}

// 	default:
// 		for i, child := range updated.Children {
// 			logger.Infof("Updating child %d", i)
// 			child, cmd := child.Update(msg)
// 			updated.Children[i] = child
// 			if cmd != nil {
// 				cmds = append(cmds, cmd)
// 			}
// 		}
// 	}
// 	cmds = append(cmds, updated.updateAlwaysUpdateChildren(msg))
// 	return updated, tea.Batch(cmds...)
// }

// // Render implements Model.
// func (a ContentArea) Render() string {
// 	return a.Style.Render(a.PanelManager.View(a.Children))
// }

// // View implements tea.Model.
// func (a ContentArea) View() tea.View {
// 	return tea.NewView(a.Render())
// }

// func (a ContentArea) updateAlwaysUpdateChildren(msg tea.Msg) tea.Cmd {
// 	var cmd tea.Cmd
// 	var cmds []tea.Cmd
// 	for i, child := range a.Children {
// 		if _, ok := any(child).(ShouldAlwaysUpdate); ok {
// 			a.Children[i], cmd = child.Update(msg)
// 			if cmd != nil {
// 				cmds = append(cmds, cmd)
// 			}
// 		}
// 	}
// 	return tea.Batch(cmds...)
// }

// func (a *ContentArea) cycleView(i int) tea.Cmd {
// 	var cmds []tea.Cmd
// 	prevFocusedChild := a.focus.Value.(int)
// 	_, cmdPrev := a.Children[prevFocusedChild].Update(LostFocusMsg{})
// 	cmds = append(cmds, cmdPrev)

// 	a.focus = a.focus.Move(i)
// 	newFocusedChild := a.focus.Value.(int)

// 	a.PanelManager.SetFocusedIndex(newFocusedChild)
// 	_, cmdNew := a.Children[newFocusedChild].Update(GotFocusMsg{})
// 	cmds = append(cmds, cmdNew)

// 	return tea.Batch(cmds...)
// }
