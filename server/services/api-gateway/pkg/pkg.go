package pkg

// TODO: Step 1 - Create service discovery utilities:
//   - ServiceRegistry interface (Register, Deregister, Discover)
//   - ConsulRegistry implementation
//   - InMemoryRegistry for development
//   - ServiceEndpoint struct (Host, Port, Health, Metadata)
// TODO: Step 2 - Implement load balancing utilities:
//   - LoadBalancer interface with different strategies
//   - RoundRobinBalancer implementation
//   - WeightedBalancer implementation
//   - HealthAwareBalancer that excludes unhealthy services
// TODO: Step 3 - Add monitoring and metrics utilities:
//   - RequestMetrics (count, latency, error rate per service)
//   - PrometheusExporter for metrics collection
//   - HealthChecker for service monitoring
// TODO: Step 4 - Create configuration utilities:
//   - RoutingConfig struct (path patterns, target services, middleware)
//   - ConfigLoader for dynamic configuration updates
//   - ConfigValidator for route rule validation

// Flow: External configuration -> pkg utilities -> proxy/handlers
// Exports: Service discovery, load balancing, monitoring interfaces
