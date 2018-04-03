#!/usr/bin/env bash
cd /tmp
for i in $(seq 0 31)
do
    jobid=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1)
    node=$(cat /dev/urandom | tr -dc '1-8' | fold -w 1 | head -n 1)
    #echo $jobid
    #echo $node
    curl "http://localhost:8080/job/start?jobid=$jobid&node=$node" && sleep 5 && curl "http://localhost:8080/job/stop?jobid=$jobid&node=$node&result=ko" &
done

    jobid=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1)
    node=$(cat /dev/urandom | tr -dc '1-8' | fold -w 1 | head -n 1)
    curl "http://localhost:8080/job/start?jobid=$jobid&node=$node" && sleep 10 && curl "http://localhost:8080/job/stop?jobid=$jobid&node=$node&result=ok&pass=yeah" &

cd -
