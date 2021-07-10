package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	_ "github.com/logrusorgru/aurora"
	"github.com/marcusolsson/tui-go"

	"github.com/rivo/tview"
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
		tui.NewPadder(20, 0, tui.NewLabel("Welcome to the FGL Database")),
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

	cPassword.OnSubmit(func(e *tui.Entry) {
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
		tui.NewPadder(20, 0, tui.NewLabel("Welcome to the FGL Database")),
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

	password.OnSubmit(func(e *tui.Entry) {
		status.SetText("Logging in...")
		creds := credentials{
			Username: username.Text(),
			Password: password.Text(),
		}
		tui.UI.Quit(ui)
		login(creds)
	})

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
	app := tview.NewApplication()
	header := tview.NewTextView().SetText("FGL Database v" + string(version))
	list := tview.NewList().
		AddItem("Announcements", "", '0', func() {
			app.Stop()
			announcementsPreMenu()
		}).
		AddItem("Lab Member Reports", "", '1', func() {
			app.Stop()
			lmReportsPreMenu()
		}).
		AddItem("Future Gadget Reports", "", '2', func() {
			app.Stop()
			fgReportsPreMenu()
		}).
		AddItem("Shell", "", '3', func() {
			app.Stop()
			shell()
		}).
		AddItem("Chat", "", '4', func() {
			app.Stop()
			chat()
		}).
		AddItem("Exit", "", 'q', func() {
			app.Stop()
			terminate()
		})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(list, 0, 10, true)

	if err := app.SetRoot(flex, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}

func announcementsPreMenu() {
	app := tview.NewApplication()
	header := tview.NewTextView().SetText("Announcements")
	list := tview.NewList().
		AddItem("View Announcements", "", '0', func() {
			app.Stop()
			listAnnouncements()
		}).
		AddItem("Make Announcement", "", '1', func() {
			app.Stop()
			makeAnnouncement("")
		}).
		AddItem("Back", "", 'b', func() {
			app.Stop()
			mainMenu("")
		})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(list, 0, 10, true)

	if err := app.SetRoot(flex, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}

func viewAnnouncement(id string, opt string) {
	clear()
	ann := getAnnouncement(id)

	status := tui.NewStatusBar(opt)

	edit := tui.NewButton("[Edit]")

	delete := tui.NewButton("[Delete]")

	back := tui.NewButton("[Back]")

	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, edit),
		tui.NewPadder(1, 0, delete),
		tui.NewPadder(1, 0, back),
	)

	annBody := tui.NewVBox(
		tui.NewLabel(wordWrap(ann.Body, 70)),
	)

	window := tui.NewVBox(
		tui.NewPadder(20, 1, tui.NewLabel(ann.Title)),
		annBody,
		tui.NewSpacer(),
		tui.NewSpacer(),
		buttons,
	)
	window.SetBorder(false)

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

	tui.DefaultFocusChain.Set(back, edit, delete)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	edit.OnActivated(func(b *tui.Button) {
		if ann.Author != currentUser {
			ui.Quit()
			s := strconv.FormatInt(ann.AID, 10)
			viewAnnouncement(s, "You cannot edit another member's announcement.")
		} else {
			ui.Quit()
			s := strconv.FormatInt(ann.AID, 10)
			editAnnouncement(s)
		}
	})

	delete.OnActivated(func(b *tui.Button) {
		if ann.Author != currentUser {
			ui.Quit()
			s := strconv.FormatInt(ann.AID, 10)
			viewAnnouncement(s, "You cannot delete another member's announcement.")
		} else {
			ui.Quit()
			s := strconv.FormatInt(ann.AID, 10)
			deleteAnnouncement(s)
			listAnnouncements()
		}
	})

	back.OnActivated(func(b *tui.Button) {
		ui.Quit()
		listAnnouncements()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func editAnnouncement(id string) {
	clear()

	ann := getAnnouncement(id)

	title := tui.NewEntry()
	title.SetFocused(true)
	title.SetText(ann.Title)

	body := tui.NewEntry()
	body.SetText(ann.Body)

	form := tui.NewGrid(1, 2)
	form.AppendRow(tui.NewLabel("Title"))
	form.AppendRow(title)
	form.AppendRow(tui.NewLabel("Body"))
	form.AppendRow(body)

	status := tui.NewStatusBar("")

	submit := tui.NewButton("[Submit]")

	cancel := tui.NewButton("[Cancel]")

	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, submit),
		tui.NewPadder(1, 0, cancel),
	)

	window := tui.NewVBox(
		tui.NewPadder(12, 0, tui.NewLabel("Edit Announcement")),
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

	tui.DefaultFocusChain.Set(title, body, submit, cancel)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	submit.OnActivated(func(b *tui.Button) {
		status.SetText("Submitting...")
		a := announcement{
			Author: ann.Author,
			Title:  title.Text(),
			Body:   body.Text(),
			AID:    ann.AID,
		}
		ui.Quit()
		s := strconv.FormatInt(ann.AID, 10)
		updateAnnouncement(s, a)
		viewAnnouncement(s, "")
	})

	cancel.OnActivated(func(b *tui.Button) {
		status.SetText("Cancelling...")
		ui.Quit()
		s := strconv.FormatInt(ann.AID, 10)
		viewAnnouncement(s, "")
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func makeAnnouncement(opt string) {
	clear()

	title := tui.NewEntry()
	title.SetFocused(true)

	body := tui.NewEntry()

	form := tui.NewGrid(1, 2)
	form.AppendRow(tui.NewLabel("Title"))
	form.AppendRow(title)
	form.AppendRow(tui.NewSpacer())
	form.AppendRow(tui.NewLabel("Body"))
	form.AppendRow(body)

	status := tui.NewStatusBar(opt)

	submit := tui.NewButton("[Submit]")

	cancel := tui.NewButton("[Cancel]")

	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, submit),
		tui.NewPadder(1, 0, cancel),
	)

	window := tui.NewVBox(
		tui.NewPadder(12, 0, tui.NewLabel("Submit a New Announcement")),
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

	tui.DefaultFocusChain.Set(title, body, submit, cancel)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	submit.OnActivated(func(b *tui.Button) {
		status.SetText("Submitting...")
		ann := announcement{
			Author: currentUser,
			Title:  title.Text(),
			Body:   body.Text(),
		}
		ui.Quit()
		sendAnnouncement(ann)
	})

	cancel.OnActivated(func(b *tui.Button) {
		status.SetText("Cancelling...")
		ui.Quit()
		announcementsPreMenu()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func listAnnouncements() {
	announcements := getAnnouncements()

	app := tview.NewApplication()
	header := tview.NewTextView().SetText("Announcements")
	list := tview.NewList()

	for i, ann := range announcements {
		a := ann
		list.InsertItem(i, ann.Title+" - "+ann.Author, "", ' ', func() {
			app.Stop()
			s := strconv.FormatInt(a.AID, 10)
			viewAnnouncement(s, "")
		})
	}

	list.InsertItem(len(announcements), "Back", "", 'b', func() {
		app.Stop()
		announcementsPreMenu()
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(list, 0, 10, true)

	if err := app.SetRoot(flex, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}

// Lab Reports

func lmReportsPreMenu() {
	app := tview.NewApplication()
	header := tview.NewTextView().SetText("Lab Member Reports")
	list := tview.NewList().
		AddItem("View Reports", "", '0', func() {
			app.Stop()
			listLmReports()
		}).
		AddItem("Submit Report", "", '1', func() {
			app.Stop()
			makeLmReport()
		}).
		AddItem("Back", "", 'b', func() {
			app.Stop()
			mainMenu("")
		})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(list, 0, 10, true)

	if err := app.SetRoot(flex, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}

func viewLmReport(id string, opt string) {
	clear()
	lr := getLmReport(id)

	status := tui.NewStatusBar(opt)

	edit := tui.NewButton("[Edit]")

	delete := tui.NewButton("[Delete]")

	back := tui.NewButton("[Back]")

	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, edit),
		tui.NewPadder(1, 0, delete),
		tui.NewPadder(1, 0, back),
	)

	lrBody := tui.NewVBox(
		tui.NewLabel(wordWrap(lr.Body, 70)),
	)

	window := tui.NewVBox(
		tui.NewPadder(20, 1, tui.NewLabel(lr.Title)),
		tui.NewPadder(0, 1, tui.NewLabel(lr.Subject)),
		lrBody,
		tui.NewSpacer(),
		tui.NewSpacer(),
		buttons,
	)
	window.SetBorder(false)

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

	tui.DefaultFocusChain.Set(back, edit, delete)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	edit.OnActivated(func(b *tui.Button) {
		if lr.Author != currentUser {
			ui.Quit()
			s := strconv.FormatInt(lr.RID, 10)
			viewLmReport(s, "You cannot edit another member's Lab Report.")
		} else {
			ui.Quit()
			s := strconv.FormatInt(lr.RID, 10)
			editLmReport(s)
		}
	})

	delete.OnActivated(func(b *tui.Button) {
		if lr.Author != currentUser {
			ui.Quit()
			s := strconv.FormatInt(lr.RID, 10)
			viewLmReport(s, "You cannot delete another member's Lab Report.")
		} else {
			ui.Quit()
			s := strconv.FormatInt(lr.RID, 10)
			deleteLmReport(s)
			listLmReports()
		}
	})

	back.OnActivated(func(b *tui.Button) {
		ui.Quit()
		listLmReports()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func makeLmReport() {
	clear()

	title := tui.NewEntry()
	title.SetFocused(true)

	subject := tui.NewEntry()
	body := tui.NewEntry()

	form := tui.NewGrid(1, 2)
	form.AppendRow(tui.NewLabel("Title"))
	form.AppendRow(title)
	form.AppendRow(tui.NewSpacer())
	form.AppendRow(tui.NewLabel("Subject"))
	form.AppendRow(subject)
	form.AppendRow(tui.NewSpacer())
	form.AppendRow(tui.NewLabel("Body"))
	form.AppendRow(body)

	status := tui.NewStatusBar("")

	submit := tui.NewButton("[Submit]")

	cancel := tui.NewButton("[Cancel]")

	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, submit),
		tui.NewPadder(1, 0, cancel),
	)

	window := tui.NewVBox(
		tui.NewPadder(12, 0, tui.NewLabel("Submit a New Lab Member Report")),
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

	tui.DefaultFocusChain.Set(title, subject, body, submit, cancel)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	submit.OnActivated(func(b *tui.Button) {
		status.SetText("Submitting...")
		lr := labReport{
			Author:  currentUser,
			Title:   title.Text(),
			Subject: subject.Text(),
			Body:    body.Text(),
		}
		ui.Quit()
		sendLmReport(lr)
	})

	cancel.OnActivated(func(b *tui.Button) {
		status.SetText("Cancelling...")
		ui.Quit()
		lmReportsPreMenu()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func editLmReport(id string) {
	clear()

	lr := getLmReport(id)

	title := tui.NewEntry()
	title.SetFocused(true)
	title.SetText(lr.Title)

	subject := tui.NewEntry()
	subject.SetText(lr.Subject)

	body := tui.NewEntry()
	body.SetText(lr.Body)

	form := tui.NewGrid(1, 2)
	form.AppendRow(tui.NewLabel("Title"))
	form.AppendRow(title)
	form.AppendRow(tui.NewSpacer())
	form.AppendRow(tui.NewLabel("Subject"))
	form.AppendRow(subject)
	form.AppendRow(tui.NewSpacer())
	form.AppendRow(tui.NewLabel("Body"))
	form.AppendRow(body)

	status := tui.NewStatusBar("")

	submit := tui.NewButton("[Submit]")

	cancel := tui.NewButton("[Cancel]")

	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, submit),
		tui.NewPadder(1, 0, cancel),
	)

	window := tui.NewVBox(
		tui.NewPadder(12, 0, tui.NewLabel("Edit Lab Report")),
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

	tui.DefaultFocusChain.Set(title, subject, body, submit, cancel)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	submit.OnActivated(func(b *tui.Button) {
		status.SetText("Submitting...")
		l := labReport{
			Author:  lr.Author,
			Title:   title.Text(),
			Subject: subject.Text(),
			Body:    body.Text(),
			RID:     lr.RID,
		}
		tui.UI.Quit(ui)
		s := strconv.FormatInt(lr.RID, 10)
		updateLmReport(s, l)
		viewLmReport(s, "")
	})

	cancel.OnActivated(func(b *tui.Button) {
		status.SetText("Cancelling...")
		tui.UI.Quit(ui)
		s := strconv.FormatInt(lr.RID, 10)
		viewLmReport(s, "")
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func listLmReports() {
	labReports := getLmReports()

	app := tview.NewApplication()
	header := tview.NewTextView().SetText("Lab Member Reports")
	list := tview.NewList()

	for i, lr := range labReports {
		l := lr
		list.InsertItem(i, l.Title+" - "+l.Author, "", ' ', func() {
			app.Stop()
			s := strconv.FormatInt(l.RID, 10)
			viewLmReport(s, "")
		})
	}

	list.InsertItem(len(labReports), "Back", "", 'b', func() {
		app.Stop()
		lmReportsPreMenu()
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(list, 0, 10, true)

	if err := app.SetRoot(flex, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}

// Future Gadget Report

func fgReportsPreMenu() {
	app := tview.NewApplication()
	header := tview.NewTextView().SetText("Future Gadget Reports")
	list := tview.NewList().
		AddItem("View Reports", "", '0', func() {
			app.Stop()
			listFgReports()
		}).
		AddItem("Submit Report", "", '1', func() {
			app.Stop()
			makeFgReport()
		}).
		AddItem("Back", "", 'b', func() {
			app.Stop()
			mainMenu("")
		})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(list, 0, 10, true)

	if err := app.SetRoot(flex, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}

func viewFgReport(id string, opt string) {
	clear()
	gr := getFgReport(id)

	status := tui.NewStatusBar(opt)

	edit := tui.NewButton("[Edit]")

	delete := tui.NewButton("[Delete]")

	back := tui.NewButton("[Back]")

	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, edit),
		tui.NewPadder(1, 0, delete),
		tui.NewPadder(1, 0, back),
	)

	grBody := tui.NewVBox(
		tui.NewLabel(wordWrap(gr.Body, 70)),
	)

	window := tui.NewVBox(
		tui.NewPadder(20, 1, tui.NewLabel(gr.Title)),
		grBody,
		tui.NewSpacer(),
		tui.NewSpacer(),
		buttons,
	)
	window.SetBorder(false)

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

	tui.DefaultFocusChain.Set(back, edit, delete)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	edit.OnActivated(func(b *tui.Button) {
		if gr.Author != currentUser {
			ui.Quit()
			s := strconv.FormatInt(gr.RID, 10)
			viewFgReport(s, "You cannot edit another member's Future Gadget Report.")
		} else {
			ui.Quit()
			s := strconv.FormatInt(gr.RID, 10)
			editFgReport(s)
		}
	})

	delete.OnActivated(func(b *tui.Button) {
		if gr.Author != currentUser {
			ui.Quit()
			s := strconv.FormatInt(gr.RID, 10)
			viewFgReport(s, "You cannot delete another member's Future Gadget Report.")
		} else {
			ui.Quit()
			s := strconv.FormatInt(gr.RID, 10)
			deleteFgReport(s)
			listFgReports()
		}
	})

	back.OnActivated(func(b *tui.Button) {
		ui.Quit()
		listFgReports()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func makeFgReport() {
	clear()

	title := tui.NewEntry()
	title.SetFocused(true)

	body := tui.NewEntry()

	form := tui.NewGrid(1, 2)
	form.AppendRow(tui.NewLabel("Title"))
	form.AppendRow(title)
	form.AppendRow(tui.NewSpacer())
	form.AppendRow(tui.NewLabel("Body"))
	form.AppendRow(body)

	status := tui.NewStatusBar("")

	submit := tui.NewButton("[Submit]")

	cancel := tui.NewButton("[Cancel]")

	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, submit),
		tui.NewPadder(1, 0, cancel),
	)

	window := tui.NewVBox(
		tui.NewPadder(12, 0, tui.NewLabel("Submit a New Future Gadget Report")),
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

	tui.DefaultFocusChain.Set(title, body, submit, cancel)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	submit.OnActivated(func(b *tui.Button) {
		status.SetText("Submitting...")
		gr := gadgetReport{
			Author: currentUser,
			Title:  title.Text(),
			Body:   body.Text(),
		}
		ui.Quit()
		sendFgReport(gr)
		listFgReports()
	})

	cancel.OnActivated(func(b *tui.Button) {
		status.SetText("Cancelling...")
		ui.Quit()
		fgReportsPreMenu()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func editFgReport(id string) {
	clear()

	gr := getFgReport(id)

	title := tui.NewEntry()
	title.SetFocused(true)
	title.SetText(gr.Title)

	body := tui.NewEntry()
	body.SetText(gr.Body)

	form := tui.NewGrid(1, 2)
	form.AppendRow(tui.NewLabel("Title"))
	form.AppendRow(title)
	form.AppendRow(tui.NewSpacer())
	form.AppendRow(tui.NewLabel("Body"))
	form.AppendRow(body)

	status := tui.NewStatusBar("")

	submit := tui.NewButton("[Submit]")

	cancel := tui.NewButton("[Cancel]")

	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, submit),
		tui.NewPadder(1, 0, cancel),
	)

	window := tui.NewVBox(
		tui.NewPadder(12, 0, tui.NewLabel("Edit Future Gadget Report")),
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

	tui.DefaultFocusChain.Set(title, body, submit, cancel)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	submit.OnActivated(func(b *tui.Button) {
		status.SetText("Submitting...")
		g := gadgetReport{
			Author: gr.Author,
			Title:  title.Text(),
			Body:   body.Text(),
			RID:    gr.RID,
		}
		ui.Quit()
		s := strconv.FormatInt(gr.RID, 10)
		updateFgReport(s, g)
		viewFgReport(s, "")
	})

	cancel.OnActivated(func(b *tui.Button) {
		status.SetText("Cancelling...")
		ui.Quit()
		s := strconv.FormatInt(gr.RID, 10)
		viewFgReport(s, "")
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func listFgReports() {
	gadgetReports := getFgReports()

	app := tview.NewApplication()
	header := tview.NewTextView().SetText("Future Gadget Reports")
	list := tview.NewList()

	for i, gr := range gadgetReports {
		g := gr
		list.InsertItem(i, g.Title+" - "+g.Author, "", ' ', func() {
			app.Stop()
			s := strconv.FormatInt(g.RID, 10)
			viewFgReport(s, "")
		})
	}

	list.InsertItem(len(gadgetReports), "Back", "", 'b', func() {
		app.Stop()
		fgReportsPreMenu()
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(list, 0, 10, true)

	if err := app.SetRoot(flex, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}

func chat() {
	clear()
	//host := "159.89.8.129"
	host := "fglteam.com"
	port := "3333"
	connect := host + ":" + port
	c, err := net.Dial("tcp", connect)
	fmt.Fprintf(c, "/join:"+currentUser+"\n")
	if err != nil {
		mainMenu(err.Error())
	}

	online := tui.NewVBox(
		tui.NewLabel("ONLINE"),
	)
	online.SetSizePolicy(tui.Expanding, tui.Expanding)

	offline := tui.NewVBox(
		tui.NewLabel("OFFLINE"),
	)
	offline.SetSizePolicy(tui.Expanding, tui.Expanding)

	sidebar := tui.NewVBox(
		online,
		tui.NewSpacer(),
		offline,
	)
	sidebar.SetBorder(true)

	history := tui.NewVBox()

	pastMessages := getMessages()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	root := tui.NewHBox(chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	var once sync.Once
	onceBody := func() {
		ui.Quit()
		c.Close()
		mainMenu("")
	}

	ui.SetKeybinding("Up", func() {
		historyScroll.Scroll(0, -1)
	})

	ui.SetKeybinding("Down", func() {
		historyScroll.Scroll(0, 1)
	})

	ui.SetKeybinding("Esc", func() {
		once.Do(onceBody)
	})

	for _, msg := range pastMessages {
		snd := msg.Author + ": " + msg.Body
		if msg.Author == "SERVER" {
			history.Append(tui.NewHBox(
				tui.NewLabel(wordWrap2(msg.Body, 80)),
				tui.NewSpacer(),
			))
		} else {
			history.Append(tui.NewHBox(
				tui.NewLabel(wordWrap2(snd, 80)),
				tui.NewSpacer(),
			))
		}
	}

	go func() {
		for {
			ui.Update(func() {})
			input.OnSubmit(func(e *tui.Entry) {
				if len(e.Text()) >= 1 {
					fmt.Fprintf(c, e.Text()+"\n")
				}
				input.SetText("")
			})

			message, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				once.Do(onceBody)
			}
			history.Append(tui.NewHBox(
				tui.NewLabel(wordWrap2(message, 80)),
				tui.NewSpacer(),
			))
		}
	}()

	if err := ui.Run(); err != nil {
		fmt.Println("WHOA BIG ERROR HERE uWu IT'S A FUCKY WUCKY")
		log.Fatal(err)
	}
}
