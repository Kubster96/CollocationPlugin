FROM ubuntu:14.04
COPY collocation-scheduler  .
COPY manifests/scheduler-config.yaml .
COPY coefficients.json .