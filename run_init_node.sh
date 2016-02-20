get_host_IP () {
    cat /etc/hosts | grep $1 | awk '{print $1}'
}

go run main.go --addr=$(get_host_IP $(hostname))
