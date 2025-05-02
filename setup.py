from setuptools import setup, find_packages

# Read README.md if it exists
try:
    with open("README.md", "r", encoding="utf-8") as f:
        long_description = f.read()
except FileNotFoundError:
    long_description = "A Python package for creating and managing Delta tables with Spark."

setup(
    name="nessi",
    version="0.1.0",
    packages=find_packages(where="src"),
    install_requires=[
        "numpy<2.0.0",
        "pandas==2.1.4",
        "pyarrow==14.0.1",
        "pyspark==3.5.0",
        "delta-spark==3.0.0",
        "pytest==8.3.5",
        "pytest-cov==6.1.1",
        "jinja2>=3.0.0",
        "fastparquet>=2023.1.0",
        "cryptography==41.0.5",
        "python-dateutil==2.8.2",
        "pytz==2023.3",
        "typing-extensions==4.8.0",
        "flake8==6.1.0",
        "click==8.1.7",
        "pathlib==1.0.1",
    ],
    python_requires=">=3.11",
    package_dir={"": "src"},
    package_data={
        "nessi": ["*.py"],
    },
    author="Nessi Team",
    author_email="nessi@example.com",
    description="A Python package for creating and managing Delta tables with Spark",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/nessi/nessi",
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3.11",
        "Operating System :: OS Independent",
    ],
) 