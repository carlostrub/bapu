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

var selector int

// AccountReturn contains fields for informations about the Gandi account
type AccountReturn struct {
	AverageCreditCost     float64   `xmlrpc:"average_credit_cost"`
	Credits               int       `xmlrpc:"credits"`
	CycleDay              int       `xmlrpc:"cycle_day"`
	DateCreditsExpiration time.Time `xmlrpc:"date_credits_expiration"`
	FullName              string    `xmlrpc:"fullname"`
	Handle                string    `xmlrpc:"handle"`
	ID                    int       `xmlrpc:"id"`
}

// DiskReturn contains fields for informations about the Disks
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

// VMReturn contains fields for informations about the virtual machines
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

func serverList(list []VMReturn) (servers [][]string) {

	servers = append(servers, []string{
		"Selected",
		"Hostname",
		"Datacenter",
		"Cores",
		"Memory",
		"State",
	})

	for i, val := range list {
		s := ""
		if selector == i {
			s = "*"
		}
		servers = append(servers, []string{
			s,
			val.Hostname,
			strconv.Itoa(val.DatacenterID),
			strconv.Itoa(val.Cores),
			strconv.Itoa(val.Memory) + "MB",
			val.State,
		})
	}

	return servers
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
	uiTitle := termui.NewPar("Bapu -- Control your Gandi Machines")
	uiTitle.Border = false
	uiTitle.Height = 3
	uiTitle.TextFgColor = termui.ColorMagenta

	// Summary
	var hostingVMCount *int
	err = api.Call("hosting.vm.count", apiKey, &hostingVMCount)
	if err != nil {
		log.Fatal(err)
	}
	var hostingAccountInfo *AccountReturn
	err = api.Call("hosting.account.info", apiKey, &hostingAccountInfo)
	if err != nil {
		log.Fatal(err)
	}

	info := *hostingAccountInfo

	uiSummary := termui.NewPar("Owner: " + info.FullName + "    Virtual Machines: " + strconv.Itoa(*hostingVMCount) + "    Remaining Credit: " + strconv.Itoa(info.Credits))
	uiSummary.Height = 3
	uiSummary.Border = true
	uiSummary.BorderLabel = "Summary"
	uiSummary.TextFgColor = termui.ColorWhite

	// List instances
	var hostingVMList *[]VMReturn
	err = api.Call("hosting.vm.list", apiKey, &hostingVMList)
	if err != nil {
		log.Fatal(err)
	}

	list := *hostingVMList

	uiTable := termui.NewTable()
	uiTable.Rows = serverList(list)
	uiTable.FgColor = termui.ColorWhite
	uiTable.BgColor = termui.ColorDefault
	uiTable.TextAlign = termui.AlignCenter
	uiTable.Seperator = false
	uiTable.Analysis()
	uiTable.SetSize()
	uiTable.BgColors[0] = termui.ColorWhite
	uiTable.FgColors[0] = termui.ColorBlack
	uiTable.Y = 20
	uiTable.X = 0
	uiTable.Border = false
	for i := 0; i < len(list); i++ {
		switch list[i].State {
		case "paused":
			uiTable.BgColors[i+1] = termui.ColorBlack
			uiTable.FgColors[i+1] = termui.ColorWhite
		case "running":
			uiTable.BgColors[i+1] = termui.ColorBlue
			uiTable.FgColors[i+1] = termui.ColorWhite
		case "halted":
			uiTable.BgColors[i+1] = termui.ColorYellow
			uiTable.FgColors[i+1] = termui.ColorBlack
		case "locked":
			uiTable.BgColors[i+1] = termui.ColorMagenta
			uiTable.FgColors[i+1] = termui.ColorGreen
		case "being_created":
			uiTable.BgColors[i+1] = termui.ColorWhite
			uiTable.FgColors[i+1] = termui.ColorBlack
		case "deleted":
			uiTable.BgColors[i+1] = termui.ColorBlack
			uiTable.FgColors[i+1] = termui.ColorRed
		}
	}

	// Commands
	uiCommands := termui.NewPar("<[S]tart>    <St[o]p>    <[R]eboot>  Virtual Machine |   <[Q]uit>")
	uiCommands.Height = 3
	uiCommands.Border = false
	uiCommands.BorderLabel = "Summary"
	uiCommands.TextFgColor = termui.ColorBlack
	uiCommands.TextBgColor = termui.ColorWhite

	// Create termui Grid system
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, uiTitle),
		),
		termui.NewRow(
			termui.NewCol(12, 0, uiSummary),
		),
		termui.NewRow(
			termui.NewCol(12, 0, uiTable),
		),
		termui.NewRow(
			termui.NewCol(12, 0, uiCommands),
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
		if selector > 0 {
			selector--
		}
		uiTable.Rows = serverList(list)
		termui.Render(termui.Body)
	})

	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		if selector < len(list)-1 {
			selector++
		}
		uiTable.Rows = serverList(list)
		termui.Render(termui.Body)
	})

	termui.Handle("/sys/kbd/s", func(termui.Event) {
		err = api.Call("hosting.vm.start", []interface{}{apiKey, list[selector].ID}, nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	termui.Handle("/sys/kbd/o", func(termui.Event) {
		err = api.Call("hosting.vm.stop", []interface{}{apiKey, list[selector].ID}, nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	termui.Handle("/sys/kbd/r", func(termui.Event) {
		err = api.Call("hosting.vm.reboot", []interface{}{apiKey, list[selector].ID}, nil)
		if err != nil {
			log.Fatal(err)
		}
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
