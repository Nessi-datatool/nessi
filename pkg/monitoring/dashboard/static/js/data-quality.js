// Initialize dashboard when document is ready
$(document).ready(function() {
    // Initialize Feather icons
    feather.replace();
    
    // Load initial data
    loadTables();
    
    // Set up event handlers
    $('#refreshBtn').click(refreshData);
    $('#exportBtn').click(exportData);
    $('#addRuleBtn').click(showAddRuleModal);
    $('#runValidationBtn').click(runValidation);
    
    // Tab change handlers
    $('#qualityTabs a').on('shown.bs.tab', function (e) {
        const tabId = $(e.target).attr('id');
        if (tabId === 'profile-tab') {
            loadProfiles();
        } else if (tabId === 'rules-tab') {
            loadRules();
        } else if (tabId === 'validation-tab') {
            loadValidationResults();
        } else if (tabId === 'trends-tab') {
            loadTrends();
        }
    });
});

// Global variables
let currentTable = '';
let qualityScoreChart = null;
let validationTrendChart = null;

// Load available tables
function loadTables() {
    // In a real implementation, this would fetch tables from the server
    // For now, we'll add some example tables
    const tables = [
        { id: 'users', name: 'Users Table' },
        { id: 'orders', name: 'Orders Table' },
        { id: 'products', name: 'Products Table' }
    ];
    
    const dropdown = $('#tableDropdown');
    dropdown.empty();
    
    tables.forEach(table => {
        dropdown.append(`
            <a class="dropdown-item table-item" href="#" data-id="${table.id}">${table.name}</a>
        `);
    });
    
    // Set up click handlers for table selection
    $('.table-item').click(function(e) {
        e.preventDefault();
        currentTable = $(this).data('id');
        $('#tableSelector').text($(this).text());
        refreshData();
    });
    
    // Select first table by default
    if (tables.length > 0) {
        currentTable = tables[0].id;
        $('#tableSelector').text(tables[0].name);
        refreshData();
    }
}

// Refresh all data based on current tab
function refreshData() {
    const activeTab = $('#qualityTabs .nav-link.active').attr('id');
    
    if (activeTab === 'profile-tab') {
        loadProfiles();
    } else if (activeTab === 'rules-tab') {
        loadRules();
    } else if (activeTab === 'validation-tab') {
        loadValidationResults();
    } else if (activeTab === 'trends-tab') {
        loadTrends();
    }
}

// Export data
function exportData() {
    alert('Export functionality will be implemented here');
}

// Load profile data
function loadProfiles() {
    if (!currentTable) return;
    
    // Show loading indicator
    $('#profileTable tbody').html('<tr><td colspan="7" class="text-center">Loading...</td></tr>');
    
    // Fetch profile data from API
    $.ajax({
        url: `/api/profiles/summary?table=${currentTable}`,
        method: 'GET',
        success: function(response) {
            if (response.success) {
                displayProfiles(response.data);
            } else {
                $('#profileTable tbody').html(`<tr><td colspan="7" class="text-center text-danger">${response.message}</td></tr>`);
            }
        },
        error: function(xhr, status, error) {
            $('#profileTable tbody').html(`<tr><td colspan="7" class="text-center text-danger">Failed to load profiles: ${error}</td></tr>`);
        }
    });
}

// Display profile data
function displayProfiles(profiles) {
    const tbody = $('#profileTable tbody');
    tbody.empty();
    
    if (profiles.length === 0) {
        tbody.html('<tr><td colspan="7" class="text-center">No profiles available</td></tr>');
        return;
    }
    
    profiles.forEach(profile => {
        const qualityScore = profile.quality_score || 'N/A';
        const qualityClass = getQualityClass(qualityScore);
        
        tbody.append(`
            <tr>
                <td>${profile.name}</td>
                <td>${profile.type}</td>
                <td>${profile.null_percent.toFixed(2)}%</td>
                <td class="${qualityClass}">${typeof qualityScore === 'number' ? qualityScore.toFixed(2) : qualityScore}</td>
                <td>${profile.anomaly_count}</td>
                <td>${profile.pattern_count}</td>
                <td>
                    <button class="btn btn-sm btn-outline-primary view-profile-btn" data-name="${profile.name}">View</button>
                </td>
            </tr>
        `);
    });
    
    // Set up click handlers for profile details
    $('.view-profile-btn').click(function() {
        const columnName = $(this).data('name');
        showProfileDetails(columnName);
    });
}

