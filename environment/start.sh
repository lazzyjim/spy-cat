#!/bin/bash
docker-compose up -d --remove-orphans

POSTGRESQL_IP=`docker inspect  --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' spy_cat.postgresql`
docker exec spy_cat.postgresql ./initialize_db.sh
echo "Docker images:"
echo "PostgreSQL IP: ${POSTGRESQL_IP}"

