#!/bin/sh

go build
for i in {1..5}
do
mkdir data/$i
./ra2
mv data/out.csv data/$i/out.csv
mv data/stats.txt data/$i/stats.txt
done

