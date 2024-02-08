package main

import (
	// "fmt"
	"log"
	"net/http"

	it "github.com/somuthink/TurtleCalc/internal"
)

func main() {
	it.CreateServers()
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", it.HtmlPage)
	mux.HandleFunc("/newprob/", it.ProblemHandler)
	mux.HandleFunc("/getprobs/", it.SendProbs)
	mux.HandleFunc("/updateopers/", it.SendOpers)
	mux.HandleFunc("/getservers/", it.SendServers)
	log.Fatal(http.ListenAndServe(":8000", mux))
}
