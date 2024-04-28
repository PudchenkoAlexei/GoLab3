#!/bin/bash
curl -X POST http://localhost:17000 -d "green"
curl -X POST http://localhost:17000 -d "figure 0.5 0.5"

max=6
for (( i=0; i <= $max; ++i ))
do
curl -X POST http://localhost:17000 -d "move 0.1 0.1"
curl -X POST http://localhost:17000 -d "update"
sleep 1
done