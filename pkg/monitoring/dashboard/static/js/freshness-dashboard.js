/**
 * Data Freshness Dashboard JavaScript
 * This file handles the freshness dashboard functionality
 */

// Initialize dashboard when document is ready
$(document).ready(function() {
    // Initialize Feather icons
    feather.replace();
    
    // Load initial data
    loadFreshnessData();
    loadSLAConfigurations();
    loadFreshnessTrends();
    
    // Set up event handlers
    setupEventHandlers();
});

/**
 * Set up event handlers for the freshness dashboard
 */
function setupEventHandlers() {
    // Refresh button
    $('#refreshBtn').click(function() {
        loadFreshnessData();
        loadSLAConfigurations();
        loadFreshnessTrends();
    });
    
    // Export button
    $('#exportBtn').click(function() {
        exportFreshnessData();
    });
    
    // Add SLA button
    $('#addSlaBtn').click(function() {
        // Reset form
        $('#slaForm')[0].reset();
        $('#slaModalLabel').text('Add SLA Configuration');
        $('#slaModal').modal('show');
    });
    
    // Save SLA button
    $('#saveSlaBtn').click(function() {
        saveSLAConfiguration();
    });
    
    // Custom frequency toggle
    $('#expectedFrequency').change(function() {
        if ($(this).val() === 'custom') {
            $('#customFrequencyGroup').show();
        } else {
            $('#customFrequencyGroup').hide();
        }
    });
    
    // Table selector change
    $('#tableSelector').change(function() {
        loadFreshnessTrends();
    });
}

/**
 * Load freshness data from the API
 */
function loadFreshnessData() {
    // Show loading indicator
    $('#freshnessTableBody').html('<tr><td colspan="7" class="text-center">Loading...</td></tr>');
    
    // Fetch freshness data from API
    $.ajax({
        url: '/api/freshness/status',
        method: 'GET',
        success: function(data) {
            displayFreshnessData(data);
        },
        error: function(xhr, status, error) {
            $('#freshnessTableBody').html('<tr><td colspan="7" class="text-center text-danger">Error loading freshness data: ' + error + '</td></tr>');
            console.error('Error loading freshness data:', error);
        }
    });
}

/**
 * Display freshness data in the table
 */
function displayFreshnessData(data) {
    if (!data || data.length === 0) {
        $('#freshnessTableBody').html('<tr><td colspan="7" class="text-center">No freshness data found</td></tr>');
        return;
    }
    
    let html = '';
    data.forEach(function(item) {
        // Format dates
        const lastUpdate = new Date(item.last_update_time).toLocaleString();
        const nextExpected = new Date(item.next_expected_update).toLocaleString();
        
        // Format status with color
        let statusHtml = '';
        switch (item.status) {
            case 'critical':
                statusHtml = '<span class="badge badge-danger">Critical</span>';
                break;
            case 'warning':
                statusHtml = '<span class="badge badge-warning">Warning</span>';
                break;
            case 'info':
                statusHtml = '<span class="badge badge-success">Up-to-date</span>';
                break;
            default:
                statusHtml = '<span class="badge badge-secondary">Unknown</span>';
        }
        
        html += `
            <tr>
                <td>${item.table_name}</td>
                <td>${lastUpdate}</td>
                <td>${formatDuration(item.time_since_update)}</td>
                <td>${formatDuration(item.expected_frequency)}</td>
                <td>${nextExpected}</td>
                <td>${statusHtml}</td>
                <td>
                    <button class="btn btn-sm btn-info view-table-btn" data-table="${item.table_name}">View</button>
                </td>
            </tr>
        `;
    });
    
    $('#freshnessTableBody').html(html);
    
    // Add event listener for view buttons
    $('.view-table-btn').click(function() {
        const tableName = $(this).data('table');
        viewTableDetails(tableName);
    });
}

/**
 * Load SLA configurations from the API
 */
