package common

import (
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

type DataDogConfig struct {
	AgentAddr        string `env:"DD_AGENT_ADDR"       envDefault:"localhost:8126"`
	ServiceName      string `env:"DD_SERVICE"          envDefault:"my-service"`
	Environment      string `env:"DD_ENV"              envDefault:"dev"`
	Version          string `env:"DD_VERSION"          envDefault:"1.0.0"`
	Enabled          bool   `env:"DD_ENABLED"          envDefault:"false"`
	ProfilingEnabled bool   `env:"DD_PROFILING_ENABLED" envDefault:"false"`
}

func InitDataDog(cfg *DataDogConfig) error {
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
