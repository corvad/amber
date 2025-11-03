#!/bin/bash

go build -v
../vendor/maelstrom/maelstrom test -w broadcast --bin maelstrom-broadcast-a --node-count 1 --time-limit 20 --rate 10
