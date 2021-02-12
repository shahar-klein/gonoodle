set +e
T=60
B=40m
RP=5.5.5.2
INITIATOR_NS=nsInitiator
RP_PUB_IP=30.30.30.100
TCPDUMP_PKT_CNT=10000
INITIATOR_DEV=enp1s0
LOADER_DEV=enp3s0
days=2
d=0
dat=`date '+%A%d%B%Y'`
logD="/root/15D_RP_test_start_at_$dat"
NUMC=1
NUMR=1
F=500

log_info_on_rp() {
        when=$1
        log_dir="${logD}/day${when}"
        ssh $RP "mkdir -p $log_dir"
        ssh $RP "ethtool -S enp7s0 > $log_dir/enp7s0.ethtool"
        ssh $RP "ethtool -S enp8s0 > $log_dir/enp8s0.ethtool"
        ssh $RP "cat /proc/meminfo> $log_dir/meminfo"
}

tcpdump_on_initiator() {
        d=$1
        port1=$2
        port2=$3
        log_dir="${logD}/day${d}"

        ip netns exec $INITIATOR_NS mkdir -p $log_dir
        ip netns exec $INITIATOR_NS tcpdump -enni $INITIATOR_DEV -c ${TCPDUMP_PKT_CNT} -w $log_dir/tcpdumpI.${port1}.pcap udp dst port ${port1} &
        ip netns exec $INITIATOR_NS tcpdump -enni $INITIATOR_DEV -c ${TCPDUMP_PKT_CNT} -w $log_dir/tcpdumpI.${port2}.pcap udp dst port ${port2} &
}
tcpdump_on_loader() {
        d=$1
        log_dir="${logD}/day${d}"
        port1=$2
        port2=$3

        mkdir -p $log_dir
        tcpdump -enni $LOADER_DEV -c ${TCPDUMP_PKT_CNT} -w $log_dir/tcpdumpL.${port1}.pcap udp dst port ${port1} &
        tcpdump -enni $LOADER_DEV -c ${TCPDUMP_PKT_CNT} -w $log_dir/tcpdumpL.${port2}.pcap udp dst port ${port2} &
}

ssh $RP "rm -rf $log_dir"
ip netns exec $INITIATOR_NS "rm -rf $log_dir"
rm -rf $log_dir
TS=$((T/2))

d=0
log_info_on_rp $d
P1=12000
P2=12050

echo TS=$TS

for (( d=1; d<=$days; d++ ))
do
        for (( h=0; h<1; h++ ))
        do
                /root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C $NUMC -R $NUMC -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -f $F -t $T &
                sleep 4
                ip netns exec $INITIATOR_NS /root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C $NUMC -R $NUMR -M 1 -b 1k -p 12000 -L :12000 -l 1400 -f 1000 -t $T &
                sleep $TS
                tcpdump_on_initiator $d $P1 $P2
                tcpdump_on_loader $d $P1 $P2
                sleep $TS
                sleep 3
                killall -9 gonoodle
                ip netns exec $INITIATOR_NS killall -9 gonoodle
                ip netns exec $INITIATOR_NS killall -9 tcpdump
                sleep 3
                P1=$(( $P1 + 1 ))
                P2=$(( $P2 + 1 ))
        done
        log_info_on_rp $d

done







