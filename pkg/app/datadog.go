package app

import (
	"github.com/gingray/go-template/pkg/config"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func InitDataDog(cfg *config.DataDogConfig) error {
	if !cfg.Enabled {
		return nil
	}
	tracer.Start(
		tracer.WithAgentAddr(cfg.AgentAddr),
		tracer.WithServiceName(cfg.ServiceName),
		tracer.WithEnv(cfg.Environment),
		tracer.WithServiceVersion(cfg.Version),
	)

	if cfg.ProfilingEnabled {
		err := profiler.Start(
			profiler.WithService(cfg.ServiceName),
			profiler.WithEnv(cfg.Environment),
			profiler.WithVersion(cfg.Version),
			profiler.WithProfileTypes(
				profiler.CPUProfile,
				profiler.HeapProfile,
			))
		if err != nil {
			tracer.Stop()
			return err
		}
	}
	return nil
}
