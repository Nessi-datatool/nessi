from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="nessi-client",
    version="0.1.0",
    author="Nessi Team",
    author_email="info@nessi-dev.com",
    description="Python client for the Nessi monitoring system",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/nessi-dev/nessi-dev",
    packages=find_packages(),
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.7",
    install_requires=[
        "requests>=2.25.0",
        "dataclasses;python_version<'3.7'"
    ],
)
