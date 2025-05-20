/**
 * Pattern Frequency Charts for Data Quality Dashboard
 * 
 * This module provides visualization components for displaying pattern frequency
 * data in the data quality dashboard.
 */

// Initialize pattern charts when the DOM is fully loaded
document.addEventListener('DOMContentLoaded', function() {
    // Check if we're on the data quality page
    if (document.getElementById('pattern-frequency-container')) {
        initializePatternCharts();
    }
});

/**
 * Initialize pattern frequency charts
 */
function initializePatternCharts() {
    // Fetch pattern data for the selected table
    const tableId = getSelectedTableId();
    if (tableId) {
        fetchPatternData(tableId);
    }

    // Set up event listener for table selection change
    document.getElementById('table-selector').addEventListener('change', function() {
        const tableId = this.value;
        fetchPatternData(tableId);
    });
}

/**
 * Get the currently selected table ID
 * @returns {string} The selected table ID
 */
function getSelectedTableId() {
    const selector = document.getElementById('table-selector');
    return selector ? selector.value : null;
}

/**
 * Fetch pattern data for a specific table
 * @param {string} tableId - The table ID to fetch pattern data for
 */
function fetchPatternData(tableId) {
    fetch(`/api/data-quality/patterns?table_id=${encodeURIComponent(tableId)}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch pattern data');
            }
            return response.json();
        })
        .then(data => {
            renderPatternCharts(data);
        })
        .catch(error => {
            console.error('Error fetching pattern data:', error);
            showErrorMessage('Failed to load pattern data. Please try again later.');
        });
}

/**
 * Render pattern frequency charts
 * @param {Object} data - The pattern data to visualize
 */
function renderPatternCharts(data) {
    const container = document.getElementById('pattern-frequency-container');
    
    // Clear previous charts
    container.innerHTML = '';
    
    if (!data || !data.columns || data.columns.length === 0) {
        container.innerHTML = '<div class="alert alert-info">No pattern data available for this table.</div>';
        return;
    }
    
    // Create a section for each column
    data.columns.forEach(column => {
        createColumnPatternSection(container, column);
    });
}

/**
 * Create a pattern visualization section for a column
 * @param {HTMLElement} container - The container element
 * @param {Object} column - The column data
 */
function createColumnPatternSection(container, column) {
    // Create column section
    const columnSection = document.createElement('div');
    columnSection.className = 'column-pattern-section card mb-4';
    
    // Create header
    const header = document.createElement('div');
    header.className = 'card-header';
    header.innerHTML = `<h5>${column.name} <small class="text-muted">(${column.type})</small></h5>`;
    columnSection.appendChild(header);
    
    // Create body
    const body = document.createElement('div');
    body.className = 'card-body';
    
    // Check if we have pattern data
    if (!column.patterns || column.patterns.length === 0) {
        body.innerHTML = '<p class="text-muted">No patterns detected for this column.</p>';
    } else {
        // Create pattern frequency chart
        const chartContainer = document.createElement('div');
        chartContainer.className = 'pattern-chart-container';
        chartContainer.style.height = '250px';
        body.appendChild(chartContainer);
        
        // Create pattern details table
        const tableContainer = document.createElement('div');
        tableContainer.className = 'pattern-details-container mt-3';
        body.appendChild(tableContainer);
        
        // Render chart and table
        renderPatternChart(chartContainer, column);
        renderPatternTable(tableContainer, column);
    }
    
    columnSection.appendChild(body);
    container.appendChild(columnSection);
}

/**
 * Render a pattern frequency chart for a column
 * @param {HTMLElement} container - The chart container element
 * @param {Object} column - The column data
 */
function renderPatternChart(container, column) {
    // Prepare data for chart
    const labels = column.patterns.map(p => p.name);
    const data = column.patterns.map(p => p.frequency);
    const colors = generateColors(labels.length);
    
    // Create canvas for chart
    const canvas = document.createElement('canvas');
    container.appendChild(canvas);
    
    // Create chart
    new Chart(canvas, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Pattern Frequency',
                data: data,
                backgroundColor: colors,
                borderColor: colors.map(c => darkenColor(c, 0.2)),
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: {
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: 'Frequency (%)'
                    },
                    max: 100
                }
            },
            plugins: {
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            return `${context.raw.toFixed(2)}%`;
                        }
                    }
                },
                legend: {
                    display: false
                }
            }
        }
    });
}

/**
 * Render a pattern details table for a column
 * @param {HTMLElement} container - The table container element
 * @param {Object} column - The column data
 */
function renderPatternTable(container, column) {
    const table = document.createElement('table');
    table.className = 'table table-sm table-striped';
    
    // Create table header
    const thead = document.createElement('thead');
    thead.innerHTML = `
        <tr>
            <th>Pattern</th>
            <th>Description</th>
            <th>Confidence</th>
            <th>Examples</th>
        </tr>
    `;
    table.appendChild(thead);
    
    // Create table body
    const tbody = document.createElement('tbody');
    column.patterns.forEach(pattern => {
        const tr = document.createElement('tr');
        
        // Pattern name
        const tdName = document.createElement('td');
        tdName.textContent = pattern.name;
        tr.appendChild(tdName);
        
        // Pattern description
        const tdDesc = document.createElement('td');
        tdDesc.textContent = pattern.description || 'N/A';
        tr.appendChild(tdDesc);
        
        // Pattern confidence
        const tdConf = document.createElement('td');
        tdConf.textContent = `${(pattern.confidence * 100).toFixed(2)}%`;
        tr.appendChild(tdConf);
        
        // Pattern examples
        const tdExamples = document.createElement('td');
        if (pattern.examples && pattern.examples.length > 0) {
            tdExamples.textContent = pattern.examples.join(', ');
        } else {
            tdExamples.textContent = 'N/A';
        }
        tr.appendChild(tdExamples);
        
        tbody.appendChild(tr);
    });
    table.appendChild(tbody);
    
    container.appendChild(table);
}

/**
 * Generate an array of colors for chart elements
 * @param {number} count - The number of colors to generate
 * @returns {string[]} Array of color strings
 */
function generateColors(count) {
    const baseColors = [
        'rgba(54, 162, 235, 0.7)',   // Blue
        'rgba(255, 99, 132, 0.7)',   // Red
        'rgba(75, 192, 192, 0.7)',   // Green
        'rgba(255, 159, 64, 0.7)',   // Orange
        'rgba(153, 102, 255, 0.7)',  // Purple
        'rgba(255, 205, 86, 0.7)',   // Yellow
        'rgba(201, 203, 207, 0.7)'   // Grey
    ];
    
    // If we have fewer colors than needed, generate more
    if (count <= baseColors.length) {
        return baseColors.slice(0, count);
    }
    
    // Generate additional colors by varying opacity
    const colors = [...baseColors];
    const neededExtra = count - baseColors.length;
    
    for (let i = 0; i < neededExtra; i++) {
        const baseColor = baseColors[i % baseColors.length];
        const opacity = 0.4 + (0.3 * (i / neededExtra));
        colors.push(baseColor.replace(/[\d.]+\)$/, `${opacity})`));
    }
    
    return colors;
}

/**
 * Darken a color by a specified amount
 * @param {string} color - The color to darken
 * @param {number} amount - The amount to darken (0-1)
 * @returns {string} The darkened color
 */
function darkenColor(color, amount) {
    // Extract rgba values
    const rgba = color.match(/[\d.]+/g);
    if (!rgba || rgba.length < 4) return color;
    
    // Create darkened version
    const r = Math.max(0, parseInt(rgba[0]) - Math.round(parseInt(rgba[0]) * amount));
    const g = Math.max(0, parseInt(rgba[1]) - Math.round(parseInt(rgba[1]) * amount));
    const b = Math.max(0, parseInt(rgba[2]) - Math.round(parseInt(rgba[2]) * amount));
    
    return `rgba(${r}, ${g}, ${b}, ${rgba[3]})`;
}

/**
 * Show an error message in the UI
 * @param {string} message - The error message to display
 */
function showErrorMessage(message) {
    const container = document.getElementById('pattern-frequency-container');
    container.innerHTML = `<div class="alert alert-danger">${message}</div>`;
}
