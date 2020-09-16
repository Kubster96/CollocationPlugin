FROM ubuntu:14.04

RUN apt-get update
RUN apt-get -y install redis-server

COPY collocation-scheduler  .
COPY manifests/scheduler-config.yaml .
COPY coefficients-montage-soykb.json coefficients.json