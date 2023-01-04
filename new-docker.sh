#!/usr/bin/env sh

docker run -i -v "$(pwd):/src" --rm --network host --workdir /src --name=BackEndContainer wasa-photo-backend