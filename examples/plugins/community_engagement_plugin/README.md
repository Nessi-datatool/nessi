# Nessi Community Engagement Plugin

This plugin enhances Nessi with features to foster community engagement, increase adoption, and encourage contributions to the project.

## Features

- **Community Section in Reports**: Adds a dedicated community section to data quality reports
- **Feedback Collection**: Provides easy ways for users to submit feedback and suggestions
- **Contribution Guidance**: Suggests ways for users to contribute based on their experience level
- **Community Links**: Integrates links to GitHub, Slack, and other community resources
- **Quick Action Buttons**: One-click access to create issues, feature requests, and bug reports

## Installation

```bash
# Install the plugin
nessi plugins install community_engagement_plugin
```

## Usage

### Enhancing Reports with Community Features

```bash
# Generate a report with community engagement features
nessi report generate --table my_table --plugin community_engagement_plugin
```

### Collecting User Feedback

```bash
# Submit feedback directly from the CLI
nessi community feedback --type "feature_request" --text "It would be great to have..."
```

### Getting Contribution Suggestions

```bash
# Get suggestions for ways to contribute to Nessi
nessi community contribute --experience "beginner"
```

### Available Plugin Commands

| Command | Description |
|---------|-------------|
| `collect_feedback` | Collects and processes user feedback |
| `suggest_contributions` | Provides tailored contribution suggestions |
| `enhance_report` | Adds community engagement elements to reports |

## Integration with Other Plugins

This plugin works seamlessly with other Nessi plugins:

- **Social Sharing Plugin**: Combines community engagement with social sharing capabilities
- **Report Extensions Plugin**: Adds community section alongside other custom report sections
- **Quality Rules Plugin**: Encourages contributions of new quality rules

## Example: Community-Driven Quality Rule Development

1. A user identifies a need for a new quality rule while using Nessi
2. They use the community engagement plugin to submit feedback:
   ```bash
   nessi community feedback --type "quality_rule_idea" --text "We need a rule to validate date formats across different columns"
   ```

3. The plugin provides a direct link to create a GitHub issue with the appropriate template
4. The user receives suggestions for how to contribute the rule themselves
5. After implementing the rule, they can share their contribution using the social sharing plugin

## Benefits for Nessi's Growth

- **Increased Contributions**: Makes it easier for users to contribute back to the project
- **Community Building**: Fosters a sense of community among Nessi users
- **Feedback Loop**: Creates a direct channel for user feedback to improve the tool
- **Viral Growth**: Encourages users to share their work and bring others to the project
- **Documentation Improvement**: Facilitates community-driven documentation enhancements

## Contributing

We welcome contributions to enhance the community engagement capabilities! Some ideas:

- Integration with additional community platforms
- Enhanced contributor recognition features
- Gamification elements to encourage participation
- User testimonial collection and showcase

Please see our [contribution guidelines](../../docs/CONTRIBUTING.md) for more information.

## License

This plugin is released under the same license as Nessi (Apache 2.0).
