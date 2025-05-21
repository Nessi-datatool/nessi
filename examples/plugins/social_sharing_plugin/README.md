# Nessi Social Sharing Plugin

This plugin enhances Nessi's data quality reports with advanced social sharing capabilities, making it easier for users to share their reports and increasing Nessi's visibility and adoption.

## Features

- **QR Code Generation**: Automatically generates QR codes for reports, making them easy to share in presentations or printed materials
- **Enhanced Social Sharing**: Adds optimized sharing for Twitter, LinkedIn, Email, and WhatsApp
- **Embed Code**: Provides HTML code for embedding reports in websites and blogs
- **Analytics Tracking**: Optional UTM parameter support for tracking report shares
- **Customizable Hashtags**: Define custom hashtags for social media sharing

## Installation

```bash
# Install the plugin
nessi plugins install social_sharing_plugin
```

## Usage

### Enhancing Reports with Social Sharing

```bash
# Generate a report with enhanced social sharing
nessi report generate --table my_table --plugin social_sharing_plugin
```

### Customizing Social Sharing Options

```bash
# Generate a report with custom sharing options
nessi report generate --table my_table --plugin social_sharing_plugin \
  --plugin-args "title=My Data Quality Report,hashtags=dataquality,analytics,opensource"
```

### Available Plugin Arguments

| Argument | Description | Default |
|----------|-------------|--------|
| `title` | Report title for sharing | Table name from report |
| `description` | Brief description of the report | Auto-generated from quality score |
| `image_url` | Custom image URL for social previews | Nessi logo |
| `hashtags` | Comma-separated hashtags for social media | nessi,dataquality,datalake,opensource |
| `analytics_id` | ID for analytics tracking | None |

## Integration with Other Plugins

This plugin works seamlessly with other Nessi plugins:

- **Report Extensions Plugin**: Adds additional report sections while maintaining social sharing capabilities
- **Quality Rules Plugin**: Custom quality rules are reflected in shareable reports
- **Format Handler Plugin**: Works with any data format supported by Nessi

## Example: Sharing a Report in a Data Team Workflow

1. Generate a quality report with the social sharing plugin:
   ```bash
   nessi report generate --table sales_data --plugin social_sharing_plugin
   ```

2. The generated report includes:
   - QR code for easy access
   - Pre-formatted sharing links for various platforms
   - Embed code for internal dashboards

3. Share the report with stakeholders via:
   - Email (with pre-formatted message)
   - Slack (using the generated Slack-friendly format)
   - Team wiki (using the embed code)

## Contributing

We welcome contributions to enhance the social sharing capabilities! Some ideas:

- Support for additional social platforms
- Custom templates for different sharing contexts
- Integration with URL shorteners
- Enhanced analytics tracking

Please see our [contribution guidelines](../../docs/CONTRIBUTING.md) for more information.

## License

This plugin is released under the same license as Nessi (Apache 2.0).
