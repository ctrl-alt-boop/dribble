package popup

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribbler/logging"
	"github.com/ctrl-alt-boop/dribbler/ui"
)

var logger = logging.GlobalLogger()

const (
	KindConnect      Kind = "login"
	KindQueryBuilder Kind = "query"
	KindTableCell    Kind = "cell"
	KindDetails      Kind = "details"
)

// I will handle the rendering of this in three stages from the innermost it looks like this
// 1. The PopupModel will have its content type from the Popup method,
// this is then handled in the PopupModel and the PopupHandler will use its View() method.
// 2. The PopupModel's View is then inserted into a PopupStyle that is a viewport,
// where its size is either some MaxWidth, MaxHeight or the popup's content, whichever is the smallest.
// 3.* The PopupHandler will then render the Popup inside a Workspace style with the same size as the Workspace

type (
	Kind string

	PopupModel interface {
		tea.Model
		SetMaxSize(width, height int)
		GetContentSize() (int, int)
		GetContentWidth() int
		GetContentHeight() int
	}

	PopupHandler struct {
		dribbleClient                 *dribble.Client
		InnerWidth, InnerHeight       int // these are temporary until I find a good way to render this
		PopupMaxWidth, PopupMaxHeight int
		PopupWidth, PopupHeight       int
		// ContentWidth, ContentHeight int

		currentPopup PopupModel
	}
)

func (p *PopupHandler) IsOpen() bool {
	return p.currentPopup != nil
}

func NewHandler(dribbleClient *dribble.Client) *PopupHandler {
	return &PopupHandler{
		dribbleClient: dribbleClient,
	}
}

func (p *PopupHandler) Init() tea.Cmd {
	if p.currentPopup != nil {
		return p.currentPopup.Init()
	}
	return nil
}

func (p *PopupHandler) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if p.currentPopup != nil {
		popup, cmd := p.currentPopup.Update(msg)
		if pop, ok := popup.(PopupModel); ok {
			p.currentPopup = pop
			return p, cmd
		}
	}

	return p, cmd
}

func (p *PopupHandler) UpdateSize(width int, height int) {
	p.InnerWidth = width - ui.WorkspaceStyle.GetHorizontalFrameSize() - ui.PopupHandlerStyle.GetHorizontalFrameSize()
	p.InnerHeight = height - ui.WorkspaceStyle.GetVerticalFrameSize() - ui.PopupHandlerStyle.GetVerticalFrameSize()

	p.PopupMaxWidth = p.InnerWidth - ui.PopupStyle.GetHorizontalFrameSize()
	p.PopupMaxHeight = p.InnerHeight - ui.PopupStyle.GetVerticalFrameSize()
}

func (p *PopupHandler) View() string {
	if p.currentPopup != nil {
		p.currentPopup.SetMaxSize(p.PopupMaxWidth, p.PopupMaxHeight)

		return ui.PopupStyle.MaxWidth(p.PopupMaxWidth).MaxHeight(p.PopupMaxHeight).Render(p.currentPopup.View())
		// return ui.PopupStyle.MaxWidth(p.PopupMaxWidth).MaxHeight(p.PopupMaxHeight).Render(p.currentPopup.View())
		// currentPopupView := ui.PopupStyle.MaxWidth(p.PopupMaxWidth).MaxHeight(p.PopupMaxHeight).Render(p.currentPopup.View())
		// centered := lipgloss.Place(p.InnerWidth, p.InnerHeight, lipgloss.Center, lipgloss.Center, currentPopupView)

		// return ui.PopupHandlerStyle.MaxWidth(p.InnerWidth).MaxHeight(p.InnerHeight).Render(centered)
	}
	return ""
}

func (p *PopupHandler) Popup(popupType Kind, args ...any) tea.Cmd {
	// popupWidth, popupHeight := p.PopupWidth/2, p.PopupHeight/2
	switch popupType {
	case KindConnect:
		driverName := args[0].(string)
		p.currentPopup = newConnect(driverName)
	case KindQueryBuilder:
		switch arg := args[0].(type) {
		case database.Intent:
			p.currentPopup = newQueryBuilder(nil) // TODO: query builder way
		case database.OperationType:
			p.currentPopup = newTableQueryBuilder(arg, args[1].(string))
		}
	case KindTableCell:
		value := args[0].(string)
		p.currentPopup = newCellData(value)
	case KindDetails:
		value := args[0].(string)
		p.currentPopup = newDetails(value)
	default:
		p.currentPopup = nil
	}
	return p.currentPopup.Init()
}

func (p *PopupHandler) Close() {
	if p.currentPopup != nil {
		p.currentPopup = nil
	}
}