// Get CSS class based on quality score
function getQualityClass(score) {
    if (typeof score !== 'number') return '';
    if (score >= 0.9) return 'text-success';
    if (score >= 0.7) return 'text-warning';
    return 'text-danger';
}

// Show profile details modal
function showProfileDetails(columnName) {
    // Show loading indicator
    $('#profileDetailContent').html('<div class="text-center">Loading...</div>');
    $('#profileDetailModalLabel').text(`Column Profile: ${columnName}`);
    $('#profileDetailModal').modal('show');
    
    // Fetch detailed profile data
    $.ajax({
        url: `/api/profiles?table=${currentTable}`,
        method: 'GET',
        success: function(response) {
            if (response.success) {
                const profiles = response.data;
                const profile = profiles.find(p => p.Profile.Name === columnName);
                if (profile) {
                    displayProfileDetails(profile);
                } else {
                    $('#profileDetailContent').html('<div class="text-center text-danger">Profile not found</div>');
                }
            } else {
                $('#profileDetailContent').html(`<div class="text-center text-danger">${response.message}</div>`);
            }
        },
        error: function(xhr, status, error) {
            $('#profileDetailContent').html(`<div class="text-center text-danger">Failed to load profile details: ${error}</div>`);
        }
    });
}

// Display profile details in modal
function displayProfileDetails(profile) {
    const content = $('#profileDetailContent');
    content.empty();
    
    // Basic information
    content.append(`
        <div class="row">
            <div class="col-md-6">
                <h6>Basic Information</h6>
                <table class="table table-sm">
                    <tr><th>Name</th><td>${profile.Profile.Name}</td></tr>
                    <tr><th>Type</th><td>${profile.Profile.Type}</td></tr>
                    <tr><th>Row Count</th><td>${profile.Profile.RowCount}</td></tr>
                    <tr><th>Null Count</th><td>${profile.Profile.NullCount}</td></tr>
                    <tr><th>Null Percent</th><td>${profile.Profile.NullPercent.toFixed(2)}%</td></tr>
                    <tr><th>Distinct Values</th><td>${profile.Profile.Distinct}</td></tr>
                </table>
            </div>
            <div class="col-md-6">
                <h6>Quality Scores</h6>
                <table class="table table-sm">
                    <tr><th>Overall</th><td class="${getQualityClass(profile.QualityScore?.Overall)}">${profile.QualityScore?.Overall?.toFixed(2) || 'N/A'}</td></tr>
                    <tr><th>Completeness</th><td class="${getQualityClass(profile.QualityScore?.Completeness)}">${profile.QualityScore?.Completeness?.toFixed(2) || 'N/A'}</td></tr>
                    <tr><th>Consistency</th><td class="${getQualityClass(profile.QualityScore?.Consistency)}">${profile.QualityScore?.Consistency?.toFixed(2) || 'N/A'}</td></tr>
                    <tr><th>Accuracy</th><td class="${getQualityClass(profile.QualityScore?.Accuracy)}">${profile.QualityScore?.Accuracy?.toFixed(2) || 'N/A'}</td></tr>
                    <tr><th>Uniqueness</th><td class="${getQualityClass(profile.QualityScore?.Uniqueness)}">${profile.QualityScore?.Uniqueness?.toFixed(2) || 'N/A'}</td></tr>
                    <tr><th>Recommended Type</th><td>${profile.QualityScore?.RecommendedType || profile.Profile.Type}</td></tr>
                </table>
            </div>
        </div>
    `);
    
    // Statistics (for numeric columns)
    if (profile.Profile.Stats && (profile.Profile.Type === 'integer' || profile.Profile.Type === 'float')) {
        content.append(`
            <div class="row mt-3">
                <div class="col-12">
                    <h6>Statistics</h6>
                    <table class="table table-sm">
                        <tr><th>Min</th><td>${profile.Profile.Stats.Min}</td></tr>
                        <tr><th>Max</th><td>${profile.Profile.Stats.Max}</td></tr>
                        <tr><th>Mean</th><td>${profile.Profile.Stats.Mean.toFixed(2)}</td></tr>
                        <tr><th>Standard Deviation</th><td>${profile.Profile.Stats.StdDev.toFixed(2)}</td></tr>
                        <tr><th>Quartiles (25%, 50%, 75%)</th><td>${profile.Profile.Stats.Quartiles.map(q => q.toFixed(2)).join(', ')}</td></tr>
                    </table>
                </div>
            </div>
        `);
    }
    
    // Patterns
    if (profile.Profile.Patterns && profile.Profile.Patterns.length > 0) {
        content.append(`
            <div class="row mt-3">
                <div class="col-12">
                    <h6>Detected Patterns</h6>
                    <ul class="list-group">
                        ${profile.Profile.Patterns.map(pattern => `<li class="list-group-item">${pattern}</li>`).join('')}
                    </ul>
                </div>
            </div>
        `);
    }
    
    // Detailed patterns
    if (profile.DetailedPatterns && profile.DetailedPatterns.length > 0) {
        content.append(`
            <div class="row mt-3">
                <div class="col-12">
                    <h6>Pattern Analysis</h6>
                    <table class="table table-sm">
                        <thead>
                            <tr>
                                <th>Pattern</th>
                                <th>Description</th>
                                <th>Confidence</th>
                                <th>Examples</th>
                            </tr>
                        </thead>
                        <tbody>
                            ${profile.DetailedPatterns.map(pattern => `
                                <tr>
                                    <td>${pattern.Name}</td>
                                    <td>${pattern.Description}</td>
                                    <td>${(pattern.Confidence * 100).toFixed(1)}%</td>
                                    <td>${pattern.Examples.join(', ')}</td>
                                </tr>
                            `).join('')}
                        </tbody>
                    </table>
                </div>
            </div>
        `);
    }
    
    // Anomalies
    if (profile.Profile.Anomalies && profile.Profile.Anomalies.length > 0) {
        content.append(`
            <div class="row mt-3">
                <div class="col-12">
                    <h6>Detected Anomalies</h6>
                    <table class="table table-sm">
                        <thead>
                            <tr>
                                <th>Type</th>
                                <th>Value</th>
                                <th>Description</th>
                            </tr>
                        </thead>
                        <tbody>
                            ${profile.Profile.Anomalies.map(anomaly => `
                                <tr>
                                    <td>${anomaly.Type}</td>
                                    <td>${anomaly.Value}</td>
                                    <td>${anomaly.Description}</td>
                                </tr>
                            `).join('')}
                        </tbody>
                    </table>
                </div>
            </div>
        `);
    }
    
    // Distribution (if available)
    if (profile.Distribution) {
        content.append(`
            <div class="row mt-3">
                <div class="col-12">
                    <h6>Distribution Analysis</h6>
                    <div class="row">
                        <div class="col-md-6">
                            <canvas id="distributionChart"></canvas>
                        </div>
                        <div class="col-md-6">
                            <table class="table table-sm">
                                <tr><th>Entropy</th><td>${profile.Distribution.Entropy.toFixed(2)}</td></tr>
                                ${profile.Distribution.Skewness ? `<tr><th>Skewness</th><td>${profile.Distribution.Skewness.toFixed(2)}</td></tr>` : ''}
                                ${profile.Distribution.Kurtosis ? `<tr><th>Kurtosis</th><td>${profile.Distribution.Kurtosis.toFixed(2)}</td></tr>` : ''}
                                ${profile.Distribution.IsNormal !== undefined ? `<tr><th>Normal Distribution</th><td>${profile.Distribution.IsNormal ? 'Yes' : 'No'}</td></tr>` : ''}
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        `);
        
        // Create distribution chart
        setTimeout(() => {
            const ctx = document.getElementById('distributionChart').getContext('2d');
            const labels = Object.keys(profile.Distribution.Histogram);
            const data = Object.values(profile.Distribution.Histogram);
            
            new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Distribution',
                        data: data,
                        backgroundColor: 'rgba(54, 162, 235, 0.5)',
                        borderColor: 'rgba(54, 162, 235, 1)',
                        borderWidth: 1
                    }]
                },
                options: {
                    responsive: true,
                    scales: {
                        y: {
                            beginAtZero: true
                        }
                    }
                }
            });
        }, 100);
    }
}

