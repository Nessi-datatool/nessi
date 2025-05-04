from setuptools import setup, find_packages

setup(
    name="nessi",
    version="0.1.0",
    description="NESSI - Network and System Scanner for Infrastructure",
    author="NESSI Team",
    author_email="nessi@example.com",
    packages=find_packages(where="src"),
    package_dir={"": "src"},
    install_requires=[
        "pyspark>=3.4.1",
        "delta-spark>=2.4.0",
        "pandas>=2.0.3",
        "numpy>=1.26.2",
        "pyarrow>=14.0.1",
    ],
    extras_require={
        "test": [
            "pytest>=7.4.3",
            "pytest-cov>=4.1.0",
            "pytest-html>=4.1.1",
            "pytest-benchmark>=4.0.0",
            "pytest-mock>=3.12.0",
        ],
        "monitoring": [
            "prometheus-client>=0.19.0",
            "grafana-api>=1.0.3",
        ],
    },
    python_requires=">=3.9",
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: Apache Software License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.9",
        "Topic :: Software Development :: Libraries :: Python Modules",
    ],
) 