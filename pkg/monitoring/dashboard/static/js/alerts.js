// This file has been removed in the OSS version.
// Alerting and advanced dashboards are only available in LakeDiff Enterprise.
    document.getElementById('type-filter').addEventListener('change', applyFilters);
    document.getElementById('severity-filter').addEventListener('change', applyFilters);
    document.getElementById('search-input').addEventListener('input', applyFilters);
    
    // Modal event listeners
    document.getElementById('submit-alert-btn').addEventListener('click', createAlert);
    document.getElementById('cancel-alert-btn').addEventListener('click', hideCreateAlertModal);
    document.getElementById('submit-rule-btn').addEventListener('click', createRule);
    document.getElementById('cancel-rule-btn').addEventListener('click', hideCreateRuleModal);
    
    // Close buttons for modals
    const closeButtons = document.querySelectorAll('.close');
    closeButtons.forEach(button => {
        button.addEventListener('click', function() {
            hideAllModals();
        });
    });
    
    // Alert action buttons
    document.getElementById('acknowledge-btn').addEventListener('click', acknowledgeAlert);
    document.getElementById('resolve-btn').addEventListener('click', resolveAlert);
    document.getElementById('silence-btn').addEventListener('click', silenceAlert);
    document.getElementById('close-modal-btn').addEventListener('click', hideAlertModal);
});

// Global variables
let currentAlerts = [];
let currentRules = [];
let selectedAlertId = null;

// Initialize the dashboard
function initAlertsDashboard() {
    // Load alerts and rules
    fetchAlerts();
    fetchRules();
    
    // Update counts every 30 seconds
    setInterval(updateAlertCounts, 30000);
}

// Fetch alerts from the API
function fetchAlerts() {
    fetch('/api/alerts')
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch alerts');
            }
            return response.json();
        })
        .then(data => {
            currentAlerts = data.alerts || [];
            renderAlerts(currentAlerts);
            updateAlertCounts();
        })
        .catch(error => {
            console.error('Error fetching alerts:', error);
            document.getElementById('alerts-list').innerHTML = 
                `<div class="error-message">Failed to load alerts: ${error.message}</div>`;
        });
}

// Fetch alert rules from the API
function fetchRules() {
    fetch('/api/alerts/rules')
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch alert rules');
            }
            return response.json();
        })
        .then(data => {
            currentRules = data.rules || [];
            renderRules(currentRules);
        })
        .catch(error => {
            console.error('Error fetching rules:', error);
            document.getElementById('rules-list').innerHTML = 
                `<div class="error-message">Failed to load rules: ${error.message}</div>`;
        });
}

// Render alerts in the list
function renderAlerts(alerts) {
    const alertsList = document.getElementById('alerts-list');
    
    if (alerts.length === 0) {
        alertsList.innerHTML = '<div class="empty-message">No alerts found</div>';
        return;
    }
    
    let html = '';
    alerts.forEach(alert => {
        const severityClass = alert.severity.toLowerCase();
        const statusClass = alert.status.toLowerCase();
        
        html += `
            <div class="alert-item ${severityClass} ${statusClass}" data-id="${alert.id}">
                <div class="alert-header">
                    <div class="alert-name">${alert.name}</div>
                    <div class="alert-badges">
                        <span class="badge severity ${severityClass}">${alert.severity}</span>
                        <span class="badge status ${statusClass}">${alert.status}</span>
                        <span class="badge type">${alert.type}</span>
                    </div>
                </div>
                <div class="alert-description">${alert.description}</div>
                <div class="alert-meta">
                    <div class="alert-source">Source: ${alert.source}</div>
                    <div class="alert-time">Time: ${formatTimestamp(alert.timestamp)}</div>
                </div>
            </div>
        `;
    });
    
    alertsList.innerHTML = html;
    
    // Add click event to show alert details
    const alertItems = document.querySelectorAll('.alert-item');
    alertItems.forEach(item => {
        item.addEventListener('click', function() {
            const alertId = this.getAttribute('data-id');
            showAlertDetails(alertId);
        });
    });
}

