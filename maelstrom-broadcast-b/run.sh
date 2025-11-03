#!/bin/bash

go build -v
../vendor/maelstrom/maelstrom test -w broadcast --bin maelstrom-broadcast-b --node-count 5 --time-limit 20 --rate 10
