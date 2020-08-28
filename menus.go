package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/c-bata/go-prompt"
	_ "github.com/logrusorgru/aurora"
	"github.com/marcusolsson/tui-go"
)

var (
	logo = `
	 ___ ___ _      ___       _        _                   
	| __/ __| |    |   \ __ _| |_ __ _| |__  __ _ ___ ___  
	| _| (_ | |__  | |) / _' |  _/ _' | '_ \/ _' (_-</ -_) 
	|_| \___|____| |___/\__,_|\__\__,_|_.__/\__,_/__/\___| 														   
	`
)

func registerMenu(opt string) {
	clear()

	username := tui.NewEntry()
	username.SetFocused(true)

	password := tui.NewEntry()
	password.SetEchoMode(tui.EchoModePassword)

	cPassword := tui.NewEntry()
	cPassword.SetEchoMode(tui.EchoModePassword)

	form := tui.NewGrid(0, 0)
	form.AppendRow(tui.NewLabel("Username"), tui.NewLabel("Password"), tui.NewLabel("Confirm Password"))
	form.AppendRow(username, password, cPassword)

	status := tui.NewStatusBar(opt)

	registerb := tui.NewButton("[Register]")

	loginb := tui.NewButton("[Login]")

	exitb := tui.NewButton("[Exit]")

	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, registerb),
		tui.NewPadder(1, 0, loginb),
		tui.NewPadder(1, 0, exitb),
	)

	window := tui.NewVBox(
		tui.NewPadder(10, 1, tui.NewLabel(logo)),
		tui.NewPadder(12, 0, tui.NewLabel("Welcome to the FGL Database")),
		tui.NewPadder(1, 1, form),
		buttons,
	)
	window.SetBorder(true)

	wrapper := tui.NewVBox(
		tui.NewSpacer(),
		window,
		tui.NewSpacer(),
	)
	content := tui.NewHBox(tui.NewSpacer(), wrapper, tui.NewSpacer())

	root := tui.NewVBox(
		content,
		status,
	)

	tui.DefaultFocusChain.Set(username, password, cPassword, registerb, loginb, exitb)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	registerb.OnActivated(func(b *tui.Button) {
		status.SetText("Registering...")
		if password.Text() != cPassword.Text() {
			registerMenu("Passwords Do Not Match")
		} else {
			creds := credentials{
				Username: username.Text(),
				Password: password.Text(),
			}
			tui.UI.Quit(ui)
			register(creds)
		}
	})

	loginb.OnActivated(func(b *tui.Button) {
		status.SetText("Login.")
		tui.UI.Quit(ui)
		loginMenu("")
	})

	exitb.OnActivated(func(b *tui.Button) {
		status.SetText("Exiting...")
		tui.UI.Quit(ui)
		terminate()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}

}

