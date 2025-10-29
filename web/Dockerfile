# ---- Build ----
FROM rust:1.90.0-slim AS build

RUN apt-get update && apt-get install -y curl pkg-config libssl-dev \
    && curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh

WORKDIR /app/wasm

COPY wasm/ .

RUN wasm-pack build --target web --release --out-dir /app/web/pkg

# ---- Runtime ----
FROM nginx:alpine AS runtime

COPY index.html /usr/share/nginx/html/index.html
COPY --from=build /app/web/pkg /usr/share/nginx/html/pkg

EXPOSE 80