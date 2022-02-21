#! /bin/sh

service=dp-find-insights-poc-api

binary=dp-find-insights-poc-api
dir=/home/ubuntu

current="$dir/$binary"
next="$dir/.${binary}.next"
prev="$dir/.${binary}.prev"

usage() {
    echo "usage: $0 <new-binary>|previous" 1>&2
    exit $1
}

main() {
    new=$1
    if test -z "$new"
    then
        usage 2
    fi

    if test "$new" = previous
    then
        new=$prev
    fi

    # get file ready to install
    cp -a "$new" "$next" || exit

    # save current executable (ok if there isn't a current)
    cp -a "$current" "$prev"

    # install new executable
    mv "$next" "$current" || exit

    sudo systemctl restart "$service"
    systemctl status "$service"
}

main "$@"
