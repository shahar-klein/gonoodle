FROM debian:buster-slim

RUN apt-get update \
  && apt-get install -y netcat-openbsd procps vim iputils-ping net-tools tcpdump man iperf nload\
  && apt-get clean

COPY gonoodle /
COPY ethtool_PHY_PPS.sh /
