package main


import (
	"session"
	"fmt"
	"net/http"
	"log"
	"strconv"
	"encoding/json"
)


var SessionManager = session.InitManager()


func addSession(w http.ResponseWriter, r *http.Request) {
	session := SessionManager.AddSession("localhost")
	b, _ := json.Marshal(session)
	fmt.Println("added session", string(b))
	fmt.Fprintf(w, string(b))
}


func rmSession(w http.ResponseWriter, r *http.Request) {
	uuid, ok := r.URL.Query()["uuid"]
	if (!ok) {
		fmt.Fprintf(w, "Please supply a uuid")
		return
	}
	SessionManager.RemoveSession(uuid[0])
	fmt.Println("removed session", uuid)
	fmt.Fprintf(w, "OK")
}


func listSession(w http.ResponseWriter, r *http.Request) {
	json, err := SessionManager.ToJson()
	if(err != nil) {
		fmt.Fprintf(w, "ERROR")
		return
	}
	fmt.Fprintf(w, string(json))
}


func processSession(w http.ResponseWriter, r *http.Request) {
	_uuid, uok := r.URL.Query()["uuid"]
	if (!uok) {
		fmt.Fprintf(w, "Please supply a uuid")
		return
	}
	uuid := _uuid[0]
	timeout := 10
	_timeout, tok := r.URL.Query()["timeout"]
	if tok {
		i, _ := strconv.Atoi(_timeout[0])
		timeout = i
	}
	mode := "live"
	_mode, mok := r.URL.Query()["mode"]
	if mok {
		mode = _mode[0]
	}	
	fmt.Println("Processing session with", timeout, uuid, mode)
	if mode == "live" {
		fmt.Println("Running in live mode")
		go SessionManager.ProcessSession(uuid, timeout)
	} else if mode == "batch"{
		fmt.Println("Running in batch mode")
		SessionManager.ProcessSession(uuid, timeout)
	} else {
		fmt.Fprintf(w, "Wrong mode supplied, mode should be `batch` or `live`")
	}
	fmt.Println("Done processing for", uuid)
	fmt.Fprintf(w, "OK")
}


func getSession(w http.ResponseWriter, r *http.Request) {
	_uuid, uok := r.URL.Query()["uuid"]
	if (!uok) {
		fmt.Fprintf(w, "Please supply a uuid")
		return
	}
	uuid := _uuid[0]	
	bs, ok := SessionManager.GetSession(uuid)
	if !ok {
		fmt.Fprintf(w, "Not a valid uuid")
		return
	}
	json, err := json.Marshal(bs)
	if(err != nil) {
		fmt.Fprintf(w, "ERROR")
		return
	}
	fmt.Fprintf(w, string(json))
}


func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}


func main() {
	http.HandleFunc("/create/", addSession)
	http.HandleFunc("/remove/", rmSession)
	http.HandleFunc("/list/", listSession)
	http.HandleFunc("/process/", processSession)
	http.HandleFunc("/getSession/", getSession)
	http.HandleFunc("/health/", healthCheck)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
