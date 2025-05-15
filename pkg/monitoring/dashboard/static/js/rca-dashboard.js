/**
 * Root Cause Analysis Dashboard JavaScript
 * This file handles the RCA dashboard functionality
 */

// Initialize dashboard when document is ready
$(document).ready(function() {
    // Initialize Feather icons
    feather.replace();
    
    // Load initial data
    loadRecentRcaData();
    loadAnomalyIds();
    loadRcaInsights();
    
    // Set up event handlers
    setupEventHandlers();
});

/**
 * Set up event handlers for the RCA dashboard
 */
function setupEventHandlers() {
    // Refresh button
    $('#refreshRcaBtn').click(function() {
        loadRecentRcaData();
        loadRcaInsights();
    });
    
    // Export button
    $('#exportRcaBtn').click(function() {
        exportRcaData();
    });
    
    // RCA form submission
    $('#rcaForm').submit(function(e) {
        e.preventDefault();
        runRootCauseAnalysis();
    });
    
    // Export detail button
    $('#exportDetailBtn').click(function() {
        exportRcaDetail();
    });
}

/**
 * Load recent RCA data from the API
 */
function loadRecentRcaData() {
    // Show loading indicator
    $('#rcaTableBody').html('<tr><td colspan="6" class="text-center">Loading...</td></tr>');
    
    // Fetch RCA data from API
    $.ajax({
        url: '/api/rca/recent',
        method: 'GET',
        success: function(data) {
            displayRecentRcaData(data);
        },
        error: function(xhr, status, error) {
            $('#rcaTableBody').html('<tr><td colspan="6" class="text-center text-danger">Error loading RCA data: ' + error + '</td></tr>');
            console.error('Error loading RCA data:', error);
        }
    });
}

/**
 * Display recent RCA data in the table
 */
function displayRecentRcaData(data) {
    if (!data || data.length === 0) {
        $('#rcaTableBody').html('<tr><td colspan="6" class="text-center">No recent analyses found</td></tr>');
        return;
    }
    
    let html = '';
    data.forEach(function(rca) {
        html += `
            <tr>
                <td>${rca.anomaly_id}</td>
                <td>${formatDate(rca.analysis_time)}</td>
                <td>${rca.primary_root_cause ? rca.primary_root_cause.description : 'N/A'}</td>
                <td>${rca.primary_root_cause ? (rca.primary_root_cause.confidence * 100).toFixed(1) + '%' : 'N/A'}</td>
                <td>${rca.affected_tables ? rca.affected_tables.join(', ') : 'N/A'}</td>
                <td>
                    <button class="btn btn-sm btn-info view-rca-btn" data-id="${rca.anomaly_id}">View</button>
                </td>
            </tr>
        `;
    });
    
    $('#rcaTableBody').html(html);
    
    // Add event listener for view buttons
    $('.view-rca-btn').click(function() {
        const anomalyId = $(this).data('id');
        viewRcaDetail(anomalyId);
    });
}

/**
 * Load anomaly IDs for the dropdown
 */
function loadAnomalyIds() {
    $.ajax({
        url: '/api/anomalies',
        method: 'GET',
        success: function(data) {
            const dropdown = $('#anomalyId');
            dropdown.empty();
            
            if (!data || data.length === 0) {
                dropdown.append('<option value="" disabled selected>No anomalies found</option>');
                return;
            }
            
            dropdown.append('<option value="" disabled selected>Select an anomaly</option>');
            data.forEach(function(anomaly) {
                dropdown.append(`<option value="${anomaly.id}">${anomaly.id} - ${anomaly.description || 'No description'}</option>`);
            });
        },
        error: function(xhr, status, error) {
            $('#anomalyId').html('<option value="" disabled selected>Error loading anomalies</option>');
            console.error('Error loading anomalies:', error);
        }
    });
}

/**
 * Run root cause analysis for the selected anomaly
 */
function runRootCauseAnalysis() {
    const anomalyId = $('#anomalyId').val();
    const outputFormat = $('#outputFormat').val();
    
    if (!anomalyId) {
        alert('Please select an anomaly');
        return;
    }
    
    // Show loading indicator
    const submitBtn = $('#rcaForm button[type="submit"]');
    const originalText = submitBtn.text();
    submitBtn.prop('disabled', true).text('Running Analysis...');
    
    // Run RCA via API
    $.ajax({
        url: '/api/rca/analyze',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            anomaly_id: anomalyId,
            format: outputFormat
        }),
        success: function(data) {
            submitBtn.prop('disabled', false).text(originalText);
            
            // Show the result in a modal
            viewRcaDetail(anomalyId);
            
            // Refresh the recent analyses list
            loadRecentRcaData();
        },
        error: function(xhr, status, error) {
            submitBtn.prop('disabled', false).text(originalText);
            alert('Error running analysis: ' + error);
            console.error('Error running analysis:', error);
        }
    });
}

