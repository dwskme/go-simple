package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type ToDoList struct {
	ToDos     []string
	ToDoCount int
}

func errorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getString(filename string) []string {
	var lines []string
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil
	}
	errorCheck(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	errorCheck(scanner.Err())
	return lines
}

func write(writer http.ResponseWriter, msg string) {
	_, err := writer.Write([]byte(msg))
	errorCheck(err)
}

func englishHandler(writer http.ResponseWriter, _ *http.Request) {
	write(writer, "Hello World")
}

func interactHandler(writer http.ResponseWriter, _ *http.Request) {
	todoValues := getString("todos.txt")
	fmt.Printf("%#v\n", todoValues)
	t, err := template.ParseFiles("view.html")
	errorCheck(err)
	todos := ToDoList{
		ToDoCount: len(todoValues),
		ToDos:     todoValues,
	}
	err = t.Execute(writer, todos)
	errorCheck(err)
}

func newHandler(writer http.ResponseWriter, _ *http.Request) {
	t, err := template.ParseFiles("new.html")
	errorCheck(err)
	err = t.Execute(writer, nil)
	errorCheck(err)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	todo := request.FormValue("todo")
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("todos.txt", options, os.FileMode(0600))
	errorCheck(err)
	_, err = fmt.Fprintln(file, todo)
	errorCheck(err)
	err = file.Close()
	errorCheck(err)

	http.Redirect(writer, request, "/interact", http.StatusFound)
}

func main() {
	http.HandleFunc("/", englishHandler)
	http.HandleFunc("/interact", interactHandler)
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/create", createHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
