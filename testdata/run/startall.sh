#!/bin/bash
cd n0
./start.sh > ../log/0.log &
cd ../n1
./start.sh > ../log/1.log &
cd ../n2
./start.sh > ../log/2.log &
cd ../n3
./start.sh > ../log/3.log &
cd ../n4
./start.sh > ../log/4.log &
cd ../n5
./start.sh > ../log/5.log &
cd ../n6
./start.sh > ../log/6.log &
