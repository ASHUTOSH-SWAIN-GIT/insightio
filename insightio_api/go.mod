module github.com/ASHUTOSH-SWAIN-GIT/insightio_api

go 1.25.5

require (
	github.com/ASHUTOSH-SWAIN-GIT/insightio v0.0.0-20251213034838-ebc37056f937
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.77.0
	google.golang.org/protobuf v1.36.10
)

require (
	golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
)

// Replace with local path - adjust the path based on your directory structure
replace github.com/ASHUTOSH-SWAIN-GIT/insightio => ../insightio
