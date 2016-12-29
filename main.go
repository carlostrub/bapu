package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gizak/termui"
	"github.com/kolo/xmlrpc"
	"github.com/spf13/viper"
)

func main() {
	// handle configurations for server
	viper.SetConfigName("bapu")           // no need to include file extension
	viper.AddConfigPath("/usr/local/etc") // set the path of your config file
	viper.AddConfigPath("../bapu")        // set the path of your config file

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	var apiKey string
	var api *xmlrpc.Client

	production := viper.GetBool("production.enabled")
	if production {
		apiKey = viper.GetString("production.apiKey")
		api, err = xmlrpc.NewClient("https://rpc.gandi.net/xmlrpc/", nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	development := viper.GetBool("development.enabled")
	if development {
		log.Println("Development Config found")
		api, err = xmlrpc.NewClient("https://rpc.ote.gandi.net/xmlrpc/", nil)
		if err != nil {
			log.Fatal(err)
		}
		apiKey = viper.GetString("development.apiKey")
	}

	if api == nil {
		log.Fatal("neither production nor development environment enabled in config")
	}

	// initialize termui
	err = termui.Init()
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

	termui.Render(title)

	// Count number of instances
	var hostingVMCount *int
	err = api.Call("hosting.vm.count", apiKey, &hostingVMCount)
	if err != nil {
		log.Fatal(err)
	}

	count := termui.NewPar("VM #: " + strconv.Itoa(*hostingVMCount))
	count.Border = false
	count.Height = 1
	count.TextFgColor = termui.ColorWhite
	count.Width = 20
	count.X = 1
	count.Y = 3

	termui.Render(count)

	// Define output structs
	type DiskReturn struct {
		CanSnapshot   bool      `xmlrpc:"can_snapshot"`
		DatacenterID  int       `xmlrpc:"datacenter_id"`
		DateCreated   time.Time `xmlrpc:"date_created"`
		DateUpdated   time.Time `xmlrpc:"date_updated"`
		ID            int       `xmlrpc:"id"`
		IsBootDisk    bool      `xmlrpc:"is_boot_disk"`
		KernelVersion string    `xmlrpc:"kernel_version"`
		Label         string    `xmlrpc:"label"`
		Name          string    `xmlrpc:"name"`
		Size          int       `xmlrpc:"size"`
		State         string    `xmlrpc:"state"`
		TotalSize     int       `xmlrpc:"total_size"`
		Type          string    `xmlrpc:"type"`
		Visibility    string    `xmlrpc:"visibility"`
	}
	type VMReturn struct {
		AiActive     int          `xmlrpc:"ai_active"`
		Console      int          `xmlrpc:"console"`
		ConsoleURL   string       `xmlrpc:"console_url"`
		Cores        int          `xmlrpc:"cores"`
		DatacenterID int          `xmlrpc:"datacenter_id"`
		DateCreated  time.Time    `xmlrpc:"date_created"`
		DateUpdated  time.Time    `xmlrpc:"date_updated"`
		Description  string       `xmlrpc:"description"`
		Disks        []DiskReturn `xmlrpc:"disks"`
		Farm         string       `xmlrpc:"farm"`
		FlexShares   int          `xmlrpc:"flex_shares"`
		Hostname     string       `xmlrpc:"hostname"`
		HVMState     string       `xmlrpc:"hvm_state"`
		ID           int          `xmlrpc:"id"`
		Memory       int          `xmlrpc:"memory"`
		State        string       `xmlrpc:"state"`
		VMmaxMemory  int          `xmlrpc:"vm_max_memory"`
	}

	// List instances
	var hostingVMList *[]VMReturn
	err = api.Call("hosting.vm.list", apiKey, &hostingVMList)
	if err != nil {
		log.Fatal(err)
	}

	var strs []string
	list := *hostingVMList
	for _, val := range list {
		strs = append(strs, val.Hostname+" ("+val.State+")")
	}

	ls := termui.NewList()
	ls.Items = strs
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "Servers"
	ls.Height = 20
	ls.Width = 80
	ls.Y = 5

	termui.Render(ls)

	// Quit with q
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		termui.Render(ls)
	})

	termui.Loop()
}
