# InsightIO Examples

This directory contains example code demonstrating how to use InsightIO components.

## SDK Example

The `sdk_example` demonstrates how to use the InsightIO SDK to send events to the analytics engine.

### Prerequisites

1. Start the InsightIO gRPC server:
   ```bash
   cd ../insightio
   export INSIGHTIO_API_KEY="test-api-key-123"
   go run ./cmd/server
   ```

2. Start the InsightIO API gateway:
   ```bash
   cd ../insightio_api
   go run ./cmd
   ```

### Running the Example

```bash
cd sdk_example
go run main.go
```

The example will send various types of events to demonstrate the SDK usage.

