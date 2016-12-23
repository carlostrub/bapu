package main

import "github.com/gizak/termui"

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

	// Quit with q
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		termui.Render(ls)
	})

	termui.Loop()

	//	apiKey := os.Getenv("GANDI_KEY")
	//
	//	//	api, err := xmlrpc.NewClient("https://rpc.gandi.net/xmlrpc/", nil)
	//	api, err := xmlrpc.NewClient("https://rpc.ote.gandi.net/xmlrpc/", nil)
	//	if err != nil {
	//		log.Fatal(err)
	//	}

	//	var result struct {
	//		Count int `xmlrpc:"count"`
	//	}

	// Count number of instances
	//	var paasCount *int
	//	err = api.Call("paas.count", apiKey, &paasCount)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	fmt.Println(*paasCount)

}
