package main

import (
		"fmt"
		"net"
		"os"
		"strings"
		"os/signal"
		"syscall"
		"flag"
		"sync"
		"strconv"
		log "github.com/sirupsen/logrus"
)
var wg sync.WaitGroup

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Fatal("Ctrl+C pressed in Terminal")
	}()
}

func listen(p string) {
	PORT := ":" + p
	l, err := net.Listen("tcp", PORT)
	if err != nil {
			log.Error(err)
			return
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error(err)
			return
		}
		handler(conn)
	}
}

func handler(c net.Conn) {
	remote_ip := strings.Split(fmt.Sprintf("%s",c.RemoteAddr()), ":")[0]
	remote_port := strings.Split(fmt.Sprintf("%s",c.RemoteAddr()), ":")[1]
	local_ip := strings.Split(fmt.Sprintf("%s",c.LocalAddr()), ":")[0]
	local_port := strings.Split(fmt.Sprintf("%s",c.LocalAddr()), ":")[1]

	c.Close()
	log.WithFields(log.Fields{
		"source_ip": remote_ip,
		"source_port": remote_port,
		"destination_ip": local_ip,
		"destination_port": local_port,
  	}).Warn("Connection attempted")
}

func main() {
	options := ParseOptions()
	log.SetFormatter(&log.JSONFormatter{})
	if options.output == "stdout" {
		log.SetOutput(os.Stdout)
	} else {
		file, err := os.OpenFile(options.output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
	}
	SetupCloseHandler()

	log.WithFields(log.Fields{
		"output": options.output,
		"ports": options.ports,
  	}).Info("Honey Options")

  	for _, p := range strings.Split(options.ports, ",") {
  		if strings.Contains(p, "-") {
  			start, _ := strconv.Atoi(strings.Split(p, "-")[0])
  			end, _ :=  strconv.Atoi(strings.Split(p, "-")[1])
  			for i := start; i <= end; i++ {
				wg.Add(1)
  				go listen(strconv.Itoa(i))
			}
  		} else {
  			wg.Add(1)
  			go listen(p)
  		}
  	}
  	wg.Wait()
}

type Options struct {
	ports		string
	output		string
}

func ParseOptions() *Options {
	options := &Options{}
	flag.StringVar(&options.ports, "ports", "", "Ports to listen to (comma separated and/or ranges)\n > Pro Tip: nmap -oG - -v --top-ports 25 2>/dev/null | grep Ports\n -ports 21,22,80,443\n -ports 20-25,80,81,8000-8100")
	flag.StringVar(&options.output, "output", "stdout", "Where to save the logs\n -output /tmp/honey.log")
	flag.Parse()
	options.validateOptions()
	return options
}

func (options *Options) validateOptions() {
	if options.ports == "" {
		log.Error("You must provide at least 1 port")
		usage()
	}
}

func usage() {
	fmt.Println("honey - Usage:")
	flag.PrintDefaults()
	os.Exit(0)
}
