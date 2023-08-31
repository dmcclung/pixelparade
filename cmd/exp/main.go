package main

import (
	"html/template"
	"os"
)

type User struct {
	Name   string
	Age    int
	Script template.HTML
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name:   "Dylan",
		Age:    100,
		Script: "<p>my html script</p>",
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
