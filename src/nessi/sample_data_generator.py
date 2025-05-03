"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 nessi.dev. All rights reserved.

nessi.dev is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of nessi.dev will require a license for all users.

1. PERSONAL USE LICENSE
   a) The Software is provided free of charge for personal use
   b) Personal use includes:
      - Individual data analysis and processing
      - Educational purposes
      - Non-commercial research
   c) Enterprise or consulting use requires a paid license

2. RESTRICTIONS
   You shall not:
   a) Copy, modify, adapt, translate, reverse engineer, decompile, or disassemble the Software
   b) Create derivative works based on the Software
   c) Rent, lease, loan, sell, sublicense, distribute, transmit, or otherwise transfer the Software
   d) Remove or alter any proprietary notices or labels on the Software
   e) Use the Software for enterprise or consulting purposes without a valid paid license

3. OWNERSHIP
   The Software is licensed, not sold. nessi.dev retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

"""Sample data generator module.

This module provides functionality to generate and save sample data
for testing purposes.
"""

import os
import logging
import pandas as pd
import numpy as np
from datetime import date, timedelta
import random
import argparse
from typing import Dict, Any, List, Union

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

def generate_sample_data(num_rows: int = 100) -> pd.DataFrame:
    """Generate sample data for testing.
    
    Args:
        num_rows (int): Number of rows to generate.
        
    Returns:
        pd.DataFrame: Generated sample data.
    """
    try:
        # Generate data
        data = {
            'id': list(range(1, num_rows + 1)),
            'name': [f"User_{i}" for i in range(1, num_rows + 1)],
            'age': [random.randint(20, 60) for _ in range(num_rows)],
            'salary': [random.randint(30000, 100000) for _ in range(num_rows)],
            'department': [random.choice(['HR', 'Engineering', 'Finance', 'Sales']) for _ in range(num_rows)],
            'join_date': [date(2020, 1, 1) + timedelta(days=random.randint(0, 1000)) for _ in range(num_rows)]
        }
        
        # Convert to DataFrame
        df = pd.DataFrame(data)
        
        # Add some null values (10% chance for each field except id)
        for col in ['name', 'age', 'salary', 'department', 'join_date']:
            mask = np.random.random(num_rows) < 0.1
            df.loc[mask, col] = None
        
        return df
    except Exception as e:
        logger.error(f"Failed to generate sample data: {str(e)}")
        raise RuntimeError(f"Failed to generate sample data: {str(e)}")

def save_as_parquet(data: pd.DataFrame, table_path: str) -> None:
    """Save data as Parquet.
    
    Args:
        data (pd.DataFrame): Data to save.
        table_path (str): Path where to save the data.
        
    Raises:
        RuntimeError: If the data cannot be saved.
    """
    try:
        # Validate inputs
        if not isinstance(data, pd.DataFrame):
            raise ValueError("Data must be a pandas DataFrame")
        if not table_path:
            raise ValueError("Table path cannot be empty")
        
        # Create directory if it doesn't exist
        directory = os.path.dirname(table_path)
        if directory:
            os.makedirs(directory, exist_ok=True)
        
        # Save data
        if os.path.isdir(table_path):
            # If table_path is a directory, save as data.parquet inside it
            file_path = os.path.join(table_path, "data.parquet")
        else:
            # If table_path doesn't end with .parquet, append it
            file_path = table_path if table_path.lower().endswith('.parquet') else table_path + '.parquet'
        
        data.to_parquet(file_path)
    except Exception as e:
        logger.error(f"Failed to save data as Parquet: {str(e)}")
        raise RuntimeError(f"Failed to save data as Parquet: {str(e)}")

def save_as_csv(data: pd.DataFrame, file_path: str) -> None:
    """Save data as CSV.
    
    Args:
        data (pd.DataFrame): Data to save.
        file_path (str): Path where to save the data.
        
    Raises:
        RuntimeError: If the data cannot be saved.
    """
    try:
        # Validate inputs
        if not isinstance(data, pd.DataFrame):
            raise ValueError("Data must be a pandas DataFrame")
        if not file_path:
            raise ValueError("File path cannot be empty")
        
        # Create directory if it doesn't exist
        os.makedirs(os.path.dirname(file_path), exist_ok=True)
        
        # Save data
        data.to_csv(file_path, index=False)
    except Exception as e:
        logger.error(f"Failed to save data as CSV: {str(e)}")
        raise RuntimeError(f"Failed to save data as CSV: {str(e)}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Generate sample data for testing.")
    parser.add_argument("--num-rows", type=int, default=100, help="Number of rows to generate")
    parser.add_argument("--output-path", type=str, required=True, help="Path where to save the data")
    parser.add_argument("--format", type=str, choices=['parquet', 'csv'], default='parquet', help="Output format")
    args = parser.parse_args()
    
    try:
        # Generate data
        data = generate_sample_data(args.num_rows)
        
        # Save data
        if args.format == 'parquet':
            save_as_parquet(data, args.output_path)
        else:
            save_as_csv(data, args.output_path)
        
        logger.info(f"Successfully generated and saved {args.num_rows} rows of sample data to {args.output_path}")
    except Exception as e:
        logger.error(f"Failed to generate and save sample data: {str(e)}")
        raise RuntimeError(f"Failed to generate and save sample data: {str(e)}")