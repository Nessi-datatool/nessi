from setuptools import setup, find_packages

setup(
    name="nessi-scanner",
    version="1.0.0",
    description="Data quality scanner and report generator for Nessi.dev",
    author="Nessi.dev Team",
    author_email="nessi-datatool@gmail.com",
    packages=find_packages(where="src"),
    package_dir={"": "src"},
    install_requires=[
        "delta-spark>=3.3.1",
        "pyspark>=3.5.5",
        "pandas>=2.0.0",
        "numpy>=1.24.0",
        "pyarrow>=14.0.0",
        "fastparquet>=2023.10.1",
        "plotly>=5.18.0",
        "jinja2>=3.1.0",
        "reportlab>=4.0.0",
        "weasyprint>=60.1",
        "python-docx>=1.0.0",
        "pydantic>=2.11.4",
        "pydantic_core>=2.34.1",
        "typing-extensions>=4.13.2",
        "annotated-types>=0.7.0"
    ],
    extras_require={
        "dev": [
            "pytest>=7.4.0",
            "pytest-cov>=4.1.0",
            "black>=23.7.0",
            "flake8>=6.1.0",
            "mypy>=1.5.0"
        ]
    },
    python_requires=">=3.8",
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
        "Programming Language :: Python :: 3.13",
        "Topic :: Software Development :: Libraries :: Python Modules",
        "Topic :: Database :: Database Engines/Servers",
        "Topic :: Scientific/Engineering :: Information Analysis"
    ],
    entry_points={
        "console_scripts": [
            "nessi-scan=scanner.cli:main",
        ],
    },
) 