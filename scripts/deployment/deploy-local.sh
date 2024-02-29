#!/bin/bash
docker-compose -f scripts/deployment/compose/docker-compose.yml up --force-recreate --remove-orphans --build -d
