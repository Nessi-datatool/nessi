# Nessi Badge Plugin

This plugin generates embeddable "Powered by Nessi" badges for reports and projects, helping to increase Nessi's visibility and adoption through viral growth.

## Features

- **Badge Generation**: Creates professional "Powered by Nessi" badges for embedding in projects
- **Multiple Formats**: Supports Markdown, HTML, and RST formats for maximum compatibility
- **Quality Score Badges**: Automatically generates quality score badges based on report results
- **Report Enhancement**: Adds badge section to reports with easy copy-paste functionality
- **Customization Options**: Customize badge style, colors, and text

## Installation

```bash
# Install the plugin
nessi plugins install badge_plugin
```

## Usage

### Enhancing Reports with Badges

```bash
# Generate a report with badge section
nessi report generate --table my_table --plugin badge_plugin
```

### Generating Standalone Badges

```bash
# Generate a badge in markdown format
nessi badge generate --format markdown

# Generate a badge with custom text
nessi badge generate --label "verified by" --message "nessi" --color blue

# Generate a quality score badge
nessi badge generate --quality-score 95
```

### Available Plugin Commands

| Command | Description |
|---------|-------------|
| `generate_badge` | Creates a badge with customizable options |
| `enhance_report` | Adds badge section to reports |
| `get_badge_markdown` | Returns markdown code for embedding a badge |
| `get_badge_html` | Returns HTML code for embedding a badge |

## Badge Customization Options

| Option | Description | Default |
|--------|-------------|--------|
| `label` | Text on the left side of the badge | "powered by" |
| `message` | Text on the right side of the badge | "nessi" |
| `color` | Color of the right side (hex or named color) | "3498db" (blue) |
| `style` | Badge style (flat, flat-square, plastic, etc.) | "flat" |
| `logo_width` | Width of the logo in pixels | 16 |
| `link_url` | URL when badge is clicked | Nessi GitHub repo |

## Integration with Other Plugins

This plugin works seamlessly with other Nessi plugins:

- **Social Sharing Plugin**: Combines badges with social sharing capabilities
- **Community Engagement Plugin**: Adds badges alongside community features
- **Report Extensions Plugin**: Includes badges in custom report sections

## Example: Adding Badges to Your Project

1. Generate a quality report with the badge plugin:
   ```bash
   nessi report generate --table customer_data --plugin badge_plugin
   ```

2. Copy the badge code from the report in your preferred format (Markdown, HTML, or RST)

3. Add the badge to your project's README, documentation, or website

4. When others see the badge, they'll be directed to Nessi, increasing visibility and adoption

## Benefits for Nessi's Growth

- **Increased Visibility**: Badges in project READMEs expose Nessi to new users
- **Social Proof**: Shows that respected projects are using Nessi
- **Community Building**: Creates a sense of pride in using and supporting Nessi
- **Viral Growth**: Each badge serves as a mini-advertisement for Nessi

## Contributing

We welcome contributions to enhance the badge plugin! Some ideas:

- Additional badge styles and formats
- Animated or interactive badges
- Integration with more badge services beyond shields.io
- Custom badge templates for different use cases

Please see our [contribution guidelines](../../docs/CONTRIBUTING.md) for more information.

## License

This plugin is released under the same license as Nessi (Apache 2.0).
