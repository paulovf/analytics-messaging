# Analytics Messaging

A Go library providing an abstraction layer for RabbitMQ communication, ensuring infrastructure decoupling across the Analytics platform microservices.

## Overview

This module is not a standalone web service. It is designed to be imported by other Go projects within the platform (such as `analytics-api` and `analytics-reports`) to handle message publishing and consuming consistently.

## Installation

```bash
go get [github.com/paulovf/analytics-messaging](https://github.com/paulovf/analytics-messaging)
```

## Testing

To run the unit tests in an isolated environment, you can use the provided Dockerfile:

```bash
docker build -t analytics-messaging-tests -f Dockerfile.test .
docker run --rm analytics-messaging-tests
```
