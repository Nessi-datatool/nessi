"""
Kubernetes operators for Nessi.dev.

This module provides Kubernetes operators for running Nessi.dev operations
as Kubernetes jobs.
"""

import os
import json
import logging
import time
import uuid
from typing import Dict, Any, Optional, List, Union, Tuple

from kubernetes import client, config
from kubernetes.client.rest import ApiException

logger = logging.getLogger(__name__)


class NessiK8sOperator:
    """
    Base class for Nessi Kubernetes operators.
    
    This class provides common functionality for running Nessi operations
    in Kubernetes.
    """
    
    def __init__(
        self,
        namespace: str = "default",
        image: str = "nessi/nessi:latest",
        service_account_name: Optional[str] = None,
        api_host: Optional[str] = None,
        api_key: Optional[str] = None,
        api_secret: Optional[str] = None,
        job_ttl_seconds_after_finished: int = 3600,
        job_backoff_limit: int = 3,
        job_timeout_seconds: int = 600,
        resources: Optional[Dict[str, Any]] = None,
        node_selector: Optional[Dict[str, str]] = None,
        tolerations: Optional[List[Dict[str, Any]]] = None,
        labels: Optional[Dict[str, str]] = None,
        annotations: Optional[Dict[str, str]] = None,
        env_from: Optional[List[Dict[str, Any]]] = None,
        volumes: Optional[List[Dict[str, Any]]] = None,
        volume_mounts: Optional[List[Dict[str, Any]]] = None,
        image_pull_secrets: Optional[List[Dict[str, Any]]] = None,
        config_file: Optional[str] = None,
        context: Optional[str] = None,
    ):
        """
        Initialize the Nessi Kubernetes operator.
        
        Args:
            namespace: Kubernetes namespace to run the job in
            image: Docker image to use for the job
            service_account_name: Service account to use for the job
            api_host: Nessi API host URL
            api_key: Nessi API key
            api_secret: Nessi API secret
            job_ttl_seconds_after_finished: Time to live for completed jobs
            job_backoff_limit: Number of retries for failed jobs
            job_timeout_seconds: Timeout for jobs in seconds
            resources: Resource requests and limits for the job
            node_selector: Node selector for the job
            tolerations: Tolerations for the job
            labels: Labels to add to the job
            annotations: Annotations to add to the job
            env_from: Environment variables from ConfigMaps or Secrets
            volumes: Volumes to mount in the job
            volume_mounts: Volume mounts for the job
            image_pull_secrets: Image pull secrets for the job
            config_file: Kubernetes config file path
            context: Kubernetes config context
        """
        self.namespace = namespace
        self.image = image
        self.service_account_name = service_account_name
        self.api_host = api_host
        self.api_key = api_key
        self.api_secret = api_secret
        self.job_ttl_seconds_after_finished = job_ttl_seconds_after_finished
        self.job_backoff_limit = job_backoff_limit
        self.job_timeout_seconds = job_timeout_seconds
        self.resources = resources or {
            "requests": {"cpu": "100m", "memory": "128Mi"},
            "limits": {"cpu": "500m", "memory": "512Mi"},
        }
        self.node_selector = node_selector
        self.tolerations = tolerations
        self.labels = labels or {}
        self.annotations = annotations or {}
        self.env_from = env_from or []
        self.volumes = volumes or []
        self.volume_mounts = volume_mounts or []
        self.image_pull_secrets = image_pull_secrets or []
        
        # Load Kubernetes configuration
        if config_file:
            config.load_kube_config(config_file=config_file, context=context)
        else:
            try:
                config.load_incluster_config()
                logger.info("Loaded in-cluster Kubernetes configuration")
            except config.ConfigException:
                config.load_kube_config(context=context)
                logger.info("Loaded Kubernetes configuration from default location")
        
        # Create Kubernetes API clients
        self.batch_v1_api = client.BatchV1Api()
        self.core_v1_api = client.CoreV1Api()
    
    def _create_job_object(
        self,
        name: str,
        command: List[str],
        args: List[str],
        env_vars: Optional[List[Dict[str, str]]] = None,
    ) -> client.V1Job:
        """
        Create a Kubernetes Job object.
        
        Args:
            name: Name of the job
            command: Command to run in the container
            args: Arguments for the command
            env_vars: Environment variables for the container
            
        Returns:
            Kubernetes Job object
        """
        # Set up environment variables
        env = []
        
        # Add API credentials if provided
        if self.api_host:
            env.append(client.V1EnvVar(name="NESSI_API_HOST", value=self.api_host))
        if self.api_key:
            env.append(client.V1EnvVar(name="NESSI_API_KEY", value=self.api_key))
        if self.api_secret:
            env.append(client.V1EnvVar(name="NESSI_API_SECRET", value=self.api_secret))
        
        # Add custom environment variables
        if env_vars:
            for env_var in env_vars:
                env.append(client.V1EnvVar(
                    name=env_var["name"],
                    value=env_var.get("value"),
                    value_from=env_var.get("value_from"),
                ))
        
        # Set up container
        container = client.V1Container(
            name="nessi",
            image=self.image,
            command=command,
            args=args,
            env=env,
            env_from=[client.V1EnvFromSource(**e) for e in self.env_from] if self.env_from else None,
            resources=client.V1ResourceRequirements(**self.resources) if self.resources else None,
            volume_mounts=[client.V1VolumeMount(**vm) for vm in self.volume_mounts] if self.volume_mounts else None,
        )
        
        # Set up pod template
        pod_template = client.V1PodTemplateSpec(
            metadata=client.V1ObjectMeta(
                labels={**self.labels, "app": "nessi", "job-name": name},
                annotations=self.annotations,
            ),
            spec=client.V1PodSpec(
                containers=[container],
                restart_policy="Never",
                service_account_name=self.service_account_name,
                node_selector=self.node_selector,
                tolerations=[client.V1Toleration(**t) for t in self.tolerations] if self.tolerations else None,
                volumes=[client.V1Volume(**v) for v in self.volumes] if self.volumes else None,
                image_pull_secrets=[client.V1LocalObjectReference(**s) for s in self.image_pull_secrets] if self.image_pull_secrets else None,
            ),
        )
        
        # Set up job
        job = client.V1Job(
            api_version="batch/v1",
            kind="Job",
            metadata=client.V1ObjectMeta(
                name=name,
                labels={**self.labels, "app": "nessi"},
                annotations=self.annotations,
            ),
            spec=client.V1JobSpec(
                template=pod_template,
                backoff_limit=self.job_backoff_limit,
                ttl_seconds_after_finished=self.job_ttl_seconds_after_finished,
                active_deadline_seconds=self.job_timeout_seconds,
            ),
        )
        
        return job
    
    def _create_job(self, job: client.V1Job) -> client.V1Job:
        """
        Create a Kubernetes Job.
        
        Args:
            job: Kubernetes Job object
            
        Returns:
            Created Job
        """
        try:
            return self.batch_v1_api.create_namespaced_job(
                namespace=self.namespace,
                body=job,
            )
        except ApiException as e:
            logger.error(f"Error creating job: {e}")
            raise
    
    def _delete_job(self, name: str) -> None:
        """
        Delete a Kubernetes Job.
        
        Args:
            name: Name of the job
        """
        try:
            self.batch_v1_api.delete_namespaced_job(
                name=name,
                namespace=self.namespace,
                body=client.V1DeleteOptions(
                    propagation_policy="Background",
                ),
            )
        except ApiException as e:
            logger.error(f"Error deleting job: {e}")
            raise
    
    def _get_job_status(self, name: str) -> Dict[str, Any]:
        """
        Get the status of a Kubernetes Job.
        
        Args:
            name: Name of the job
            
        Returns:
            Job status
        """
        try:
            job = self.batch_v1_api.read_namespaced_job_status(
                name=name,
                namespace=self.namespace,
            )
            
            status = {
                "name": job.metadata.name,
                "active": job.status.active or 0,
                "succeeded": job.status.succeeded or 0,
                "failed": job.status.failed or 0,
                "completion_time": job.status.completion_time,
                "start_time": job.status.start_time,
                "conditions": [],
            }
            
            if job.status.conditions:
                for condition in job.status.conditions:
                    status["conditions"].append({
                        "type": condition.type,
                        "status": condition.status,
                        "reason": condition.reason,
                        "message": condition.message,
                        "last_transition_time": condition.last_transition_time,
                    })
            
            return status
        except ApiException as e:
            logger.error(f"Error getting job status: {e}")
            raise
    
    def _wait_for_job_completion(
        self,
        name: str,
        timeout: int = 600,
        poll_interval: int = 10,
    ) -> Dict[str, Any]:
        """
        Wait for a Kubernetes Job to complete.
        
        Args:
            name: Name of the job
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds
            
        Returns:
            Job status
        """
        start_time = time.time()
        max_time = start_time + timeout
        
        while True:
            # Check timeout
            if time.time() > max_time:
                raise TimeoutError(f"Timeout reached while waiting for job completion: {name}")
            
            # Get job status
            status = self._get_job_status(name)
            
            # Check if job is complete
            if status["succeeded"] > 0:
                logger.info(f"Job completed successfully: {name}")
                return status
            elif status["failed"] > 0:
                logger.error(f"Job failed: {name}")
                # Get pod logs for debugging
                try:
                    pods = self.core_v1_api.list_namespaced_pod(
                        namespace=self.namespace,
                        label_selector=f"job-name={name}",
                    )
                    if pods.items:
                        pod_name = pods.items[0].metadata.name
                        logs = self.core_v1_api.read_namespaced_pod_log(
                            name=pod_name,
                            namespace=self.namespace,
                        )
                        logger.error(f"Pod logs for {pod_name}:\n{logs}")
                except ApiException as e:
                    logger.error(f"Error getting pod logs: {e}")
                
                raise RuntimeError(f"Job failed: {name}")
            
            # Wait for next poll
            logger.info(f"Job still running, waiting... {name}")
            time.sleep(poll_interval)
    
    def _get_job_logs(self, name: str) -> str:
        """
        Get logs from a Kubernetes Job.
        
        Args:
            name: Name of the job
            
        Returns:
            Job logs
        """
        try:
            pods = self.core_v1_api.list_namespaced_pod(
                namespace=self.namespace,
                label_selector=f"job-name={name}",
            )
            
            if not pods.items:
                return "No pods found for job"
            
            pod_name = pods.items[0].metadata.name
            logs = self.core_v1_api.read_namespaced_pod_log(
                name=pod_name,
                namespace=self.namespace,
            )
            
            return logs
        except ApiException as e:
            logger.error(f"Error getting job logs: {e}")
            return f"Error getting logs: {e}"


