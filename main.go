package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/c-bata/go-prompt"
)

var (
	history     = prompt.NewHistory()
	currentUser = ""
	authCode    = ""
	mClear      map[string]func()
	currentDir  = directory{
		Name: "~",
	}
	version     = "0.4"
	server      string
	checkUpdate bool
)

// Version object to more easily read version
type Version struct {
	Version string `json:"version"`
}

func main() {
	//setName()
	loadConfig()
	if checkUpdate {
		update()
	}
	fglInit()
}

func setName() {
	cmd := exec.Command("title", "FGL-Client")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func loadConfig() {
	server = "www.fglteam.com"
	checkUpdate = false
}

func init() {
	mClear = make(map[string]func()) //Initialize it
	mClear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	mClear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func fglInit() {
	loginMenu("")
	//chat()
}

func update() {
	clear()

	fmt.Println()
	fmt.Println()
	fmt.Println("Checking For Update...")
	fmt.Println()
	fmt.Println()

	cVersion := Version{Version: version}
	cv, _ := json.Marshal(cVersion)
	resp, err := http.Post(baseURL("/recversion"), "application/json", bytes.NewBuffer(cv))

	if err != nil {
		fmt.Println("Unable to connect to server.")
		fmt.Println()
		fmt.Println()
		fmt.Println("Press enter to close.")
		t := prompt.Input("", emptyCompleter)
		fmt.Print(t)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusUpgradeRequired:
		fmt.Println("Update Available. Press Enter To Download Updater..")
		fmt.Println()
		fmt.Println()
		t := prompt.Input("", emptyCompleter)
		fmt.Print(t)
		getUpdater()
	case http.StatusOK:
		fmt.Println("You Have The Latest Version. Press Enter To Continue.")
		fmt.Println()
		fmt.Println()
		t := prompt.Input("", emptyCompleter)
		fmt.Print(t)
		fglInit()
	case http.StatusBadRequest:
		loginMenu("Internal Server Error")
	default:
		fmt.Println(http.StatusText(resp.StatusCode))
	}
}

func getUpdater() {
	clear()
	fmt.Println()
	fmt.Println()
	fmt.Println("Obtaining Updater...")
	fmt.Println()

	fileURL := server + "/?file=fgl-updater.exe"
	filePath := "fgl-updater.exe"

	resp, err := http.Get(fileURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return
	}
	//defer out.Close()

	_, err = io.Copy(out, resp.Body)

	fmt.Println("Success")
	fmt.Println()

	up := prompt.Input("Unfortunately, due to limitations of Windows, you'll have to manually run the updater. Press enter to exit the application, then double click fgl-updater.exe.", emptyCompleter)
	fmt.Println(up)
	out.Close()
	terminate()
}

func terminate() {
	clear()
	os.Clearenv()
	os.Exit(0)
}
