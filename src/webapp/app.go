package main

import (
    "fmt"
    "net/http"
)

func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)

    http.HandleFunc("/sendTransaction", handleTransaction)
    http.HandleFunc("/stake", handleStake)
    http.HandleFunc("/viewLastBlock", handleViewLastBlock)
    http.HandleFunc("/checkBalance", handleCheckBalance)

    fmt.Println("Server started at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

func handleTransaction(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }
    fmt.Fprintf(w, "Transaction processed")
}

func handleStake(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }
    fmt.Fprintf(w, "Stake processed")
}

func handleViewLastBlock(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Last block details")
}

func handleCheckBalance(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Balance details")
}

