#!/bin/bash
set -e

mkdir -p bin
javac src/GenerateKeys.java -d bin
java -cp bin GenerateKeys
