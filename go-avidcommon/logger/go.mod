module github.com/rkvst/avidcommon/logger

go 1.21

replace jitsuin.com/avid/correlationid => ../correlationid

require (
	github.com/KimMachineGun/automemlimit v0.2.6
	go.uber.org/automaxprocs v1.5.3
	go.uber.org/zap v1.25.0
	google.golang.org/grpc v1.57.0
	jitsuin.com/avid/correlationid v0.0.0-00010101000000-000000000000
)

require (
	github.com/cilium/ebpf v0.9.1 // indirect
	github.com/containerd/cgroups/v3 v3.0.1 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/godbus/dbus/v5 v5.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/opencontainers/runtime-spec v1.0.2 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
