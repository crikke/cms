#!/bin/bash

for f in $(find ./testdata -name *.json)
do
    mongoimport --username $1 --password $2 --uri $3 --collection $(basename $f .json)
done