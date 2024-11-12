package main

import (
  "fmt"
  "net/http"
)

func main(){
  http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request){
    fmt.Fprint(res, "Welcome to my website!")
  })

  fmt.Println("Server running locally on port 8000!")
  http.ListenAndServe(":8000", nil)
}
