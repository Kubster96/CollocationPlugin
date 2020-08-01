go build -o app . &&
  docker build . -t collocation-scheduler || exit 1