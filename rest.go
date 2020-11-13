package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func login(creds credentials) {
	credentials, _ := json.Marshal(creds)
	resp, err := http.Post(baseURL("/login"), "application/json", bytes.NewBuffer(credentials))
	if err != nil {
		loginMenu("Internal Server Error")
	}
	defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusNotFound:
		loginMenu("Incorrect Username or Password")
	case http.StatusUnauthorized:
		loginMenu("Incorrect Username Or Password")
	case http.StatusInternalServerError:
		loginMenu("Internal Server Error")
	case http.StatusOK:
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
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		loginMenu("Internal Server Error")
	case http.StatusInternalServerError:
		loginMenu("Internal Server Error")
	case http.StatusCreated:
		loginMenu("")
	default:
		log.Fatal(resp.StatusCode)
	}
}

func sendAnnouncement(announcement announcement) {
	ann, _ := json.Marshal(announcement)

	client := &http.Client{}

	req, err := http.NewRequest("POST", baseURL("/makeann"), bytes.NewBuffer(ann))
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("username", currentUser)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusBadRequest:
		mainMenu("Bad Request")
	case http.StatusInternalServerError:
		mainMenu("Internal Server Error")
	case http.StatusCreated:
		listAnnouncements()
	}
}

func updateAnnouncement(id string, ann announcement) {
	client := &http.Client{}

	json, err := json.Marshal(ann)
	if err != nil {
		mainMenu(err.Error())
	}

	req, err := http.NewRequest(http.MethodPut, baseURL("/updateannouncement/"+id), bytes.NewBuffer(json))
	if err != nil {
		mainMenu(err.Error())
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		mainMenu(err.Error())
	}

	defer resp.Body.Close()
}

func deleteAnnouncement(id string) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", baseURL("/deleteannouncement/"+id), nil)
	if err != nil {
		mainMenu(err.Error())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		mainMenu(err.Error())
	}
	defer resp.Body.Close()
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

func updateLmReport(id string, lr labReport) {
	client := &http.Client{}

	json, err := json.Marshal(lr)
	if err != nil {
		mainMenu(err.Error())
	}

	req, err := http.NewRequest(http.MethodPut, baseURL("/updatelabreport/"+id), bytes.NewBuffer(json))
	if err != nil {
		mainMenu(err.Error())
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		mainMenu(err.Error())
	}

	defer resp.Body.Close()
}

func deleteLmReport(id string) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", baseURL("/deletelabreport/"+id), nil)
	if err != nil {
		mainMenu(err.Error())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		mainMenu(err.Error())
	}
	defer resp.Body.Close()
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

func updateFgReport(id string, gr gadgetReport) {
	client := &http.Client{}

	json, err := json.Marshal(gr)
	if err != nil {
		mainMenu(err.Error())
	}

	req, err := http.NewRequest(http.MethodPut, baseURL("/updategadgetreport/"+id), bytes.NewBuffer(json))
	if err != nil {
		mainMenu(err.Error())
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		mainMenu(err.Error())
	}

	defer resp.Body.Close()
}

func deleteFgReport(id string) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", baseURL("/deletegadgetreport/"+id), nil)
	if err != nil {
		mainMenu(err.Error())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		mainMenu(err.Error())
	}
	defer resp.Body.Close()
}

func sendMessage(msg message) {
	m, _ := json.Marshal(msg)
	resp, err := http.Post(baseURL("/newmessage"), "application/json", bytes.NewBuffer(m))
	if err != nil {
		mainMenu("Internal Server Error")
	}
	defer resp.Body.Close()
}

// GETs

func getMessages() []message {
	messages := make([]message, 0)

	resp, err := http.Get(baseURL("/messages"))
	if err != nil {
		mainMenu("GET: Internal Server Error")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &messages)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		mainMenu("GET: Internal Server Error")
	}

	return messages
}

func getAnnouncement(id string) announcement {
	ann := announcement{}

	annr, err := http.Get(baseURL("/announcement/" + id))
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

	return ann
}

func getAnnouncements() []announcement {
	announcements := make([]announcement, 0)

	client := &http.Client{}

	req, err := http.NewRequest("GET", baseURL("/announcements"), nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("username", currentUser)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &announcements)

	defer resp.Body.Close()

	return announcements
}

func getLmReport(id string) labReport {
	lr := labReport{}

	lrr, err := http.Get(baseURL("/labreport/" + id))
	if err != nil {
		listLmReports()
	}
	lrBody, err := ioutil.ReadAll(lrr.Body)
	if err != nil {
		listLmReports()
	}
	asdf := json.Unmarshal(lrBody, &lr)
	if asdf != nil {
		listLmReports()
	}
	defer lrr.Body.Close()

	return lr
}

func getLmReports() []labReport {
	labReports := make([]labReport, 0)

	resp, err := http.Get(baseURL("/labreports"))
	if err != nil {
		mainMenu("GET: Internal Server Error")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &labReports)

	defer resp.Body.Close()

	return labReports
}

func getFgReport(id string) gadgetReport {
	gr := gadgetReport{}

	grr, err := http.Get(baseURL("/gadgetreport/" + id))
	if err != nil {
		listFgReports()
	}
	grBody, err := ioutil.ReadAll(grr.Body)
	if err != nil {
		listFgReports()
	}
	asdf := json.Unmarshal(grBody, &gr)
	if asdf != nil {
		listFgReports()
	}
	defer grr.Body.Close()

	return gr
}

func getFgReports() []gadgetReport {
	gadgetReports := make([]gadgetReport, 0)

	resp, err := http.Get(baseURL("/gadgetreports"))
	if err != nil {
		mainMenu("GET: Internal Server Error")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &gadgetReports)

	defer resp.Body.Close()

	return gadgetReports
}
