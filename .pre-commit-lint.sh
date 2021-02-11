if [ "$(git branch --show-current)" == "main" ]; then
  echo "made it"
    golangci-lint run
fi
