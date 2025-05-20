"""
Nessi.dev flows for Prefect.
"""

from nessi_prefect.flows.data_quality_flow import data_quality_flow
from nessi_prefect.flows.data_pipeline_flow import data_pipeline_flow

__all__ = [
    'data_quality_flow',
    'data_pipeline_flow',
]
