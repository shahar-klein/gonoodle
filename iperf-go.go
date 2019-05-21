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
	C := parser.Int("C", "conns", &argparse.Options{Help: "Num concurrent connections", Default: 100})
	r := parser.Int("r", "threads", &argparse.Options{Help: "Num threads to use - default(1)", Default: 1})
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
	self.numConns = *C
	self.numThreads = *r
	self.rampRate = *R
	if self.rampRate > self.numConns {
		self.rampRate = self.numConns
	}
	self.rampRate = self.rampRate / self.numThreads
	if self.rampRate == 0 {
		self.rampRate = 1
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

type connectionManager struct {
	id		int
	config		*Config
	numConns	int
	conns		[]Connection
	secondNotOver	bool
}

func (self *connectionManager) initLoad() {
	fmt.Println("Init:", self.id)
	self.numConns = self.config.numConns/self.config.numThreads
	fmt.Println("numConns:", self.numConns)
	//self.conns = make([]Connection, self.numConns)
	for i:=0; i<self.numConns; i++ {
		c := Connection{id: i,
				daddr: self.config.daddr,
				dport: self.config.port,
				byteBWPerSec: self.config.bwPerConn,
				isActive: false,
				socketMode: self.config.socketMode,
				msgSize: self.config.msgSize,
				msg: &self.config.sendBuff}
		self.conns = append(self.conns, c)
		fmt.Println("size conns:", len(self.conns))
	}
}

func (self *connectionManager) runLoad() {
	//fmt.Println("RUN:", self.id)
	//time.Sleep(time.Microsecond * 2000000)
	//fmt.Println("Done:", self.id)
	fmt.Println("RUN:", self.id, "size conns:", len(self.conns))
	for {
		self.secondNotOver = true
		duration := time.Duration(1) * time.Second
		time.AfterFunc(duration, self.secondOver)

		connsLeftToCreate := self.config.rampRate
		for i, _ := range self.conns {
			self.conns[i].zeroCounters()
		}

		for self.secondNotOver == true {
			for i, _ := range self.conns {
				//fmt.Println(i, "isActive:", c.isActive)
				if self.conns[i].isActive {
					self.conns[i].send()
				} else if connsLeftToCreate > 0 {
					self.conns[i].connect()
					self.conns[i].send()
					connsLeftToCreate -= 1
				}
			}
		}


	}

}

func (self *connectionManager) secondOver() {
	self.secondNotOver = false
}

func (self connectionManager) dump() {
        fmt.Println("CM id=", self.id, "\nConfig: port:", self.config.port)
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
}

func (self *Connection) zeroCounters() {
	fmt.Println("Sent:", self.byteSent)
	self.byteSent = 0
}

func (self *Connection) dump() {
	fmt.Println("Connection id:", self.id, "Dial to:", net.JoinHostPort(self.daddr, strconv.Itoa(self.dport)))
}

func (self *Connection) connect() {
	var err error
	fmt.Println("Connecting")
	self.conn, err = net.Dial(self.socketMode, net.JoinHostPort(self.daddr, strconv.Itoa(self.dport)))
	//Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	fmt.Println("connect:", self.id, err)
	if err != nil {
		fmt.Println("Error connect:", self.id, err)
		os.Exit(1)
	}
	self.isActive = true
	fmt.Println("setting: isActive", self.isActive)
}

func (self *Connection) send() {
	//fmt.Println("Send", self.byteSent, self.byteBWPerSec)
	if self.byteSent < self.byteBWPerSec {
		sent, err := self.conn.Write(*self.msg)
		if err != nil {
			fmt.Println("Error sent:", self.id, err)
		} else {
			self.byteSent += sent
		}
	}
}


func runClient(config *Config) {
	CM := []connectionManager{}
	for i:=0; i<config.numThreads; i++ {
		cm := connectionManager{id: i, config: config}
		CM = append(CM, cm)
		cm.initLoad()
		go cm.runLoad()
	}
}


type Server struct {
	id		int
	ch		chan net.Conn
}

func (self Server) worker() {
	conns := []net.Conn{}
	buf := make([]byte, 8*1024)

	for {
		fmt.Printf("worker: %d: ", self.id)
		for i:=0; i<len(conns); i++ {
			fmt.Printf("%d.", conns[i])
			conns[i].Read(buf)
		}
		fmt.Printf("\n")
		newConn := <-self.ch
		conns = append(conns, newConn)
	}

}

func runServer(config *Config) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(config.port))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	servers := []Server{}
	nextServer := 0
	for i:=0; i<config.numThreads; i++ {
		s := Server{id: i, ch: make(chan net.Conn)}
		servers = append(servers, s)
		go s.worker()
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Got connected from ", conn.RemoteAddr().String())
		servers[nextServer].ch <- conn
		nextServer++
		nextServer = nextServer % config.numThreads


	}
}


func main() {

	config := new(Config)
	config.parse(os.Args)
	config.dump()

	if config.server {
		runServer(config)
	} else {
		runClient(config)
	}

	//var wg sync.WaitGroup


	//wg.Add(1)
	//wg.Wait()

	select{}





/*
func main() {
	fmt.Printf("hello, world\n")
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
        	fmt.Printf("Some error %v", err)
        	return
    	}
	b := make([]byte, 1000)
  	for i := range b {
    		b[i] = charset[seededRand.Intn(len(charset))]
  	} 
    	for {
		//msg := strconv.Itoa(i)
        	//i++
        	//buf := []byte(msg)
        	conn.Write(b)
        	//if err != nil {
            	//	fmt.Println("ERROR:", err)
        	//}
        	//time.Sleep(time.Microsecond * 10000)
    	}
*/
}
