#!/bin/bash

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

source_archive=imposcg_$(git rev-list HEAD -1).zip

echo source_archive = "\"$source_archive\"" > ./terraform.tfvars

terraform apply