/**
 * View RCA details for a specific anomaly
 */
function viewRcaDetail(anomalyId) {
    // Show loading indicator
    $('#rcaDetailContent').html('<div class="text-center"><div class="spinner-border" role="status"><span class="sr-only">Loading...</span></div></div>');
    $('#rcaDetailModal').modal('show');
    
    // Fetch RCA details from API
    $.ajax({
        url: `/api/rca/${anomalyId}`,
        method: 'GET',
        success: function(data) {
            displayRcaDetail(data);
        },
        error: function(xhr, status, error) {
            $('#rcaDetailContent').html('<div class="alert alert-danger">Error loading RCA details: ' + error + '</div>');
            console.error('Error loading RCA details:', error);
        }
    });
}

/**
 * Display RCA details in the modal
 */
function displayRcaDetail(data) {
    if (!data) {
        $('#rcaDetailContent').html('<div class="alert alert-warning">No RCA data found</div>');
        return;
    }
    
    // Update modal title
    $('#rcaDetailModalLabel').text(`Root Cause Analysis: ${data.anomaly_id}`);
    
    // Build HTML content
    let html = `
        <div class="rca-detail">
            <div class="card mb-3">
                <div class="card-header bg-primary text-white">
                    <h5 class="mb-0">Primary Root Cause</h5>
                </div>
                <div class="card-body">
                    <h5>${data.primary_root_cause ? data.primary_root_cause.description : 'N/A'}</h5>
                    <div class="progress mb-2">
                        <div class="progress-bar" role="progressbar" style="width: ${data.primary_root_cause ? (data.primary_root_cause.confidence * 100).toFixed(1) + '%' : '0%'}" 
                            aria-valuenow="${data.primary_root_cause ? (data.primary_root_cause.confidence * 100).toFixed(1) : '0'}" aria-valuemin="0" aria-valuemax="100">
                            ${data.primary_root_cause ? (data.primary_root_cause.confidence * 100).toFixed(1) + '%' : '0%'}
                        </div>
                    </div>
                    <p>${data.primary_root_cause ? data.primary_root_cause.details : 'No details available'}</p>
                </div>
            </div>
    `;
    
    // Other causes
    if (data.other_causes && data.other_causes.length > 0) {
        html += `
            <div class="card mb-3">
                <div class="card-header">
                    <h5 class="mb-0">Other Potential Causes</h5>
                </div>
                <div class="card-body">
                    <div class="table-responsive">
                        <table class="table table-sm">
                            <thead>
                                <tr>
                                    <th>Description</th>
                                    <th>Confidence</th>
                                    <th>Details</th>
                                </tr>
                            </thead>
                            <tbody>
        `;
        
        data.other_causes.forEach(function(cause) {
            html += `
                <tr>
                    <td>${cause.description}</td>
                    <td>
                        <div class="progress">
                            <div class="progress-bar" role="progressbar" style="width: ${(cause.confidence * 100).toFixed(1)}%" 
                                aria-valuenow="${(cause.confidence * 100).toFixed(1)}" aria-valuemin="0" aria-valuemax="100">
                                ${(cause.confidence * 100).toFixed(1)}%
                            </div>
                        </div>
                    </td>
                    <td>${cause.details || 'No details'}</td>
                </tr>
            `;
        });
        
        html += `
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
        `;
    }
    
    // Affected tables
    if (data.affected_tables && data.affected_tables.length > 0) {
        html += `
            <div class="card mb-3">
                <div class="card-header">
                    <h5 class="mb-0">Affected Tables</h5>
                </div>
                <div class="card-body">
                    <ul class="list-group">
        `;
        
        data.affected_tables.forEach(function(table) {
            html += `<li class="list-group-item">${table}</li>`;
        });
        
        html += `
                    </ul>
                </div>
            </div>
        `;
    }
    
    // Recommended actions
    if (data.recommended_actions && data.recommended_actions.length > 0) {
        html += `
            <div class="card mb-3">
                <div class="card-header">
                    <h5 class="mb-0">Recommended Actions</h5>
                </div>
                <div class="card-body">
                    <ol class="list-group list-group-numbered">
        `;
        
        data.recommended_actions.forEach(function(action) {
            html += `<li class="list-group-item">${action}</li>`;
        });
        
        html += `
                    </ol>
                </div>
            </div>
        `;
    }
    
    // Related anomalies
    if (data.related_anomalies && data.related_anomalies.length > 0) {
        html += `
            <div class="card mb-3">
                <div class="card-header">
                    <h5 class="mb-0">Related Anomalies</h5>
                </div>
                <div class="card-body">
                    <ul class="list-group">
        `;
        
        data.related_anomalies.forEach(function(anomaly) {
            html += `<li class="list-group-item"><a href="#" class="view-related-anomaly" data-id="${anomaly}">${anomaly}</a></li>`;
        });
        
        html += `
                    </ul>
                </div>
            </div>
        `;
    }
    
    // Grafana dashboard link
    if (data.grafana_dashboard_url) {
        html += `
            <div class="card mb-3">
                <div class="card-header">
                    <h5 class="mb-0">Monitoring Dashboard</h5>
                </div>
                <div class="card-body">
                    <a href="${data.grafana_dashboard_url}" target="_blank" class="btn btn-info">
                        <i data-feather="external-link"></i> View in Grafana
                    </a>
                </div>
            </div>
        `;
    }
    
    // Analysis metadata
    html += `
        <div class="card">
            <div class="card-header">
                <h5 class="mb-0">Analysis Metadata</h5>
            </div>
            <div class="card-body">
                <dl class="row">
                    <dt class="col-sm-3">Anomaly ID</dt>
                    <dd class="col-sm-9">${data.anomaly_id}</dd>
                    
                    <dt class="col-sm-3">Analysis Time</dt>
                    <dd class="col-sm-9">${formatDate(data.analysis_time)}</dd>
                </dl>
            </div>
        </div>
    </div>
    `;
    
    // Set the HTML content
    $('#rcaDetailContent').html(html);
    
    // Initialize Feather icons in the modal
    feather.replace();
    
    // Add event listeners for related anomalies
    $('.view-related-anomaly').click(function(e) {
        e.preventDefault();
        const relatedId = $(this).data('id');
        viewRcaDetail(relatedId);
    });
}

