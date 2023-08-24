// Package grpchealth package provides server implementing Check rpc that meets https://github.com/grpc/grpc/blob/master/doc/health-checking.md
package grpchealth

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rkvst/go-rkvstcommon/logger"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	livenessServiceName  = "liveness"
	readinessServiceName = "readiness"
)

type HealthCheckingService struct {
	grpc_health_v1.UnimplementedHealthServer
	sync.RWMutex
	healthStatus map[string]grpc_health_v1.HealthCheckResponse_ServingStatus
}

func New() HealthCheckingService {
	return HealthCheckingService{
		healthStatus: map[string]grpc_health_v1.HealthCheckResponse_ServingStatus{
			livenessServiceName:  grpc_health_v1.HealthCheckResponse_SERVING,
			readinessServiceName: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		},
	}
}

func (s *HealthCheckingService) serving(service string) {
	s.Lock()
	defer s.Unlock()
	s.healthStatus[service] = grpc_health_v1.HealthCheckResponse_SERVING
	logger.Sugar.Infof("Health set to 'SERVING': %s", service)
}

func (s *HealthCheckingService) notServing(service string) {
	s.Lock()
	defer s.Unlock()
	s.healthStatus[service] = grpc_health_v1.HealthCheckResponse_NOT_SERVING
	logger.Sugar.Infof("Health set to 'NOT_SERVING': %s", service)
}

// Dead - changes status of service to dead
func (s *HealthCheckingService) Dead() {
	s.notServing(livenessServiceName)
}

// Live - changes status of service to alive
func (s *HealthCheckingService) Live() {
	s.serving(livenessServiceName)
}

// NotReady - changes status of service to not ready
func (s *HealthCheckingService) NotReady() {
	s.notServing(readinessServiceName)
}

// Ready - changes status of service to ready
func (s *HealthCheckingService) Ready() {
	s.serving(readinessServiceName)
}

// Check implements `service Health`.
func (s *HealthCheckingService) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (
	*grpc_health_v1.HealthCheckResponse, error) {
	s.RLock()
	defer s.RUnlock()

	// logger.Sugar.Debugf("Health Check for '%s'", in.Service)
	if in.Service == "" {
		for _, v := range s.healthStatus {
			// logger.Sugar.Debugf("Health Check for '%s'-> '%s'", in.Service, v.String())
			if v != grpc_health_v1.HealthCheckResponse_SERVING {
				logger.Sugar.Infof("Health Check '%s' is NOT SERVING: '%s'", in.Service, v.String())
				return &grpc_health_v1.HealthCheckResponse{
					Status: v,
				}, nil
			}
		}
		logger.Sugar.Infof("Health Check '%s' is SERVING", in.Service)
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		}, nil
	}
	if stat, ok := s.healthStatus[in.Service]; ok {
		// logger.Sugar.Debugf("Health Check '%s' -> '%s'", in.Service, stat.String())
		logger.Sugar.Infof("Health Check '%s' is `%s'", in.Service, stat)
		return &grpc_health_v1.HealthCheckResponse{
			Status: stat,
		}, nil
	}
	err := status.Error(codes.NotFound, "unknown service: "+in.Service)

	logger.Sugar.Infof("Health Check failed: %v", err)
	return nil, err
}

func (s *HealthCheckingService) Watch(in *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	logger.Sugar.Infof("Health Check watch not supported")
	return status.Error(codes.Unimplemented, "watch not supported")
}
