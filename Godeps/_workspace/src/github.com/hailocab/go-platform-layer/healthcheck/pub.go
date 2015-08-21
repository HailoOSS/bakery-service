package healthcheck

import (
	"github.com/hailocab/go-platform-layer/client"
)

// pubLastSample pings this healthcheck sample out into the ether
func pubLastSample(hc *HealthCheck, ls *Sample) {
	client.Pub("com.hailocab.monitor.healthcheck", healthCheckSampleToProto(hc, ls))
}
