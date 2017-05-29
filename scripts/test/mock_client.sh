#!/usr/bin/env bash

function register(){
    curl -v -H "Content-Type: application/json" \
        -d '{ "auth_string" : "'$1'", "hostname" : "'$(hostname)'"}' \
    "0.0.0.0:8080/register"
}

function post(){
    echo "$3"
    curl -v -H "Content-Type: application/json" \
        -d '{ "uuid" : "'$1'", "auth_string" : "'$2'", "hostname" : "'$(hostname)'"}' \
        "0.0.0.0:8080/$3"
}

function usage(){
str="Usage: $0 [-r(egister) -l(ogin) -L(ogout) -p(ing) -s(status)] <uuid> <auth_string>
     e.g: $0 -r 5a09aea6-69d7-43be-ba2b-6d8ef562f97d 2BE5501D15C15A53\n"
    printf "$str"
}

if [ -z "$1" ]; then
    usage
    exit -1
fi

while getopts "r:l:L:p:s:" opt; do
    case "$opt" in
        r) register $1
            ;;
        l) post $1 $2 "login"
            ;;
        L) post $1 $2 "logoff"
            ;;
        p) post $1 $2 "ping"
            ;;
        s) post $1 $2 "status"
            ;;
    esac
done

