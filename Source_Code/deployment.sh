docker compose down
docker pull apache/rocketmq:5.3.2
docker pull docker:27-cli
docker pull minmuslin/icw_activity_classification:latest
docker pull minmuslin/icw_activity_reasoning:latest
docker pull minmuslin/icw_activity_summary:latest
docker pull minmuslin/icw_core_api:latest
docker pull minmuslin/icw_core_biz:latest
docker pull minmuslin/icw_core_web:latest
docker compose up -d
