package main

import (
	"context"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

const (
	size          = 10
	uiListItemLen = 20
)

func main() {
	// @TODO: add clear all items, add delete item
	d := make([]string, size)
	ctx := context.Background()

	ch := clipboard.Watch(ctx, clipboard.FmtText)

	myWindow, list := makeWindow(d)

	go func() {
		current := -1
		for data := range ch {
			copied := string(data)

			if copied == "" || copied == "\n" || copied == "\t" {
				continue
			}

			exists := false

			if exists {
				continue
			}

			for _, v := range d {
				if v == copied {
					exists = true
					break
				}
			}

			if current == size-1 {
				current = 0
			} else {
				current++
			}

			d[current] = copied
			list.Refresh()
		}
	}()

	myWindow.ShowAndRun()
}

func substrLongText(str string) string {
	str = regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(str), " ")
	if len(str) > uiListItemLen {
		return str[0:uiListItemLen-1] + "..."
	}

	return str
}

func makeWindow(d []string) (fyne.Window, *widget.List) {
	myApp := app.NewWithID("clip-app")
	myWindow := myApp.NewWindow("Memory clipboard")

	list := widget.NewList(
		func() int {
			return len(d)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).Wrapping = fyne.TextTruncate
			o.(*widget.Label).SetText(substrLongText(d[i]))
		})

	list.OnSelected = func(id widget.ListItemID) {
		if d[id] == "" {
			return
		}

		clipboard.Write(clipboard.FmtText, []byte(d[id]))
		myApp.SendNotification(fyne.NewNotification("copied to clipboard", d[id]))
	}

	list.BaseWidget.MinSize()

	b := widget.NewButton("clear", func() {
		d = make([]string, size)
		list.Refresh()
	})

	butContainer := container.NewGridWrap(fyne.NewSize(300, 30), b)
	listContainer := container.NewGridWrap(fyne.NewSize(300, 270), list)

	content := container.NewVBox(butContainer, listContainer)
	myWindow.SetContent(content)

	return myWindow, list
}