function loadSLAConfigurations() {
    // Show loading indicator
    $('#slaTableBody').html('<tr><td colspan="7" class="text-center">Loading...</td></tr>');
    
    // Fetch SLA configurations from API
    $.ajax({
        url: '/api/freshness/sla',
        method: 'GET',
        success: function(data) {
            displaySLAConfigurations(data);
            updateTableSelector(data);
        },
        error: function(xhr, status, error) {
            $('#slaTableBody').html('<tr><td colspan="7" class="text-center text-danger">Error loading SLA configurations: ' + error + '</td></tr>');
            console.error('Error loading SLA configurations:', error);
        }
    });
}

/**
 * Display SLA configurations in the table
 */
function displaySLAConfigurations(data) {
    if (!data || data.length === 0) {
        $('#slaTableBody').html('<tr><td colspan="7" class="text-center">No SLA configurations found</td></tr>');
        return;
    }
    
    let html = '';
    data.forEach(function(item) {
        // Format enabled status
        const enabledHtml = item.enabled ? 
            '<span class="badge badge-success">Yes</span>' : 
            '<span class="badge badge-secondary">No</span>';
        
        html += `
            <tr>
                <td>${item.table_name}</td>
                <td>${item.table_path}</td>
                <td>${formatDuration(item.expected_frequency)}</td>
                <td>${item.warning_threshold}%</td>
                <td>${item.critical_threshold}%</td>
                <td>${enabledHtml}</td>
                <td>
                    <button class="btn btn-sm btn-primary edit-sla-btn" data-table="${item.table_name}">Edit</button>
                    <button class="btn btn-sm btn-danger delete-sla-btn" data-table="${item.table_name}">Delete</button>
                </td>
            </tr>
        `;
    });
    
    $('#slaTableBody').html(html);
    
    // Add event listeners for edit and delete buttons
    $('.edit-sla-btn').click(function() {
        const tableName = $(this).data('table');
        editSLAConfiguration(tableName);
    });
    
    $('.delete-sla-btn').click(function() {
        const tableName = $(this).data('table');
        deleteSLAConfiguration(tableName);
    });
}

/**
 * Update the table selector dropdown with available tables
 */
function updateTableSelector(data) {
    if (!data || data.length === 0) {
        return;
    }
    
    // Get current selection
    const currentSelection = $('#tableSelector').val();
    
    // Clear existing options (except "All Tables")
    $('#tableSelector option:not(:first)').remove();
    
    // Add options for each table
    data.forEach(function(item) {
        $('#tableSelector').append(`<option value="${item.table_name}">${item.table_name}</option>`);
    });
    
    // Restore selection if possible
    if (currentSelection && currentSelection !== 'all') {
        $('#tableSelector').val(currentSelection);
    }
}

/**
 * Load freshness trends from the API
 */
function loadFreshnessTrends() {
    const tableName = $('#tableSelector').val();
    const url = tableName === 'all' ? 
        '/api/freshness/trends' : 
        `/api/freshness/trends?table=${tableName}`;
    
    // Fetch freshness trends from API
    $.ajax({
        url: url,
        method: 'GET',
        success: function(data) {
            displayFreshnessTrends(data);
        },
        error: function(xhr, status, error) {
            console.error('Error loading freshness trends:', error);
        }
    });
}

/**
 * Display freshness trends in charts
 */
