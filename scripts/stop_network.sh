#!/bin/bash
docker-compose -f docker-compose-mrtgexchg.yaml down

sleep 10

./scripts/cleanup.sh
