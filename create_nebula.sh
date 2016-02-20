#!/bin/bash

NebulaSize=10

get_host_IP () {
    cat /etc/hosts | grep $1 | awk '{print $1}'
}

for i in `seq 0 $(($NebulaSize - 1))`
do
    echo "Creating node$i ..."
    #docker kill node$i
    docker rm node$i > /dev/null

    # Create the first node of nebula cluster
    if [ $i -eq 0 ]; then
        docker run -h node$i --name=node$i -t v1ct0r/nebula:v1 /bin/bash ./run_init_node.sh &
        continue
    fi

    # Wait for 3 seconds
    sleep 3

    # Create a normal node
    docker run -h node$i --name=node$i --link=node0 -t v1ct0r/nebula:v1 /bin/bash ./run_node.sh &
done

exit
