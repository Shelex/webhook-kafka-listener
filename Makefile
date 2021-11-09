NAME=webhook-listener
ROOT=github.com/Shelex/${NAME}
GO111MODULE=on
SHELL=/bin/bash

.PHONY: start
start:
	docker-compose up

.PHONY: simulation
simulation:
	k6 run simulation.js

.PHONY: open
open-web:
	open http://localhost:8080