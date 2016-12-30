package main

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gizak/termui"
	"github.com/kolo/xmlrpc"
	"github.com/spf13/viper"
)

// LoadAPI returns the api and apiKey according to the settings defined in the
// configuration file, respectively.
func LoadAPI() (api *xmlrpc.Client, apiKey string, err error) {
	viper.SetConfigName("bapu")

	homePath := os.Getenv("HOME")
	viper.AddConfigPath(homePath)
	viper.AddConfigPath("/usr/local/etc")
	viper.AddConfigPath("/etc")

	err = viper.ReadInConfig()
	if err != nil {
		return api, apiKey, err
	}

	production := viper.GetBool("production.enabled")
	if production {
		apiKey = viper.GetString("production.apiKey")
		api, err = xmlrpc.NewClient("https://rpc.gandi.net/xmlrpc/", nil)
		if err != nil {
			return api, apiKey, err
		}
	}

	development := viper.GetBool("development.enabled")
	if development {
		log.Println("Development Config found")
		api, err = xmlrpc.NewClient("https://rpc.ote.gandi.net/xmlrpc/", nil)
		if err != nil {
			return api, apiKey, err
		}
		apiKey = viper.GetString("development.apiKey")
	}

	if api == nil {
		return api, apiKey, errors.New("neither production nor development environment enabled in config")
	}

	return api, apiKey, nil
}

func main() {
	// Load API
	api, apiKey, err := LoadAPI()
	if err != nil {
		log.Fatal(err)
	}

	// initialize termui
	err = termui.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termui.Close()

	// Title
	uiTitle := termui.NewPar("Bapu")
	uiTitle.Border = false
	uiTitle.TextFgColor = termui.ColorMagenta

	// Count number of instances
	var hostingVMCount *int
	err = api.Call("hosting.vm.count", apiKey, &hostingVMCount)
	if err != nil {
		log.Fatal(err)
	}

	uiCount := termui.NewPar("VM #: " + strconv.Itoa(*hostingVMCount))
	uiCount.Border = false
	uiCount.TextFgColor = termui.ColorWhite

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
		strs = append(strs, "["+strconv.Itoa(val.ID)+"] "+val.Hostname+" ("+val.State+")")
	}

	uiList := termui.NewList()
	uiList.Items = strs
	uiList.ItemFgColor = termui.ColorYellow
	uiList.BorderLabel = "Servers"
	uiList.Height = len(strs) + 2

	// Create termui Grid system
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(3, 0, uiTitle),
		),
		termui.NewRow(
			termui.NewCol(2, 0, uiCount),
		),
		termui.NewRow(
			termui.NewCol(10, 0, uiList),
		),
	)

	// calculate layout
	termui.Body.Align()
	termui.Render(termui.Body)

	// Quit with q
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		termui.Body.Align()
		termui.Render(termui.Body)
	})

	termui.Handle("/timer/1s", func(e termui.Event) {
		t := e.Data.(termui.EvtTimer)
		// t is a EvtTimer
		if t.Count%2 == 0 {
			termui.Body.Align()
			termui.Render(termui.Body)
		}
	})

	termui.Loop()
}
