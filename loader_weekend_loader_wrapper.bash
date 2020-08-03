#!/bin/bash



# Slowly go up to 1000 at 20mb and stay there for an hour or so
#./gonoodle -u -c 10.8.51.71 --rp loader -C 1000 -R  10 -M 10 -b 20m -p 5000 -L 10.8.51.71:8000 -l 1400 -t 40

#sleep 60

# 200 sessions for 10 hours 40mb
#./gonoodle -u -c 10.8.51.71 --rp loader -C 200 -R 20 -M 10 -b 40m -p 5000 -L 10.8.51.71:8000 -l 1400 -t 360 &

T=21600
B=20m
# day 1
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 200 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 800 -R 10 -M 10 -b $B -p 12200 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 300 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 700 -R 10 -M 10 -b $B -p 12300 -L 5.5.50.0:47998 -l 1400 -t $T
# day 2
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 400 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 600 -R 10 -M 10 -b $B -p 12400 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 500 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 500 -R 10 -M 10 -b $B -p 12500 -L 5.5.50.0:47998 -l 1400 -t $T
# day 3
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 400 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 600 -R 10 -M 10 -b $B -p 12400 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 500 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 500 -R 10 -M 10 -b $B -p 12500 -L 5.5.50.0:47998 -l 1400 -t $T
# day 4
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 200 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 800 -R 10 -M 10 -b $B -p 12200 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 300 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 700 -R 10 -M 10 -b $B -p 12300 -L 5.5.50.0:47998 -l 1400 -t $T
# day 5
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 400 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 600 -R 10 -M 10 -b $B -p 12400 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 500 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 500 -R 10 -M 10 -b $B -p 12500 -L 5.5.50.0:47998 -l 1400 -t $T
# day 6
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 200 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 800 -R 10 -M 10 -b $B -p 12200 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 300 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 700 -R 10 -M 10 -b $B -p 12300 -L 5.5.50.0:47998 -l 1400 -t $T
# day 7
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 400 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 600 -R 10 -M 10 -b $B -p 12400 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 500 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 500 -R 10 -M 10 -b $B -p 12500 -L 5.5.50.0:47998 -l 1400 -t $T
# day 8
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 400 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 600 -R 10 -M 10 -b $B -p 12400 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 500 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 500 -R 10 -M 10 -b $B -p 12500 -L 5.5.50.0:47998 -l 1400 -t $T
# day 9
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 200 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 800 -R 10 -M 10 -b $B -p 12200 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 300 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 700 -R 10 -M 10 -b $B -p 12300 -L 5.5.50.0:47998 -l 1400 -t $T
# day 10
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 400 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 600 -R 10 -M 10 -b $B -p 12400 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 500 -R 10 -M 10 -b $B -p 12000 -L 5.5.50.0:47998 -l 1400 -t $T
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 500 -R 10 -M 10 -b $B -p 12500 -L 5.5.50.0:47998 -l 1400 -t $T


