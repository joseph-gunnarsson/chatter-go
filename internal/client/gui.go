package client

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/joseph-gunnarsson/chatter-go/internal/errors"
	"github.com/jroimartin/gocui"
)

func newGUI(client *Client) {
	g, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		errors.HandleError(errors.ClientGUILoadError{Message: err.Error()})
	}
	defer g.Close()

	g.SetManagerFunc(generateLayout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		errors.HandleError(errors.ClientGUILoadError{Message: err.Error()})
	}

	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, writeMessages); err != nil {
		errors.HandleError(errors.ClientGUILoadError{Message: err.Error()})
	}

	go client.readMessage(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func generateLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("messages", 0, 0, maxX-1, maxY-3); err != nil {
		v.Title = "Chat Messages"
		v.Autoscroll = true
		v.Wrap = true
	}

	if v, err := g.SetView("input", 0, maxY-3, maxX-1, maxY-1); err != nil {
		v.Title = "Type your message"
		v.Editable = true
		v.Wrap = true
	}

	g.SetCurrentView("input")

	return nil
}

func writeMessages(g *gocui.Gui, v *gocui.View) error {
	message := strings.TrimPrefix(v.Buffer(), " ")
	if message == "" {
		message = "\n"
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	formattedMessage := fmt.Sprintf("[%s] you: %s", timestamp, message)

	messageView, err := g.View("messages")

	if err != nil {

	}

	if !strings.HasPrefix(message, "/") {
		_, err = messageView.Write([]byte(formattedMessage))
	}
	messages <- message

	if err != nil {
		errors.HandleError(errors.ClientErrorExit{Message: err.Error()})
	}

	v.SetCursor(0, 0)
	v.Clear()
	return nil
}

func (c *Client) readMessage(g *gocui.Gui) error {
	buf := make([]byte, 1024)
	messageView, _ := g.View("messages")
	for {
		n, err := c.conn.Read(buf)

		if err != nil {
			errors.HandleError(err)
			break
		}

		go g.Update(func(g *gocui.Gui) error {
			requestString := strings.TrimSpace(string(buf[:n]))
			messageView.Write([]byte(requestString + "\n"))
			return nil
		})

	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
