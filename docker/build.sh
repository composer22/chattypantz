#!/bin/bash
docker build -t composer22/chattypantz_build .
docker run -v /var/run/docker.sock:/var/run/docker.sock -v $(which docker):$(which docker) -ti --name chattypantz_build composer22/chattypantz_build
docker rm chattypantz_build
docker rmi composer22/chattypantz_build
