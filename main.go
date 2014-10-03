package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
)

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

var tcp = flag.String("tcp", "127.0.0.1:9090", "listen host:port for FCGI")

func getGroups(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()

	headers.Add("Content-Type", "application/json")
	fmt.Fprint(w, Response{"success": true,
		"groups": listCgroups(),
	})
}

func getPIDs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	headers := w.Header()

	headers.Add("Content-Type", "application/json")
	fmt.Fprint(w, Response{"success": true,
		"cgroup": vars["cgroup"],
		"pids":   listPids(vars["cgroup"]),
	})
}

func putPID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	headers := w.Header()
	r.ParseForm()

	for _, pid := range r.Form["pid"] {
		addPid(vars["cgroup"], pid)
	}

	headers.Add("Content-Type", "application/json")
	fmt.Fprint(w, Response{"success": true,
		"cgroup": vars["cgroup"],
		"pids":   listPids(vars["cgroup"]),
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/cgroups", getGroups).Methods("GET")
	r.HandleFunc("/cgroups/{cgroup}", getPIDs).Methods("GET")
	r.HandleFunc("/cgroups/{cgroup}", putPID).Methods("PUT")

	flag.Parse()
	var err error

	listener, err := net.Listen("tcp", *tcp)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	err = fcgi.Serve(listener, r)
}