// Render alert rules in the list
function renderRules(rules) {
    const rulesList = document.getElementById('rules-list');
    
    if (rules.length === 0) {
        rulesList.innerHTML = '<div class="empty-message">No rules found</div>';
        return;
    }
    
    let html = '';
    rules.forEach(rule => {
        const severityClass = rule.severity.toLowerCase();
        const statusClass = rule.enabled ? 'enabled' : 'disabled';
        
        html += `
            <div class="rule-item ${severityClass} ${statusClass}" data-id="${rule.id}">
                <div class="rule-header">
                    <div class="rule-name">${rule.name}</div>
                    <div class="rule-badges">
                        <span class="badge severity ${severityClass}">${rule.severity}</span>
                        <span class="badge status ${statusClass}">${rule.enabled ? 'Enabled' : 'Disabled'}</span>
                    </div>
                </div>
                <div class="rule-description">${rule.description}</div>
                <div class="rule-condition">
                    ${rule.metric} ${rule.comparison_operator} ${rule.threshold}
                </div>
                <div class="rule-actions">
                    <button class="btn small edit-rule" data-id="${rule.id}">Edit</button>
                    <button class="btn small ${rule.enabled ? 'disable-rule' : 'enable-rule'}" data-id="${rule.id}">
                        ${rule.enabled ? 'Disable' : 'Enable'}
                    </button>
                    <button class="btn small delete-rule" data-id="${rule.id}">Delete</button>
                </div>
            </div>
        `;
    });
    
    rulesList.innerHTML = html;
    
    // Add click events for rule actions
    document.querySelectorAll('.edit-rule').forEach(btn => {
        btn.addEventListener('click', function(e) {
            e.stopPropagation();
            const ruleId = this.getAttribute('data-id');
            editRule(ruleId);
        });
    });
    
    document.querySelectorAll('.enable-rule, .disable-rule').forEach(btn => {
        btn.addEventListener('click', function(e) {
            e.stopPropagation();
            const ruleId = this.getAttribute('data-id');
            const isEnabled = this.classList.contains('disable-rule');
            toggleRuleStatus(ruleId, !isEnabled);
        });
    });
    
    document.querySelectorAll('.delete-rule').forEach(btn => {
        btn.addEventListener('click', function(e) {
            e.stopPropagation();
            const ruleId = this.getAttribute('data-id');
            deleteRule(ruleId);
        });
    });
}

// Update alert counts
function updateAlertCounts() {
    const activeAlerts = currentAlerts.filter(alert => alert.status === 'active');
    const criticalAlerts = activeAlerts.filter(alert => alert.severity === 'critical');
    const warningAlerts = activeAlerts.filter(alert => alert.severity === 'warning');
    const infoAlerts = activeAlerts.filter(alert => alert.severity === 'info');
    
    document.getElementById('active-alerts-count').textContent = activeAlerts.length;
    document.getElementById('critical-alerts-count').textContent = criticalAlerts.length;
    document.getElementById('warning-alerts-count').textContent = warningAlerts.length;
    document.getElementById('info-alerts-count').textContent = infoAlerts.length;
}

// Apply filters to the alerts list
function applyFilters() {
    const statusFilter = document.getElementById('status-filter').value;
    const typeFilter = document.getElementById('type-filter').value;
    const severityFilter = document.getElementById('severity-filter').value;
    const searchQuery = document.getElementById('search-input').value.toLowerCase();
    
    let filteredAlerts = [...currentAlerts];
    
    // Apply status filter
    if (statusFilter !== 'all') {
        filteredAlerts = filteredAlerts.filter(alert => alert.status === statusFilter);
    }
    
    // Apply type filter
    if (typeFilter !== 'all') {
        filteredAlerts = filteredAlerts.filter(alert => alert.type === typeFilter);
    }
    
    // Apply severity filter
    if (severityFilter !== 'all') {
        filteredAlerts = filteredAlerts.filter(alert => alert.severity === severityFilter);
    }
    
    // Apply search query
    if (searchQuery) {
        filteredAlerts = filteredAlerts.filter(alert => 
            alert.name.toLowerCase().includes(searchQuery) || 
            alert.description.toLowerCase().includes(searchQuery) ||
            alert.source.toLowerCase().includes(searchQuery)
        );
    }
    
    renderAlerts(filteredAlerts);
}