class NessiDataQualityOperator(NessiK8sOperator):
    """
    Kubernetes operator for running Nessi data quality checks.
    """
    
    def run_quality_check(
        self,
        table_name: str,
        rules: Optional[List[Dict[str, Any]]] = None,
        profile: bool = True,
        output_format: str = "json",
        output_path: Optional[str] = None,
        wait_for_completion: bool = True,
        timeout: int = 600,
        poll_interval: int = 10,
        job_name: Optional[str] = None,
        env_vars: Optional[List[Dict[str, str]]] = None,
    ) -> Dict[str, Any]:
        """
        Run a data quality check as a Kubernetes job.
        
        Args:
            table_name: Name of the table to check
            rules: List of quality rules to apply
            profile: Whether to generate a profile
            output_format: Output format (json, yaml, text)
            output_path: Path to write output to
            wait_for_completion: Whether to wait for job completion
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds
            job_name: Name for the job (generated if not provided)
            env_vars: Additional environment variables
            
        Returns:
            Job status or job name if not waiting for completion
        """
        # Generate job name if not provided
        if not job_name:
            job_name = f"nessi-quality-{str(uuid.uuid4())[:8]}"
        
        # Set up command and arguments
        command = ["nessi"]
        args = ["quality", "check", table_name]
        
        if rules:
            # Write rules to a file and mount it
            rules_file = f"/tmp/rules-{job_name}.json"
            with open(rules_file, "w") as f:
                json.dump(rules, f)
            
            # Add rules file to volumes
            self.volumes.append({
                "name": "rules",
                "host_path": {
                    "path": rules_file,
                },
            })
            self.volume_mounts.append({
                "name": "rules",
                "mount_path": "/rules.json",
                "sub_path": f"rules-{job_name}.json",
            })
            
            args.extend(["--rules", "/rules.json"])
        
        if profile:
            args.append("--profile")
        
        if output_format:
            args.extend(["--format", output_format])
        
        if output_path:
            args.extend(["--output", output_path])
        
        # Create and run the job
        job = self._create_job_object(
            name=job_name,
            command=command,
            args=args,
            env_vars=env_vars,
        )
        
        created_job = self._create_job(job)
        logger.info(f"Created job: {job_name}")
        
        if wait_for_completion:
            return self._wait_for_job_completion(
                name=job_name,
                timeout=timeout,
                poll_interval=poll_interval,
            )
        else:
            return {"name": job_name, "status": "created"}
    
    def get_quality_results(
        self,
        job_name: str,
        wait_for_completion: bool = True,
        timeout: int = 600,
        poll_interval: int = 10,
    ) -> Dict[str, Any]:
        """
        Get results of a data quality check job.
        
        Args:
            job_name: Name of the job
            wait_for_completion: Whether to wait for job completion
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds
            
        Returns:
            Job status and logs
        """
        if wait_for_completion:
            status = self._wait_for_job_completion(
                name=job_name,
                timeout=timeout,
                poll_interval=poll_interval,
            )
        else:
            status = self._get_job_status(job_name)
        
        logs = self._get_job_logs(job_name)
        
        return {
            "status": status,
            "logs": logs,
        }


