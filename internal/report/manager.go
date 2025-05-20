package report

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/jung-kurt/gofpdf"
)

// ReportFormat represents the format of a report
type ReportFormat string

const (
	PDF  ReportFormat = "pdf"
	CSV  ReportFormat = "csv"
	HTML ReportFormat = "html"
	JSON ReportFormat = "json"
)

// Report represents a generated report
type Report struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Format      ReportFormat `json:"format"`
	Content     interface{}  `json:"content"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// ReportTemplate represents a report template
type ReportTemplate struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Format      ReportFormat `json:"format"`
	Template    string       `json:"template"`
	Parameters  []string     `json:"parameters"`
}

// ReportManager handles report generation
type ReportManager struct {
	templates map[string]ReportTemplate
	reports   map[string]Report
	outputDir string
}

// NewManager creates a new report manager
func NewManager(outputDir string) (*ReportManager, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	return &ReportManager{
		templates: make(map[string]ReportTemplate),
		reports:   make(map[string]Report),
		outputDir: outputDir,
	}, nil
}

// AddTemplate adds a new report template
func (m *ReportManager) AddTemplate(template ReportTemplate) error {
	if _, exists := m.templates[template.ID]; exists {
		return fmt.Errorf("template %s already exists", template.ID)
	}

	m.templates[template.ID] = template
	return nil
}

// RemoveTemplate removes a report template
func (m *ReportManager) RemoveTemplate(id string) error {
	if _, exists := m.templates[id]; !exists {
		return fmt.Errorf("template %s does not exist", id)
	}

	delete(m.templates, id)
	return nil
}

// GetTemplates returns all report templates
func (m *ReportManager) GetTemplates() []ReportTemplate {
	templates := make([]ReportTemplate, 0, len(m.templates))
	for _, template := range m.templates {
		templates = append(templates, template)
	}
	return templates
}

// GenerateReport generates a report from a template
func (m *ReportManager) GenerateReport(ctx context.Context, templateID string, parameters map[string]interface{}) (*Report, error) {
	template, exists := m.templates[templateID]
	if !exists {
		return nil, fmt.Errorf("template %s does not exist", templateID)
	}

	// Validate parameters
	for _, param := range template.Parameters {
		if _, ok := parameters[param]; !ok {
			return nil, fmt.Errorf("missing required parameter: %s", param)
		}
	}

	// Generate report based on format
	var content interface{}
	var err error

	switch template.Format {
	case PDF:
		content, err = m.generatePDF(template, parameters)
	case CSV:
		content, err = m.generateCSV(template, parameters)
	case HTML:
		content, err = m.generateHTML(template, parameters)
	case JSON:
		content, err = m.generateJSON(template, parameters)
	default:
		return nil, fmt.Errorf("unsupported report format: %s", template.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate report: %w", err)
	}

	report := &Report{
		ID:          fmt.Sprintf("report_%d", time.Now().UnixNano()),
		Name:        template.Name,
		Description: template.Description,
		Format:      template.Format,
		Content:     content,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.reports[report.ID] = *report
	return report, nil
}

// generatePDF generates a PDF report
func (m *ReportManager) generatePDF(template ReportTemplate, parameters map[string]interface{}) (interface{}, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, template.Name)
	pdf.Ln(20)

	// Add content based on template
	// TODO: Implement template parsing and content generation
	// This would involve parsing the template string and filling in the parameters

	return pdf, nil
}

// generateCSV generates a CSV report
func (m *ReportManager) generateCSV(template ReportTemplate, parameters map[string]interface{}) (interface{}, error) {
	// Check if data parameter exists and is a slice or map
	data, ok := parameters["data"]
	if !ok {
		return nil, fmt.Errorf("missing data parameter for CSV generation")
	}

	// Create a CSV buffer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Handle different data types
	switch d := data.(type) {
	case [][]string:
		// Data is already in CSV format
		for _, row := range d {
			if err := writer.Write(row); err != nil {
				return nil, fmt.Errorf("failed to write CSV row: %w", err)
			}
		}
	case []map[string]interface{}:
		// Data is a slice of maps, extract headers from first item
		if len(d) == 0 {
			return nil, fmt.Errorf("empty data for CSV generation")
		}

		// Extract headers from first map
		headers := make([]string, 0, len(d[0]))
		for k := range d[0] {
			headers = append(headers, k)
		}

		// Write headers
		if err := writer.Write(headers); err != nil {
			return nil, fmt.Errorf("failed to write CSV headers: %w", err)
		}

		// Write data rows
		for _, item := range d {
			row := make([]string, len(headers))
			for i, header := range headers {
				val, ok := item[header]
				if !ok {
					row[i] = ""
				} else {
					row[i] = fmt.Sprintf("%v", val)
				}
			}
			if err := writer.Write(row); err != nil {
				return nil, fmt.Errorf("failed to write CSV row: %w", err)
			}
		}
	default:
		return nil, fmt.Errorf("unsupported data type for CSV generation")
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	return buf.String(), nil
}

// generateHTML generates an HTML report
func (m *ReportManager) generateHTML(template ReportTemplate, parameters map[string]interface{}) (interface{}, error) {
	// Create a template from the template string
	tmpl, err := template.New("report").Parse(template.Template)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template with the parameters
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, parameters); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// generateJSON generates a JSON report
func (m *ReportManager) generateJSON(template ReportTemplate, parameters map[string]interface{}) (interface{}, error) {
	// For JSON reports, we simply return the parameters as JSON
	jsonData, err := json.MarshalIndent(parameters, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(jsonData), nil
}

// SaveReport saves a report to a file
func (m *ReportManager) SaveReport(report *Report) error {
	// Create output file
	filename := filepath.Join(m.outputDir, fmt.Sprintf("%s.%s", report.ID, report.Format))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Write report content to file
	switch report.Format {
	case PDF:
		// Handle PDF content
		if pdf, ok := report.Content.(*gofpdf.Fpdf); ok {
			return pdf.OutputFileAndClose(filename)
		}
		return fmt.Errorf("invalid content type for PDF report")
	case CSV, HTML, JSON:
		// Write string content to file
		if content, ok := report.Content.(string); ok {
			_, err := file.WriteString(content)
			return err
		}
		return fmt.Errorf("invalid content type for %s report", report.Format)
	default:
		return fmt.Errorf("unsupported report format: %s", report.Format)
	}
}

// GetReport returns a report by ID
func (m *ReportManager) GetReport(id string) (*Report, error) {
	report, exists := m.reports[id]
	if !exists {
		return nil, fmt.Errorf("report %s does not exist", id)
	}

	return &report, nil
}

// GetReports returns all reports
func (m *ReportManager) GetReports() []Report {
	reports := make([]Report, 0, len(m.reports))
	for _, report := range m.reports {
		reports = append(reports, report)
	}
	return reports
}

// DeleteReport deletes a report
func (m *ReportManager) DeleteReport(id string) error {
	report, exists := m.reports[id]
	if !exists {
		return fmt.Errorf("report %s does not exist", id)
	}

	// Delete the report file
	filename := filepath.Join(m.outputDir, fmt.Sprintf("%s.%s", report.ID, report.Format))
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete report file: %w", err)
	}

	delete(m.reports, id)
	return nil
}

// GenerateTableReport generates a report from a table
func (m *ReportManager) GenerateTableReport(ctx context.Context, templateID string, record arrow.Record) (*Report, error) {
	template, exists := m.templates[templateID]
	if !exists {
		return nil, fmt.Errorf("template %s does not exist", templateID)
	}

	// Convert record to table data
	table := make([][]string, int(record.NumRows())+1)

	// Add headers
	headers := make([]string, int(record.NumCols()))
	for i := 0; i < int(record.NumCols()); i++ {
		headers[i] = record.ColumnName(i)
	}
	table[0] = headers

	// Add data
	for i := 0; i < int(record.NumRows()); i++ {
		row := make([]string, int(record.NumCols()))
		for j := 0; j < int(record.NumCols()); j++ {
			col := record.Column(j)
			if col.IsNull(i) {
				row[j] = "NULL"
			} else {
				row[j] = col.ValueStr(i)
			}
		}
		table[i+1] = row
	}

	// Generate report based on format
	var content interface{}
	var err error

	switch template.Format {
	case PDF:
		content, err = m.generateTablePDF(template, table)
	case CSV:
		content, err = m.generateTableCSV(template, table)
	case HTML:
		content, err = m.generateTableHTML(template, table)
	case JSON:
		content, err = m.generateTableJSON(template, table)
	default:
		return nil, fmt.Errorf("unsupported report format: %s", template.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate table report: %w", err)
	}

	report := &Report{
		ID:          fmt.Sprintf("report_%d", time.Now().UnixNano()),
		Name:        template.Name,
		Description: template.Description,
		Format:      template.Format,
		Content:     content,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.reports[report.ID] = *report
	return report, nil
}

// generateTablePDF generates a PDF report from a table
func (m *ReportManager) generateTablePDF(template ReportTemplate, table [][]string) (interface{}, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, template.Name)
	pdf.Ln(20)

	// Add table
	pdf.SetFont("Arial", "", 12)
	colWidth := 190.0 / float64(len(table[0]))
	for _, row := range table {
		for _, cell := range row {
			pdf.Cell(colWidth, 10, cell)
		}
		pdf.Ln(10)
	}

	return pdf, nil
}

// generateTableCSV generates a CSV report from a table
func (m *ReportManager) generateTableCSV(template ReportTemplate, table [][]string) (interface{}, error) {
	// Create a CSV buffer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write headers
	if err := writer.Write(table[0]); err != nil {
		return nil, fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// Write data rows
	for _, row := range table[1:] {
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	return buf.String(), nil
}

// generateTableHTML generates an HTML report from a table
func (m *ReportManager) generateTableHTML(template ReportTemplate, table [][]string) (interface{}, error) {
	// Create an HTML buffer
	var buf bytes.Buffer

	// Write HTML header
	buf.WriteString("<html><body><table border='1'>")

	// Write headers
	buf.WriteString("<tr>")
	for _, header := range table[0] {
		buf.WriteString(fmt.Sprintf("<th>%s</th>", header))
	}
	buf.WriteString("</tr>")

	// Write data rows
	for _, row := range table[1:] {
		buf.WriteString("<tr>")
		for _, cell := range row {
			buf.WriteString(fmt.Sprintf("<td>%s</td>", cell))
		}
		buf.WriteString("</tr>")
	}

	// Write HTML footer
	buf.WriteString("</table></body></html>")

	return buf.String(), nil
}

// generateTableJSON generates a JSON report from a table
func (m *ReportManager) generateTableJSON(template ReportTemplate, table [][]string) (interface{}, error) {
	// Create a JSON buffer
	var buf bytes.Buffer

	// Write JSON header
	buf.WriteString("[")

	// Write data rows
	for i, row := range table[1:] {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("{")
		for j, cell := range row {
			buf.WriteString(fmt.Sprintf("\"%s\":\"%s\"", table[0][j], cell))
			if j < len(row)-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString("}")
	}

	// Write JSON footer
	buf.WriteString("]")

	return buf.String(), nil
}
