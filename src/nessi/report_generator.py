"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 nessi.dev. All rights reserved.

nessi.dev is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of nessi.dev will require a license for all users.
"""

from typing import Dict, Any, List
import json
import os
from datetime import datetime
import pandas as pd
import plotly.graph_objects as go
import plotly.express as px
import numpy as np
from plotly.subplots import make_subplots

class ReportGenerator:
    """Generates reports from scan results."""
    
    def __init__(self):
        """Initialize the report generator."""
        self.template_dir = os.path.join(os.path.dirname(__file__), 'templates')
        
    def generate_report(self, data: Dict, template_name: str = 'default', 
                       output_format: str = 'html') -> str:
        """Generate a report with advanced visualizations.
        
        Args:
            data: Data to include in the report
            template_name: Name of the template to use
            output_format: Output format (html or pdf)
            
        Returns:
            Generated report as a string
        """
        # Generate visualizations
        visualizations = self._generate_visualizations(data)
        
        # Load template
        template = self._load_template(template_name)
        
        # Render template with data and visualizations
        report = template.render(
            data=data,
            visualizations=visualizations,
            timestamp=datetime.now().isoformat()
        )
        
        # Convert to PDF if requested
        if output_format == 'pdf':
            return self._convert_to_pdf(report)
        
        return report

    def _generate_visualizations(self, data: Dict) -> Dict:
        """Generate advanced visualizations from data."""
        visualizations = {}
        
        # Generate partition heatmap
        if 'partition_data' in data:
            visualizations['partition_heatmap'] = self._generate_partition_heatmap(
                data['partition_data']
            )
        
        # Generate value distributions
        if 'value_data' in data:
            visualizations['value_distributions'] = self._generate_value_distributions(
                data['value_data']
            )
        
        # Generate time series analysis
        if 'time_series_data' in data:
            visualizations['time_series'] = self._generate_time_series(
                data['time_series_data']
            )
        
        # Generate correlation matrix
        if 'correlation_data' in data:
            visualizations['correlation_matrix'] = self._generate_correlation_matrix(
                data['correlation_data']
            )
        
        return visualizations

    def _generate_partition_heatmap(self, data: Dict) -> str:
        """Generate a partition heatmap visualization."""
        df = pd.DataFrame(data)
        
        # Create heatmap
        fig = go.Figure(data=go.Heatmap(
            z=df.values,
            x=df.columns,
            y=df.index,
            colorscale='Viridis',
            showscale=True
        ))
        
        # Update layout
        fig.update_layout(
            title='Partition Distribution Heatmap',
            xaxis_title='Partition Values',
            yaxis_title='Partition Keys',
            height=600
        )
        
        return self._plot_to_html(fig)

    def _generate_value_distributions(self, data: Dict) -> str:
        """Generate value distribution visualizations."""
        df = pd.DataFrame(data)
        
        # Create subplots
        fig = make_subplots(rows=len(df.columns), cols=1)
        
        for i, col in enumerate(df.columns, 1):
            # Add histogram
            fig.add_trace(
                go.Histogram(x=df[col], name=col),
                row=i, col=1
            )
            
            # Add box plot
            fig.add_trace(
                go.Box(x=df[col], name=col),
                row=i, col=1
            )
        
        # Update layout
        fig.update_layout(
            title='Value Distributions',
            height=300 * len(df.columns),
            showlegend=False
        )
        
        return self._plot_to_html(fig)

    def _generate_time_series(self, data: Dict) -> str:
        """Generate time series visualizations."""
        df = pd.DataFrame(data)
        df['timestamp'] = pd.to_datetime(df['timestamp'])
        
        # Create time series plot
        fig = go.Figure()
        
        for col in df.columns:
            if col != 'timestamp':
                fig.add_trace(
                    go.Scatter(
                        x=df['timestamp'],
                        y=df[col],
                        name=col,
                        mode='lines+markers'
                    )
                )
        
        # Add trend lines
        for col in df.columns:
            if col != 'timestamp':
                z = np.polyfit(range(len(df)), df[col], 1)
                p = np.poly1d(z)
                fig.add_trace(
                    go.Scatter(
                        x=df['timestamp'],
                        y=p(range(len(df))),
                        name=f'{col} Trend',
                        line=dict(dash='dash')
                    )
                )
        
        # Update layout
        fig.update_layout(
            title='Time Series Analysis',
            xaxis_title='Time',
            yaxis_title='Value',
            height=600
        )
        
        return self._plot_to_html(fig)

    def _generate_correlation_matrix(self, data: Dict) -> str:
        """Generate correlation matrix visualization."""
        df = pd.DataFrame(data)
        corr = df.corr()
        
        # Create correlation heatmap
        fig = go.Figure(data=go.Heatmap(
            z=corr.values,
            x=corr.columns,
            y=corr.index,
            colorscale='RdBu',
            zmin=-1,
            zmax=1,
            showscale=True
        ))
        
        # Update layout
        fig.update_layout(
            title='Correlation Matrix',
            height=600
        )
        
        return self._plot_to_html(fig)

    def _plot_to_html(self, fig: go.Figure) -> str:
        """Convert Plotly figure to HTML."""
        return fig.to_html(full_html=False, include_plotlyjs='cdn')

    def _convert_to_pdf(self, html: str) -> str:
        """Convert HTML report to PDF."""
        # This is a placeholder - implement actual PDF conversion
        return html

    def _load_template(self, template_name: str) -> str:
        """Load a template from the template directory."""
        template_path = os.path.join(self.template_dir, f'{template_name}.html')
        with open(template_path, 'r') as f:
            template = f.read()
        return template

    def _generate_html_report(self, results: Dict[str, Any]) -> str:
        """Generate an HTML report."""
        template_path = os.path.join(self.template_dir, 'report_template.html')
        with open(template_path, 'r') as f:
            template = f.read()
            
        # Format the template with results
        report = template.format(
            timestamp=datetime.now().isoformat(),
            results=json.dumps(results, indent=2)
        )
        
        return report
        
    def _generate_json_report(self, results: Dict[str, Any]) -> str:
        """Generate a JSON report."""
        return json.dumps(results, indent=2)
        
    def _generate_text_report(self, results: Dict[str, Any]) -> str:
        """Generate a text report."""
        report = []
        report.append(f"Report generated at: {datetime.now().isoformat()}")
        report.append("\nScan Results:")
        report.append(json.dumps(results, indent=2))
        return "\n".join(report)

    def customize_report(self, template_name: str, customizations: Dict) -> None:
        """Customize report template with user-defined options.
        
        Args:
            template_name: Name of the template to customize
            customizations: Dictionary of customization options
        """
        template_path = os.path.join(self.template_dir, f'{template_name}.html')
        
        # Load template
        with open(template_path, 'r') as f:
            template_content = f.read()
        
        # Apply customizations
        if 'styles' in customizations:
            template_content = self._apply_styles(template_content, customizations['styles'])
        
        if 'layout' in customizations:
            template_content = self._apply_layout(template_content, customizations['layout'])
        
        if 'filters' in customizations:
            template_content = self._apply_filters(template_content, customizations['filters'])
        
        # Save customized template
        custom_template_path = os.path.join(
            self.template_dir,
            f'{template_name}_custom.html'
        )
        with open(custom_template_path, 'w') as f:
            f.write(template_content)

    def export_report(self, report: str, format: str, options: Dict = None) -> str:
        """Export report in various formats with customization options.
        
        Args:
            report: Report content to export
            format: Export format (pdf, docx, xlsx, csv)
            options: Export options
            
        Returns:
            Path to exported file
        """
        if format == 'pdf':
            return self._export_pdf(report, options or {})
        elif format == 'docx':
            return self._export_docx(report, options or {})
        elif format == 'xlsx':
            return self._export_xlsx(report, options or {})
        elif format == 'csv':
            return self._export_csv(report, options or {})
        else:
            raise ValueError(f"Unsupported export format: {format}")

    def _apply_styles(self, template: str, styles: Dict) -> str:
        """Apply custom styles to template."""
        # Add custom CSS
        css = []
        for selector, properties in styles.items():
            css.append(f"{selector} {{")
            for prop, value in properties.items():
                css.append(f"    {prop}: {value};")
            css.append("}")
        
        # Insert styles in head
        style_tag = f"<style>\n{' '.join(css)}\n</style>"
        return template.replace("</head>", f"{style_tag}\n</head>")

    def _apply_layout(self, template: str, layout: Dict) -> str:
        """Apply custom layout to template."""
        # Modify template structure based on layout
        if 'sections' in layout:
            for section in layout['sections']:
                template = template.replace(
                    f"<!-- {section['name']} -->",
                    self._build_section(section)
                )
        
        return template

    def _apply_filters(self, template: str, filters: Dict) -> str:
        """Apply custom filters to template."""
        # Add custom Jinja2 filters
        for name, func in filters.items():
            self.env.filters[name] = func
        
        return template

    def _export_pdf(self, report: str, options: Dict) -> str:
        """Export report to PDF."""
        from weasyprint import HTML, CSS
        
        # Create PDF
        html = HTML(string=report)
        css = CSS(string=options.get('css', ''))
        
        # Generate PDF
        pdf_path = os.path.join(
            self.output_dir,
            f"report_{datetime.now().strftime('%Y%m%d_%H%M%S')}.pdf"
        )
        html.write_pdf(pdf_path, stylesheets=[css])
        
        return pdf_path

    def _export_docx(self, report: str, options: Dict) -> str:
        """Export report to DOCX."""
        from docx import Document
        
        # Create document
        doc = Document()
        
        # Add content
        for line in report.split('\n'):
            if line.strip():
                doc.add_paragraph(line)
        
        # Save document
        docx_path = os.path.join(
            self.output_dir,
            f"report_{datetime.now().strftime('%Y%m%d_%H%M%S')}.docx"
        )
        doc.save(docx_path)
        
        return docx_path

    def _export_xlsx(self, report: str, options: Dict) -> str:
        """Export report to XLSX."""
        import pandas as pd
        
        # Convert report to DataFrame
        data = []
        for line in report.split('\n'):
            if line.strip():
                data.append(line.split('\t'))
        
        df = pd.DataFrame(data)
        
        # Save to Excel
        xlsx_path = os.path.join(
            self.output_dir,
            f"report_{datetime.now().strftime('%Y%m%d_%H%M%S')}.xlsx"
        )
        df.to_excel(xlsx_path, index=False)
        
        return xlsx_path

    def _export_csv(self, report: str, options: Dict) -> str:
        """Export report to CSV."""
        import pandas as pd
        
        # Convert report to DataFrame
        data = []
        for line in report.split('\n'):
            if line.strip():
                data.append(line.split('\t'))
        
        df = pd.DataFrame(data)
        
        # Save to CSV
        csv_path = os.path.join(
            self.output_dir,
            f"report_{datetime.now().strftime('%Y%m%d_%H%M%S')}.csv"
        )
        df.to_csv(csv_path, index=False)
        
        return csv_path

    def _build_section(self, section: Dict) -> str:
        """Build HTML section from configuration."""
        html = []
        html.append(f"<div class='section {section.get('class', '')}'>")
        
        if 'title' in section:
            html.append(f"<h2>{section['title']}</h2>")
        
        if 'content' in section:
            html.append(f"<div class='content'>{section['content']}</div>")
        
        html.append("</div>")
        return '\n'.join(html)

    def create_interactive_dashboard(self, data: Dict[str, Any], 
                                   dashboard_config: Dict[str, Any]) -> Dict[str, Any]:
        """Create an interactive dashboard from data.
        
        Args:
            data: Data to visualize
            dashboard_config: Dashboard configuration
            
        Returns:
            Dictionary with dashboard creation results
        """
        try:
            # Create dashboard layout
            layout = self._create_dashboard_layout(dashboard_config)
            
            # Generate visualizations
            visualizations = self._generate_visualizations(data, dashboard_config)
            
            # Combine into dashboard
            dashboard = self._combine_dashboard(layout, visualizations)
            
            return {
                "created": True,
                "dashboard": dashboard,
                "config": dashboard_config
            }
            
        except Exception as e:
            logger.error(f"Error creating interactive dashboard: {str(e)}")
            raise

    def _create_dashboard_layout(self, config: Dict[str, Any]) -> Dict[str, Any]:
        """Create dashboard layout from configuration."""
        try:
            layout = {
                "title": config.get("title", "Interactive Dashboard"),
                "panels": []
            }
            
            for panel in config.get("panels", []):
                layout["panels"].append({
                    "title": panel.get("title"),
                    "type": panel.get("type", "graph"),
                    "position": panel.get("position", {}),
                    "size": panel.get("size", {"width": 12, "height": 8})
                })
            
            return layout
            
        except Exception as e:
            logger.error(f"Error creating dashboard layout: {str(e)}")
            raise

    def _combine_dashboard(self, layout: Dict[str, Any], 
                          visualizations: Dict[str, Any]) -> Dict[str, Any]:
        """Combine layout and visualizations into dashboard."""
        try:
            dashboard = layout.copy()
            
            for panel in dashboard["panels"]:
                panel_id = panel["title"].lower().replace(" ", "_")
                if panel_id in visualizations:
                    panel["visualization"] = visualizations[panel_id]
            
            return dashboard
            
        except Exception as e:
            logger.error(f"Error combining dashboard: {str(e)}")
            raise

    def version_report(self, report_id: str, version_notes: str) -> Dict[str, Any]:
        """Create a new version of a report.
        
        Args:
            report_id: ID of the report to version
            version_notes: Notes about the version changes
            
        Returns:
            Dictionary with versioning results
        """
        try:
            # Get current report
            current_report = self._get_report(report_id)
            
            # Create new version
            new_version = {
                "id": f"{report_id}_v{len(current_report['versions']) + 1}",
                "timestamp": datetime.now().isoformat(),
                "notes": version_notes,
                "content": current_report["content"]
            }
            
            # Update report versions
            current_report["versions"].append(new_version)
            self._update_report(report_id, current_report)
            
            return {
                "versioned": True,
                "report_id": report_id,
                "new_version": new_version["id"]
            }
            
        except Exception as e:
            logger.error(f"Error versioning report: {str(e)}")
            raise

    def _get_report(self, report_id: str) -> Dict[str, Any]:
        """Get report from storage."""
        try:
            # In a real implementation, this would fetch from a database
            # For now, return a mock report
            return {
                "id": report_id,
                "content": {},
                "versions": []
            }
            
        except Exception as e:
            logger.error(f"Error getting report: {str(e)}")
            raise

    def _update_report(self, report_id: str, report: Dict[str, Any]) -> None:
        """Update report in storage."""
        try:
            # In a real implementation, this would update a database
            pass
            
        except Exception as e:
            logger.error(f"Error updating report: {str(e)}")
            raise

    def get_report_versions(self, report_id: str) -> List[Dict[str, Any]]:
        """Get all versions of a report.
        
        Args:
            report_id: ID of the report
            
        Returns:
            List of report versions
        """
        try:
            report = self._get_report(report_id)
            return report["versions"]
            
        except Exception as e:
            logger.error(f"Error getting report versions: {str(e)}")
            raise

    def compare_versions(self, report_id: str, version1: str, 
                        version2: str) -> Dict[str, Any]:
        """Compare two versions of a report.
        
        Args:
            report_id: ID of the report
            version1: First version to compare
            version2: Second version to compare
            
        Returns:
            Dictionary with comparison results
        """
        try:
            report = self._get_report(report_id)
            
            v1 = next(v for v in report["versions"] if v["id"] == version1)
            v2 = next(v for v in report["versions"] if v["id"] == version2)
            
            # Compare content
            differences = self._compare_content(v1["content"], v2["content"])
            
            return {
                "compared": True,
                "report_id": report_id,
                "versions": [version1, version2],
                "differences": differences
            }
            
        except Exception as e:
            logger.error(f"Error comparing versions: {str(e)}")
            raise

    def _compare_content(self, content1: Dict[str, Any], 
                        content2: Dict[str, Any]) -> Dict[str, Any]:
        """Compare content of two report versions."""
        try:
            differences = {
                "added": [],
                "removed": [],
                "modified": []
            }
            
            # Compare keys
            keys1 = set(content1.keys())
            keys2 = set(content2.keys())
            
            # Find added keys
            for key in keys2 - keys1:
                differences["added"].append({
                    "key": key,
                    "value": content2[key]
                })
            
            # Find removed keys
            for key in keys1 - keys2:
                differences["removed"].append({
                    "key": key,
                    "value": content1[key]
                })
            
            # Find modified values
            for key in keys1 & keys2:
                if content1[key] != content2[key]:
                    differences["modified"].append({
                        "key": key,
                        "old_value": content1[key],
                        "new_value": content2[key]
                    })
            
            return differences
            
        except Exception as e:
            logger.error(f"Error comparing content: {str(e)}")
            raise 