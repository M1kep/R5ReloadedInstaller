package gui

import (
	"R5ReloadedInstaller/internal/app"
	giouiapp "gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

type WindowManager struct {
	Window      *giouiapp.Window
	PageTable   map[string]app.Page
	PageHistory []string
	currentPage app.Page
	ops         op.Ops
}

func NewWindowManager(title string, startPage app.Page) *WindowManager {
	w := giouiapp.NewWindow(
		giouiapp.Title(title),
		giouiapp.Size(unit.Dp(400), unit.Dp(600)),
	)

	windowManager := &WindowManager{
		Window:      w,
		PageTable:   map[string]app.Page{},
		PageHistory: []string{},
	}
	windowManager.SetCurrentPage(startPage)
	return windowManager
}

func (w *WindowManager) SetCurrentPage(page app.Page) {
	// Unload the current page
	if w.currentPage != nil {
		w.currentPage.OnPageUnload()
	}

	// Load the new page
	w.currentPage = page
	w.currentPage.OnPageLoad(w)
	//w.currentPage.SetWindowManager()
	w.PageTable[w.currentPage.ID()] = w.currentPage
	w.PageHistory = append(w.PageHistory, w.currentPage.ID())
}

func (w *WindowManager) SetCurrentPageByID(id string) {
	w.SetCurrentPage(w.PageTable[id])
}

func (w *WindowManager) KnowsPage(id string) bool {
	_, ok := w.PageTable[id]
	return ok
}

func (w *WindowManager) GetCurrentPage() app.Page {
	return w.currentPage
}

func (w *WindowManager) PreviousPage() {
	if len(w.PageHistory) < 2 {
		return
	}

	w.PageHistory = w.PageHistory[:len(w.PageHistory)-1]
	w.SetCurrentPage(w.PageTable[w.PageHistory[len(w.PageHistory)-1]])
}

func (w *WindowManager) Run() error {
	for {
		e := <-w.Window.Events()
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&w.ops, e)
			w.currentPage.HandleUserInteractions()
			w.currentPage.Layout(gtx)
			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
}
