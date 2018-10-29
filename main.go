package main
/*
This is the main entry point for the program

In here we only define the endpoints and their handlers
	GET /create/ 			  creates new session and returns the session object
	GET /remove/?uuid=... 	  removes a session given a mandatory uuid in the querystring
	GET /list/ 				  lists all available sessions
	GET /process/?uuid=... 	  starts listening on the session queue given a mandatory uuid in the querystring
								by processing we mean consume all messages and write them to the database
	GET /getSession/?uuid=... gets the session object given a mandatory uuid in the querystring
	GET /health/ 			  Just checks that the server is alive
*/


import (
	"session"
	"fmt"
	"net/http"
	"log"
	"strconv"
	"encoding/json"
)


// Instantiate the ServerManager object
// It is really important to note that the server should only have one ServerManager instance
// It is not straight forward to implement a singletone struct like the Java world, so we just create one instance like this
var SessionManager = session.InitManager()


// GET /create/ handler
// takes no input and returns a session object as json
func addSession(w http.ResponseWriter, r *http.Request) {
	session := SessionManager.AddSession("localhost")
	b, _ := json.Marshal(session)
	fmt.Fprintf(w, string(b))
}


// GET /remove/?uuid=... handler
// takes uuid as a mandatory input and removes a session
// returns either the removed session object as json or ERROR as a string in case that the session doesn't exist
// We should note here that in case the session is under processing then we block until the whole session is processed
func rmSession(w http.ResponseWriter, r *http.Request) {
	uuid, ok := r.URL.Query()["uuid"]
	if (!ok) {
		fmt.Fprintf(w, "Please supply a uuid")
		return
	}
	session, success := SessionManager.RemoveSession(uuid[0])
	if success {
		fmt.Fprintf(w, string(session))
	} else {
		fmt.Fprintf(w, "ERROR")
	}
}


// GET /list/ handler
// takes no input and returns all sessions available in the SessionManager
func listSession(w http.ResponseWriter, r *http.Request) {
	json, err := SessionManager.ToJson()
	if(err != nil) {
		fmt.Fprintf(w, "ERROR")
		return
	}
	fmt.Fprintf(w, string(json))
}


// GET /process/?uuid=...&mode=live|batch&timeout=[0...] handler
// takes a mandatory uuid input and returns either OK as string in case of success or error message
// takes optional mode argument in the query string. mode should be either live or batch. default is live.
// takes optional timeout argument in the query string. default is 10 seconds. Timeout is the duration in seconds
// 		where the session should terminate after x seconds after the session has started processing
// tells the session manager to start listening on this session queue
// In case the client wishes to go in live mode then we don't block and return immediately OK
// And in the case of batch mode we block until we process all the messages and timeout
// by processing we mean consuming the messages from rabbitmq queue, parsing them and writing them to the database
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
	// now if we are going to live mode then we spawn a new goroutine and return immediately
	// otherwise if we are in batch mode then we block until we process all messages
	// else we return an error message stating we don't know how to handle the given mode
	if mode == "live" {
		go SessionManager.ProcessSession(uuid, timeout) // don't block
	} else if mode == "batch"{
		SessionManager.ProcessSession(uuid, timeout) // block
	} else {
		fmt.Fprintf(w, "Wrong mode supplied, mode should be `batch` or `live`") // undefined mode is given
		return
	}
	fmt.Fprintf(w, "OK")
}


// GET /getSession handler
// takes a mandatory UUID input in the query string
// returns either the session object or error message
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


// GET /health/ handler
// returns a dummy text response OK
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
	fmt.Println("listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
