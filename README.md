# Checking out go windmill library with kafka

Dead simple http listener of webhooks that produces messages to kafka

## Setup

- Requires docker engine
- Run:

```bash
   make start
```

- API will accept POST requests to `http://localhost:9000/webhooks` url

## Simulate webhooks

- Requires [k6.io](https://k6.io/docs/getting-started/installation/)
- Run:

```bash
    make simulation
```

## Check out UI for Kafka

- Run:

```bash
    make open
```

- kafka dashboard will be available at `http://localhost:8080`
