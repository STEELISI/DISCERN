module metaOS

go 1.21.3

require (
	DataSorcerers/helpers v0.0.0-00010101000000-000000000000
	FusionBridge v0.0.0-00010101000000-000000000000
	OS v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.5.3
	google.golang.org/grpc v1.58.2
)

require (
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230711160842-782d3b101e98 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace FusionBridge => ../../../FusionBridge

replace OS => ./

replace DataSorcerers/helpers => ../../helpers
