go build -o app ./cmd &&
  docker build . -t collocation-scheduler || exit 1