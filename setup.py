from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="nessi",
    version="1.0.0",
    author="Nessi Data Tool",
    author_email="nessi.datatool@gmail.com",
    description="A powerful data processing and analysis tool",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/Nessi-datatool/nessi",
    packages=find_packages(),
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "Intended Audience :: Science/Research",
        "License :: Other/Proprietary License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
    ],
    python_requires=">=3.8",
    install_requires=[
        "pandas>=2.0.0",
        "numpy>=1.24.0",
        "pyarrow>=14.0.0",
        "delta-spark>=3.0.0",
        "click>=8.0.0",
        "pathlib>=1.0.1",
        "cryptography>=41.0.0",
    ],
    entry_points={
        "console_scripts": [
            "nessi=src.cli:cli",
        ],
    },
) 