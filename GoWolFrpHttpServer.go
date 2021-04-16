package main

import (
	"encoding/json"
	"fmt"
	"github.com/tignioj/go-wol-frp-http/wol"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)
func main()  {
	//2C:F0:5D:54:34:D0
	args := os.Args
	var port = "9999"
	if len(args) >1 {
		port = args[1]
	}

	log.Println("Starting wol server: http://localhost:" + port)
	addr := ":" + port
	http.HandleFunc("/", indexView)
	http.HandleFunc("/wake", handleWake)
	http.HandleFunc("/wakeAll", handleWakeAll)
	http.HandleFunc("/del", handleDel)
	http.HandleFunc("/add", handleAdd)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleWakeAll(w http.ResponseWriter, r *http.Request) {
	//TODO
	w.Write([]byte("Developing this function..."))
}

func handleDel(w http.ResponseWriter, r *http.Request) {
	//TODO
	w.Write([]byte("Developing this function..."))
}
var templates = template.Must(template.ParseFiles("index.html"))

func checkPost(w http.ResponseWriter, r *http.Request) bool{
	if strings.ToUpper(r.Method) == "POST" {
		if err := r.ParseForm(); err != nil {
			return false
		}
	} else {
		w.WriteHeader(403)
		w.Write([]byte("403:Not allow!"))
	}
	return false
}

type Device struct {
	DeviceName string `json:"device_name"`
	MacAddr string `json:"mac_address"`
	BroadCastIP string `json:"broad_cast_ip"`
}
func handleAdd(w http.ResponseWriter, r *http.Request) {
	//TODO
	w.Write([]byte("Developing this function..."))
	if checkPost(w, r) {
		//macAddr := r.FormValue("mac")
		//deviceName := r.FormValue("name")
		//f, err := ioutil.ReadFile("db.json")
	}
}

func handleWake(w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		macAddr := r.FormValue("macAddr")
		broadCastIP := r.FormValue("broadCastIP")
		err1 := wol.WakeCmd(macAddr, broadCastIP)
		if err1 != nil {
			w.Write([]byte("Send Error:" + err1.Error()))
			return
		}
		w.Write([]byte("Magic packet sent successfully to:" + macAddr))
	} else {
		w.WriteHeader(403)
		w.Write([]byte("403:Not allow!"))
	}
}



func indexView(w http.ResponseWriter, r *http.Request) {
	var m = make(map[string][]Device)
	devices := readFile("devices.json")
	obj := json.Unmarshal(devices, &m)
	fmt.Println(obj)

	err := templates.ExecuteTemplate(w, "index.html", m["devices"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func readFile(fileName string) []byte {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("read", err.Error())
	}
	return f
}