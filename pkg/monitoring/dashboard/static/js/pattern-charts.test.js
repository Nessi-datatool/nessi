/**
 * Tests for pattern-charts.js
 * 
 * This file contains unit tests for the pattern frequency chart functionality.
 * Run with: npm test
 */

// Mock the DOM elements
document.body.innerHTML = `
  <div id="pattern-frequency-container"></div>
  <select id="table-selector">
    <option value="table1">Table 1</option>
    <option value="table2">Table 2</option>
  </select>
`;

// Mock Chart.js
global.Chart = class Chart {
  constructor(canvas, config) {
    this.canvas = canvas;
    this.config = config;
    this.type = config.type;
    this.data = config.data;
    this.options = config.options;
  }
};

// Mock fetch API
global.fetch = jest.fn();

// Import the module under test
const fs = require('fs');
const path = require('path');
const patternChartsPath = path.join(__dirname, 'pattern-charts.js');
const patternChartsCode = fs.readFileSync(patternChartsPath, 'utf8');
eval(patternChartsCode);

describe('Pattern Frequency Charts', () => {
  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
    
    // Reset DOM
    document.getElementById('pattern-frequency-container').innerHTML = '';
    
    // Mock successful fetch response
    global.fetch.mockImplementation(() => 
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockPatternData)
      })
    );
  });
  
  // Mock pattern data for testing
  const mockPatternData = {
    columns: [
      {
        name: 'email',
        type: 'string',
        patterns: [
          {
            name: 'email',
            description: 'Email address format',
            frequency: 95.5,
            confidence: 0.955,
            examples: ['user@example.com', 'test@test.com']
          }
        ]
      },
      {
        name: 'id',
        type: 'integer',
        patterns: [
          {
            name: 'integer',
            description: 'Integer number format',
            frequency: 100.0,
            confidence: 1.0,
            examples: ['1', '42', '1000']
          }
        ]
      }
    ]
  };
  
  test('fetchPatternData makes correct API call', async () => {
    // Call the function
    fetchPatternData('test_table');
    
    // Check if fetch was called with the correct URL
    expect(global.fetch).toHaveBeenCalledWith('/api/data-quality/patterns?table_id=test_table');
  });
  
  test('renderPatternCharts creates sections for each column', async () => {
    // Call the function
    renderPatternCharts(mockPatternData);
    
    // Check if container has sections for each column
    const container = document.getElementById('pattern-frequency-container');
    const sections = container.querySelectorAll('.column-pattern-section');
    
    expect(sections.length).toBe(2);
    
    // Check if headers contain column names
    const headers = container.querySelectorAll('.card-header h5');
    expect(headers[0].textContent).toContain('email');
    expect(headers[1].textContent).toContain('id');
  });
  
  test('createColumnPatternSection creates chart and table', () => {
    // Get container
    const container = document.getElementById('pattern-frequency-container');
    
    // Call the function with mock column data
    createColumnPatternSection(container, mockPatternData.columns[0]);
    
    // Check if chart container exists
    const chartContainer = container.querySelector('.pattern-chart-container');
    expect(chartContainer).not.toBeNull();
    
    // Check if pattern details table exists
    const tableContainer = container.querySelector('.pattern-details-container');
    expect(tableContainer).not.toBeNull();
    
    // Check if table has correct data
    const table = tableContainer.querySelector('table');
    const rows = table.querySelectorAll('tbody tr');
    expect(rows.length).toBe(1);
    
    // Check row content
    const cells = rows[0].querySelectorAll('td');
    expect(cells[0].textContent).toBe('email');
    expect(cells[1].textContent).toBe('Email address format');
    expect(cells[2].textContent).toBe('95.50%');
    expect(cells[3].textContent).toBe('user@example.com, test@test.com');
  });
  
  test('renderPatternChart creates chart with correct data', () => {
    // Create a chart container
    const container = document.createElement('div');
    
    // Call the function
    renderPatternChart(container, mockPatternData.columns[0]);
    
    // Check if canvas was created
    const canvas = container.querySelector('canvas');
    expect(canvas).not.toBeNull();
    
    // Check if Chart was instantiated
    expect(Chart).toHaveBeenCalledWith(canvas, expect.any(Object));
    
    // Get the chart instance
    const chartInstance = Chart.mock.instances[0];
    
    // Check chart type
    expect(chartInstance.type).toBe('bar');
    
    // Check chart data
    expect(chartInstance.data.labels).toEqual(['email']);
    expect(chartInstance.data.datasets[0].data).toEqual([95.5]);
  });
  
  test('generateColors returns correct number of colors', () => {
    // Test with small number
    const colors3 = generateColors(3);
    expect(colors3.length).toBe(3);
    
    // Test with large number
    const colors10 = generateColors(10);
    expect(colors10.length).toBe(10);
    
    // Check that all colors are rgba strings
    for (const color of colors10) {
      expect(color).toMatch(/^rgba\(\d+, \d+, \d+, [0-9.]+\)$/);
    }
  });
  
  test('darkenColor darkens color correctly', () => {
    // Test darkening a color
    const original = 'rgba(100, 150, 200, 0.7)';
    const darkened = darkenColor(original, 0.2);
    
    // Check format
    expect(darkened).toMatch(/^rgba\(\d+, \d+, \d+, [0-9.]+\)$/);
    
    // Parse the values
    const originalValues = original.match(/\d+/g).map(Number);
    const darkenedValues = darkened.match(/\d+/g).map(Number);
    
    // Check that RGB values are lower
    expect(darkenedValues[0]).toBeLessThan(originalValues[0]);
    expect(darkenedValues[1]).toBeLessThan(originalValues[1]);
    expect(darkenedValues[2]).toBeLessThan(originalValues[2]);
    
    // Check that alpha is unchanged
    expect(darkenedValues[3]).toBe(originalValues[3]);
  });
  
  test('showErrorMessage displays error correctly', () => {
    // Call the function
    showErrorMessage('Test error message');
    
    // Check if error message is displayed
    const container = document.getElementById('pattern-frequency-container');
    const alert = container.querySelector('.alert-danger');
    
    expect(alert).not.toBeNull();
    expect(alert.textContent).toBe('Test error message');
  });
  
  test('handles empty pattern data', () => {
    // Create empty data
    const emptyData = {
      columns: []
    };
    
    // Call the function
    renderPatternCharts(emptyData);
    
    // Check if info message is displayed
    const container = document.getElementById('pattern-frequency-container');
    const alert = container.querySelector('.alert-info');
    
    expect(alert).not.toBeNull();
    expect(alert.textContent).toContain('No pattern data available');
  });
  
  test('handles column with no patterns', () => {
    // Create column with no patterns
    const columnWithNoPatterns = {
      name: 'empty_column',
      type: 'string',
      patterns: []
    };
    
    // Get container
    const container = document.getElementById('pattern-frequency-container');
    
    // Call the function
    createColumnPatternSection(container, columnWithNoPatterns);
    
    // Check if message is displayed
    const message = container.querySelector('.text-muted');
    
    expect(message).not.toBeNull();
    expect(message.textContent).toContain('No patterns detected');
  });
  
  test('handles fetch error', async () => {
    // Mock fetch error
    global.fetch.mockImplementation(() => 
      Promise.resolve({
        ok: false,
        status: 500
      })
    );
    
    // Call the function
    await fetchPatternData('test_table');
    
    // Check if error message is displayed
    const container = document.getElementById('pattern-frequency-container');
    const alert = container.querySelector('.alert-danger');
    
    expect(alert).not.toBeNull();
    expect(alert.textContent).toContain('Failed to load pattern data');
  });
});
