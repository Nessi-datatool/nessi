#!/usr/bin/env python3
"""
read_parquet.py - Reads a Parquet file and outputs Arrow IPC stream to stdout.

Usage:
    python3 read_parquet.py <file_path>

This script is intended to be called from Go for robust Parquet reading,
aligning with the modular Go-Python bridge described in IMPLEMENTATION.md.
"""
import sys
import pyarrow.parquet as pq
import pyarrow as pa

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 read_parquet.py <file_path>", file=sys.stderr)
        sys.exit(1)
    file_path = sys.argv[1]
    try:
        table = pq.read_table(file_path)
    except Exception as e:
        print(f"Failed to read Parquet: {e}", file=sys.stderr)
        sys.exit(2)
    # Output as Arrow IPC stream to stdout
    with pa.ipc.new_stream(sys.stdout.buffer, table.schema) as writer:
        writer.write_table(table)

if __name__ == "__main__":
    main()