// Load validation rules
function loadRules() {
    // Show loading indicator
    $('#rulesTable tbody').html('<tr><td colspan="6" class="text-center">Loading...</td></tr>');
    
    // Fetch rules from API
    $.ajax({
        url: '/api/rules',
        method: 'GET',
        success: function(response) {
            if (response.success) {
                displayRules(response.data);
            } else {
                $('#rulesTable tbody').html(`<tr><td colspan="6" class="text-center text-danger">${response.message}</td></tr>`);
            }
        },
        error: function(xhr, status, error) {
            $('#rulesTable tbody').html(`<tr><td colspan="6" class="text-center text-danger">Failed to load rules: ${error}</td></tr>`);
        }
    });
}

// Display validation rules
function displayRules(rules) {
    const tbody = $('#rulesTable tbody');
    tbody.empty();
    
    if (rules.length === 0) {
        tbody.html('<tr><td colspan="6" class="text-center">No rules available</td></tr>');
        return;
    }
    
    rules.forEach(rule => {
        const severityClass = getSeverityClass(rule.Metadata.Severity);
        
        tbody.append(`
            <tr>
                <td>${rule.Metadata.ID}</td>
                <td>${rule.Metadata.Name}</td>
                <td>${getRuleType(rule)}</td>
                <td class="${severityClass}">${rule.Metadata.Severity}</td>
                <td>${rule.Metadata.Description}</td>
                <td>
                    <button class="btn btn-sm btn-outline-primary view-rule-btn" data-id="${rule.Metadata.ID}">View</button>
                    <button class="btn btn-sm btn-outline-danger delete-rule-btn" data-id="${rule.Metadata.ID}">Delete</button>
                </td>
            </tr>
        `);
    });
    
    // Set up click handlers for rule details
    $('.view-rule-btn').click(function() {
        const ruleId = $(this).data('id');
        showRuleDetails(ruleId);
    });
    
    // Set up click handlers for rule deletion
    $('.delete-rule-btn').click(function() {
        const ruleId = $(this).data('id');
        if (confirm(`Are you sure you want to delete rule ${ruleId}?`)) {
            deleteRule(ruleId);
        }
    });
}

