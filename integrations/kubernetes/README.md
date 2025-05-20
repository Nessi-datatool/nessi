# Nessi Kubernetes Integration

This package provides integration between Nessi.dev and Kubernetes, allowing you to run Nessi data quality operations as Kubernetes jobs.

## Features

- Run data quality checks, profiling, and validation as Kubernetes jobs
- Configurable resource allocation, node selection, and tolerations
- Support for environment variables, ConfigMaps, and Secrets
- Automatic job cleanup with configurable TTL
- Comprehensive logging and error handling
- Support for waiting for job completion or asynchronous operation

## Installation

```bash
pip install nessi-k8s
```

## Requirements

- Python 3.7+
- Kubernetes Python client (`kubernetes`)
- Access to a Kubernetes cluster
- Nessi.dev API credentials (optional)

## Configuration

The Kubernetes integration can be configured in multiple ways:

1. **Configuration file** - YAML or JSON file with all settings
2. **Environment variables** - Prefixed with `NESSI_K8S_`
3. **Direct parameters** - Passed to the client constructor

### Configuration File Example

```yaml
# nessi-k8s-config.yaml
namespace: nessi
image: nessi/nessi:latest
service_account_name: nessi-service-account
api_host: http://nessi-api:8080
api_key: your-api-key
api_secret: your-api-secret
job_ttl_seconds_after_finished: 3600
job_backoff_limit: 3
job_timeout_seconds: 600
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 512Mi
node_selector:
  node-role: data-processing
tolerations:
  - key: dedicated
    operator: Equal
    value: data-processing
    effect: NoSchedule
labels:
  app: nessi
  environment: production
annotations:
  description: Nessi data quality job
env_from:
  - config_map_ref:
      name: nessi-config
  - secret_ref:
      name: nessi-secrets
```

### Environment Variables

You can use environment variables prefixed with `NESSI_K8S_` to configure the integration:

```bash
export NESSI_K8S_NAMESPACE=nessi
export NESSI_K8S_IMAGE=nessi/nessi:latest
export NESSI_K8S_SERVICE_ACCOUNT_NAME=nessi-service-account
export NESSI_K8S_API_HOST=http://nessi-api:8080
export NESSI_K8S_API_KEY=your-api-key
export NESSI_K8S_API_SECRET=your-api-secret
export NESSI_K8S_JOB_TTL_SECONDS=3600
export NESSI_K8S_JOB_BACKOFF_LIMIT=3
export NESSI_K8S_JOB_TIMEOUT_SECONDS=600
export NESSI_K8S_RESOURCES='{"requests": {"cpu": "100m", "memory": "128Mi"}, "limits": {"cpu": "500m", "memory": "512Mi"}}'
export NESSI_K8S_NODE_SELECTOR='{"node-role": "data-processing"}'
export NESSI_K8S_TOLERATIONS='[{"key": "dedicated", "operator": "Equal", "value": "data-processing", "effect": "NoSchedule"}]'
export NESSI_K8S_LABELS='{"app": "nessi", "environment": "production"}'
export NESSI_K8S_ANNOTATIONS='{"description": "Nessi data quality job"}'
```

## Usage

### Basic Usage

```python
from nessi_k8s.client import NessiK8sClient

# Create client
client = NessiK8sClient(
    namespace="nessi",
    image="nessi/nessi:latest",
    service_account_name="nessi-service-account",
    api_host="http://nessi-api:8080",
    api_key="your-api-key",
    api_secret="your-api-secret",
)

# Run a data quality check
result = client.run_quality_check(
    table_name="/path/to/table",
    rules=[
        {"name": "not_null_check", "rule_type": "not_null", "column": "id"},
        {"name": "unique_check", "rule_type": "unique", "column": "id"},
    ],
    profile=True,
    output_format="json",
    wait_for_completion=True,
)

print(f"Job status: {result}")
```

### Using a Configuration File

```python
from nessi_k8s.client import NessiK8sClient

# Create client from configuration file
client = NessiK8sClient(config_file="nessi-k8s-config.yaml")

# Run a data profile
result = client.run_profile(
    table_name="/path/to/table",
    columns=["id", "name", "email"],
    sample_size=1000,
    wait_for_completion=True,
)

print(f"Job status: {result}")
```

### Asynchronous Operation

```python
from nessi_k8s.client import NessiK8sClient

# Create client
client = NessiK8sClient(
    namespace="nessi",
    image="nessi/nessi:latest",
)

# Run a data validation asynchronously
result = client.run_validation(
    table_name="/path/to/table",
    rules=[
        {"name": "not_null_check", "rule_type": "not_null", "column": "id"},
        {"name": "unique_check", "rule_type": "unique", "column": "id"},
    ],
    wait_for_completion=False,
)

print(f"Job created: {result['name']}")

# Later, check the results
job_result = client.get_validation_results(
    job_name=result["name"],
    wait_for_completion=True,
)

print(f"Job status: {job_result['status']}")
print(f"Job logs: {job_result['logs']}")
```

## Command Line Interface

The package includes a command-line interface for running Nessi operations in Kubernetes:

```bash
# Run a data quality check
python -m nessi_k8s.cli quality-check --table /path/to/table --rules rules.json --namespace nessi

# Run a data profile
python -m nessi_k8s.cli profile --table /path/to/table --columns id,name,email --namespace nessi

# Run a data validation
python -m nessi_k8s.cli validate --table /path/to/table --rules rules.json --namespace nessi

# Create a default configuration file
python -m nessi_k8s.cli create-config --output nessi-k8s-config.yaml
```

## Kubernetes Deployment

See the [examples/kubernetes-deployment.yaml](examples/kubernetes-deployment.yaml) file for an example of how to deploy Nessi in a Kubernetes cluster.

## Examples

Check out the [examples](examples) directory for more usage examples:

- [quality_check_example.py](examples/quality_check_example.py) - Example of running a data quality check
- [kubernetes-deployment.yaml](examples/kubernetes-deployment.yaml) - Example Kubernetes deployment

## License

This project is licensed under the MIT License - see the LICENSE file for details.
