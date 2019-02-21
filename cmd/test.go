package main

import (
    "fmt"
    "net/http"
)
func home_h(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w,"Hello world \n")
}
func main(){
    http.HandleFunc("/",home_h)
    http.ListenAndServe(":8080",nil)
    // println("Hoge")
}
