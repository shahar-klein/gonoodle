package main
/*
TODO: 
	burst
	Local address 
        global send buff?
        raw packet
        time to run
        time to run per session
        serve epoll?
        adaptive send over second
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
)

const charset = "abcdefghijklmnopqrstuvwxyz" +

  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


func stringToBytes(s string) (int) {

	unit := s[len(s)-1:]
	if  ! strings.Contains("kmgKMG", unit) {
		fmt.Println("Error parsing bandwidth, no unit")
	}
	numS := s[:len(s)-1]
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
	server		bool
	socketMode	string
	port		int
	numConns	int
	rampRate	int // per thread
	numThreads	int
	bwPerConn	int
	msgSize		int
	sendBuff	[]byte
}

func (self Config) dump() {
	fmt.Println("Config:", "server:", self.server, "\nclient:", self.daddr, "\nport:", self.port, "\nnum ports:", self.numConns,
			"\nRamp:", self.rampRate, "\nnum threads:", self.numThreads, "\nBW per conn:", self.bwPerConn, "\nmsg size:", self.msgSize)
}

func (self *Config) parse(args []string) {
	parser := argparse.NewParser("Noodle", "iperf with goodies")
	s := parser.Flag("s", "server", &argparse.Options{Help: "Server mode"})
	c := parser.String("c", "client", &argparse.Options{Help: "Client mode"})
	u := parser.Flag("u", "udp", &argparse.Options{Help: "UDP mode", Default: false})
	p := parser.Int("p", "port", &argparse.Options{Help: "port to listen on/connect to", Required: false, Default: 10005})
	P := parser.Int("P", "conns", &argparse.Options{Help: "number of parallel connections to run", Default: 100})
	R := parser.Int("R", "ramp", &argparse.Options{Help: "Ramp up connections per second", Default: 100})
	b := parser.String("b", "bandwidth", &argparse.Options{Help: "Banwidth per connection in kmgKMG", Default: "1m"})
	B := parser.String("B", "total-bandwidth", &argparse.Options{Help: "Total Banwidth in kmgKMG bits. Overrides -b"})
	l := parser.Int("l", "msg size", &argparse.Options{Help: "length(in bytes) of buffer in bytes to read or write", Default: 1440})

	err := parser.Parse(args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	if *s == true {
		self.server = true
	}
	self.daddr = *c
	self.socketMode = "tcp"
	if *u == true {
		self.socketMode = "udp"
	}
	self.port = *p
	self.numConns = *P
	self.rampRate = *R
	if self.rampRate > self.numConns {
		self.rampRate = self.numConns
	}

	// BW per conn in bytes
	self.bwPerConn = stringToBytes(*b)
	// -B overides -b
	if strings.ContainsAny(*B, "kmgKMG") {
		totalBW := stringToBytes(*B)
		self.bwPerConn = totalBW/self.numConns
	}
	self.bwPerConn /= 8

	self.msgSize = *l
	// be polite
	if self.msgSize > self.bwPerConn {
		self.msgSize = self.bwPerConn
	}

	self.sendBuff = make([]byte, self.msgSize)
	for i := range self.sendBuff {
                self.sendBuff[i] = charset[seededRand.Intn(len(charset))]
        }

}


type Connection struct {
	id		int
	conn		net.Conn
	daddr		string
	dport		int
	byteSent	int
	byteBWPerSec	int
	isActive	bool
	socketMode	string
	msgSize		int
	msg		*[]byte
	secondNotOver	bool
}

func (self *Connection) dump() {
	fmt.Println("Connection id:", self.id, "Dial to:", net.JoinHostPort(self.daddr, strconv.Itoa(self.dport)))
}

func (self *Connection) connect() {
	var err error
	fmt.Println("Connecting")
	self.conn, err = net.Dial(self.socketMode, net.JoinHostPort(self.daddr, strconv.Itoa(self.dport)))
	fmt.Println("connect:", self.id, err)
	if err != nil {
		fmt.Println("Error connect:", self.id, err)
		os.Exit(1)
	}
	self.isActive = true
}

func (self *Connection) send() {
	ticker := time.NewTicker(1000 * time.Millisecond)
	for range ticker.C {
		self.byteSent = 0
		for self.byteSent < self.byteBWPerSec {
			sent, err := self.conn.Write(*self.msg)
			if err != nil {
				fmt.Println("Error sent:", self.id, err)
			} else {
				self.byteSent += sent
			}
		}
	}
}


func (self *Connection) run() {
	self.connect()
	fmt.Println("RUN:", self.id)
	self.send()
}

func runClient(config *Config) {
	ticker := time.NewTicker(1000 * time.Millisecond)
	created := 0
	for range ticker.C {
		if created >= config.numConns {
			// good position to collect stats
			break
		}
		for i:=0; i<config.rampRate; i++ {
			c := Connection{id: i,
				daddr: config.daddr,
				dport: config.port,
				byteBWPerSec: config.bwPerConn,
				isActive: false,
				socketMode: config.socketMode,
				msgSize: config.msgSize,
				msg: &config.sendBuff}
			created++
			go c.run()
	    }
        }
}

type Server struct {
	id		int
	ch		chan net.Conn
}

func read(s net.Conn) {
	fmt.Println("in Read:")
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

func main() {

	config := new(Config)
	config.parse(os.Args)
	config.dump()

	if config.server {
		if config.socketMode == "tcp" {
			runTCPServer(config)
		} else {
			runUDPServer(config)
		}
	} else {
		runClient(config)
	}

	select {}
}
