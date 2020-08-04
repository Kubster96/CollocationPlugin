FROM ubuntu:14.04
COPY app .
COPY manifests/scheduler-config.yaml .
COPY coefficients.json .