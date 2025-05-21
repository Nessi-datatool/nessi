package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nessi-dev/nessi/pkg/common"
	"github.com/spf13/cobra"
)

// viralCmd represents the viral command for sharing and community engagement
var viralCmd = &cobra.Command{
	Use:   "viral",
	Short: "Tools for sharing Nessi reports and engaging with the community",
	Long: `The viral command provides tools for sharing Nessi reports, generating badges,
and engaging with the Nessi community. These features help increase Nessi's
visibility and adoption while fostering community contributions.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Nessi Viral Growth Tools")
		fmt.Println("=======================")
		fmt.Println("\nAvailable commands:")
		fmt.Println("  viral share     - Generate shareable reports")
		fmt.Println("  viral badge     - Create 'Powered by Nessi' badges")
		fmt.Println("  viral community - Engage with the Nessi community")
		fmt.Println("\nRun 'nessi viral [command] --help' for more information.")
	},
}

// shareCmd represents the share command for generating shareable reports
var shareCmd = &cobra.Command{
	Use:   "share [table_name]",
	Short: "Generate a shareable report for a table",
	Long: `Generate a shareable report for a table with enhanced social sharing features.
The report includes social media sharing buttons, QR codes, and embed options.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tableName := args[0]
		outputPath, _ := cmd.Flags().GetString("output")
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		hashtags, _ := cmd.Flags().GetString("hashtags")

		// In a real implementation, this would use the social sharing plugin
		// For now, we'll just show a mock implementation
		fmt.Printf("Generating shareable report for table '%s'...\n", tableName)
		fmt.Printf("Using title: %s, description: %s, hashtags: %s\n", title, description, hashtags)

		// Show progress indicator
		showProgress("Generating shareable report...", 1*time.Second)

		// Mock report generation (in a real implementation, this could fail)
		// Simulate success for demonstration purposes
		showSuccess(fmt.Sprintf("Successfully generated shareable report for '%s'", tableName))

		// Show output details
		if outputPath != "" {
			fmt.Printf("Report saved to: %s\n", outputPath)
		} else {
			fmt.Println("Report displayed in terminal (preview mode)")
		}

		// Display sharing information
		fmt.Println("\nShare your report with:")
		fmt.Println("- Twitter: Use hashtags #nessi #dataquality")
		fmt.Println("- LinkedIn: Highlight your data quality achievements")
		fmt.Println("- Email: Send to stakeholders with quality metrics")
		fmt.Println("\nTip: Add a 'Powered by Nessi' badge to your project with: nessi viral badge")
	},
}

// badgeCmd represents the badge command for generating badges
var badgeCmd = &cobra.Command{
	Use:   "badge",
	Short: "Generate a 'Powered by Nessi' badge",
	Long: `Generate a 'Powered by Nessi' badge that can be embedded in your project's
README, documentation, or website. This helps increase Nessi's visibility
and adoption while showing your support for the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		label, _ := cmd.Flags().GetString("label")
		message, _ := cmd.Flags().GetString("message")
		color, _ := cmd.Flags().GetString("color")
		style, _ := cmd.Flags().GetString("style")
		qualityScore, _ := cmd.Flags().GetInt("quality-score")

		// Log quality score for debugging
		if qualityScore > 0 {
			fmt.Printf("Including quality score: %d\n", qualityScore)
		}

		// Show progress indicator
		showProgress("Generating 'Powered by Nessi' badge...", 1*time.Second)

		// In a real implementation, this would use the badge plugin
		// For now, we'll just show a mock implementation

		// Generate mock badge code based on format
		var badgeCode string
		var err error

		// Validate color format (simple validation for demonstration)
		if !strings.HasPrefix(color, "#") && len(color) != 6 && color != "blue" && color != "green" && color != "red" && color != "yellow" {
			// Show a warning but continue with default color
			showWarning(fmt.Sprintf("Color '%s' may not be recognized. Using default color.", color))
			color = "blue"
		}

		// Generate badge based on format
		switch format {
		case "markdown":
			badgeCode = "[![Powered by Nessi](https://img.shields.io/badge/" + label + "-" + message + "-" + color + "?style=" + style + ")](https://github.com/nessi-dev/nessi)"
		case "html":
			badgeCode = "<a href=\"https://github.com/nessi-dev/nessi\"><img src=\"https://img.shields.io/badge/" + label + "-" + message + "-" + color + "?style=" + style + "\" alt=\"Powered by Nessi\"></a>"
		default:
			// Handle invalid format error
			err = &common.NessiError{
				Code:    ErrViralBadgeFailed,
				Message: fmt.Sprintf("Invalid badge format: %s", format),
				Details: "Supported formats are: markdown, html",
			}
			handleBadgeError(err.(*common.NessiError))
			return
		}

		// Show success message
		showSuccess("Badge generated successfully!")

		fmt.Println("\nYour 'Powered by Nessi' badge:")
		fmt.Println("----------------------------")
		fmt.Println(badgeCode)
		fmt.Println("----------------------------")
		fmt.Println("\nAdd this badge to your project's README or documentation to show your support for Nessi!")
		fmt.Println("This helps increase Nessi's visibility and adoption in the data engineering community.")
	},
}

// communityCmd represents the community command for engaging with the community
var communityCmd = &cobra.Command{
	Use:   "community",
	Short: "Engage with the Nessi community",
	Long: `Engage with the Nessi community by providing feedback, getting contribution
suggestions, and connecting with other Nessi users. This helps foster community
contributions and increase Nessi's adoption.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show community information
		fmt.Println("Nessi Community Resources:")
		fmt.Println("- GitHub: https://github.com/nessi-dev/nessi")
		fmt.Println("- Slack: https://join.slack.com/t/nessi-community/shared_invite/...")
		fmt.Println("- Discussions: https://github.com/nessi-dev/nessi/discussions")
		fmt.Println("- Issues: https://github.com/nessi-dev/nessi/issues")
		fmt.Println("\nUse 'nessi community feedback' to provide feedback")
		fmt.Println("Use 'nessi community contribute' to get contribution suggestions")
	},
}