class NessiProfileOperator(NessiK8sOperator):
    """
    Kubernetes operator for running Nessi data profiling.
    """
    
    def run_profile(
        self,
        table_name: str,
        columns: Optional[List[str]] = None,
        sample_size: Optional[int] = None,
        output_format: str = "json",
        output_path: Optional[str] = None,
        wait_for_completion: bool = True,
        timeout: int = 600,
        poll_interval: int = 10,
        job_name: Optional[str] = None,
        env_vars: Optional[List[Dict[str, str]]] = None,
    ) -> Dict[str, Any]:
        """
        Run a data profile as a Kubernetes job.
        
        Args:
            table_name: Name of the table to profile
            columns: List of columns to profile
            sample_size: Number of rows to sample
            output_format: Output format (json, yaml, text)
            output_path: Path to write output to
            wait_for_completion: Whether to wait for job completion
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds
            job_name: Name for the job (generated if not provided)
            env_vars: Additional environment variables
            
        Returns:
            Job status or job name if not waiting for completion
        """
        # Generate job name if not provided
        if not job_name:
            job_name = f"nessi-profile-{str(uuid.uuid4())[:8]}"
        
        # Set up command and arguments
        command = ["nessi"]
        args = ["profile", "run", table_name]
        
        if columns:
            args.extend(["--columns", ",".join(columns)])
        
        if sample_size:
            args.extend(["--sample", str(sample_size)])
        
        if output_format:
            args.extend(["--format", output_format])
        
        if output_path:
            args.extend(["--output", output_path])
        
        # Create and run the job
        job = self._create_job_object(
            name=job_name,
            command=command,
            args=args,
            env_vars=env_vars,
        )
        
        created_job = self._create_job(job)
        logger.info(f"Created job: {job_name}")
        
        if wait_for_completion:
            return self._wait_for_job_completion(
                name=job_name,
                timeout=timeout,
                poll_interval=poll_interval,
            )
        else:
            return {"name": job_name, "status": "created"}
    
    def get_profile_results(
        self,
        job_name: str,
        wait_for_completion: bool = True,
        timeout: int = 600,
        poll_interval: int = 10,
    ) -> Dict[str, Any]:
        """
        Get results of a data profile job.
        
        Args:
            job_name: Name of the job
            wait_for_completion: Whether to wait for job completion
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds
            
        Returns:
            Job status and logs
        """
        if wait_for_completion:
            status = self._wait_for_job_completion(
                name=job_name,
                timeout=timeout,
                poll_interval=poll_interval,
            )
        else:
            status = self._get_job_status(job_name)
        
        logs = self._get_job_logs(job_name)
        
        return {
            "status": status,
            "logs": logs,
        }


