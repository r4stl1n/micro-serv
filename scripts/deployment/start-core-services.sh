#!/bin/bash
docker-compose -f scripts/deployment/compose/docker-compose.yml up nats mongodb redis
