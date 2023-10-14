#####################################################
#  generates gRPC server from provided proto files  #
#####################################################

function generate() {
  if [ -z "$1" ]; then
      echo "❌️ proto-path is empty"
      exit 1
  fi

  if [ -z "$2" ]; then
        echo "❌ .proto filename not provided"
        exit 1
  fi

  mkdir -p "$2"

  echo "⏳  generating gRPC server..."
  protoc --proto_path="$1" \
    --go_out="$2" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$2" \
    --go-grpc_opt=paths=source_relative \
    "$3"

  echo "⏳  downloading dependencies..."
  go mod tidy

  echo "🎉 server was generated"
}

generate "$1" "$2" "$3"