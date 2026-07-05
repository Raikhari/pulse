#!/bin/bash

curl localhost:8080/metrics | jq 'map(. + {time: (.timestamp | strftime("%Y-%m-%d %H:%M:%S"))})'
