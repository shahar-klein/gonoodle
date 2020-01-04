noodle is a c10k iperf like load tool with some goodies like bw throttle per session
and session ramp up throttle

```

usage: Noodle [-h|--help] [-s|--server] [-c|--client "<value>"] [-u|--udp]
              [-p|--port <integer>] [-L|--local "<value>"] [-C|--conns
              <integer>] [-R|--ramp <integer>] [-b|--bandwidth "<value>"]
              [-B|--total-bandwidth "<value>"] [-r|--burst "<value>"] [-l|--msg
              size <integer>] [-t|--time <integer>] [-M|--cms <integer>]
              [--rp "<value>"] [-T|--stime <integer>] [-i|--report interval
              <integer>]

              iperf with goodies

Arguments:

  -h  --help             Print help information
  -s  --server           Server mode
  -c  --client           <host> Client mode
  -u  --udp              UDP mode. Default: false
  -p  --port             port to listen on/connect to. Default: 10005
  -L  --local            [ip | ip:port ] Local address to bind as the first
                         port, use 0:port for start port. Default: 0:0
  -C  --conns            number of concurrent connections to run. Default: 100
  -R  --ramp             Ramp up connections per second. Default: 100
  -b  --bandwidth        Banwidth per connection in kmgKMG. Default: 1m
  -B  --total-bandwidth  Total Banwidth in kmgKMG bits. Overrides -b
  -r  --burst            burst in percentage from avarage low:high. Default:
                         100:100
  -l  --msg size         length(in bytes) of buffer in bytes to read or write.
                         Default: 1440
  -t  --time             time in seconds to transmit. Default: 10
  -M  --cms              number of connection managers. Default: 0
      --rp               RP mode <loader|initiator>, UDP only
  -T  --stime            session time in seconds. After T seconds the session
                         closes and re-opens immediately. 0 means don't close
                         till the process ends. Default: 0
  -i  --report interval  report interval. -1 means report only at the end. -2
                         means no report. Default: -1
TODOs/Features ideas: 
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




```
