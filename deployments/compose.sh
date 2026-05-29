#!/usr/bin/env bash

set -euo pipefail

COMPOSE_FILE="$(dirname "$0")/docker-compose.dev.yaml"

case "${1:-}" in
  up)
    docker compose -f "$COMPOSE_FILE" up -d
    ;;
  down)
    docker compose -f "$COMPOSE_FILE" down
    ;;
  logs)
    docker compose -f "$COMPOSE_FILE" logs -f "${@:2}"
    ;;
  restart)
    docker compose -f "$COMPOSE_FILE" restart "${@:2}"
    ;;
  *)
    echo "Usage: $0 {up|down|logs|restart}"
    exit 1
    ;;
esac