func loginMenu(opt string) {
	clear()

	username := tui.NewEntry()
	username.SetFocused(true)

	password := tui.NewEntry()
	password.SetEchoMode(tui.EchoModePassword)

	form := tui.NewGrid(0, 0)
	form.AppendRow(tui.NewLabel("User"), tui.NewLabel("Password"))
	form.AppendRow(username, password)

	status := tui.NewStatusBar(opt)

	loginb := tui.NewButton("[Login]")

	registerb := tui.NewButton("[Register]")

	exitb := tui.NewButton("[Exit]")

	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, loginb),
		tui.NewPadder(1, 0, registerb),
		tui.NewPadder(1, 0, exitb),
	)

	window := tui.NewVBox(
		tui.NewPadder(10, 1, tui.NewLabel(logo)),
		tui.NewPadder(12, 0, tui.NewLabel("Welcome to the FGL Database")),
		tui.NewPadder(1, 1, form),
		buttons,
	)
	window.SetBorder(true)

	wrapper := tui.NewVBox(
		tui.NewSpacer(),
		window,
		tui.NewSpacer(),
	)
	content := tui.NewHBox(tui.NewSpacer(), wrapper, tui.NewSpacer())

	root := tui.NewVBox(
		content,
		status,
	)

	tui.DefaultFocusChain.Set(username, password, loginb, registerb, exitb)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	loginb.OnActivated(func(b *tui.Button) {
		status.SetText("Logging in...")
		creds := credentials{
			Username: username.Text(),
			Password: password.Text(),
		}
		tui.UI.Quit(ui)
		login(creds)
	})

	registerb.OnActivated(func(b *tui.Button) {
		status.SetText("Register")
		tui.UI.Quit(ui)
		registerMenu("")
	})

	exitb.OnActivated(func(b *tui.Button) {
		status.SetText("Exiting...")
		tui.UI.Quit(ui)
		terminate()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func mainMenu(opt string) {
	clear()
	fmt.Println(opt)
	//fmt.Println(Red("FGL Database v" + version))
	fmt.Println()
	fmt.Println("[0] Announcements")
	fmt.Println("[1] Lab Member Reports")
	fmt.Println("[2] Future Gadget Reports")
	fmt.Println("[3] Shell")
	fmt.Println("[4] Chat")
	fmt.Println("[5] Exit")
	fmt.Println()
	fmt.Println()

	c := prompt.Choose(makeMainPrompt(), []string{})

	switch c {
	case "0":
		announcementsPreMenu()
	case "1":
		lmReportsPreMenu()
	case "2":
		fgReportsPreMenu()
	case "3":
		shell()
	case "4":
		chat()
	case "5":
		terminate()
	default:
		mainMenu("")
	}
}

// Announcement

func announcementsPreMenu() {
	clear()
	fmt.Println()
	fmt.Println("Announcements")
	fmt.Println()
	fmt.Println("[0] View Announcements")
	fmt.Println("[1] Make Announcement")
	fmt.Println("[2] Back")
	fmt.Println()
	fmt.Println()

	c := prompt.Choose(makeMainPrompt(), []string{})

	switch c {
	case "0":
		listAnnouncements()
	case "1":
		makeAnnouncement()
	case "2":
		mainMenu("")
	default:
		announcementsPreMenu()
	}
}

func viewAnnouncement(ann announcement) {
	clear()
	fmt.Println()
	fmt.Println()
	//fmt.Println(Red(ann.Author))
	fmt.Println(ann.Title)
	fmt.Println()
	fmt.Println(wordWrap(ann.Body, 70))
	fmt.Println()
	fmt.Println()

	ap := prompt.Choose(makeMainPrompt(), []string{})

	switch ap {
	default:
		listAnnouncements()
	}
}

func makeAnnouncement() {
	clear()
	fmt.Print("\nMake Announcement\n\n")
	title := prompt.Input("Title: ", emptyCompleter)
	body := prompt.Input("Body: ", emptyCompleter)

	ann := announcement{
		Author: currentUser,
		Title:  title,
		Body:   body,
	}

	sendAnnouncement(ann)
}

func listAnnouncements() {
	clear()

	announcements := make([]announcement, 0)

	resp, err := http.Get(baseURL("/announcements"))
	if err != nil {
		mainMenu("GET: Internal Server Error")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &announcements)

	defer resp.Body.Close()

	choices := make([]string, 0)

	fmt.Println()
	fmt.Println()
	for i, ann := range announcements {
		fmt.Println("[" + strconv.Itoa(i) + "] " + ann.Title + " - " + ann.Author)
		choices = append(choices, strconv.Itoa(i))
	}
	fmt.Println()
	fmt.Println()
	a := prompt.Choose(makeMainPrompt(), []string{})

	switch a {
	case "":
		mainMenu("")
	case "back":
		mainMenu("")
	case "exit":
		mainMenu("")
	default:
		fmt.Println(a)
		ann := announcement{}

		aint, _ := strconv.ParseInt(a, 10, 32)
		aint++
		astr := strconv.Itoa(int(aint))

		annr, err := http.Get(baseURL("/announcement/" + astr))
		if err != nil {
			listAnnouncements()
		}
		abody, err := ioutil.ReadAll(annr.Body)
		if err != nil {
			listAnnouncements()
		}
		asdf := json.Unmarshal(abody, &ann)
		if asdf != nil {
			listAnnouncements()
		}
		defer annr.Body.Close()
		viewAnnouncement(ann)
	}
}

// Lab Member Report

func lmReportsPreMenu() {
	clear()
	fmt.Println()
	fmt.Println("Lab Member Reports")
	fmt.Println()
	fmt.Println("[0] View Reports")
	fmt.Println("[1] Submit Report")
	fmt.Println("[2] Back")
	fmt.Println()
	fmt.Println()

	c := prompt.Choose(makeMainPrompt(), []string{})

	switch c {
	case "0":
		listLmReports()
	case "1":
		makeLmReport()
	case "2":
		mainMenu("")
	default:
		lmReportsPreMenu()
	}
}

func viewLmReport(lr labReport) {
	clear()
	fmt.Println()
	fmt.Println()
	//fmt.Println(Red(lr.Author))
	fmt.Println(lr.Title)
	fmt.Println()
	fmt.Println(lr.Subject)
	fmt.Println()
	fmt.Println(wordWrap(lr.Body, 70))
	fmt.Println()
	fmt.Println()

	ap := prompt.Choose(makeMainPrompt(), []string{})

	switch ap {
	default:
		listLmReports()
	}
}

func makeLmReport() {
	clear()
	fmt.Print("\nMake Lab Member Report\n\n")
	title := prompt.Input("Title: ", emptyCompleter)
	subject := prompt.Input("Subject: ", emptyCompleter)
	body := prompt.Input("Body: ", emptyCompleter)

	lr := labReport{
		Author:  currentUser,
		Title:   title,
		Subject: subject,
		Body:    body,
	}

	sendLmReport(lr)
}

func listLmReports() {
	clear()

	lmr := make([]labReport, 0)

	resp, err := http.Get(baseURL("/labreports"))
	if err != nil {
		mainMenu("GET: Internal Server Error")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &lmr)

	defer resp.Body.Close()

	choices := make([]string, 0)

	fmt.Println()
	fmt.Println()
	for i, lm := range lmr {
		fmt.Println("[" + strconv.Itoa(i) + "] " + lm.Title + " - " + lm.Author)
		choices = append(choices, strconv.Itoa(i))
	}
	fmt.Println()
	fmt.Println()
	a := prompt.Choose(makeMainPrompt(), []string{})

	switch a {
	case "":
		mainMenu("")
	case "back":
		mainMenu("")
	case "exit":
		mainMenu("")
	default:
		fmt.Println(a)
		lm := labReport{}

		lint, _ := strconv.ParseInt(a, 10, 32)
		lint++
		lstr := strconv.Itoa(int(lint))

		lmrr, err := http.Get(baseURL("/labreport/" + lstr))
		if err != nil {
			listLmReports()
		}
		lbody, err := ioutil.ReadAll(lmrr.Body)
		if err != nil {
			listLmReports()
		}
		asdf := json.Unmarshal(lbody, &lm)
		if asdf != nil {
			listLmReports()
		}
		defer lmrr.Body.Close()
		viewLmReport(lm)
	}
}

// Future Gadget Report

func fgReportsPreMenu() {
	clear()
	fmt.Println()
	fmt.Println("Future Gadget Reports")
	fmt.Println()
	fmt.Println("[0] View Reports")
	fmt.Println("[1] Submit Report")
	fmt.Println("[2] Back")
	fmt.Println()
	fmt.Println()

	c := prompt.Choose(makeMainPrompt(), []string{})

	switch c {
	case "0":
		listFgReports()
	case "1":
		makeFgReport()
	case "2":
		mainMenu("")
	default:
		fgReportsPreMenu()
	}
}

func viewFgReport(fg gadgetReport) {
	clear()
	fmt.Println()
	fmt.Println()
	//fmt.Println(Red(fg.Author))
	fmt.Println(fg.Title)
	fmt.Println()
	fmt.Println(wordWrap(fg.Body, 70))
	fmt.Println()
	fmt.Println()

	ap := prompt.Choose(makeMainPrompt(), []string{})

	switch ap {
	default:
		listFgReports()
	}
}

func makeFgReport() {
	clear()
	fmt.Print("\nMake Future Gadget Report\n\n")
	title := prompt.Input("Title: ", emptyCompleter)
	body := prompt.Input("Body: ", emptyCompleter)

	fg := gadgetReport{
		Author: currentUser,
		Title:  title,
		Body:   body,
	}

	sendFgReport(fg)
}

func listFgReports() {
	clear()

	fgr := make([]gadgetReport, 0)

	resp, err := http.Get(baseURL("/gadgetreports"))
	if err != nil {
		mainMenu("GET: Internal Server Error")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &fgr)

	defer resp.Body.Close()

	choices := make([]string, 0)

	fmt.Println()
	fmt.Println()
	for i, fg := range fgr {
		fmt.Println("[" + strconv.Itoa(i) + "] " + fg.Title + " - " + fg.Author)
		choices = append(choices, strconv.Itoa(i))
	}
	fmt.Println()
	fmt.Println()
	a := prompt.Choose(makeMainPrompt(), []string{})

	switch a {
	case "":
		mainMenu("")
	case "back":
		mainMenu("")
	case "exit":
		mainMenu("")
	default:
		fmt.Println(a)
		fg := gadgetReport{}

		fint, _ := strconv.ParseInt(a, 10, 32)
		fint++
		fstr := strconv.Itoa(int(fint))

		fgrr, err := http.Get(baseURL("/gadgetreport/" + fstr))
		if err != nil {
			listFgReports()
		}
		fbody, err := ioutil.ReadAll(fgrr.Body)
		if err != nil {
			listFgReports()
		}
		asdf := json.Unmarshal(fbody, &fg)
		if asdf != nil {
			listFgReports()
		}
		defer fgrr.Body.Close()
		viewFgReport(fg)
	}
}

func chat() {
	clear()
	fmt.Println()
	fmt.Println()
	fmt.Println("Not implemented yet, sorry! Press enter to return to the main menu.")
	fmt.Println()
	fmt.Println()
	p := prompt.Input("", emptyCompleter)
	fmt.Println(p)
	mainMenu("")
}
