package main

import (
	"github.com/gorilla/mux"
	"github.com/yuuhikaze/rgxr/handlers"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	// Operation endpoints
	r.HandleFunc("/union", handlers.BooleanHandler).Methods("POST")
	r.HandleFunc("/intersection", handlers.BooleanHandler).Methods("POST")
	r.HandleFunc("/concatenation", handlers.ConcatenationHandler).Methods("POST")
	r.HandleFunc("/complement", handlers.ComplementHandler).Methods("GET")
	r.HandleFunc("/minimize-dfa", handlers.MinimizeDFAHandler).Methods("GET")
	r.HandleFunc("/fa-to-regex", handlers.FAToRegexHandler).Methods("GET")
	r.HandleFunc("/regex-to-nfa", handlers.RegexToNFAHandler).Methods("POST")
	r.HandleFunc("/nfa-to-dfa", handlers.NFAToDFAHandler).Methods("GET")

	// Storage endpoints
	r.HandleFunc("/tex/{uuid}", handlers.GetTeXHandler).Methods("GET")
	r.HandleFunc("/svg/{uuid}", handlers.GetSVGHandler).Methods("GET")
	r.HandleFunc("/render", handlers.RenderHandler).Methods("POST")

	log.Println("Backend running on :8080")
	log.Println("Available endpoints:")
	log.Println("  POST /union - Union of multiple FAs")
	log.Println("  POST /intersection - Intersection of multiple FAs")
	log.Println("  POST /concatenation - Concatenation of multiple FAs")
	log.Println("  GET  /complement?uuid=<uuid> - Complement of FA")
	log.Println("  GET  /minimize-dfa?uuid=<uuid> - Minimize DFA")
	log.Println("  GET  /fa-to-regex?uuid=<uuid> - Convert FA to regex")
	log.Println("  POST /regex-to-nfa - Convert regex to NFA")
	log.Println("  GET  /nfa-to-dfa?uuid=<uuid> - Convert NFA to DFA")
	log.Println("  POST /render - Render FA to SVG")
	log.Println("  GET  /tex/{uuid} - Get saved TeX file")
	log.Println("  GET  /svg/{uuid} - Get saved SVG file")

	log.Fatal(http.ListenAndServe(":8080", r))
}
