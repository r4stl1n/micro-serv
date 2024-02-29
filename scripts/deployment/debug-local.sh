#!/bin/bash

if [ $# -eq 1 ]
then

if [ $1 = "all" ]; then
  for dir in scripts/deployment/docker/*     # list directories in the form "/tmp/dirname/"
  do
      dir=${dir%*/}      # remove the trailing "/"
      file="${dir##*/}"    # print everything after the final "/"
      name="${file%%.*}"
      docker build . -t micro-$name -f $dir
  done
else
  docker build . -t micro-$1 -f scripts/deployment/docker/$1.docker
fi

fi

docker-compose -f scripts/deployment/compose/docker-compose.yml up --force-recreate --remove-orphans --build
