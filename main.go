package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
)

var port string

func main()  {
	args := os.Args

	if args == nil || len(args) < 1 {
		fmt.Println("COMMAND PORT")
		return
	}
	port = args[1]

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	//server开始监听
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}
	return ip
}


func CrossDomain(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	log.Println("request domain ", r.Host, "URL: ", r.URL)
	w.Header().Set("Access-Control-Allow-Origin", "*")                                                                                              //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie,Action, Module") //header的类型
	w.Header().Set("content-type", "application/json")                                                                                              //返回数据格式是json
	//header("Access-Control-Allow-Credentials : true");
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	return w
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	CrossDomain(w, r)
	ip, err := externalIP()
   if err != nil {
         fmt.Println(err)
     }

	w.Write([]byte("hello world~! server ip addr: "+ip.String()+":"+port+"\n"))
}
