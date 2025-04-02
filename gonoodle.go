package main
/*
TODO: 
	verbosity/report/jitter
	burst/sporadic jitter
	pps
        raw packet
        time to run
        time to run per session
        serve epoll? not sure I need it in go.
        adaptive send over second/add 100 milli rsend resolution within the second
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
	"encoding/binary"
	"log"
	"bufio"
        "golang.org/x/sys/unix"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +

  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

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
        echo            bool
	trace		bool
	socketMode	string
	lPort		int
	port		int
	sport		int
	numConns	int
	numCM		int
	numConnsCM	int // sessions per CM
	rampRate	int
	rampRateCM      int // per CM
        rampFactor      int
	bwPerConn	int
	bwPerCycLo	int
	bwPerCycHi	int
	cycleInMiliSec  float64
	sendPerConPerCycle  float64
	msgSize		int
	timeToRun	int
	sessionTime	int
	tosVal          int
	sendBuff	[]byte
	rpMode		string
	multiIpFile	string
	reportInterval	int
	randomIPs	[]string
        debugInc        bool
}

func (self *Config) dump() {
	if self.server {
		fmt.Println("Config:", strings.ToUpper(self.socketMode), "Server\nListen to:", self.port)
	} else {
		fmt.Println("Config:",  strings.ToUpper(self.socketMode), "client, Calling", self.daddr, "\b:" , self.port, "\nnum conns:", self.numConns,
			"\nRamp:", self.rampRate, "\nrampRateCM:", self.rampRateCM, "\nrampFactor:", self.rampFactor, "\nBW per conn:", self.bwPerConn, "\nCycle in mili:", self.cycleInMiliSec, "\nBytes send per Cycle:", self.sendPerConPerCycle,  "\nmsg size:", self.msgSize, "\nsTime:", self.sessionTime,
			"\nrpMode:", self.rpMode, "\nnumCM:", self.numCM, "\nnumConnsCM:", self.numConnsCM)
	}
}

func ReadInts(f string) ([]string, error) {
    file, err := os.Open(f)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanWords)
    var result []string
    for scanner.Scan() {
        result = append(result, scanner.Text())
    }
    return result, scanner.Err()
}

func (self *Config) parse(args []string) {
	parser := argparse.NewParser("Noodle", "iperf with goodies")
	s := parser.Flag("s", "server", &argparse.Options{Help: "Server mode"})
	e := parser.Flag("e", "echo", &argparse.Options{Help: "Server echo mode"})
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
	tos := parser.Int("", "tos", &argparse.Options{Help: "tos in Dec", Default: 0})
	tr := parser.Flag("", "tr", &argparse.Options{Help: "trace mode", Default: false})
	f := parser.Int("f", "frequency", &argparse.Options{Help: "frequency in msec. Sessions will transmit every f millisec.", Default: 100})
	M := parser.Int("M", "cms", &argparse.Options{Help: "number of connection managers", Default: 0})
	RP := parser.String("", "rp", &argparse.Options{Help: "RP mode <loader_multi|loader|initiator>, UDP only"})
	RF := parser.String("", "rpips", &argparse.Options{Help: "multi ips file"})
	T := parser.Int("T", "stime", &argparse.Options{Help: "session time in seconds. After T seconds the session closes and re-opens immediately. 0 means don't close till the process ends", Default: 0})
        D := parser.Flag("D", "msginc", &argparse.Options{Help: "Msg payload inc int", Default: false})
	i := parser.Int("i", "report interval", &argparse.Options{Help: "report interval. -1 means report only at the end. -2 means no report", Default: -1})


	if len(args) < 2 {
		fmt.Print(parser.Usage(0))
		os.Exit(1)
	}

	err := parser.Parse(args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	self.timeToRun = *t
	self.sessionTime = *T
        self.tosVal = *tos

	self.reportInterval = *i
	if self.reportInterval == -1 {
		self.reportInterval = (self.timeToRun-1)
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

	if *e == true {
		self.echo = true
	}

	self.daddr = *c

        if *tr == true {
                self.trace = true
        }



	self.socketMode = "tcp"
	if *u == true {
		self.socketMode = "udp"
	}

	if *D == true {
		self.debugInc = true
	}

	self.multiIpFile= *RF

	self.rpMode = *RP
	if self.rpMode == "loader_multi" {
		if self.multiIpFile == "" {
			fmt.Println("\nFor the loader_multi option you must specify an ips file with --rpips\n")
			fmt.Print(parser.Usage(err))
			os.Exit(1)
		}
		self.randomIPs, _ = ReadInts(self.multiIpFile)
	}

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

	self.cycleInMiliSec = float64(*f)
	// bytes to send per connection per cycle b/f/1000
	self.sendPerConPerCycle = float64(self.bwPerConn)*self.cycleInMiliSec/1000


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
	self.bwPerCycHi = int(float64(self.sendPerConPerCycle)*float64(high)/100)
	self.bwPerCycLo = int(float64(self.sendPerConPerCycle)*float64(low)/100)

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
	self.rampRateCM = self.rampRate/self.numCM
	if self.rampRateCM < 1 {
		self.rampRateCM = 1
	}
        self.rampFactor = self.numCM / self.rampRate
        if self.rampFactor < 1 {
                self.rampFactor = 1
        }

}


type Connection struct {
	id		int
	thrId		int
	conn		net.Conn
	daddr		string
	saddr		string
	dport		int
	sport		int
	byteSent	int
	byteBWPerCycle	int
	byteBWPerCycLo	int
	byteBWPerCycHi	int
	isActive	bool
	isWaiting	bool
	isReady		bool
	socketMode	string
	msgSize		int
	msg		*[]byte
	sessionTime	int
	rpMode		string
        tosVal          int
        debugInc        bool
        msgCount        int
        started         time.Time
}

func (self *Connection)setTosUDP(tos int) error {
	sc, err := self.conn.(*net.UDPConn).SyscallConn()
	if err != nil {
		return err
	}
	var serr error
	err = sc.Control(func(fd uintptr) {
		serr = unix.SetsockoptInt(int(fd), unix.SOL_IP, unix.IP_TOS, tos)
	})
	if err != nil {
		return err
	}
	return serr
}

func (self *Connection)setTosTCP(tos int) error {
	sc, err := self.conn.(*net.UDPConn).SyscallConn()
	if err != nil {
		return err
	}
	var serr error
	err = sc.Control(func(fd uintptr) {
		serr = unix.SetsockoptInt(int(fd), unix.SOL_IP, unix.IP_TOS, tos)
	})
	if err != nil {
		return err
	}
	return serr
}


func (self *Connection) dump() {
	fmt.Println("Connection id:", self.id, "Thread id:", self.thrId, "Dial to:", net.JoinHostPort(self.daddr, strconv.Itoa(self.dport)), "From:", net.JoinHostPort(self.saddr, strconv.Itoa(self.sport)), self.isActive, self.isReady)
}

func (self *Connection) connect() bool {
	var err error
	//fmt.Println("Connection: Low:", self.byteBWPerSecLo, "high:", self.byteBWPerSecHi)

        self.started = time.Now()

	if self.socketMode == "tcp" {
		dAddr, _  := net.ResolveTCPAddr(self.socketMode, net.JoinHostPort(self.daddr, strconv.Itoa(self.dport)))
		sAddr, _  := net.ResolveTCPAddr(self.socketMode, net.JoinHostPort(self.saddr, strconv.Itoa(self.sport)))
		self.conn, err = net.DialTCP(self.socketMode, sAddr, dAddr)
		if err != nil {
			fmt.Println(self.id, err)
			return false
		}
                if self.tosVal > 0 {
                        self.setTosTCP(self.tosVal)
                }
		self.conn.(*net.TCPConn).SetLinger(0)
		self.conn.(*net.TCPConn).SetNoDelay(true)
	} else {
		dAddr, _ := net.ResolveUDPAddr(self.socketMode, net.JoinHostPort(self.daddr, strconv.Itoa(self.dport)))
		sAddr, _ := net.ResolveUDPAddr(self.socketMode, net.JoinHostPort(self.saddr, strconv.Itoa(self.sport)))
		self.conn, err = net.DialUDP(self.socketMode, sAddr, dAddr)
		if err != nil {
			fmt.Println(self.id, err)
			return false
		}
                if self.tosVal > 0 {
                        self.setTosUDP(self.tosVal)
                }

	}
	self.isActive = true
	return true
}

func (self *Connection) zero() {

	if self.byteBWPerCycHi != self.byteBWPerCycLo {
		self.byteBWPerCycle = rand.Intn(self.byteBWPerCycHi-self.byteBWPerCycLo+1) + self.byteBWPerCycLo
	}

	self.byteSent  = 0
}

func (self *Connection) waitForInitiator() {
	if self.rpMode == "loader" || self.rpMode == "loader_multi" {
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
        self.isWaiting = false
	self.isReady = true

}

func (self *Connection) send() {

        if self.sessionTime > 0 {
                if time.Since(self.started) >= time.Duration(self.sessionTime) * time.Second {
                        fmt.Println("handle:", self.sport, self.isReady)
                        self.started = time.Now()
                        self.conn.Close()
                        self.connect()
                }
        }

        if self.isReady != true {
                if self.isWaiting != true {
                        self.isWaiting = true
                                //fmt.Println("Go")
                                go self.waitForInitiator()
                }
        }

	if self.byteSent < self.byteBWPerCycle && self.isActive == true && self.isReady == true {
                if self.debugInc == true {
                        //fmt.Println("Sending", self.msgCount, strconv.Itoa(self.msgCount))
                        s := fmt.Sprintf("This is msg No. %d", self.msgCount)
                        sent, err := self.conn.Write([]byte(s))
			if err != nil {
				fmt.Println("Error sent:", self.id, err)
			} else {
				self.byteSent += sent
			}
                        self.msgCount += 1
                } else {
                        sent, err := self.conn.Write(*self.msg)
			if err != nil {
				fmt.Println("Error sent:", self.id, err)
			} else {
				self.byteSent += sent
			}
                }
	}
}

func runCM(config *Config, startSport, startDport, needToCreate, id int, ch chan string) {
	totalCreated := 0
	conCreatedThisCycle := 0
	cyclesForReport := 0
	conns := make([]Connection, 0)
	sentTilReport := 0
	cyclesTillreport := config.reportInterval*1000 / int(config.cycleInMiliSec)
	if cyclesTillreport < 1 {
		cyclesTillreport = 1
	}

	sPort := config.sport
	if sPort != 0 {
		if config.rpMode != "loader_multi" {
                        sPort = startSport
		}
	}
        dPort := config.port
	sAddr := uint32(0)
	if config.rpMode != "" {
                dPort = startDport
	}
        secondStarted := time.Now()
        for { // This is the main load loop per thread(CM)
		cycleOver := false
		duration := time.Duration(config.cycleInMiliSec) * time.Millisecond
		f := func() {
			cycleOver = true
			cyclesForReport += 1
                        if time.Since(secondStarted) >= time.Duration(config.rampFactor) * time.Second {
                               conCreatedThisCycle = 0
                               secondStarted = time.Now()
                        }
		}
		time.AfterFunc(duration, f)

		//TODO reporting - consider go routine for that
		for i:=0; i<len(conns); i++ {
			sentTilReport += conns[i].byteSent
			// can add here report per session - not recommended
			conns[i].zero()
		}
		if cyclesTillreport == cyclesForReport {
			go func (reportBytes int) {
				ch <- fmt.Sprint("CM-", id, " Sent ", humanRead(reportBytes), ". over the last ", cyclesTillreport*int(config.cycleInMiliSec), " mili-seconds")
			}(sentTilReport)
			sentTilReport = 0
			cyclesForReport = 0
		}
		for cycleOver != true {
			for i:=0; i<len(conns); i++ {
				conns[i].send()

			}
			for i:=0; i<config.rampRateCM; i++ {
				if totalCreated >= needToCreate {
					break
				}
				if conCreatedThisCycle >= config.rampRateCM {
					break
				}
				if config.rpMode == "loader_multi" {
					randomIP := config.randomIPs[(id*config.numConnsCM)+totalCreated]
					//fmt.Println("Got:", randomIP)
					sAddr = ip2int(net.ParseIP(randomIP))
				}
                                if config.trace == true {
                                        fmt.Println("NewConn: ", time.Now(), id, sPort, dPort)
                                }
				c := Connection{id: totalCreated,
					thrId: id,
					daddr: config.daddr,
					dport: dPort,
					sport: sPort,
					saddr: int2ip(sAddr).String(),
					byteBWPerCycle: int(config.sendPerConPerCycle),
					byteBWPerCycLo: config.bwPerCycLo,
					byteBWPerCycHi: config.bwPerCycHi,
					isActive: false,
					isReady: false,
					isWaiting: false,
					socketMode: config.socketMode,
					msgSize: config.msgSize,
					sessionTime: config.sessionTime,
					rpMode: config.rpMode,
                                        tosVal: config.tosVal,
                                        debugInc: config.debugInc,
					msgCount: 1,
					msg: &config.sendBuff}
				if sPort != 0 {
					sPort++
				}
				if config.rpMode != "" {
					dPort++
				}
				if config.rpMode == "loader_multi" {
					sPort--

				}
				if c.connect() {
					conns = append(conns, c)
					totalCreated++
					conCreatedThisCycle++
				}
			}

		} // second is over
	}

}

type Server struct {
	id		int
	ch		chan net.Conn
}

func readEcho(s net.Conn) {
    defer s.Close() // Ensure connection is closed when function exits

    buf := make([]byte, 8*1024)
    for {
        read, err := s.Read(buf)
        if err != nil {
            fmt.Println("Error:", err)
            break
        }
        if read == 0 {
            break
        }

        _, writeErr := s.Write(buf[:read]) // Echo back the data
        if writeErr != nil {
            fmt.Println("Write Error:", writeErr)
            break
        }
    }
}


func read(s net.Conn) {
	var err error
	read := 0
	buf := make([]byte, 8*1024)
	for {
		read, err = s.Read(buf)
		// fmt.Println("Read", buf)
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
        fmt.Println("UDP server")
	l, err := net.ListenPacket(config.socketMode, ":"+strconv.Itoa(config.port))
	if err != nil {
		fmt.Println("Error UDP: ", err.Error())
	}

	defer l.Close()
        if config.echo {
                buffer := make([]byte, 1024)
                for {
		       _, addr, err := l.ReadFrom(buffer)
		       if err != nil {
		               fmt.Println("Error reading from UDP:", err)
			       continue
		       }
		       // Send back 1 byte of zero
		       _, err = l.WriteTo([]byte{0x00}, addr)
		       if err != nil {
			       fmt.Println("Error sending response:", err)
		       }
                }
        } else {
                for {
                }
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
		//fmt.Println(A, "Got connected from ", conn.RemoteAddr().String())
                if config.echo {
                        go readEcho(conn)
                } else {
		        go read(conn)
                }
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
                extra := config.numConns - (config.numConnsCM * config.numCM)
                fmt.Println("Extra=", extra)
                needCreate := 0
                portsSoFar := 0
                sleepFactor := 1000 / config.rampRate
                fmt.Println("sleepFactor:", sleepFactor)
		for i:=0; i<config.numCM; i++ {
			ch := make(chan string)
			reportChans = append(reportChans, ch)
                        startSrcPort := config.sport + portsSoFar
                        startDstPort := config.port + portsSoFar
                        if extra > 0 {
                                needCreate = config.numConnsCM + 1
                                extra -= 1
                        } else {
                                needCreate = config.numConnsCM
                        }
                        portsSoFar += needCreate
                        time.Sleep(time.Duration(sleepFactor) * time.Millisecond)
			go runCM(config, startSrcPort, startDstPort, needCreate, i, reportChans[i])
		}
		go reporter(reportChans)
	}

	select {
	case <-time.After(time.Duration(config.timeToRun) * time.Second):
		fmt.Println("Done.")
		break
	}
}
