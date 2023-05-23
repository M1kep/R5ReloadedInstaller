package interaction

import "gioui.org/widget"

type Manager struct {
	buttonHandlers map[*widget.Clickable]func()
	boolHandlers   map[*widget.Bool]func()
}

func NewManager() *Manager {
	return &Manager{
		buttonHandlers: make(map[*widget.Clickable]func()),
		boolHandlers:   make(map[*widget.Bool]func()),
	}
}

func (m *Manager) RegisterClickableHandler(clickable *widget.Clickable, handler func()) {
	m.buttonHandlers[clickable] = handler
}

func (m *Manager) RegisterBoolHandler(bool *widget.Bool, handler func()) {
	m.boolHandlers[bool] = handler
}

func (m *Manager) RunHandlers() {
	// Checkboxes come first as they may change values referenced by other handlers
	for bool, handler := range m.boolHandlers {
		if bool.Changed() {
			handler()
		}
	}

	for clickable, handler := range m.buttonHandlers {
		if clickable.Clicked() {
			handler()
		}
	}
}
