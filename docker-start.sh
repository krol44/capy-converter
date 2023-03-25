#!/bin/bash
docker build . -t capy-converter-image
docker rm -f capy-converter
docker run -d \
  -p 127.0.0.1:3003:3003 \
  -e MAX_FILE_SIZE_MB="100" \
  --restart=always --log-opt max-size=5m --name=capy-converter capy-converter-image