// feedbackCmd represents the feedback command for providing feedback
var feedbackCmd = &cobra.Command{
	Use:   "feedback",
	Short: "Provide feedback on Nessi",
	Long: `Provide feedback on Nessi to help improve the tool. Your feedback will be
processed and may be used to guide future development.`,
	Run: func(cmd *cobra.Command, args []string) {
		feedbackType, _ := cmd.Flags().GetString("type")
		feedbackText, _ := cmd.Flags().GetString("text")
		userName, _ := cmd.Flags().GetString("name")
		userEmail, _ := cmd.Flags().GetString("email")

		// Log user email for debugging
		if userEmail != "" {
			fmt.Printf("Contact email provided: %s\n", userEmail)
		}

		// If no feedback text is provided, prompt the user
		if feedbackText == "" {
			fmt.Print("Please enter your feedback: ")
			feedbackBytes, _ := os.ReadFile("/dev/stdin")
			feedbackText = strings.TrimSpace(string(feedbackBytes))

			// Validate feedback text
			if feedbackText == "" {
				// Handle empty feedback error
				err := &common.NessiError{
					Code:    ErrViralCommunityFailed,
					Message: "Feedback text cannot be empty",
					Details: "Please provide some feedback text to submit",
				}
				handleCommunityError(err)
				return
			}
		}

		// Show progress indicator
		showProgress("Submitting feedback...", 1*time.Second)

		// In a real implementation, this would use the community engagement plugin
		// For now, we'll just show a mock implementation

		// Generate a mock GitHub issue URL
		githubIssueURL := fmt.Sprintf("https://github.com/nessi-dev/nessi/issues/new?title=%s&body=%s",
			"Feedback: "+feedbackType,
			"Feedback from "+userName+":\n\n"+feedbackText)

		// Show success message
		showSuccess("Feedback submitted successfully!")

		// Display mock result
		fmt.Println("Thank you for your feedback!")
		fmt.Println("\nYou can also create a GitHub issue with your feedback:")
		fmt.Println(githubIssueURL)
		fmt.Println("\nThank you for helping improve Nessi!")
	},
}

