package main

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

func main() {
	size := 10
	d := make([]string, size)
	ctx := context.Background()

	ch := clipboard.Watch(ctx, clipboard.FmtText)

	myWindow, list := makeWindow(d)

	go func() {
		current := -1
		for data := range ch {
			if current == size-1 {
				current = 0
			} else {
				current++
			}

			copied := string(data)
			exists := false

			for _, v := range d {
				if v == copied {
					exists = true
					break
				}
			}

			if !exists {
				d[current] = copied
				list.Refresh()
			}
		}
	}()

	myWindow.ShowAndRun()
}

func makeWindow(d []string) (fyne.Window, *widget.List) {
	myApp := app.NewWithID("clip-app")
	myWindow := myApp.NewWindow("Memory clipboard")

	myWindow.Resize(fyne.NewSize(300, 300))

	list := widget.NewList(
		func() int {
			return len(d)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(d[i])
		})

	list.OnSelected = func(id widget.ListItemID) {
		clipboard.Write(clipboard.FmtText, []byte(d[id]))
		myApp.SendNotification(fyne.NewNotification("copied to clipboard", d[id]))
	}

	myWindow.SetContent(list)

	return myWindow, list
}
