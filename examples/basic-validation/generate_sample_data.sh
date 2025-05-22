#!/bin/bash
# Script to generate sample data for Nessi validation example

echo "Generating sample Delta Lake table..."

# Create directories
mkdir -p ./data/sample_table

# Check if Python is available
if ! command -v python3 &> /dev/null; then
    echo "Error: Python 3 is required but not found."
    exit 1
fi

# Create Python script to generate sample data
cat > ./generate_data.py << 'EOF'
import os
import random
import datetime
import pandas as pd
from faker import Faker
from delta import configure_spark_with_delta_pip
from pyspark.sql import SparkSession

# Initialize Faker
fake = Faker()

# Create sample data
def generate_sample_data(num_rows=1000):
    data = []
    
    # Generate current date for reference
    now = datetime.datetime.now()
    
    for i in range(num_rows):
        # Introduce some nulls and invalid data for demonstration
        user_id = None if random.random() < 0.05 else f"user_{i}"
        name = fake.name()
        
        # Age with some out of range values
        age = random.randint(10, 130) if random.random() < 0.1 else random.randint(18, 80)
        
        # Email with some invalid formats
        email = f"invalid-email-{i}" if random.random() < 0.1 else fake.email()
        
        # Dates
        signup_date = fake.date_between(start_date="-5y", end_date="today")
        
        # Some very old last logins
        if random.random() < 0.1:
            last_login = fake.date_time_between(start_date="-5y", end_date="-2y")
        else:
            last_login = fake.date_time_between(start_date="-1y", end_date="now")
        
        # Active status
        active = random.choice([True, False])
        
        # Score with some out of range values
        score = random.uniform(-10, 110) if random.random() < 0.1 else random.uniform(0, 100)
        
        data.append({
            "user_id": user_id,
            "name": name,
            "age": age,
            "email": email,
            "signup_date": signup_date,
            "last_login": last_login,
            "active": active,
            "score": score
        })
    
    return pd.DataFrame(data)

# Main function
def main():
    print("Creating sample data...")
    df = generate_sample_data()
    
    print("Initializing Spark with Delta Lake...")
    # Configure Spark with Delta Lake
    builder = SparkSession.builder.appName("NessiSampleData") \
        .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
        .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog")
    
    spark = configure_spark_with_delta_pip(builder).getOrCreate()
    
    # Convert pandas DataFrame to Spark DataFrame
    spark_df = spark.createDataFrame(df)
    
    # Write to Delta format
    print("Writing data to Delta format...")
    table_path = os.path.abspath("./data/sample_table")
    spark_df.write.format("delta").mode("overwrite").save(table_path)
    
    print(f"Sample Delta table created at: {table_path}")
    print("Schema:")
    spark_df.printSchema()
    
    # Clean up Spark session
    spark.stop()

if __name__ == "__main__":
    main()
EOF

# Create requirements file
cat > ./requirements.txt << 'EOF'
delta-spark>=2.0.0
pyspark>=3.1.2
pandas>=1.3.0
faker>=8.0.0
EOF

# Install dependencies
echo "Installing required Python packages..."
pip3 install -r ./requirements.txt

# Run the Python script
echo "Running data generation script..."
python3 ./generate_data.py

# Create rules file
cat > ./rules.yaml << 'EOF'
rules:
  - name: "user_id_not_null"
    description: "User ID should not be null"
    column: "user_id"
    type: "not_null"
    
  - name: "valid_email"
    description: "Email should be in valid format"
    column: "email"
    type: "regex"
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
    
  - name: "age_range"
    description: "Age should be between 18 and 120"
    column: "age"
    type: "range"
    min: 18
    max: 120
    
  - name: "score_range"
    description: "Score should be between 0 and 100"
    column: "score"
    type: "range"
    min: 0
    max: 100
    
  - name: "recent_login"
    description: "Last login should be within the last year"
    column: "last_login"
    type: "recency"
    max_age: "365d"
EOF

# Create reports directory
mkdir -p ./reports

echo "Setup complete! You can now run:"
echo "  nessi validate ./data/sample_table --rules ./rules.yaml"
echo "  nessi report ./data/sample_table --rules ./rules.yaml --format html --output ./reports"
