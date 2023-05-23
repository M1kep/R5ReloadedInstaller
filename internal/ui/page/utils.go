package page

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func checkBoxListBuilder(checkBoxes []map[*widget.Bool]string) []layout.FlexChild {
	var checkBoxList []layout.FlexChild
	for _, checkBox := range checkBoxes {
		for bool, label := range checkBox {
			checkBoxList = append(checkBoxList, layout.Rigid(
				func(gtx layout.Context) layout.Dimensions {
					chkBox := material.CheckBox(theme, bool, label)
					return chkBox.Layout(gtx)
				},
			))
		}
	}
	return checkBoxList
}
