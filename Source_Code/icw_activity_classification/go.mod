module icw_activity_classification

go 1.26.0

require icw_common v0.0.0

replace icw_common => ../icw_common

require (
	golang.org/x/image v0.39.0
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/sqids/sqids-go v0.4.1 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)
