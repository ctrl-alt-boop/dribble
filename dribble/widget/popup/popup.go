package popup

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/dribble/ui"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

var logger = logging.NewLogger("popup.log")

const (
	KindConnect      Kind = "login"
	KindQueryOptions Kind = "query"
	KindTableCell    Kind = "cell"
)

// I will handle the rendering of this in three stages from the innermost it looks like this
// 1. The PopupModel will have its content type from the Popup method,
// this is then handled in the PopupModel and the PopupHandler will use its View() method.
// 2. The PopupModel's View is then inserted into a PopupStyle that is a viewport,
// where its size is either some MaxWidth, MaxHeight or the popup's content, whichever is the smallest.
// 3.* The PopupHandler will then render the Popup inside a Workspace style with the same size as the Workspace

// * When a solution for actual modals/popups has been found I will adapt the last step to render it as a modal/popup

type (
	Kind string

	PopupModel interface {
		tea.Model
		Exec() tea.Cmd
		Cancel() tea.Cmd
		SetMaxSize(width, height int)
		GetContentSize() (int, int)
		GetContentWidth() int
		GetContentHeight() int
	}

	PopupHandler struct {
		goolDb                        *gooldb.GoolDb
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

func NewHandler(gool *gooldb.GoolDb) *PopupHandler {
	return &PopupHandler{
		goolDb: gool,
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
	case KindQueryOptions:
		p.currentPopup = newQueryOptions(nil)
	case KindTableCell:
		value := args[0].(string)
		p.currentPopup = newCellData(value)
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
