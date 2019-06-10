package main
/*
TODO: 
	verbosity
	burst
	pps
        raw packet
        time to run
        time to run per session
        serve epoll? not sure I need it in go.
        adaptive send over second
	specific output device
        mixed sessions

*/

import (
	"fmt"
	"net"
	"time"
	"math/rand"
	"github.com/akamensky/argparse"
	"os"
	"strings"
	"strconv"
	"runtime"
	"reflect"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +

  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func humanRead(bytes int) (string){


	switch {
		case bytes > 1000000000:
			return fmt.Sprint(bytes/1000000000, " Gbytes")
		case bytes > 1000000:
			return fmt.Sprint(bytes/1000000, " Mbytes")
		case bytes > 1000:
			return fmt.Sprint(bytes/1000, " Kbytes")
	}

	return fmt.Sprint(bytes, " bytes")

}

func stringToBytes(s string) (int) {

	unit := s[len(s)-1:]
	numS := ""
	if  ! strings.Contains("kmgKMG", unit) {
		numS = s
	} else {
		numS = s[:len(s)-1]
	}
	f := 1
	switch {
                case unit == "k" || unit == "K":
                        f = 1000
                case unit == "m" || unit == "M":
                        f = 1000000
                case unit == "g" || unit == "G":
                        f = 1000000000

        }

	num, err := strconv.Atoi(numS)
	if err != nil {
		fmt.Println("Error parsing bandwidth")
	}
	return num*f

}

var seededRand *rand.Rand = rand.New(
  rand.NewSource(time.Now().UnixNano()))

type Config struct {
	daddr		string
	saddr		string
	server		bool
	socketMode	string
	lPort		int
	port		int
	sport		int
	numConns	int
	numCM		int
	numConnsCM	int // sessions per CM
	rampRate	int // per CM
	bwPerConn	int
	bwPerConnLo	int
	bwPerConnHi	int
	msgSize		int
	timeToRun	int
	sessionTime	int
	sendBuff	[]byte
	rpMode		string
	reportInterval	int
}

func (self *Config) dump() {
	if self.server {
		fmt.Println("Config:", strings.ToUpper(self.socketMode), "Server\nListen to:", self.port)
	} else {
		fmt.Println("Config:",  strings.ToUpper(self.socketMode), "client, Calling", self.daddr, "\b:" , self.port, "\nnum conns:", self.numConns,
			"\nRamp:", self.rampRate, "\nBW per conn:", self.bwPerConn, "\nmsg size:", self.msgSize, "\nsTime:", self.sessionTime,
			"\nrpMode:", self.rpMode, "\nnumCM:", self.numCM, "\nnumConnsCM:", self.numConnsCM)
	}
}

func (self *Config) parse(args []string) {
	parser := argparse.NewParser("Noodle", "iperf with goodies")
	s := parser.Flag("s", "server", &argparse.Options{Help: "Server mode"})
	c := parser.String("c", "client", &argparse.Options{Help: "<host> Client mode"})
	u := parser.Flag("u", "udp", &argparse.Options{Help: "UDP mode", Default: false})
	p := parser.Int("p", "port", &argparse.Options{Help: "port to listen on/connect to", Required: false, Default: 10005})
	L := parser.String("L", "local", &argparse.Options{Help: "[ip | ip:port ] Local address to bind as the first port, use 0:port for start port", Default: "0:0"})
	C := parser.Int("C", "conns", &argparse.Options{Help: "number of concurrent connections to run", Default: 100})
	R := parser.Int("R", "ramp", &argparse.Options{Help: "Ramp up connections per second", Default: 100})
	b := parser.String("b", "bandwidth", &argparse.Options{Help: "Banwidth per connection in kmgKMG", Default: "1m"})
	B := parser.String("B", "total-bandwidth", &argparse.Options{Help: "Total Banwidth in kmgKMG bits. Overrides -b"})
	r := parser.String("r", "burst", &argparse.Options{Help: "burst in percentage from avarage low:high", Default: "100:100"})
	l := parser.Int("l", "msg size", &argparse.Options{Help: "length(in bytes) of buffer in bytes to read or write", Default: 1440})
	t := parser.Int("t", "time", &argparse.Options{Help: "time in seconds to transmit", Default: 10})
	M := parser.Int("M", "cms", &argparse.Options{Help: "number of connection managers", Default: 0})
	RP := parser.String("", "rp", &argparse.Options{Help: "RP mode <loader|initiator>, UDP only"})
	T := parser.Int("T", "stime", &argparse.Options{Help: "session time in seconds. After T seconds the session closes and re-opens immediately. 0 means don't close till the process ends", Default: 0})
	i := parser.Int("i", "report interval", &argparse.Options{Help: "report interval. -1 means report only at the end. -2 means no report", Default: -1})



	err := parser.Parse(args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	self.timeToRun = *t
	self.sessionTime = *T

	self.reportInterval = *i
	if self.reportInterval == -1 {
		self.reportInterval = self.timeToRun-1
	}

	lAddr := strings.Split(*L, ":")
	if len(lAddr) == 1 {
		//address only
		self.saddr = lAddr[0]
	} else { //A:P or :P
		if lAddr[0] == "" { //:P
			self.saddr = ""
		} else {
			self.saddr = lAddr[0]
		}
		self.sport, err = strconv.Atoi(lAddr[1])
		if err != nil {
			fmt.Println("Error parsing local port")
			fmt.Print(parser.Usage(err))
			os.Exit(1)
		}

	}

	if *s == true {
		self.server = true
	}
	self.daddr = *c


	self.socketMode = "tcp"
	if *u == true {
		self.socketMode = "udp"
	}

	self.rpMode = *RP

	if self.rpMode != "" && self.socketMode != "udp" {
		fmt.Println("\n\nError: RP mode can run in UDP only\n")
		os.Exit(1)
	}

	self.port = *p
	self.numConns = *C

	// BW per conn in bytes
	self.bwPerConn = stringToBytes(*b)
	// -B overides -b
	if *B != "" {
		totalBW := stringToBytes(*B)
		self.bwPerConn = totalBW/self.numConns
	}
	self.bwPerConn /= 8
	if self.bwPerConn < 1 {
		self.bwPerConn = 1
	}

	low, err := strconv.Atoi(strings.Split(*r, ":")[0])
	if err != nil {
		fmt.Println("Error parsing burst numbers")
		os.Exit(11)
	}
	high, err := strconv.Atoi(strings.Split(*r, ":")[1])
	if err != nil {
		fmt.Println("Error parsing burst numbers")
		os.Exit(11)
	}
	if low > 100 {
		fmt.Println("Low burst should be less than 100")
		os.Exit(11)
	}
	if high < 100 {
		fmt.Println("high burst should be more than 100")
		os.Exit(11)
	}
	self.bwPerConnHi = int(float64(self.bwPerConn)*float64(high)/100)
	self.bwPerConnLo = int(float64(self.bwPerConn)*float64(low)/100)

	self.msgSize = *l
	// be polite
	if self.msgSize > self.bwPerConn {
		self.msgSize = self.bwPerConn
	}

	self.sendBuff = make([]byte, self.msgSize)
	for i := range self.sendBuff {
                self.sendBuff[i] = charset[seededRand.Intn(len(charset))]
        }

	if *M == 0 {
		self.numCM = runtime.NumCPU() / 4
	} else {
		self.numCM = *M
	}
	if self.numConns < self.numCM {
		self.numCM = 1
	}
	self.numConnsCM = self.numConns/self.numCM

	self.rampRate = *R
	if self.rampRate > self.numConns {
		self.rampRate = self.numConns
	}
	self.rampRate = self.numConns/self.numCM
	if self.rampRate < 1 {
		self.rampRate = 1
	}

}


type Connection struct {
	id		int
	conn		net.Conn
	daddr		string
	saddr		string
	dport		int
	sport		int
	byteSent	int
	byteBWPerSec	int
	byteBWPerSecLo	int
	byteBWPerSecHi	int
	isActive	bool
	isReady		bool
	socketMode	string
	msgSize		int
	msg		*[]byte
	secondNotOver	bool
	sessionTime	int
	rpMode		string
}

func (self *Connection) dump() {
	fmt.Println("Connection id:", self.id, "Dial to:", net.JoinHostPort(self.daddr, strconv.Itoa(self.dport)), "From:", net.JoinHostPort(self.saddr, strconv.Itoa(self.sport)), self.isActive, self.isReady)
}

func (self *Connection) connect() {
	var err error
	//fmt.Println("Connection: Low:", self.byteBWPerSecLo, "high:", self.byteBWPerSecHi)

	if self.socketMode == "tcp" {
		dAddr, _  := net.ResolveTCPAddr(self.socketMode, net.JoinHostPort(self.daddr, strconv.Itoa(self.dport)))
		sAddr, _  := net.ResolveTCPAddr(self.socketMode, net.JoinHostPort(self.saddr, strconv.Itoa(self.sport)))
		self.conn, err = net.DialTCP(self.socketMode, sAddr, dAddr)
		if err != nil {
			fmt.Println(self.id, err)
			return
		}
		self.conn.(*net.TCPConn).SetLinger(0)
		self.conn.(*net.TCPConn).SetNoDelay(true)
	} else {
		dAddr, _ := net.ResolveUDPAddr(self.socketMode, net.JoinHostPort(self.daddr, strconv.Itoa(self.dport)))
		sAddr, _ := net.ResolveUDPAddr(self.socketMode, net.JoinHostPort(self.saddr, strconv.Itoa(self.sport)))
		self.conn, err = net.DialUDP(self.socketMode, sAddr, dAddr)
		if err != nil {
			fmt.Println(self.id, err)
			return
		}

	}
	self.isActive = true
}

func (self *Connection) zero() {
	self.byteBWPerSec = seededRand.Intn(self.byteBWPerSecHi-self.byteBWPerSecLo+1) + self.byteBWPerSecLo
	self.byteSent  = 0
}

func (self *Connection) waitForInitiator() {
	if self.rpMode == "loader" {
		buffer := make([]byte, 100)
		self.conn.(*net.UDPConn).ReadFrom(buffer)
		/*
			nRead, addr, err := self.conn.(*net.UDPConn).ReadFrom(buffer)
			if err != nil {
				fmt.Println("Error Read:", err)
			}
			fmt.Println("Got read from", addr, "read:", nRead)
		*/
	}
	self.isReady = true

}

func (self *Connection) send() {
	if self.isReady != true {
		self.waitForInitiator()
	}
	if self.byteSent < self.byteBWPerSec && self.isActive == true && self.isReady == true {
	//if self.byteSent < self.byteBWPerSec && self.isActive == true {
		sent, err := self.conn.Write(*self.msg)
		if err != nil {
			fmt.Println("Error sent:", self.id, err)
		} else {
			self.byteSent += sent
		}
	}
}

func runCM(config *Config, id int, ch chan string) {
	needToCreate := config.numConnsCM
	if id == 0 && config.numConnsCM*config.numCM < config.numConns {
		needToCreate += config.numConns - (config.numConnsCM*config.numCM)
	}
	totalCreated := 0
	conns := make([]Connection, 0)
	sentTilReport := 0
	reportInterval := config.reportInterval

	sPort := config.sport
	if sPort != 0 {
		sPort = config.sport + id*(config.numConnsCM)
	}
        dPort := config.port
	if config.rpMode != "" {
		dPort = config.port + id*(config.numConnsCM)
	}
        for {
		secondOver := false
		duration := time.Duration(1) * time.Second
		f := func() {
			secondOver = true
			reportInterval--
		}
		time.AfterFunc(duration, f)

		for i:=0; i<len(conns); i++ {
			sentTilReport += conns[i].byteSent
			// can add here report per session - not recommended
			conns[i].zero()
		}
		if reportInterval == 0 {
			go func (reportBytes int) {
				ch <- fmt.Sprint("CM-", id, " Sent ", humanRead(reportBytes), ". over the last ", config.reportInterval, " seconds")
			}(sentTilReport)
			sentTilReport = 0
			reportInterval = config.reportInterval
		}
		for secondOver != true {
			for i:=0; i<len(conns); i++ {
				conns[i].send()
			}
			for i:=0; i<config.rampRate; i++ {
				// although it seems right, don't take this if out of the for
				if totalCreated >= needToCreate {
					break
				}
				c := Connection{id: totalCreated,
					daddr: config.daddr,
					dport: dPort,
					sport: sPort,
					saddr: config.saddr,
					byteBWPerSec: config.bwPerConn,
					byteBWPerSecLo: config.bwPerConnLo,
					byteBWPerSecHi: config.bwPerConnHi,
					isActive: false,
					isReady: false,
					socketMode: config.socketMode,
					msgSize: config.msgSize,
					sessionTime: config.sessionTime,
					rpMode: config.rpMode,
					msg: &config.sendBuff}
				if sPort != 0 {
					sPort++
				}
				if config.rpMode != "" {
					dPort++
				}
				conns = append(conns, c)
				conns[totalCreated].connect()
				totalCreated++
			}

		} // second is over
	}

}

type Server struct {
	id		int
	ch		chan net.Conn
}

func read(s net.Conn) {
	var err error
	read := 0
	buf := make([]byte, 8*1024)
	for {
		read, err = s.Read(buf)
		if read == 0 {
			s.Close()
			break
		}
		if err != nil {
			fmt.Println("Error:", err, "Read:", read)

		}

	}
}

func runUDPServer(config *Config) {
	//net.ListenUDP(config.socketMode, &net.UDPAddr{IP:[]byte{0,0,0,0},Port:config.port,Zone:""})
	//l, err := net.ListenUDP(config.socketMode, &net.UDPAddr{IP:[]byte{0,0,0,0},Port:config.port,Zone:""})
	l, err := net.ListenPacket(config.socketMode, ":"+strconv.Itoa(config.port))

	defer l.Close()

	if err != nil {
		fmt.Println("Error UDP: ", err.Error())
	}
	for {
	}

}

func runTCPServer(config *Config) {
	A := 0
	l, err := net.Listen("tcp", ":"+strconv.Itoa(config.port))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		A++
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println(A, "Got connected from ", conn.RemoteAddr().String())
		go read(conn)
	}
}

func reporter(reportChans []chan string) {
	cases := make([]reflect.SelectCase, len(reportChans))
	for i, ch := range reportChans {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	for {
		//chosen, value, ok := reflect.Select(cases)
		//fmt.Printf("Read from channel %#v and received %s\n", reportChans[chosen], value.String())
		_, value, _ := reflect.Select(cases)
		fmt.Printf("%s\n", value.String())
	}
}

func main() {

	config := new(Config)
	config.parse(os.Args)
	config.dump()
	var reportChans = []chan string{}

	if config.server {
		if config.socketMode == "tcp" {
			runTCPServer(config)
		} else {
			runUDPServer(config)
		}
	} else {
		for i:=0; i<config.numCM; i++ {
			ch := make(chan string)
			reportChans = append(reportChans, ch)
			go runCM(config, i, reportChans[i])
		}
		go reporter(reportChans)
	}

	select {
	case <-time.After(time.Duration(config.timeToRun) * time.Second):
		fmt.Println("Done.")
		break
	}
}
