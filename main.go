package main

import (
	"flag"
	"fmt"
	"github.com/rakyll/portmidi"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
)

var (
	buttons []int
)

func commandButton(button, value int) {
	fmt.Print("command:", button)
	switch button {
	case 12:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 13:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 14:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 15:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 16:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 17:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 18:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 19:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	}
}

func noteButton(button, value int) {
	fmt.Print("note:", button)
	switch button {
	case 36:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 37:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 38:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 39:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 40:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 41:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 42:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	case 43:
		if value == 0 {
			fmt.Print(" Up")
		} else {
			fmt.Print(" Down:", value)
		}
	}
}

func rangeButton(button, value int) {
	fmt.Print("range:", button)
	switch button {
	case 70:
		fmt.Print(" v:", value)
	case 71:
		fmt.Print(" v:", value)
	case 72:
		fmt.Print(" v:", value)
	case 73:
		fmt.Print(" v:", value)
	case 74:
		fmt.Print(" v:", value)
	case 75:
		fmt.Print(" v:", value)
	case 76:
		fmt.Print(" v:", value)
	case 77:
		fmt.Print(" v:", value)
	}
}

func programButton(button int) {
	switch button {
	case 0:
		fmt.Print("program 0")
	case 1:
		fmt.Print("program 1")
	case 2:
		fmt.Print("program 2")
	case 3:
		fmt.Print("program 3")
	case 4:
		fmt.Print("program 4")
	case 5:
		fmt.Print("program 5")
	case 6:
		fmt.Print("program 6")
	case 7:
		fmt.Print("program 7")
	}
}

func buttonActions(button, value int) {
	fmt.Println("")
	if button >= 0 && button <= 7 {
		programButton(button)
		return
	}
	if button >= 12 && button <= 19 {
		commandButton(button, value)
		return
	}
	if button >= 36 && button <= 43 {
		noteButton(button, value)
		return
	}
	if button >= 70 && button <= 77 {
		rangeButton(button, value)
		return
	}
	fmt.Println("button:", button, "value:", value)
	
}

func logButtons(buttonInt int) {
	for _, button := range buttons {
		if buttonInt == button {
			return
		}
	}
	buttons = append(buttons, buttonInt)
}

func listen(id int) {
	fmt.Printf("Listening to midi device : %d\n", id)
	in, err := portmidi.NewInputStream(portmidi.DeviceID(id), 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()
	ch := in.Listen()
	for {
		event := <-ch
		/*
			mj, err := json.MarshalIndent(event, "", "  ")
			if err != nil {
				log.Printf("Error serializing event: %s", err)
			}
			fmt.Println("\n", string(mj))
		*/
		// logButtons(int(event.Data1)) // map buttons // mode 2 data1 is value?
		buttonActions(int(event.Data1), int(event.Data2))
	}
}

func listDevices(ss string) (inputDevice, outputDevice int) {
	devices := portmidi.CountDevices()
	fmt.Printf("Found %d MiDi Devices\n", devices)
	for i := 0; i < devices; i++ {
		di := portmidi.Info(portmidi.DeviceID(i))
		if strings.Contains(strings.ToLower(di.Name), strings.ToLower(ss)) {
			if di.IsInputAvailable {
				inputDevice = i
			}
			if di.IsOutputAvailable {
				outputDevice = i
			}
			fmt.Printf("Device %d:\n Name: %s\n isOpened: %v\n Interface %+v\n isInput: %v\n isOutput: %v\n", i, di.Name, di.IsOpened, di.Interface, di.IsInputAvailable, di.IsOutputAvailable)
		}
	}
	return
}

func cleanup() {
	sort.Ints(buttons)
	fmt.Printf("\nfound buttons: %v\n", buttons)
	fmt.Println("\nClosing portmidi device")
	if err := portmidi.Terminate(); err != nil {
		log.Printf("error terminating portmidi: %v\n", err)
	}
}

func main() {
	targetDevice := flag.String("device", "LPD8", "Target device")
	flag.Parse()
	if err := portmidi.Initialize(); err != nil {
		log.Printf("Error initializing portmidi: %v\n", err)
		return
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()
	defer func() {
		cleanup()
	}()
	inputDevice, _ := listDevices(*targetDevice)
	listen(inputDevice)
}