function displayFreshnessTrends(data) {
    if (!data || !data.history || data.history.length === 0) {
        return;
    }
    
    // Prepare data for charts
    const labels = data.history.map(item => new Date(item.timestamp).toLocaleString());
    const datasets = [];
    
    // Group data by table
    const tableData = {};
    data.history.forEach(item => {
        if (!tableData[item.table_name]) {
            tableData[item.table_name] = {
                timestamps: [],
                timeSinceUpdate: []
            };
        }
        
        tableData[item.table_name].timestamps.push(item.timestamp);
        tableData[item.table_name].timeSinceUpdate.push(item.time_since_update / 3600); // Convert to hours
    });
    
    // Create datasets for each table
    Object.keys(tableData).forEach((tableName, index) => {
        const color = getChartColor(index);
        datasets.push({
            label: tableName,
            data: tableData[tableName].timeSinceUpdate,
            borderColor: color,
            backgroundColor: color + '20', // Add transparency
            fill: false
        });
    });
    
    // Create freshness history chart
    const ctx1 = document.getElementById('freshnessHistoryChart').getContext('2d');
    if (window.freshnessHistoryChart) {
        window.freshnessHistoryChart.destroy();
    }
    
    window.freshnessHistoryChart = new Chart(ctx1, {
        type: 'line',
        data: {
            labels: labels,
            datasets: datasets
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                x: {
                    title: {
                        display: true,
                        text: 'Time'
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: 'Hours Since Update'
                    },
                    beginAtZero: true
                }
            },
            plugins: {
                title: {
                    display: true,
                    text: 'Freshness History'
                },
                tooltip: {
                    mode: 'index',
                    intersect: false
                }
            }
        }
    });
    
    // Create SLA compliance chart
    const ctx2 = document.getElementById('slaComplianceChart').getContext('2d');
    if (window.slaComplianceChart) {
        window.slaComplianceChart.destroy();
    }
    
    // Prepare compliance data
    const complianceData = {
        labels: ['Up-to-date', 'Warning', 'Critical'],
        datasets: [{
            data: [
                data.compliance.info_count || 0,
                data.compliance.warning_count || 0,
                data.compliance.critical_count || 0
            ],
            backgroundColor: [
                '#28a745', // Green for up-to-date
                '#ffc107', // Yellow for warning
                '#dc3545'  // Red for critical
            ]
        }]
    };
    
    window.slaComplianceChart = new Chart(ctx2, {
        type: 'doughnut',
        data: complianceData,
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: {
                    display: true,
                    text: 'SLA Compliance'
                },
                legend: {
                    position: 'bottom'
                }
            }
        }
    });
}

/**
 * View table details
 */
function viewTableDetails(tableName) {
    // Fetch table details from API
    $.ajax({
        url: `/api/freshness/status?table=${tableName}`,
        method: 'GET',
        success: function(data) {
            if (!data || data.length === 0) {
                $('#tableDetailsContent').html('<div class="alert alert-warning">No details found for table ' + tableName + '</div>');
                $('#viewGrafanaBtn').hide();
                $('#tableDetailsModal').modal('show');
                return;
            }
            
            const tableData = data[0];
            
            // Format dates
            const lastUpdate = new Date(tableData.last_update_time).toLocaleString();
            const nextExpected = new Date(tableData.next_expected_update).toLocaleString();
            
            // Format status with color
            let statusHtml = '';
            switch (tableData.status) {
                case 'critical':
                    statusHtml = '<span class="badge badge-danger">Critical</span>';
                    break;
                case 'warning':
                    statusHtml = '<span class="badge badge-warning">Warning</span>';
                    break;
                case 'info':
                    statusHtml = '<span class="badge badge-success">Up-to-date</span>';
                    break;
                default:
                    statusHtml = '<span class="badge badge-secondary">Unknown</span>';
            }
            
            // Build HTML content
            let html = `
                <div class="table-details">
                    <h4>${tableData.table_name}</h4>
                    <p class="text-muted">${tableData.table_path}</p>
                    
                    <div class="row mt-4">
                        <div class="col-md-6">
                            <h5>Freshness Status</h5>
                            <table class="table table-sm">
                                <tr>
                                    <th>Status</th>
                                    <td>${statusHtml}</td>
                                </tr>
                                <tr>
                                    <th>Last Update</th>
                                    <td>${lastUpdate}</td>
                                </tr>
                                <tr>
                                    <th>Time Since Update</th>
                                    <td>${formatDuration(tableData.time_since_update)}</td>
                                </tr>
                                <tr>
                                    <th>Next Expected Update</th>
                                    <td>${nextExpected}</td>
                                </tr>
                            </table>
                        </div>
                        
                        <div class="col-md-6">
                            <h5>SLA Configuration</h5>
                            <table class="table table-sm">
                                <tr>
                                    <th>Expected Frequency</th>
                                    <td>${formatDuration(tableData.expected_frequency)}</td>
                                </tr>
                                <tr>
                                    <th>Warning Threshold</th>
                                    <td>${tableData.sla_config.warning_threshold}%</td>
                                </tr>
                                <tr>
                                    <th>Critical Threshold</th>
                                    <td>${tableData.sla_config.critical_threshold}%</td>
                                </tr>
                                <tr>
                                    <th>Enabled</th>
                                    <td>${tableData.sla_config.enabled ? 'Yes' : 'No'}</td>
                                </tr>
                            </table>
                        </div>
                    </div>
            `;
            
            // Add description if available
            if (tableData.sla_config.description) {
                html += `
                    <div class="row mt-3">
                        <div class="col-12">
                            <h5>Description</h5>
                            <p>${tableData.sla_config.description}</p>
                        </div>
                    </div>
                `;
            }
            
            // Add tags if available
            if (tableData.sla_config.tags && tableData.sla_config.tags.length > 0) {
                html += `
                    <div class="row mt-3">
                        <div class="col-12">
                            <h5>Tags</h5>
                            <p>
                `;
                
                tableData.sla_config.tags.forEach(function(tag) {
                    html += `<span class="badge badge-info mr-1">${tag}</span>`;
                });
                
                html += `
                            </p>
                        </div>
                    </div>
                `;
            }
            
            html += '</div>';
            
            // Update modal content
            $('#tableDetailsContent').html(html);
            $('#tableDetailsModalLabel').text(`Table Details: ${tableData.table_name}`);
            

            
            // Show modal
            $('#tableDetailsModal').modal('show');
        },
        error: function(xhr, status, error) {
            $('#tableDetailsContent').html('<div class="alert alert-danger">Error loading table details: ' + error + '</div>');
            $('#viewGrafanaBtn').hide();
            $('#tableDetailsModal').modal('show');
            console.error('Error loading table details:', error);
        }
    });
}

