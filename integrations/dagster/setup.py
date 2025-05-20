from setuptools import setup, find_packages

setup(
    name="nessi-dagster",
    version="0.1.0",
    description="Dagster integration for Nessi.dev",
    author="Nessi.dev Team",
    author_email="info@nessi.dev",
    packages=find_packages(),
    install_requires=[
        "dagster>=1.0.0",
        "dagster-graphql>=1.0.0",
        "nessi-client>=0.1.0",
        "requests>=2.25.0",
        "pyyaml>=5.4.0",
    ],
    python_requires=">=3.7",
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
    ],
)
