package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gizak/termui"
	"github.com/kolo/xmlrpc"
)

func main() {

	// initialize termui
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	// Title
	title := termui.NewPar("Bapu")
	title.Border = false
	title.Height = 1
	title.TextFgColor = termui.ColorMagenta
	title.Width = 10
	title.X = 1
	title.Y = 1

	// List
	strs := []string{
		"[0] github.com/gizak/termui",
		"[1] [你好，世界](fg-blue)",
		"[2] [こんにちは世界](fg-red)",
		"[3] [color output](fg-white,bg-green)",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] nsf/termbox-go"}

	ls := termui.NewList()
	ls.Items = strs
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "Servers"
	ls.Height = 20
	ls.Width = 80
	ls.Y = 4

	termui.Render(title, ls)

	apiKey := os.Getenv("BAPU_GANDI_KEY")
	development := false
	if os.Getenv("BAPU_DEVELOPMENT") != "" {
		development = true
	}

	api, err := xmlrpc.NewClient("https://rpc.gandi.net/xmlrpc/", nil)
	if development {
		api, err = xmlrpc.NewClient("https://rpc.ote.gandi.net/xmlrpc/", nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Count number of instances
	var paasCount *int
	err = api.Call("paas.count", apiKey, &paasCount)
	if err != nil {
		log.Fatal(err)
	}

	// Count of Servers
	count := termui.NewPar(strconv.Itoa(*paasCount))
	count.Border = false
	count.Height = 1
	count.TextFgColor = termui.ColorMagenta
	count.Width = 10
	count.X = 1
	count.Y = 2

	termui.Render(count)

	// Quit with q
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		termui.Render(ls)
	})

	termui.Loop()
}