/**
 * Edit SLA configuration
 */
function editSLAConfiguration(tableName) {
    // Fetch SLA configuration from API
    $.ajax({
        url: `/api/freshness/sla?table=${tableName}`,
        method: 'GET',
        success: function(data) {
            if (!data || data.length === 0) {
                alert('No SLA configuration found for table ' + tableName);
                return;
            }
            
            const slaConfig = data[0];
            
            // Populate form
            $('#tableName').val(slaConfig.table_name);
            $('#tablePath').val(slaConfig.table_path);
            
            // Set frequency
            let frequencyValue = 'custom';
            const frequency = slaConfig.expected_frequency;
            
            if (frequency === 3600) {
                frequencyValue = 'hourly';
            } else if (frequency === 86400) {
                frequencyValue = 'daily';
            } else if (frequency === 604800) {
                frequencyValue = 'weekly';
            } else if (frequency === 2592000) {
                frequencyValue = 'monthly';
            }
            
            $('#expectedFrequency').val(frequencyValue);
            
            if (frequencyValue === 'custom') {
                $('#customFrequencyGroup').show();
                $('#customFrequency').val(formatDurationForInput(frequency));
            } else {
                $('#customFrequencyGroup').hide();
            }
            
            $('#warningThreshold').val(slaConfig.warning_threshold);
            $('#criticalThreshold').val(slaConfig.critical_threshold);
            $('#description').val(slaConfig.description || '');
            $('#tags').val(slaConfig.tags ? slaConfig.tags.join(',') : '');

            $('#enabled').prop('checked', slaConfig.enabled);
            
            // Update modal title
            $('#slaModalLabel').text('Edit SLA Configuration');
            
            // Show modal
            $('#slaModal').modal('show');
        },
        error: function(xhr, status, error) {
            alert('Error loading SLA configuration: ' + error);
            console.error('Error loading SLA configuration:', error);
        }
    });
}

/**
 * Save SLA configuration
 */
