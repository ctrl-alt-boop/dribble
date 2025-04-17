package popup

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

var logger = logging.NewLogger("popup.log")

const (
	KindConnect      Kind = "login"
	KindQueryOptions Kind = "query"
	KindTableCell    Kind = "cell"
)

type (
	Kind string

	PopupModel interface {
		tea.Model
		Exec() tea.Cmd
		Cancel() tea.Cmd
		SetSize(width, height int)
	}

	PopupHandler struct {
		goolDb        *gooldb.GoolDb
		width, height int

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

	// switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	panelWidth := msg.Width/ui.PanelWidthRatio - ui.BorderThicknessDouble
	// 	p.width, p.height = msg.Width-panelWidth-ui.BorderThicknessDouble-1, msg.Height-5
	// 	if p.currentPopup != nil {
	// 		p.currentPopup.SetSize(p.width, p.height)
	// 		return p, nil
	// 	}
	// }

	return p, cmd
}

func (p *PopupHandler) View() string {
	if p.currentPopup != nil {
		popupStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(p.width).Height(p.height)

		return popupStyle.Render(p.currentPopup.View())
	}
	return ""
}

func (p *PopupHandler) Popup(popupType Kind, args ...any) tea.Cmd {
	switch popupType {
	case KindConnect:
		driverName := args[0].(string)
		p.currentPopup = NewConnect(driverName, p.width/2, p.height/2)
	case KindQueryOptions:
		p.currentPopup = &QueryOptions{}
	case KindTableCell:
		value := args[0].(string)
		p.currentPopup = &TableCell{
			Value: value,
		}
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
