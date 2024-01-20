package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Serving static HTML files
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)

    // Handling specific routes
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
    // Add logic to handle transaction
    fmt.Fprintf(w, "Transaction processed")
}

func handleStake(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }
    // Add logic to handle stake
    fmt.Fprintf(w, "Stake processed")
}

func handleViewLastBlock(w http.ResponseWriter, r *http.Request) {
    // Add logic to view last block
    fmt.Fprintf(w, "Last block details")
}

func handleCheckBalance(w http.ResponseWriter, r *http.Request) {
    // Add logic to check balance
    fmt.Fprintf(w, "Balance details")
}

