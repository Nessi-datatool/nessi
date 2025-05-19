package dashboard

import (
	"net/http"
	
	"github.com/nessi-dev/nessi/pkg/security"
)

// handleDataQualityDashboard handles the data quality dashboard page
func (d *Dashboard) handleDataQualityDashboard(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated when authentication is enabled
	if d.authManager != nil {
		// Check if user is authenticated
		_, ok := security.UserFromContext(r.Context())
		if !ok {
			// Redirect to login page
			http.Redirect(w, r, "/login?redirect=/data-quality", http.StatusFound)
			return
		}
	}

	// Render the data quality dashboard template
	err := d.templates.ExecuteTemplate(w, "data_quality.html", map[string]interface{}{
		"Title":       "Data Quality Dashboard",
		"Description": "Monitor and manage data quality rules and profiles",
	})
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
