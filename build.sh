#!/bin/bash
rm main
rm function.zip
GOOS=linux go build -o main ./cmd/stm
zip function.zip main