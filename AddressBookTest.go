package main

import "github.com/gorilla/mux"
import (
"AddressBook"
"log"
"net/http"
)

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/address/{phone}", AddressBook.GetAddress).Methods("GET")
	router.HandleFunc("/address/{phone}", AddressBook.CreateAddress).Methods("PUT")
	router.HandleFunc("/address/{phone}", AddressBook.UpdateAddress).Methods("POST")
	router.HandleFunc("/address/{phone}", AddressBook.DeleteAddress).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))
}
