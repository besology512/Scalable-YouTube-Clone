#!/bin/bash

# Tree structure as input
read -r -d '' tree << EOM
.
├── api-gateway
│   └── gateway-config/
├── auth-service
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   ├── proto/
│   └── go.mod
├── deployments
│   ├── docker-compose.yml
│   └── k8s/
├── README.md
├── recommendation-service
│   ├── app/
│   ├── models/
│   ├── scripts/
│   └── requirements.txt
├── shared-libs
│   └── proto/
├── streaming-service
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   ├── proto/
│   └── go.mod
└── video-upload-service
    ├── cmd/
    ├── internal/
    ├── pkg/
    ├── proto/
    └── go.mod
EOM

# Create folders and files
while IFS= read -r line; do
    # Clean formatting characters
    clean_line=$(echo "$line" | sed 's/^[^a-zA-Z0-9]*//')
    
    # If line ends with '/', it's a directory
    if [[ "$clean_line" == */ ]]; then
        mkdir -p "$clean_line"
    elif [[ -n "$clean_line" ]]; then
        # Otherwise it's a file
        mkdir -p "$(dirname "$clean_line")"
        touch "$clean_line"
    fi
done <<< "$tree"
