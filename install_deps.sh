#!/bin/bash

mkdir -p vendor
cd vendor
wget https://github.com/jepsen-io/maelstrom/releases/download/v0.2.4/maelstrom.tar.bz2
rm -rf maelstrom
tar xf maelstrom.tar.bz2 -v --overwrite -C .
rm maelstrom.tar.bz2*
echo -e "\nInstalled maelstrom to vendor/maelstrom"