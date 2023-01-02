#!/usr/bin/env sh

sudo docker run -i -v "$(pwd):/src" --network host --workdir /src wasa-photo-backend