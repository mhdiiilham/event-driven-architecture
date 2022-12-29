# NOP
## The NOise Processing System
Nop is a system for receiving and processing noise data from external sensors.

## System Should:
- Receive The Telemetry
- Validate Data
- Notify client of new data

## System Requirements
- Handle load -> Use Streaming Engine
- Validate data -> First thing that should happen when receiving the data
- Unknown number of clients -> Classing event requirements
- No sync users' command -> No synchronous actions required

## What needed to be installed
- Go 1.18+
- VSCode (or any code editor)
- RabbitMQ (the channel, this will included in the docker-compose.yaml file)
