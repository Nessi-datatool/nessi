#!/usr/bin/env python3
"""
generate_sample_parquet.py - Generates a small Parquet file for testing Go↔Python integration.

Schema: id (int32), name (string), value (float64)
"""
import pandas as pd
import pyarrow as pa
import pyarrow.parquet as pq
import argparse
import random
import string

NAMES = ["Alice", "Bob", "Charlie", "Diana", "Eve", "Frank", "Grace", "Heidi", "Ivan", "Judy"]

def random_string(length=8):
    return ''.join(random.choices(string.ascii_letters, k=length))

def make_default_schema(rows):
    now = pd.Timestamp.now()
    return pd.DataFrame({
        "id": list(range(1, rows+1)),
        "name": [random.choice(NAMES) for _ in range(rows)],
        "value": [round(random.uniform(0, 100), 2) for _ in range(rows)],
        "timestamp": [now + pd.Timedelta(seconds=random.randint(-100000, 100000)) for _ in range(rows)],
        "decimal": [round(random.uniform(0, 100), 4) for _ in range(rows)],
    })

def make_wide_schema(rows):
    now = pd.Timestamp.now()
    # Nested struct column for person, and list[str] for tags
    df = pd.DataFrame({
        "id": list(range(1, rows+1)),
        "name": [random.choice(NAMES) for _ in range(rows)],
        "value": [round(random.uniform(0, 100), 2) for _ in range(rows)],
        "score": [random.uniform(0, 1) for _ in range(rows)],
        "flag": [random.choice([True, False]) for _ in range(rows)],
        "description": [random_string(16) for _ in range(rows)],
        "timestamp": [now + pd.Timedelta(seconds=random.randint(-100000, 100000)) for _ in range(rows)],
        "decimal": [round(random.uniform(0, 100), 4) for _ in range(rows)],
    })
    # Add nested struct column (person)
    df["person"] = [{"name": random.choice(NAMES), "age": random.randint(18, 80)} for _ in range(rows)]
    # Add list<string> column (tags)
    tag_pool = ["alpha", "beta", "gamma", "delta", "omega", "test", "prod"]
    df["tags"] = [[random.choice(tag_pool) for _ in range(random.randint(1, 4))] for _ in range(rows)]
    return df

def main():
    parser = argparse.ArgumentParser(description="Generate a sample Parquet file for Go↔Python integration testing.")
    parser.add_argument("--output", required=True, help="Output Parquet file path.")
    parser.add_argument("--rows", type=int, default=10, help="Number of rows to generate.")
    parser.add_argument("--schema", choices=["default", "wide"], default="default", help="Schema type.")
    parser.add_argument("--all-null", action="store_true", help="Generate all-null columns.")
    parser.add_argument("--empty", action="store_true", help="Generate an empty Parquet file.")
    parser.add_argument("--large", action="store_true", help="Generate a large Parquet file (100,000 rows).")
    args = parser.parse_args()

    if args.empty:
        df = pd.DataFrame()
        schema = pa.schema([])
    else:
        row_count = 100000 if args.large else args.rows
        if args.schema == 'default':
            df = make_default_schema(row_count)
            schema = pa.schema([
                ("id", pa.int32()),
                ("name", pa.string()),
                ("value", pa.float64()),
                ("timestamp", pa.timestamp('ns')),
                ("decimal", pa.decimal128(10, 4)),
            ])
        elif args.schema == 'wide':
            df = make_wide_schema(row_count)
            schema = pa.schema([
                ("id", pa.int32()),
                ("name", pa.string()),
                ("value", pa.float64()),
                ("score", pa.float32()),
                ("flag", pa.bool_()),
                ("description", pa.string()),
                ("timestamp", pa.timestamp('ns')),
                ("decimal", pa.decimal128(10, 4)),
                ("person", pa.struct([('name', pa.string()), ('age', pa.int32())])),
                ("tags", pa.list_(pa.string())),
            ])
        else:
            raise ValueError(f"Unsupported schema: {args.schema}")

        if args.all_null:
            for col in df.columns:
                df[col] = None

    table = pa.Table.from_pandas(df, schema=schema, preserve_index=False)
    pq.write_table(table, args.output)
    print(f"Sample Parquet file '{args.output}' created with {len(df)} rows and schema '{args.schema}'.")

if __name__ == "__main__":
    main()
