#!/bin/bash



# Slowly go up to 1000 at 20mb and stay there for an hour or so
#./gonoodle -u -c 10.8.51.71 --rp loader -C 1000 -R  10 -M 10 -b 20m -p 5000 -L 10.8.51.71:8000 -l 1400 -t 40

#sleep 60

# 200 sessions for 10 hours 40mb
#./gonoodle -u -c 10.8.51.71 --rp loader -C 200 -R 20 -M 10 -b 40m -p 5000 -L 10.8.51.71:8000 -l 1400 -t 360 &


# Slowly go up to 1000 at 20mb and stay there for an hour or so
/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 10 -M 10 -b 20m -p 5000 -L 5.5.50.0:47998 -l 1400 -t 4000

#let sessions age out
sleep 120


/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 100 -M 10 -b 20m -p 5000 -L 5.5.50.0:47998 -l 1400 -t 36000


sleep 120

/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 100 -M 10 -b 20m -p 5000 -L 5.5.50.0:47998 -l 1400 -t 3600

sleep 120

/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 100 -M 10 -b 20m -p 5000 -L 5.5.50.0:47998 -l 1400 -t 72000 

sleep 120

/root/ws/git/gonoodle/gonoodle -u -c 5.5.5.2 --rp loader_multi -C 1000 -R 100 -M 10 -b 20m -p 5000 -L 5.5.50.0:47998 -l 1400 -t 72000 




