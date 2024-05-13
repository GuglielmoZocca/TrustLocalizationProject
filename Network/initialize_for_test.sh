#!/bin/bash

#Start a network

./network.sh up -bft

#Create a channel

./network.sh createChannel

#Deploy a contract

./network.sh deployCC -ccn Position -ccp ../Chaincode_dir/ -ccl go  -cccg ../Chaincode_dir/collections_config.json