// Get rule type based on rule object
function getRuleType(rule) {
    // In a real implementation, this would determine the rule type from the rule object
    // For now, we'll use a simple heuristic
    if (rule.Metadata.Config.fields && Array.isArray(rule.Metadata.Config.fields)) {
        return 'Null Check';
    } else if (rule.Metadata.Config.fields && typeof rule.Metadata.Config.fields === 'object') {
        return 'Range Check';
    } else if (rule.Metadata.Config.pattern) {
        return 'Regex';
    } else if (rule.Metadata.Config.values) {
        return 'Enum';
    } else if (rule.Metadata.Config.min_length !== undefined) {
        return 'Length';
    } else if (rule.Metadata.Config.format) {
        return 'Date Format';
    }
    return 'Custom';
}

// Get CSS class based on severity
function getSeverityClass(severity) {
    switch (severity.toLowerCase()) {
        case 'error':
        case 'critical':
            return 'text-danger';
        case 'warning':
            return 'text-warning';
        case 'info':
            return 'text-info';
        default:
            return '';
    }
}

// Show rule details modal
function showRuleDetails(ruleId) {
    // In a real implementation, this would fetch rule details from the server
    // For now, we'll just show a placeholder
    $('#ruleDetailModalLabel').text(`Rule Details: ${ruleId}`);
    $('#ruleDetailContent').html(`
        <div class="text-center">
            <p>Detailed information for rule ${ruleId} would be displayed here.</p>
        </div>
    `);
    $('#ruleDetailModal').modal('show');
}

// Show add rule modal
function showAddRuleModal() {
    alert('Add rule functionality will be implemented here');
}

// Delete rule
function deleteRule(ruleId) {
    alert(`Delete rule ${ruleId} functionality will be implemented here`);
}

// Load validation results
function loadValidationResults() {
    if (!currentTable) return;
    
    // Show loading indicator
    $('#validationTable tbody').html('<tr><td colspan="6" class="text-center">Loading...</td></tr>');
    
    // In a real implementation, this would fetch validation results from the server
    // For now, we'll just show placeholder data
    setTimeout(() => {
        const results = [
            {
                rule_id: 'null_check_1',
                rule_name: 'Required Fields Check',
                severity: 'error',
                passed: true,
                error_count: 0,
                timestamp: new Date().toISOString(),
                details: []
            },
            {
                rule_id: 'range_check_1',
                rule_name: 'Age Range Check',
                severity: 'warning',
                passed: false,
                error_count: 2,
                timestamp: new Date().toISOString(),
                details: ['Value 150 is outside range [0, 120]', 'Value -5 is outside range [0, 120]']
            }
        ];
        
        displayValidationResults(results);
    }, 500);
}