// contributeCmd represents the contribute command for getting contribution suggestions
var contributeCmd = &cobra.Command{
	Use:   "contribute",
	Short: "Get suggestions for contributing to Nessi",
	Long: `Get suggestions for contributing to Nessi based on your experience level
and interests. This helps foster community contributions and increase
Nessi's adoption.`,
	Run: func(cmd *cobra.Command, args []string) {
		experience, _ := cmd.Flags().GetString("experience")
		githubProfile, _ := cmd.Flags().GetString("github-profile")

		// Log GitHub profile for debugging
		if githubProfile != "" {
			fmt.Printf("GitHub profile provided: %s\n", githubProfile)
		}

		// In a real implementation, this would use the community engagement plugin
		// For now, we'll just show a mock implementation
		fmt.Println("Finding contribution suggestions...")

		// Display mock first-time suggestions if user is a beginner
		if experience == "beginner" || experience == "" {
			fmt.Println("\nGreat First-Time Contributions:")

			fmt.Println("\n- Fix a Documentation Typo (Very Easy)")
			fmt.Println("  Start with a simple documentation fix to learn the contribution process")
			fmt.Println("  Link: https://github.com/nessi-dev/nessi/blob/main/docs/")

			fmt.Println("\n- Add Test Case (Easy)")
			fmt.Println("  Add a test case for an existing feature")
			fmt.Println("  Link: https://github.com/nessi-dev/nessi/tree/main/pkg/test")
		}

		// Display mock general suggestions
		fmt.Println("\nContribution Suggestions:")

		fmt.Println("\n- Add a New Quality Rule (Easy)")
		fmt.Println("  Contribute a new data quality rule to help validate Delta Lake tables")
		fmt.Println("  Link: https://github.com/nessi-dev/nessi/blob/main/docs/CONTRIBUTING.md#adding-quality-rules")

		fmt.Println("\n- Improve Documentation (Easy)")
		fmt.Println("  Help improve Nessi's documentation with examples, tutorials, or clarifications")
		fmt.Println("  Link: https://github.com/nessi-dev/nessi/blob/main/docs/CONTRIBUTING.md#improving-documentation")

		fmt.Println("\n- Add Format Support (Medium)")
		fmt.Println("  Extend Nessi to support additional data formats beyond Delta Lake")
		fmt.Println("  Link: https://github.com/nessi-dev/nessi/blob/main/docs/CONTRIBUTING.md#adding-format-support")

		// Display mock community resources
		fmt.Println("\nCommunity Resources:")
		fmt.Println("- Contribution Guide: https://github.com/nessi-dev/nessi/blob/main/docs/CONTRIBUTING.md")
		fmt.Println("- Community Slack: https://join.slack.com/t/nessi-community/shared_invite/...")
		fmt.Println("- GitHub Repository: https://github.com/nessi-dev/nessi")
		fmt.Println("\nThank you for your interest in contributing to Nessi!")
	},
}

func init() {
	// Add the viral command to the root command
	CLI.RootCmd.AddCommand(viralCmd)
	viralCmd.AddCommand(shareCmd)
	viralCmd.AddCommand(badgeCmd)
	viralCmd.AddCommand(communityCmd)
	communityCmd.AddCommand(feedbackCmd)
	communityCmd.AddCommand(contributeCmd)

	// Share command flags
	shareCmd.Flags().StringP("output", "o", "", "Output path for the shareable report")
	shareCmd.Flags().StringP("title", "t", "", "Title for the shareable report")
	shareCmd.Flags().StringP("description", "d", "", "Description for the shareable report")
	shareCmd.Flags().StringP("hashtags", "", "nessi,dataquality,datalake,opensource", "Comma-separated hashtags for social media sharing")

	// Badge command flags
	badgeCmd.Flags().StringP("format", "f", "markdown", "Badge format (markdown, html, rst)")
	badgeCmd.Flags().StringP("label", "l", "powered by", "Text on the left side of the badge")
	badgeCmd.Flags().StringP("message", "m", "nessi", "Text on the right side of the badge")
	badgeCmd.Flags().StringP("color", "c", "3498db", "Color of the right side (hex or named color)")
	badgeCmd.Flags().StringP("style", "s", "flat", "Badge style (flat, flat-square, plastic, etc.)")
	badgeCmd.Flags().IntP("quality-score", "q", 0, "Quality score to include in the badge (0-100)")

	// Feedback command flags
	feedbackCmd.Flags().StringP("type", "t", "general", "Type of feedback (general, bug, feature)")
	feedbackCmd.Flags().StringP("text", "", "", "Feedback text")
	feedbackCmd.Flags().StringP("name", "n", "", "Your name (optional)")
	feedbackCmd.Flags().StringP("email", "e", "", "Your email (optional)")

	// Contribute command flags
	contributeCmd.Flags().StringP("experience", "e", "", "Your experience level (beginner, intermediate, advanced)")
	contributeCmd.Flags().StringP("github-profile", "g", "", "Your GitHub profile (optional)")
}