// Show alert details in modal
function showAlertDetails(alertId) {
    const alert = currentAlerts.find(a => a.id === alertId);
    if (!alert) return;
    
    selectedAlertId = alertId;
    
    const modalTitle = document.getElementById('alert-modal-title');
    const modalBody = document.getElementById('alert-modal-body');
    
    modalTitle.innerHTML = `
        <span class="badge severity ${alert.severity.toLowerCase()}">${alert.severity}</span>
        ${alert.name}
    `;
    
    // Format labels and annotations
    let labelsHtml = '';
    if (alert.labels && Object.keys(alert.labels).length > 0) {
        labelsHtml = '<div class="alert-labels"><h4>Labels</h4><ul>';
        for (const [key, value] of Object.entries(alert.labels)) {
            labelsHtml += `<li><strong>${key}:</strong> ${value}</li>`;
        }
        labelsHtml += '</ul></div>';
    }
    
    let annotationsHtml = '';
    if (alert.annotations && Object.keys(alert.annotations).length > 0) {
        annotationsHtml = '<div class="alert-annotations"><h4>Annotations</h4><ul>';
        for (const [key, value] of Object.entries(alert.annotations)) {
            annotationsHtml += `<li><strong>${key}:</strong> ${value}</li>`;
        }
        annotationsHtml += '</ul></div>';
    }
    
    // Format threshold if available
    let thresholdHtml = '';
    if (alert.threshold !== undefined && alert.value !== undefined) {
        thresholdHtml = `
            <div class="alert-threshold">
                <h4>Threshold</h4>
                <p>Value: ${alert.value} ${alert.comparison_operator} Threshold: ${alert.threshold}</p>
            </div>
        `;
    }
    
    modalBody.innerHTML = `
        <div class="alert-details">
            <p class="alert-description">${alert.description}</p>
            
            <div class="alert-meta-details">
                <div class="meta-item">
                    <strong>Status:</strong> 
                    <span class="badge status ${alert.status.toLowerCase()}">${alert.status}</span>
                </div>
                <div class="meta-item">
                    <strong>Type:</strong> ${alert.type}
                </div>
                <div class="meta-item">
                    <strong>Source:</strong> ${alert.source}
                </div>
                <div class="meta-item">
                    <strong>Created:</strong> ${formatTimestamp(alert.timestamp)}
                </div>
                <div class="meta-item">
                    <strong>Last Updated:</strong> ${formatTimestamp(alert.last_updated)}
                </div>
            </div>
            
            ${thresholdHtml}
            ${labelsHtml}
            ${annotationsHtml}
        </div>
    `;
    
    // Update action buttons based on alert status
    const acknowledgeBtn = document.getElementById('acknowledge-btn');
    const resolveBtn = document.getElementById('resolve-btn');
    const silenceBtn = document.getElementById('silence-btn');
    
    if (alert.status === 'active') {
        acknowledgeBtn.style.display = 'inline-block';
        resolveBtn.style.display = 'inline-block';
        silenceBtn.style.display = 'inline-block';
    } else if (alert.status === 'acknowledged') {
        acknowledgeBtn.style.display = 'none';
        resolveBtn.style.display = 'inline-block';
        silenceBtn.style.display = 'inline-block';
    } else {
        acknowledgeBtn.style.display = 'none';
        resolveBtn.style.display = 'none';
        silenceBtn.style.display = 'none';
    }
    
    // Show the modal
    document.getElementById('alert-modal').style.display = 'block';
}

// Format timestamp for display
function formatTimestamp(timestamp) {
    if (!timestamp) return 'N/A';
    
    const date = new Date(timestamp);
    return date.toLocaleString();
}

// Refresh alerts
function refreshAlerts() {
    document.getElementById('alerts-list').innerHTML = '<div class="loading">Refreshing alerts...</div>';
    fetchAlerts();
    fetchRules();
}

// Show create alert modal
function showCreateAlertModal() {
    // Reset form
    document.getElementById('create-alert-form').reset();
    
    // Show modal
    document.getElementById('create-alert-modal').style.display = 'block';
}

// Hide create alert modal
function hideCreateAlertModal() {
    document.getElementById('create-alert-modal').style.display = 'none';
}

// Show create rule modal
function showCreateRuleModal() {
    // Reset form
    document.getElementById('create-rule-form').reset();
    
    // Show modal
    document.getElementById('create-rule-modal').style.display = 'block';
}

// Hide create rule modal
function hideCreateRuleModal() {
    document.getElementById('create-rule-modal').style.display = 'none';
}

// Hide alert details modal
function hideAlertModal() {
    document.getElementById('alert-modal').style.display = 'none';
    selectedAlertId = null;
}

// Hide all modals
function hideAllModals() {
    document.getElementById('alert-modal').style.display = 'none';
    document.getElementById('create-alert-modal').style.display = 'none';
    document.getElementById('create-rule-modal').style.display = 'none';
    selectedAlertId = null;
}

// Create a new alert
function createAlert() {
    const name = document.getElementById('alert-name').value;
    const description = document.getElementById('alert-description').value;
    const type = document.getElementById('alert-type').value;
    const severity = document.getElementById('alert-severity').value;
    const source = document.getElementById('alert-source').value;
    
    // Parse labels
    const labelsText = document.getElementById('alert-labels').value;
    const labels = {};
    if (labelsText) {
        labelsText.split('\n').forEach(line => {
            const [key, value] = line.split('=').map(s => s.trim());
            if (key && value) {
                labels[key] = value;
            }
        });
    }
    
    // Parse annotations
    const annotationsText = document.getElementById('alert-annotations').value;
    const annotations = {};
    if (annotationsText) {
        annotationsText.split('\n').forEach(line => {
            const [key, value] = line.split('=').map(s => s.trim());
            if (key && value) {
                annotations[key] = value;
            }
        });
    }
    
    // Create alert object
    const alert = {
        name,
        description,
        type,
        severity,
        source,
        labels,
        annotations
    };
    
    // Send to API
    fetch('/api/alerts', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(alert)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to create alert');
        }
        return response.json();
    })
    .then(data => {
        hideCreateAlertModal();
        refreshAlerts();
    })
    .catch(error => {
        console.error('Error creating alert:', error);
        alert(`Failed to create alert: ${error.message}`);
    });
}