// Display validation results
function displayValidationResults(results) {
    const tbody = $('#validationTable tbody');
    tbody.empty();
    
    if (results.length === 0) {
        tbody.html('<tr><td colspan="6" class="text-center">No validation results available</td></tr>');
        return;
    }
    
    results.forEach(result => {
        const severityClass = getSeverityClass(result.severity);
        const statusClass = result.passed ? 'text-success' : 'text-danger';
        const statusText = result.passed ? 'Passed' : 'Failed';
        
        tbody.append(`
            <tr>
                <td>${result.rule_name}</td>
                <td class="${severityClass}">${result.severity}</td>
                <td class="${statusClass}">${statusText}</td>
                <td>${result.error_count}</td>
                <td>${new Date(result.timestamp).toLocaleString()}</td>
                <td>
                    ${result.details.length > 0 
                        ? `<button class="btn btn-sm btn-outline-info view-details-btn" data-toggle="tooltip" title="${result.details.join('\n')}">View</button>` 
                        : '-'}
                </td>
            </tr>
        `);
    });
    
    // Initialize tooltips
    $('[data-toggle="tooltip"]').tooltip();
}

// Run validation
function runValidation() {
    if (!currentTable) return;
    
    // Show loading indicator
    $('#validationTable tbody').html('<tr><td colspan="6" class="text-center">Running validation...</td></tr>');
    
    // In a real implementation, this would send a request to run validation
    // For now, we'll just reload the validation results after a delay
    setTimeout(loadValidationResults, 1000);
}

// Load trends data
function loadTrends() {
    // In a real implementation, this would fetch trend data from the server
    // For now, we'll just show placeholder charts
    
    // Quality score trend chart
    const qualityScoreCtx = document.getElementById('qualityScoreChart').getContext('2d');
    if (qualityScoreChart) {
        qualityScoreChart.destroy();
    }
    
    const dates = [];
    const scores = [];
    
    // Generate 7 days of data
    const now = new Date();
    for (let i = 6; i >= 0; i--) {
        const date = new Date(now);
        date.setDate(date.getDate() - i);
        dates.push(date.toLocaleDateString());
        
        // Random score between 0.7 and 1.0
        scores.push(0.7 + Math.random() * 0.3);
    }
    
    qualityScoreChart = new Chart(qualityScoreCtx, {
        type: 'line',
        data: {
            labels: dates,
            datasets: [{
                label: 'Overall Quality Score',
                data: scores,
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                borderColor: 'rgba(75, 192, 192, 1)',
                borderWidth: 2,
                tension: 0.1
            }]
        },
        options: {
            responsive: true,
            scales: {
                y: {
                    min: 0,
                    max: 1,
                    ticks: {
                        callback: function(value) {
                            return value.toFixed(1);
                        }
                    }
                }
            },
            plugins: {
                title: {
                    display: true,
                    text: 'Quality Score Trend'
                }
            }
        }
    });
    
    // Validation trend chart
    const validationTrendCtx = document.getElementById('validationTrendChart').getContext('2d');
    if (validationTrendChart) {
        validationTrendChart.destroy();
    }
    
    const passedRules = [];
    const failedRules = [];
    
    // Generate 7 days of data
    for (let i = 0; i < 7; i++) {
        // Random values
        const passed = 8 + Math.floor(Math.random() * 3);
        const failed = Math.floor(Math.random() * 3);
        
        passedRules.push(passed);
        failedRules.push(failed);
    }
    
    validationTrendChart = new Chart(validationTrendCtx, {
        type: 'bar',
        data: {
            labels: dates,
            datasets: [
                {
                    label: 'Passed Rules',
                    data: passedRules,
                    backgroundColor: 'rgba(75, 192, 192, 0.5)',
                    borderColor: 'rgba(75, 192, 192, 1)',
                    borderWidth: 1
                },
                {
                    label: 'Failed Rules',
                    data: failedRules,
                    backgroundColor: 'rgba(255, 99, 132, 0.5)',
                    borderColor: 'rgba(255, 99, 132, 1)',
                    borderWidth: 1
                }
            ]
        },
        options: {
            responsive: true,
            scales: {
                y: {
                    beginAtZero: true,
                    stacked: false
                },
                x: {
                    stacked: false
                }
            },
            plugins: {
                title: {
                    display: true,
                    text: 'Validation Results Trend'
                }
            }
        }
    });
}
