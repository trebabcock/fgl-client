package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func login(creds credentials) {
	credentials, _ := json.Marshal(creds)
	resp, err := http.Post(baseURL("/login"), "application/json", bytes.NewBuffer(credentials))
	if err != nil {
		loginMenu("Internal Server Error")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case http.StatusNotFound:
		loginMenu("Incorrect Username or Password")
	case http.StatusUnauthorized:
		loginMenu("Incorrect Username Or Password")
	case http.StatusInternalServerError:
		loginMenu("Internal Server Error")
	case http.StatusOK:
		os.Setenv("FGL_TOKEN", string(body))
		currentUser = creds.Username
		mainMenu("")
	}
}

func register(creds credentials) {
	credentials, err := json.Marshal(creds)
	if err != nil {
		fmt.Println(err.Error())
	}
	resp, err := http.Post(baseURL("/register"), "application/json", bytes.NewBuffer(credentials))
	fmt.Println(string(credentials))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		loginMenu("Internal Server Error")
	case http.StatusInternalServerError:
		loginMenu("Internal Server Error")
	case http.StatusCreated:
		loginMenu("")
	}
}

func sendAnnouncement(announcement announcement) {
	ann, _ := json.Marshal(announcement)
	resp, err := http.Post(baseURL("/makeann"), "application/json", bytes.NewBuffer(ann))
	if err != nil {
		mainMenu("Internal Server Error")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusBadRequest:
		mainMenu("Internal Server Error")
	case http.StatusInternalServerError:
		mainMenu("Internal Server Error")
	case http.StatusCreated:
		listAnnouncements()
	}
}

func sendLmReport(lr labReport) {
	lmr, _ := json.Marshal(lr)
	resp, err := http.Post(baseURL("/makelabreport"), "application/json", bytes.NewBuffer(lmr))
	if err != nil {
		mainMenu("Internal Server Error")
	}
	defer resp.Body.Close()
	fmt.Println(http.StatusText(resp.StatusCode))
	switch resp.StatusCode {
	case http.StatusBadRequest:
		mainMenu("Internal Server Error")
	case http.StatusInternalServerError:
		mainMenu("Internal Server Error")
	case http.StatusCreated:
		listLmReports()
	case http.StatusNotFound:
		listLmReports()
	default:
		fmt.Println("what the fuck")
	}
}

func sendFgReport(fg gadgetReport) {
	fgr, _ := json.Marshal(fg)
	resp, err := http.Post(baseURL("/makegadgetreport"), "application/json", bytes.NewBuffer(fgr))
	if err != nil {
		mainMenu("Internal Server Error")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusBadRequest:
		mainMenu("Internal Server Error")
	case http.StatusInternalServerError:
		mainMenu("Internal Server Error")
	case http.StatusNotFound:
		listFgReports()
	case http.StatusCreated:
		listFgReports()
	}
}
