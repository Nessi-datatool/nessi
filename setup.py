from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="nessi",
    version="0.1.0",
    author="Nessi Team",
    author_email="nessi@example.com",
    description="A Python-based data processing and analysis tool built with PySpark and Delta Lake",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/Nessi-datatool/nessi",
    packages=find_packages(),
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "Intended Audience :: Science/Research",
        "License :: OSI Approved :: Apache Software License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Topic :: Scientific/Engineering :: Information Analysis",
        "Topic :: Software Development :: Libraries :: Python Modules",
    ],
    install_requires=[
        "pyspark>=3.5.0",
        "delta-spark>=2.4.0",
        "pandas>=2.0.0",
        "pytest>=7.0.0",
    ],
    python_requires=">=3.8",
    project_urls={
        "Bug Reports": "https://github.com/Nessi-datatool/nessi/issues",
        "Source": "https://github.com/Nessi-datatool/nessi",
    },
) 