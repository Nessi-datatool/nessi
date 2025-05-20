#!/usr/bin/env python
"""
ML-based anomaly detection for Nessi.dev

This script provides anomaly detection capabilities using machine learning
algorithms from scikit-learn. It can detect outliers in time series data
using various algorithms including Isolation Forest, One-Class SVM, and
Local Outlier Factor.
"""

import json
import sys
import numpy as np
import pandas as pd
from sklearn.ensemble import IsolationForest
from sklearn.svm import OneClassSVM
from sklearn.neighbors import LocalOutlierFactor
import matplotlib.pyplot as plt
from io import BytesIO
import base64

def detect_anomalies(data, algorithm="isolation_forest", contamination=0.05):
    """
    Detect anomalies in the given data using the specified algorithm.
    
    Args:
        data: List or numpy array of values
        algorithm: One of "isolation_forest", "one_class_svm", or "local_outlier_factor"
        contamination: Expected proportion of outliers in the data
        
    Returns:
        Dictionary with anomaly indices, scores, and visualization
    """
    X = np.array(data).reshape(-1, 1)
    
    if algorithm == "isolation_forest":
        model = IsolationForest(contamination=contamination, random_state=42)
        scores = model.fit_predict(X)
    elif algorithm == "one_class_svm":
        model = OneClassSVM(nu=contamination)
        scores = model.fit_predict(X)
    elif algorithm == "local_outlier_factor":
        model = LocalOutlierFactor(n_neighbors=20, contamination=contamination)
        scores = model.fit_predict(X)
    else:
        raise ValueError(f"Unknown algorithm: {algorithm}")
    
    # Convert scores to anomaly flags (-1 for anomalies, 1 for normal)
    anomaly_indices = np.where(scores == -1)[0].tolist()
    
    # Generate visualization
    visualization = generate_visualization(data, anomaly_indices)
    
    return {
        "anomaly_indices": anomaly_indices,
        "anomaly_values": [data[i] for i in anomaly_indices],
        "visualization": visualization
    }

def generate_visualization(data, anomaly_indices):
    """Generate a visualization of the data with anomalies highlighted"""
    plt.figure(figsize=(10, 6))
    plt.plot(data, 'b-', label='Data')
    plt.scatter(anomaly_indices, [data[i] for i in anomaly_indices], 
                color='red', label='Anomalies')
    plt.title('Anomaly Detection Results')
    plt.xlabel('Index')
    plt.ylabel('Value')
    plt.legend()
    plt.grid(True)
    
    # Save plot to a base64 encoded string
    buffer = BytesIO()
    plt.savefig(buffer, format='png')
    buffer.seek(0)
    image_png = buffer.getvalue()
    buffer.close()
    plt.close()
    
    return base64.b64encode(image_png).decode('utf-8')

def seasonal_decomposition(data, period=7):
    """
    Decompose time series into trend, seasonal, and residual components.
    
    Args:
        data: List or numpy array of values
        period: The period of the seasonal component
        
    Returns:
        Dictionary with trend, seasonal, and residual components
    """
    try:
        from statsmodels.tsa.seasonal import seasonal_decompose
        
        # Convert to pandas Series
        series = pd.Series(data)
        
        # Perform decomposition
        result = seasonal_decompose(series, model='additive', period=period)
        
        # Generate visualizations
        trend_viz = generate_component_visualization(result.trend, "Trend Component")
        seasonal_viz = generate_component_visualization(result.seasonal, "Seasonal Component")
        residual_viz = generate_component_visualization(result.resid, "Residual Component")
        
        return {
            "trend": result.trend.dropna().tolist(),
            "seasonal": result.seasonal.dropna().tolist(),
            "residual": result.resid.dropna().tolist(),
            "visualizations": {
                "trend": trend_viz,
                "seasonal": seasonal_viz,
                "residual": residual_viz
            }
        }
    except ImportError:
        return {
            "error": "statsmodels package is required for seasonal decomposition"
        }

def generate_component_visualization(component, title):
    """Generate visualization for a time series component"""
    plt.figure(figsize=(10, 4))
    plt.plot(component)
    plt.title(title)
    plt.grid(True)
    
    # Save plot to a base64 encoded string
    buffer = BytesIO()
    plt.savefig(buffer, format='png')
    buffer.seek(0)
    image_png = buffer.getvalue()
    buffer.close()
    plt.close()
    
    return base64.b64encode(image_png).decode('utf-8')

def main():
    """Main function to process command line arguments"""
    if len(sys.argv) < 2:
        print(json.dumps({"error": "No arguments provided"}))
        return
    
    try:
        args = json.loads(sys.argv[1])
        
        if "action" not in args:
            print(json.dumps({"error": "No action specified"}))
            return
        
        if args["action"] == "detect_anomalies":
            if "data" not in args:
                print(json.dumps({"error": "No data provided"}))
                return
            
            algorithm = args.get("algorithm", "isolation_forest")
            contamination = float(args.get("contamination", 0.05))
            
            result = detect_anomalies(args["data"], algorithm, contamination)
            print(json.dumps(result))
            
        elif args["action"] == "seasonal_decomposition":
            if "data" not in args:
                print(json.dumps({"error": "No data provided"}))
                return
            
            period = int(args.get("period", 7))
            
            result = seasonal_decomposition(args["data"], period)
            print(json.dumps(result))
            
        else:
            print(json.dumps({"error": f"Unknown action: {args['action']}"}))
            
    except Exception as e:
        print(json.dumps({"error": str(e)}))

if __name__ == "__main__":
    main()
