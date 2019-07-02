package main

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/jacobsa/go-serial/serial"
)

var mainwin *ui.Window
var portListEntry *ui.Entry
var sendButton *ui.Button
var infoLabel *ui.Label
var info1 = "Fill in the fields, choose the port and press the \"PROGRAM\" button."
var paramError int

type paramdata struct {
	deveuientry *ui.Entry
	appeuientry *ui.Entry
	appkeyentry *ui.Entry
}

var Data paramdata
var port io.ReadWriteCloser

func checkPorts() {

}

func forTesting() string {
	return "test_OK"
}

func checkConfigState(c []byte, s string) bool {
	if strings.Contains(string(c), s) {
		//fmt.Println(string(c))
		return true
	}
	return false
}

func sendCommand(p io.ReadWriteCloser, s string) {
	b := []byte(s)
	p.Write(b)
}

func waitSignalProgram(p string, d paramdata) {
	fmt.Println("TURN ON device")
	showInfo("TURN ON device!")
	sendButton.Disable()
	// Set up options.
	options := serial.OpenOptions{
		PortName:        p,
		BaudRate:        57600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	// Make sure to close it later.
	defer port.Close()

	cIndex := 0
	stringBuf := make([]byte, 64)
	for {

		charBuf := make([]byte, 1)

		_, err := port.Read(charBuf)
		if err != nil {
			fmt.Println("Error")
		}

		if charBuf[0] != byte(0) {
			//fmt.Printf("%d %c\r\n", charBuf, charBuf)    \r=10  \n=13
			//fmt.Printf("%s\r\n", stringBuf)
			stringBuf[cIndex] = charBuf[0]
			cIndex++

			if charBuf[0] == byte(10) {
				fmt.Printf("%s\r\n", stringBuf)

				if checkConfigState(stringBuf, ">CONF\r\n") {
					fmt.Println("**** ENTER PROGRAM MODE ****")
					sendCommand(port, "config\r\n")
				}
				if checkConfigState(stringBuf, ">DEVEUI\r\n") {
					fmt.Println("**** DEVEUI ****")
					time.Sleep(50 * time.Millisecond)
					strdeveui := Data.deveuientry.Text()
					sendCommand(port, strdeveui+"\r\n")
				}
				if checkConfigState(stringBuf, ">APPEUI\r\n") {
					fmt.Println("**** APPEUI ****")
					time.Sleep(50 * time.Millisecond)
					strappeui := Data.appeuientry.Text()
					sendCommand(port, strappeui+"\r\n")
				}
				if checkConfigState(stringBuf, ">APPKEY\r\n") {
					fmt.Println("**** APPKEY ****")
					time.Sleep(50 * time.Millisecond)
					strappkey := Data.appkeyentry.Text()
					sendCommand(port, strappkey+"\r\n")
				}
				if checkConfigState(stringBuf, ">FINISH\r\n") {
					fmt.Println("**** FINISH ****")
					showInfo("DEVICE SUCCESFULLY PROGRAMMED!")
					sendButton.Enable()
					return
				}
				cIndex = 0
				stringBuf = make([]byte, 64)
			}

		}

	}

}

func showError(s string) {
	infoLabel.SetText("[ERROR]  " + s)
}

func showInfo(s string) {
	infoLabel.SetText("[INFO]  " + s)
}

func makeBasicControlsPage() ui.Control {
	fullBox := ui.NewVerticalBox()
	fullBox.SetPadded(true)

	/*------------- group 1 ----------------*/
	group := ui.NewGroup("Join Parameters")
	group.SetMargined(true)
	fullBox.Append(group, false)

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)

	Data.deveuientry = ui.NewEntry()
	Data.deveuientry.SetText("0000000000000000")
	Data.appeuientry = ui.NewEntry()
	Data.appeuientry.SetText("1111111111111111")
	Data.appkeyentry = ui.NewEntry()
	Data.appkeyentry.SetText("22222222222222223333333333333333")

	Data.deveuientry.OnChanged(func(*ui.Entry) {
		showInfo(info1)
	})
	Data.appeuientry.OnChanged(func(*ui.Entry) {
		showInfo(info1)
	})

	Data.appkeyentry.OnChanged(func(*ui.Entry) {
		showInfo(info1)
	})

	entryForm.Append("DEVEUI  ", Data.deveuientry, false)
	entryForm.Append("APPEUI  ", Data.appeuientry, false)
	entryForm.Append("APPKEY  ", Data.appkeyentry, false)

	/*------------- group 3 ----------------*/

	portBox := ui.NewHorizontalBox()
	portBox.SetPadded(true)
	fullBox.Append(portBox, false)

	portBox.Append(ui.NewLabel("PORT:"), false)

	sendButton = ui.NewButton("PROGRAM")
	sendButton.Disable()

	portListEntry := ui.NewEntry()
	portListEntry.OnChanged(func(*ui.Entry) {
		re := regexp.MustCompile("COM[0-9]") // Check compatible port on windows
		if re.MatchString(portListEntry.Text()) {
			fmt.Println("port ready")
			sendButton.Enable()
		} else {
			sendButton.Disable()
		}
	})
	portBox.Append(portListEntry, false)

	sendButton.OnClicked(func(*ui.Button) {
		fmt.Println("Send button pressed!")
		lendeveui := len(Data.deveuientry.Text())
		lenappeui := len(Data.appeuientry.Text())
		lenappkey := len(Data.appkeyentry.Text())

		if lendeveui != 16 || lenappeui != 16 || lenappkey != 32 {
			fmt.Println("[ERROR] invalid len parameters.")
			showError("invalid len parameters.")
		} else {
			fmt.Println(portListEntry.Text())
			showInfo("Attempting to program device!")
			//TODO: CONNECT TO SERIAL
			go waitSignalProgram(portListEntry.Text(), Data)

		}

	})

	space := ui.NewLabel("    ")
	portBox.Append(space, false)
	portBox.Append(sendButton, true)

	/*------------- group 3 ----------------*/
	fullBox.Append(ui.NewVerticalSeparator(), false)
	infoLabel = ui.NewLabel(" ")
	showInfo(info1)
	fullBox.Append(infoLabel, false)
	return fullBox
}

func setupUI() {

	paramError = 1
	checkPorts()

	mainwin = ui.NewWindow("LoRaWAN Node Configurator v0.1", 480, 240, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	tab.Append("OTAA parameters", makeBasicControlsPage())
	tab.SetMargined(0, true)

	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
	fmt.Println(".")
}
