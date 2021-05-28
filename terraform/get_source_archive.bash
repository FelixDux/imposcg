#!/bin/bash

fgrep source_archive "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/terraform.tfvars | cut -d'=' -f2 | xargs