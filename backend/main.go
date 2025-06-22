package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/yuuhikaze/rgxr/handlers"
)

func main() {
    r := mux.NewRouter()

    // Existing endpoints
    r.HandleFunc("/union", handlers.UnionHandler).Methods("POST")
    r.HandleFunc("/convert", handlers.ConvertHandler).Methods("POST")
    r.HandleFunc("/save-image", handlers.SaveImageHandler).Methods("POST")

    // New endpoints for retrieving saved files
    r.HandleFunc("/tex/{uuid}", handlers.GetTeXHandler).Methods("GET")
    r.HandleFunc("/svg/{uuid}", handlers.GetSVGHandler).Methods("GET")

    // Extended endpoints
    r.HandleFunc("/intersection", handlers.IntersectionHandler).Methods("POST")
    r.HandleFunc("/nfa-to-dfa", handlers.NFAToDFAHandler).Methods("GET")
    r.HandleFunc("/fa-to-regex", handlers.FAToRegexHandler).Methods("GET")
    r.HandleFunc("/regex-to-nfa", handlers.RegexToNFAHandler).Methods("POST")
    r.HandleFunc("/minimize-dfa", handlers.MinimizeDFAHandler).Methods("GET")
    r.HandleFunc("/complement", handlers.ComplementHandler).Methods("GET")
    r.HandleFunc("/concatenation", handlers.ConcatenationHandler).Methods("POST")

    log.Println("Backend running on :8080")
    log.Println("Available endpoints:")
    log.Println("  POST /union - Union of multiple FAs")
    log.Println("  POST /intersection - Intersection of multiple FAs")
    log.Println("  POST /concatenation - Concatenation of multiple FAs")
    log.Println("  GET  /nfa-to-dfa?uuid=<uuid> - Convert NFA to DFA")
    log.Println("  GET  /fa-to-regex?uuid=<uuid> - Convert FA to regex")
    log.Println("  POST /regex-to-nfa - Convert regex to NFA")
    log.Println("  GET  /minimize-dfa?uuid=<uuid> - Minimize DFA")
    log.Println("  GET  /complement?uuid=<uuid> - Complement of FA")
    log.Println("  POST /convert - Convert FA to DOT/TikZ/SVG")
    log.Println("  POST /save-image - Save SVG image")
    log.Println("  GET  /tex/{uuid} - Get saved TeX file")
    log.Println("  GET  /svg/{uuid} - Get saved SVG file")
    
    log.Fatal(http.ListenAndServe(":8080", r))
}
