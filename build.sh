go build -o collocation-scheduler ./cmd &&
docker build . -t kubster96/scheduler:collocation-scheduler-latest &&
docker push kubster96/scheduler:collocation-scheduler-latest