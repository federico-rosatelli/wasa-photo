#!/usr/bin/env sh

sudo docker run -i -v "$(pwd):/src" --network host --workdir "$(pwd)" wasa-photo-backend