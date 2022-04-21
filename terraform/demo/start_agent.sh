#! /bin/bash
cd hooks-docker-example
docker build --rm --no-cache -t local-agent .

export TFC_AGENT_TOKEN=xxx

export TFC_AGENT_NAME=agent

docker run -e TFC_AGENT_TOKEN -e TFC_AGENT_NAME local-agent
