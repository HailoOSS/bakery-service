package gossie

var (
	invalidEndpoint    = "localhost:9999"
	localEndpoint      = "localhost:19160"
	localEndpointPool  = []string{localEndpoint}
	localEndpointsPool = []string{localEndpoint}

	keyspace = "TestGossie"

	standardTimeout = 3000
	shortTimeout    = 1000

	poolOptions = PoolOptions{Size: 50, Timeout: standardTimeout}
)
