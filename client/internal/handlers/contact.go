package handlers

import (
	"context"
	"log"

	"github.com/rivo/tview"
	"grpcexi/client/internal/services"
	pb "grpcexi/protos"
)

type ContactHandler struct {
	app     *tview.Application
	service *services.ContactService
}

func NewContactHandler(app *tview.Application, service *services.ContactService) *ContactHandler {
	return &ContactHandler{
		app:     app,
		service: service,
	}
}

func (h *ContactHandler) SaveContact(
	name string,
	phones []*pb.PhoneNumber,
	formText *tview.TextView,
	clearInputs func(),
) {
	contact := &pb.Contact{Name: name, Phones: phones}

	go func() {
		err := h.service.AddContact(context.Background(), contact)

		h.app.QueueUpdateDraw(func() {
			if err != nil {
				log.Printf("AddContact error: %v", err)
				formText.SetText("Failed to save contact: " + err.Error())
			} else {
				formText.SetText("Contact saved successfully!")
				clearInputs()
			}
		})
	}()
}

func (h *ContactHandler) StreamContacts(table *tview.Table) {
	go func() {
		ch, err := h.service.StreamContacts(context.Background())
		if err != nil {
			log.Printf("StreamContacts error: %v", err)
			h.app.QueueUpdateDraw(func() {
				table.SetCell(1, 0, tview.NewTableCell("Failed to load contacts: "+err.Error()))
			})
			return
		}

		row := 1
		for c := range ch {
			for _, p := range c.Phones {
				r := row
				h.app.QueueUpdateDraw(func() {
					table.SetCell(r, 0, tview.NewTableCell(c.Name))
					table.SetCell(r, 1, tview.NewTableCell(p.Digits))
					table.SetCell(r, 2, tview.NewTableCell(p.Type.String()))
				})
				row++
			}
		}
	}()
}
