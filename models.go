package main

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type announcement struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type labReport struct {
	Author  string `json:"author"`
	Title   string `json:"title"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type gadgetReport struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
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