// Create a new alert rule
function createRule() {
    const name = document.getElementById('rule-name').value;
    const description = document.getElementById('rule-description').value;
    const metric = document.getElementById('rule-metric').value;
    const threshold = parseFloat(document.getElementById('rule-threshold').value);
    const operator = document.getElementById('rule-operator').value;
    const severity = document.getElementById('rule-severity').value;
    
    // Parse labels
    const labelsText = document.getElementById('rule-labels').value;
    const labels = {};
    if (labelsText) {
        labelsText.split('\n').forEach(line => {
            const [key, value] = line.split('=').map(s => s.trim());
            if (key && value) {
                labels[key] = value;
            }
        });
    }
    
    // Create rule object
    const rule = {
        name,
        description,
        metric,
        threshold,
        comparison_operator: operator,
        severity,
        labels,
        enabled: true
    };
    
    // Send to API
    fetch('/api/alerts/rules', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(rule)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to create rule');
        }
        return response.json();
    })
    .then(data => {
        hideCreateRuleModal();
        fetchRules();
    })
    .catch(error => {
        console.error('Error creating rule:', error);
        alert(`Failed to create rule: ${error.message}`);
    });
}

// Edit an existing rule
function editRule(ruleId) {
    // Implement rule editing functionality
    alert('Edit rule functionality will be implemented in a future update');
}

// Toggle rule enabled status
function toggleRuleStatus(ruleId, enabled) {
    fetch(`/api/alerts/rules/${ruleId}/${enabled ? 'enable' : 'disable'}`, {
        method: 'POST'
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`Failed to ${enabled ? 'enable' : 'disable'} rule`);
        }
        return response.json();
    })
    .then(data => {
        fetchRules();
    })
    .catch(error => {
        console.error(`Error ${enabled ? 'enabling' : 'disabling'} rule:`, error);
        alert(`Failed to ${enabled ? 'enable' : 'disable'} rule: ${error.message}`);
    });
}

// Delete a rule
function deleteRule(ruleId) {
    if (!confirm('Are you sure you want to delete this rule?')) {
        return;
    }
    
    fetch(`/api/alerts/rules/${ruleId}`, {
        method: 'DELETE'
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to delete rule');
        }
        return response.json();
    })
    .then(data => {
        fetchRules();
    })
    .catch(error => {
        console.error('Error deleting rule:', error);
        alert(`Failed to delete rule: ${error.message}`);
    });
}

// Acknowledge an alert
function acknowledgeAlert() {
    if (!selectedAlertId) return;
    
    fetch(`/api/alerts/${selectedAlertId}/acknowledge`, {
        method: 'POST'
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to acknowledge alert');
        }
        return response.json();
    })
    .then(data => {
        hideAlertModal();
        refreshAlerts();
    })
    .catch(error => {
        console.error('Error acknowledging alert:', error);
        alert(`Failed to acknowledge alert: ${error.message}`);
    });
}

// Resolve an alert
function resolveAlert() {
    if (!selectedAlertId) return;
    
    fetch(`/api/alerts/${selectedAlertId}/resolve`, {
        method: 'POST'
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to resolve alert');
        }
        return response.json();
    })
    .then(data => {
        hideAlertModal();
        refreshAlerts();
    })
    .catch(error => {
        console.error('Error resolving alert:', error);
        alert(`Failed to resolve alert: ${error.message}`);
    });
}

// Silence an alert
function silenceAlert() {
    if (!selectedAlertId) return;
    
    const duration = prompt('Enter silence duration in hours:', '24');
    if (!duration) return;
    
    const hours = parseInt(duration);
    if (isNaN(hours) || hours <= 0) {
        alert('Please enter a valid positive number for duration');
        return;
    }
    
    fetch(`/api/alerts/${selectedAlertId}/silence`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            duration_hours: hours,
            reason: 'Silenced from dashboard'
        })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to silence alert');
        }
        return response.json();
    })
    .then(data => {
        hideAlertModal();
        refreshAlerts();
    })
    .catch(error => {
        console.error('Error silencing alert:', error);
        alert(`Failed to silence alert: ${error.message}`);
    });
}
