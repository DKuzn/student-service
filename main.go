package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var fullName string
var port int

type Student struct {
	ProcId   int    `json:"proc_id"`
	FullName string `json:"full_name"`
}

var rootCmd = &cobra.Command{
	Use:  "student",
	Long: "An example service to get info about the student.",
	Run:  serveStatic,
}

func init() {
	rootCmd.Flags().StringVarP(&fullName, "fullname", "n", "Full Name", "Full name of the student.")
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "The port is used to listen requests.")
}

func serveStatic(cmd *cobra.Command, args []string) {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/student", func(w http.ResponseWriter, r *http.Request) {
		student := Student{ProcId: os.Getpid(), FullName: fullName}
		json.NewEncoder(w).Encode(student)
	}).Methods("GET")

	methods := handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	log.Printf("Student Service is running on port: %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), handlers.CORS(origins, methods)(r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