/**
 * Load RCA insights data and display charts
 */
function loadRcaInsights() {
    $.ajax({
        url: '/api/rca/insights',
        method: 'GET',
        success: function(data) {
            displayRcaInsights(data);
        },
        error: function(xhr, status, error) {
            console.error('Error loading RCA insights:', error);
        }
    });
}

/**
 * Display RCA insights with charts
 */
function displayRcaInsights(data) {
    if (!data) {
        return;
    }
    
    // Root cause distribution chart
    if (data.root_cause_distribution) {
        const ctx1 = document.getElementById('rootCauseDistributionChart').getContext('2d');
        new Chart(ctx1, {
            type: 'pie',
            data: {
                labels: Object.keys(data.root_cause_distribution),
                datasets: [{
                    data: Object.values(data.root_cause_distribution),
                    backgroundColor: [
                        '#4e73df', '#1cc88a', '#36b9cc', '#f6c23e', '#e74a3b',
                        '#5a5c69', '#858796', '#6610f2', '#fd7e14', '#20c9a6'
                    ]
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Root Cause Distribution'
                    },
                    legend: {
                        position: 'bottom'
                    }
                }
            }
        });
    }
    
    // Affected tables chart
    if (data.affected_tables_count) {
        const ctx2 = document.getElementById('affectedTablesChart').getContext('2d');
        new Chart(ctx2, {
            type: 'bar',
            data: {
                labels: Object.keys(data.affected_tables_count),
                datasets: [{
                    label: 'Number of Incidents',
                    data: Object.values(data.affected_tables_count),
                    backgroundColor: '#4e73df'
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Most Affected Tables'
                    },
                    legend: {
                        display: false
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        ticks: {
                            precision: 0
                        }
                    }
                }
            }
        });
    }
    
    // Common root causes
    if (data.common_root_causes) {
        let html = '<div class="list-group">';
        data.common_root_causes.forEach(function(cause, index) {
            html += `
                <div class="list-group-item">
                    <div class="d-flex w-100 justify-content-between">
                        <h5 class="mb-1">${index + 1}. ${cause.description}</h5>
                        <small>${cause.count} occurrences</small>
                    </div>
                    <p class="mb-1">${cause.details || 'No additional details'}</p>
                </div>
            `;
        });
        html += '</div>';
        $('#commonRootCauses').html(html);
    }
}

/**
 * Export RCA data to CSV or JSON
 */
function exportRcaData() {
    const format = prompt('Export format (csv or json):', 'csv');
    if (!format || (format !== 'csv' && format !== 'json')) {
        alert('Invalid format. Please specify csv or json.');
        return;
    }
    
    window.location.href = `/api/rca/export?format=${format}`;
}

/**
 * Export RCA detail to JSON
 */
function exportRcaDetail() {
    const anomalyId = $('#rcaDetailModalLabel').text().split(':')[1].trim();
    window.location.href = `/api/rca/${anomalyId}/export`;
}

/**
 * Format date string to a readable format
 */
function formatDate(dateString) {
    if (!dateString) return 'N/A';
    
    const date = new Date(dateString);
    return date.toLocaleString();
}
