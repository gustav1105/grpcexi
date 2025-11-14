package page

import (
	"fmt"
	"github.com/rivo/tview"
	"grpcexi/client/internal/components/layout"
	"grpcexi/client/internal/handlers"
	pb "grpcexi/protos"
)

type Home struct {
	*tview.Flex
	app *tview.Application
}

func NewHome(app *tview.Application, handler *handlers.ContactHandler) *Home {
	var phoneNumbers []*pb.PhoneNumber

	sidebar := layout.NewSidebar("Menu")
	formText := tview.NewTextView().SetText("Create a new contact").SetSize(1, 40)
	nameInput := tview.NewInputField().SetLabel("Name: ")
	phoneInput := tview.NewInputField().SetLabel("Phone: ")
	options := []string{
		pb.PhoneType_CELL.String(),
		pb.PhoneType_WORK.String(),
		pb.PhoneType_HOME.String(),
	}
	phoneTypeInput := tview.NewDropDown().
		SetLabel("Phone Type: ").
		SetOptions(options, nil).
		SetCurrentOption(0)
	phoneList := tview.NewList()

	form := tview.NewForm().
		AddFormItem(formText).
		AddFormItem(nameInput).
		AddFormItem(phoneInput).
		AddFormItem(phoneTypeInput).
		AddButton("Add Phone", func() {
			digits := phoneInput.GetText()
			index, _ := phoneTypeInput.GetCurrentOption()
			phoneType := pb.PhoneType(index)

			phoneNumbers = append(phoneNumbers, &pb.PhoneNumber{
				Digits: digits,
				Type:   phoneType,
			})

			phoneList.InsertItem(-1, fmt.Sprintf("%s (%s)", digits, phoneType.String()), "", 0, nil)
			phoneInput.SetText("")
			phoneTypeInput.SetCurrentOption(0)
		}).AddButton("Save", func() {
		handler.SaveContact(
			nameInput.GetText(),
			phoneNumbers,
			formText,
			func() {
				nameInput.SetText("")
				phoneInput.SetText("")
				phoneTypeInput.SetCurrentOption(0)
				phoneNumbers = []*pb.PhoneNumber{}
				phoneList.Clear()
			},
		)
	}).SetButtonsAlign(tview.AlignRight)

	content := tview.NewFlex().SetDirection(tview.FlexRow).SetDirection(tview.FlexColumn)
	content.SetBorder(true)
	content.SetTitle("Contacts")

	flex := tview.NewFlex().
		AddItem(sidebar, 20, 1, true).
		AddItem(content, 0, 1, false)

	sidebar.SetItems([]string{"Create Contact", "View Contacts"}, func(index int, text string) {
		switch index {
		case 0:
			content.Clear()
			spacer := tview.NewBox()
			innerFlex := tview.NewFlex().SetDirection(tview.FlexColumn).SetDirection(tview.FlexRow)
			innerFlex.AddItem(form, 0, 1, true)
			innerFlex.AddItem(phoneList, 0, 1, false)
			content.AddItem(spacer, 24, 1, false)
			content.AddItem(innerFlex, 0, 1, true)
			content.AddItem(spacer, 28, 1, false)
			app.SetFocus(form)
		case 1:
			content.Clear()
			table := tview.NewTable()
			table.SetCell(0, 0, tview.NewTableCell("Name"))
			table.SetCell(0, 1, tview.NewTableCell("Phone"))
			table.SetCell(0, 2, tview.NewTableCell("Type"))
			content.AddItem(table, 0, 1, true)
			app.SetFocus(table)

			handler.StreamContacts(table)
		}
	})

	return &Home{Flex: flex, app: app}
}
