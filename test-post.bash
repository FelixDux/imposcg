#!/bin/bash

url1="https://4nese8424h.execute-api.eu-west-2.amazonaws.com/"
url2="api/iteration/data/"

if [ $# -gt 0 ]
then
    url="$url1$1/$url2"

    echo "$url1$1/swagger/index.html"
else
    url="$url1$url2"
fi

curl -X POST $url -H "accept: application/json" -H "Content-Type: application/x-www-form-urlencoded" -H "x-api-key header:" -d "frequency=2.8&offset=0&r=0.8&maxPeriods=100&phi=0&v=0&numIterations=1000" -v -L