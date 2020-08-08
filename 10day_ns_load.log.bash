#!/bin/bash



# Slowly go up to 1000 at 20mb and stay there for an hour or so
#./gonoodle -u -c 10.8.51.71 --rp loader -C 1000 -R  10 -M 10 -b 20m -p 5000 -L 10.8.51.71:8000 -l 1400 -t 40

#sleep 60

# 200 sessions for 10 hours 40mb
#./gonoodle -u -c 10.8.51.71 --rp loader -C 200 -R 20 -M 10 -b 40m -p 5000 -L 10.8.51.71:8000 -l 1400 -t 360 &


# ip netns  add nsInitiator
# ip link set enp1s0 netns nsInitiator
# ip netns exec nsInitiator ip add add 30.30.30.20/16 dev enp1s0
# ip netns exec nsInitiator ip link set enp1s0 up
# ip netns exec nsInitiator ip r add default  dev enp1s0





set -x 
set +e
T=21600
B=20m
RP=5.5.5.2

d=0

log_info_on_rp() {
	when=$1
	ssh $RP "mkdir /root/day${when}"
	ssh $RP "ethtool -S enp7s0 > /root/day${when}/enp7s0.ethtool"
	ssh $RP "ethtool -S enp8s0 > /root/day${when}/enp8s0.ethtool"
	ssh $RP "cat /proc/meminfo> /root/day${when}/meminfo"
}

log_info_on_rp $d
d=$((d+1))

# day 1
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 800 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 300 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 700 -R 40 -M 1 -b 1k -p 12300 -L :12300 -l 1400 -t $T
sleep 10
killall -9 gonoodle
log_info_on_rp $d
d=$((d+1))

# day 2
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12400 -L :12400 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12600 -L :12600 -l 1400 -t $T
sleep 10
killall -9 gonoodle
log_info_on_rp $d
d=$((d+1))

# day 3
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 800 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 300 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 700 -R 40 -M 1 -b 1k -p 12300 -L :12300 -l 1400 -t $T
sleep 10
killall -9 gonoodle
log_info_on_rp $d
d=$((d+1))

# day 4
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12400 -L :12400 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12600 -L :12600 -l 1400 -t $T
sleep 10
killall -9 gonoodle
log_info_on_rp $d
d=$((d+1))

# day 5
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 800 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 300 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 700 -R 40 -M 1 -b 1k -p 12300 -L :12300 -l 1400 -t $T
sleep 10
killall -9 gonoodle
log_info_on_rp $d
d=$((d+1))

# day 6
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12400 -L :12400 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12600 -L :12600 -l 1400 -t $T
sleep 10
killall -9 gonoodle
log_info_on_rp $d
d=$((d+1))

# day 7
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 800 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 300 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 700 -R 40 -M 1 -b 1k -p 12300 -L :12300 -l 1400 -t $T
sleep 10
killall -9 gonoodle
log_info_on_rp $d
d=$((d+1))

# day 8
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12400 -L :12400 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12600 -L :12600 -l 1400 -t $T
sleep 10
log_info_on_rp $d
d=$((d+1))

# day 9
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 800 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 300 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 700 -R 40 -M 1 -b 1k -p 12300 -L :12300 -l 1400 -t $T
sleep 10
killall -9 gonoodle
log_info_on_rp $d
d=$((d+1))

# day 10
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12400 -L :12400 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12600 -L :12600 -l 1400 -t $T
sleep 10
killall -9 gonoodle

log_info_on_rp $d
d=$((d+1))

# day 11
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 800 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 300 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 700 -R 40 -M 1 -b 1k -p 12300 -L :12300 -l 1400 -t $T
sleep 10
killall -9 gonoodle
log_info_on_rp $d
d=$((d+1))

# day 12
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12400 -L :12400 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12600 -L :12600 -l 1400 -t $T
sleep 10
log_info_on_rp $d
d=$((d+1))

# day 13
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 800 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 300 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 700 -R 40 -M 1 -b 1k -p 12300 -L :12300 -l 1400 -t $T
sleep 10
killall -9 gonoodle
log_info_on_rp $d
d=$((d+1))

# day 14
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12400 -L :12400 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 40 -M 1 -b 1k -p 12600 -L :12600 -l 1400 -t $T
sleep 10
killall -9 gonoodle

log_info_on_rp $d
d=$((d+1))

# day 15
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 200 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 800 -R 40 -M 1 -b 1k -p 12200 -L :12200 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 300 -R 20 -M 1 -b 1k -p 12000 -L :12000 -l 1400 -t $T
sleep 10
killall -9 gonoodle
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 1000 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T &
sleep 2
ip netns exec nsInitiator /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 700 -R 40 -M 1 -b 1k -p 12300 -L :12300 -l 1400 -t $T
sleep 10
killall -9 gonoodle
log_info_on_rp $d
