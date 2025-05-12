document.addEventListener('DOMContentLoaded', function() {
    // Navigation
    const navLinks = document.querySelectorAll('nav a');
    const sections = document.querySelectorAll('main section');
    
    navLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            
            // Update active link
            navLinks.forEach(l => l.classList.remove('active'));
            this.classList.add('active');
            
            // Show active section
            const targetId = this.getAttribute('href').substring(1);
            sections.forEach(section => {
                section.classList.remove('active-section');
                if (section.id === targetId) {
                    section.classList.add('active-section');
                }
            });
        });
    });
    
    // Metrics Chart
    let metricsChart = null;
    const metricsChartCtx = document.getElementById('metrics-chart').getContext('2d');
    
    // Initialize metrics chart
    function initMetricsChart() {
        if (metricsChart) {
            metricsChart.destroy();
        }
        
        metricsChart = new Chart(metricsChartCtx, {
            type: 'line',
            data: {
                labels: [],
                datasets: [{
                    label: 'Metric Value',
                    data: [],
                    borderColor: '#3498db',
                    backgroundColor: 'rgba(52, 152, 219, 0.1)',
                    borderWidth: 2,
                    tension: 0.1,
                    fill: true
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: {
                        type: 'time',
                        time: {
                            unit: 'hour',
                            displayFormats: {
                                hour: 'MMM d, HH:mm'
                            }
                        },
                        title: {
                            display: true,
                            text: 'Time'
                        }
                    },
                    y: {
                        beginAtZero: true,
                        title: {
                            display: true,
                            text: 'Value'
                        }
                    }
                },
                plugins: {
                    tooltip: {
                        mode: 'index',
                        intersect: false
                    },
                    legend: {
                        position: 'top'
                    },
                    title: {
                        display: true,
                        text: 'Metric History'
                    }
                }
            }
        });
    }
    
    // Load metrics data
    async function loadMetricsData() {
        const metricName = document.getElementById('metric-selector').value;
        const timeRange = document.getElementById('time-range').value;
        
        // Calculate start time based on selected range
        const end = new Date();
        let start = new Date();
        
        switch (timeRange) {
            case '1h':
                start.setHours(end.getHours() - 1);
                break;
            case '6h':
                start.setHours(end.getHours() - 6);
                break;
            case '24h':
                start.setDate(end.getDate() - 1);
                break;
            case '7d':
                start.setDate(end.getDate() - 7);
                break;
            case '30d':
                start.setDate(end.getDate() - 30);
                break;
        }
        
        try {
            // Use fetchWithAuth for authenticated requests
            const response = await fetchWithAuth(`/api/metrics?metric=${metricName}&start=${start.toISOString()}&end=${end.toISOString()}`);
            
            if (!response.ok) {
                if (response.status === 401) {
                    // Unauthorized, redirect to login
                    window.location.href = '/login';
                    return;
                }
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const data = await response.json();
            updateMetricsChart(data, metricName);
            updateMetricsTable(data);
        } catch (error) {
            console.error('Error loading metrics data:', error);
            alert('Failed to load metrics data. See console for details.');
        }
    }
    
    // Load alerts data
    async function loadAlertsData() {
        const severity = document.getElementById('alert-severity').value;
        
        try {
            // Use fetchWithAuth for authenticated requests
            const response = await fetchWithAuth(`/api/alerts${severity !== 'all' ? `?severity=${severity}` : ''}`);
            
            if (!response.ok) {
                if (response.status === 401) {
                    // Unauthorized, redirect to login
                    window.location.href = '/login';
                    return;
                }
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const data = await response.json();
            updateAlertsTable(data, severity);
        } catch (error) {
            console.error('Error loading alerts data:', error);
        }
    }
    
    // Update metrics chart with new data
    function updateMetricsChart(data, metricName) {
        if (!data || data.length === 0) {
            // No data available
            if (metricsChart) {
                metricsChart.data.labels = [];
                metricsChart.data.datasets[0].data = [];
                metricsChart.update();
            } else {
                initMetricsChart();
            }
            return;
        }
        
        // Group data by label combinations
        const datasetsByLabels = {};
        
        data.forEach(item => {
            // Create a key from labels
            const labelKey = Object.entries(item.labels || {})
                .map(([k, v]) => `${k}=${v}`)
                .sort()
                .join(';');
            
            if (!datasetsByLabels[labelKey]) {
                datasetsByLabels[labelKey] = {
                    label: labelKey || metricName,
                    data: [],
                    borderColor: getRandomColor(),
                    backgroundColor: 'rgba(52, 152, 219, 0.1)',
                    borderWidth: 2,
                    tension: 0.1,
                    fill: false
                };
            }
            
            datasetsByLabels[labelKey].data.push({
                x: new Date(item.timestamp),
                y: item.value
            });
        });
        
        // Sort data points by timestamp
        Object.values(datasetsByLabels).forEach(dataset => {
            dataset.data.sort((a, b) => a.x - b.x);
        });
        
        // Create or update chart
        if (metricsChart) {
            metricsChart.data.datasets = Object.values(datasetsByLabels);
            metricsChart.options.plugins.title.text = `${metricName} History`;
            metricsChart.update();
        } else {
            initMetricsChart();
            metricsChart.data.datasets = Object.values(datasetsByLabels);
            metricsChart.options.plugins.title.text = `${metricName} History`;
            metricsChart.update();
        }
    }
    
    // Update metrics table with new data
    function updateMetricsTable(data) {
        const tableBody = document.querySelector('#metrics-table tbody');
        tableBody.innerHTML = '';
        
        if (!data || data.length === 0) {
            const row = document.createElement('tr');
            row.innerHTML = '<td colspan="3">No data available</td>';
            tableBody.appendChild(row);
            return;
        }
        
        // Sort data by timestamp (newest first)
        data.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));
        
        // Show only the latest 10 entries
        const latestData = data.slice(0, 10);
        
        latestData.forEach(item => {
            const row = document.createElement('tr');
            
            // Format timestamp
            const timestamp = new Date(item.timestamp);
            const formattedTime = timestamp.toLocaleString();
            
            // Format labels
            const labelsText = Object.entries(item.labels || {})
                .map(([k, v]) => `${k}: ${v}`)
                .join(', ');
            
            row.innerHTML = `
                <td>${formattedTime}</td>
                <td>${labelsText || '-'}</td>
                <td>${item.value.toFixed(2)}</td>
            `;
            
            tableBody.appendChild(row);
        });
    }
    
    // Update alerts table with new data
    function updateAlertsTable(data, severity) {
        const tableBody = document.querySelector('#alerts-table tbody');
        tableBody.innerHTML = '';
        
        if (!data || data.length === 0) {
            const row = document.createElement('tr');
            row.innerHTML = '<td colspan="5">No alerts available</td>';
            tableBody.appendChild(row);
            return;
        }
        
        // Filter by severity if not 'all'
        let filteredData = data;
        if (severity !== 'all') {
            filteredData = data.filter(alert => 
                alert.severity.toLowerCase() === severity.toLowerCase()
            );
        }
        
        // Sort by timestamp (newest first)
        filteredData.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));
        
        filteredData.forEach(alert => {
            const row = document.createElement('tr');
            
            // Format timestamp
            const timestamp = new Date(alert.timestamp);
            const formattedTime = timestamp.toLocaleString();
            
            // Format metadata
            const metadataText = Object.entries(alert.metadata || {})
                .map(([k, v]) => `${k}: ${v}`)
                .join(', ');
            
            row.innerHTML = `
                <td>${formattedTime}</td>
                <td>${alert.name}</td>
                <td class="severity-${alert.severity.toLowerCase()}">${alert.severity}</td>
                <td>${alert.message}</td>
                <td>${metadataText || '-'}</td>
            `;
            
            tableBody.appendChild(row);
        });
    }
    
    // Settings form submission
    document.getElementById('settings-form').addEventListener('submit', function(e) {
        e.preventDefault();
        
        const settings = {
            retentionPeriod: parseInt(document.getElementById('retention-period').value, 10),
            snapshotInterval: parseInt(document.getElementById('snapshot-interval').value, 10),
            alertCooldown: parseInt(document.getElementById('alert-cooldown').value, 10)
        };
        
        // TODO: Implement settings update via API
        console.log('Settings updated:', settings);
        alert('Settings updated successfully!');
    });
    
    // Export metrics function
    async function exportMetrics() {
        const metricName = document.getElementById('metric-selector').value;
        const timeRange = document.getElementById('time-range').value;
        const format = document.getElementById('export-format').value;
        
        // Calculate start time based on selected range
        const end = new Date();
        let start = new Date();
        
        switch (timeRange) {
            case '1h':
                start.setHours(end.getHours() - 1);
                break;
            case '6h':
                start.setHours(end.getHours() - 6);
                break;
            case '24h':
                start.setDate(end.getDate() - 1);
                break;
            case '7d':
                start.setDate(end.getDate() - 7);
                break;
            case '30d':
                start.setDate(end.getDate() - 30);
                break;
        }
        
        // Create export URL
        const exportUrl = `/api/export?metric=${metricName}&format=${format}&start=${start.toISOString()}&end=${end.toISOString()}`;
        
        // Check if authenticated
        if (isAuthenticated()) {
            // Create a link with authentication token
            const a = document.createElement('a');
            a.href = exportUrl;
            a.target = '_blank';
            a.style.display = 'none';
            document.body.appendChild(a);
            
            // Set up authentication for the download
            const authHeader = `Bearer ${getAuthToken()}`;
            
            // Use fetch with authentication to download the file
            try {
                const response = await fetchWithAuth(exportUrl);
                
                if (!response.ok) {
                    if (response.status === 401) {
                        // Unauthorized, redirect to login
                        window.location.href = '/login';
                        return;
                    }
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                
                // Get the blob from the response
                const blob = await response.blob();
                
                // Create a download link
                const url = window.URL.createObjectURL(blob);
                a.href = url;
                a.download = `${metricName}_export.${format}`;
                a.click();
                
                // Clean up
                window.URL.revokeObjectURL(url);
                document.body.removeChild(a);
            } catch (error) {
                console.error('Error exporting metrics:', error);
                alert('Failed to export metrics. Please try again.');
            }
        } else {
            // Not authenticated, redirect to login
            window.location.href = '/login';
        }
    }
    
    // Refresh and export buttons
    document.getElementById('refresh-btn').addEventListener('click', loadMetricsData);
    document.getElementById('refresh-alerts-btn').addEventListener('click', loadAlertsData);
    document.getElementById('export-btn').addEventListener('click', exportMetrics);
    
    // Metric selector change
    document.getElementById('metric-selector').addEventListener('change', loadMetricsData);
    document.getElementById('time-range').addEventListener('change', loadMetricsData);
    
    // Alert severity filter change
    document.getElementById('alert-severity').addEventListener('change', loadAlertsData);
    
    // Helper function to generate random colors
    function getRandomColor() {
        const colors = [
            '#3498db', // blue
            '#e74c3c', // red
            '#2ecc71', // green
            '#f39c12', // orange
            '#9b59b6', // purple
            '#1abc9c', // teal
            '#34495e', // dark blue
            '#e67e22', // dark orange
            '#27ae60', // dark green
            '#d35400'  // dark orange
        ];
        
        return colors[Math.floor(Math.random() * colors.length)];
    }
    
    // Initial load
    initMetricsChart();
    loadMetricsData();
    loadAlertsData();
});