class NessiValidationOperator(NessiK8sOperator):
    """
    Kubernetes operator for running Nessi data validation.
    """
    
    def run_validation(
        self,
        table_name: str,
        rules: List[Dict[str, Any]],
        output_format: str = "json",
        output_path: Optional[str] = None,
        wait_for_completion: bool = True,
        timeout: int = 600,
        poll_interval: int = 10,
        job_name: Optional[str] = None,
        env_vars: Optional[List[Dict[str, str]]] = None,
    ) -> Dict[str, Any]:
        """
        Run a data validation as a Kubernetes job.
        
        Args:
            table_name: Name of the table to validate
            rules: List of validation rules
            output_format: Output format (json, yaml, text)
            output_path: Path to write output to
            wait_for_completion: Whether to wait for job completion
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds
            job_name: Name for the job (generated if not provided)
            env_vars: Additional environment variables
            
        Returns:
            Job status or job name if not waiting for completion
        """
        # Generate job name if not provided
        if not job_name:
            job_name = f"nessi-validation-{str(uuid.uuid4())[:8]}"
        
        # Set up command and arguments
        command = ["nessi"]
        args = ["validation", "run", table_name]
        
        # Write rules to a file and mount it
        rules_file = f"/tmp/rules-{job_name}.json"
        with open(rules_file, "w") as f:
            json.dump(rules, f)
        
        # Add rules file to volumes
        self.volumes.append({
            "name": "rules",
            "host_path": {
                "path": rules_file,
            },
        })
        self.volume_mounts.append({
            "name": "rules",
            "mount_path": "/rules.json",
            "sub_path": f"rules-{job_name}.json",
        })
        
        args.extend(["--rules", "/rules.json"])
        
        if output_format:
            args.extend(["--format", output_format])
        
        if output_path:
            args.extend(["--output", output_path])
        
        # Create and run the job
        job = self._create_job_object(
            name=job_name,
            command=command,
            args=args,
            env_vars=env_vars,
        )
        
        created_job = self._create_job(job)
        logger.info(f"Created job: {job_name}")
        
        if wait_for_completion:
            return self._wait_for_job_completion(
                name=job_name,
                timeout=timeout,
                poll_interval=poll_interval,
            )
        else:
            return {"name": job_name, "status": "created"}
    
    def get_validation_results(
        self,
        job_name: str,
        wait_for_completion: bool = True,
        timeout: int = 600,
        poll_interval: int = 10,
    ) -> Dict[str, Any]:
        """
        Get results of a data validation job.
        
        Args:
            job_name: Name of the job
            wait_for_completion: Whether to wait for job completion
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds
            
        Returns:
            Job status and logs
        """
        if wait_for_completion:
            status = self._wait_for_job_completion(
                name=job_name,
                timeout=timeout,
                poll_interval=poll_interval,
            )
        else:
            status = self._get_job_status(job_name)
        
        logs = self._get_job_logs(job_name)
        
        return {
            "status": status,
            "logs": logs,
        }