function saveSLAConfiguration() {
    // Get form data
    const tableName = $('#tableName').val();
    const tablePath = $('#tablePath').val();
    const frequencyType = $('#expectedFrequency').val();
    const warningThreshold = parseInt($('#warningThreshold').val());
    const criticalThreshold = parseInt($('#criticalThreshold').val());
    const description = $('#description').val();
    const tags = $('#tags').val() ? $('#tags').val().split(',').map(tag => tag.trim()) : [];

    const enabled = $('#enabled').is(':checked');
    
    // Validate form
    if (!tableName || !tablePath) {
        alert('Table name and path are required');
        return;
    }
    
    if (warningThreshold >= criticalThreshold) {
        alert('Warning threshold must be less than critical threshold');
        return;
    }
    
    // Parse frequency
    let frequency;
    if (frequencyType === 'custom') {
        const customFrequency = $('#customFrequency').val();
        if (!customFrequency) {
            alert('Custom frequency is required');
            return;
        }
        
        // Send to server for parsing
        frequency = customFrequency;
    } else {
        frequency = frequencyType;
    }
    
    // Create SLA configuration
    const slaConfig = {
        table_name: tableName,
        table_path: tablePath,
        expected_frequency: frequency,
        warning_threshold: warningThreshold,
        critical_threshold: criticalThreshold,
        description: description,
        tags: tags,

        enabled: enabled
    };
    
    // Save SLA configuration
    $.ajax({
        url: '/api/freshness/sla',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(slaConfig),
        success: function() {
            // Hide modal
            $('#slaModal').modal('hide');
            
            // Reload SLA configurations
            loadSLAConfigurations();
            
            // Show success message
            alert('SLA configuration saved successfully');
        },
        error: function(xhr, status, error) {
            alert('Error saving SLA configuration: ' + error);
            console.error('Error saving SLA configuration:', error);
        }
    });
}

/**
 * Delete SLA configuration
 */
function deleteSLAConfiguration(tableName) {
    if (!confirm(`Are you sure you want to delete the SLA configuration for table ${tableName}?`)) {
        return;
    }
    
    // Delete SLA configuration
    $.ajax({
        url: `/api/freshness/sla?table=${tableName}`,
        method: 'DELETE',
        success: function() {
            // Reload SLA configurations
            loadSLAConfigurations();
            
            // Show success message
            alert('SLA configuration deleted successfully');
        },
        error: function(xhr, status, error) {
            alert('Error deleting SLA configuration: ' + error);
            console.error('Error deleting SLA configuration:', error);
        }
    });
}

/**
 * Export freshness data
 */
function exportFreshnessData() {
    const format = prompt('Export format (csv or json):', 'csv');
    if (!format || (format !== 'csv' && format !== 'json')) {
        alert('Invalid format. Please specify csv or json.');
        return;
    }
    
    window.location.href = `/api/freshness/export?format=${format}`;
}

/**
 * Format duration in seconds to a human-readable string
 */
function formatDuration(seconds) {
    if (typeof seconds !== 'number') {
        return 'N/A';
    }
    
    if (seconds < 60) {
        return `${Math.round(seconds)} seconds`;
    } else if (seconds < 3600) {
        return `${Math.round(seconds / 60)} minutes`;
    } else if (seconds < 86400) {
        return `${Math.round(seconds / 3600)} hours`;
    } else if (seconds < 604800) {
        return `${Math.round(seconds / 86400)} days`;
    } else if (seconds < 2592000) {
        return `${Math.round(seconds / 604800)} weeks`;
    } else {
        return `${Math.round(seconds / 2592000)} months`;
    }
}

/**
 * Format duration in seconds to a string suitable for input
 */
function formatDurationForInput(seconds) {
    if (typeof seconds !== 'number') {
        return '';
    }
    
    if (seconds < 60) {
        return `${seconds}s`;
    } else if (seconds < 3600) {
        return `${Math.round(seconds / 60)}m`;
    } else if (seconds < 86400) {
        return `${Math.round(seconds / 3600)}h`;
    } else if (seconds < 604800) {
        return `${Math.round(seconds / 86400)}d`;
    } else {
        return `${Math.round(seconds / 86400)}d`;
    }
}

/**
 * Get chart color by index
 */
function getChartColor(index) {
    const colors = [
        '#4e73df', // Primary blue
        '#1cc88a', // Success green
        '#36b9cc', // Info teal
        '#f6c23e', // Warning yellow
        '#e74a3b', // Danger red
        '#5a5c69', // Dark gray
        '#858796', // Secondary gray
        '#6610f2', // Purple
        '#fd7e14', // Orange
        '#20c9a6'  // Light teal
    ];
    
    return colors[index % colors.length];
}
