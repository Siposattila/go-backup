#!/bin/sh

PWD=$(pwd -P)

echo "Generating dummy data for dokcer test env..."

mkdir -p $PWD/docker/data/small
for i in $(seq 1 5)
do
    head -c 5M </dev/urandom > $PWD/docker/data/small/dummy${i}.bin
done

echo "Finished generating small dummy data..."

mkdir -p $PWD/docker/data/medium
for i in $(seq 1 5)
do
    head -c 50M </dev/urandom > $PWD/docker/data/medium/dummy${i}.bin
done

echo "Finished generating medium dummy data..."

mkdir -p $PWD/docker/data/large
for i in $(seq 1 5)
do
    head -c 500M </dev/urandom > $PWD/docker/data/large/dummy${i}.bin
done

echo "Finished generating large dummy data..."
echo "Finished generating dummy data! :)"
