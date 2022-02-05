#!/bin/bash

for f in $(find . -name *.json)
do
    mongoimport --uri "mongodb://0.0.0.0/cms" --collection $(basename $f .json) --file $f --jsonArray --drop
done