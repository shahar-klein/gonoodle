#!/bin/bash

set -x

for i in {1..10} ; do
	date
	/root/ws/git/gonoodle/gonoodle -u -c 30.30.30.100 --rp initiator -C 1000 -R  1000  -M 1 -b 1k -p 12000 -L :7000 -l 1400 -t 86400
done

echo "Done"
