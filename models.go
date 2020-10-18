package main

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

type announcement struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	AID    int64  `json:"aid"`
}

type labReport struct {
	Author  string `json:"author"`
	Title   string `json:"title"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	RID     int64  `json:"rid"`
}

type gadgetReport struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	RID    int64  `json:"rid"`
}

type directory struct {
	Name     string
	Path     string
	Parent   *directory
	Children []directory
	Files    []file
}

type file struct {
	Name      string
	Path      string
	Directory directory
	Contents  string
}
