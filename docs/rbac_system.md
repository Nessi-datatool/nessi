

## Overview

Nessi.dev implements a comprehensive Role-Based Access Control (RBAC) system that provides fine-grained access management for both the CLI and API interfaces. This document outlines the implementation details and usage of the RBAC system.

## Core Components

### Roles

The RBAC system defines the following built-in roles:

- **Admin**: Full access to all features and resources
- **User**: Access to most features but limited administrative capabilities
- **Viewer**: Read-only access to resources
- **Custom**: User-defined roles with specific permissions

### Permissions

Permissions define what actions can be performed on resources:

- **Read**: Ability to view resources
- **Write**: Ability to create or modify resources
- **Execute**: Ability to run operations
- **Admin**: Ability to perform administrative operations

### Resources

The system controls access to the following resource types:

- **Table**: Delta tables and their metadata
- **Rule**: Validation rules
- **User**: User accounts
- **Webhook**: Webhook configurations
- **Plugin**: Plugin management
- **Config**: System configuration
- **Report**: Generated reports

## Usage

### CLI Commands

```bash
# List all users and their roles
nessi rbac users list

# Assign a role to a user
nessi rbac users assign --user user1 --role admin

# Create a custom role
nessi rbac roles create --name custom-role --permissions "table:read,table:write,rule:read"

# Check access for a user
nessi rbac check-access --user user1 --resource table --permission write

# List all teams
nessi rbac teams list

# Create a team
nessi rbac teams create --name "Data Quality Team" --description "Team responsible for data quality"

# Add a user to a team
nessi rbac teams add-member --team "Data Quality Team" --user user1

# Assign a role to a team
nessi rbac teams assign-role --team "Data Quality Team" --role user
```

### API Endpoints

The RBAC system exposes the following API endpoints:

- `GET /api/v1/rbac/users` - List all users and their roles
- `POST /api/v1/rbac/users/{id}/role` - Assign a role to a user
- `GET /api/v1/rbac/roles` - List all roles
- `POST /api/v1/rbac/roles` - Create a custom role
- `GET /api/v1/rbac/teams` - List all teams
- `POST /api/v1/rbac/teams` - Create a team
- `PUT /api/v1/rbac/teams/{id}/members` - Add a user to a team
- `PUT /api/v1/rbac/teams/{id}/role` - Assign a role to a team

## Configuration

The RBAC system can be configured in the `config/config.yaml` file:

```yaml
security:
  rbac:
    enabled: true
    default_role: viewer
    custom_roles_path: "/path/to/custom/roles"
```

## Implementation Details

### Access Control

The RBAC system uses an access control list (ACL) to determine if a user has permission to access a resource. Each entry in the ACL consists of a role, resource, and permission.

### Team-Based Permissions

Teams provide a way to manage permissions for groups of users. When a user is added to a team, they inherit the team's role. This simplifies permission management for organizations with many users.

### Custom Roles

Custom roles allow for fine-grained control over permissions. A custom role is defined by a set of access control entries that specify what resources the role can access and what actions it can perform.

### Integration with Authentication

The RBAC system integrates with the authentication system to ensure that users are properly authenticated before their permissions are checked.

## Security Considerations

- The RBAC system is designed to be secure by default, with new users assigned the most restrictive role
- All permission checks are performed server-side to prevent client-side bypassing
- Failed permission checks are logged for audit purposes
- The system supports the principle of least privilege by allowing fine-grained permission assignment

## Best Practices

1. **Follow the Principle of Least Privilege**: Assign users the minimum permissions they need to perform their tasks
2. **Use Teams for Group Permissions**: Organize users into teams based on their roles in the organization
3. **Audit Permissions Regularly**: Regularly review user and team permissions to ensure they are appropriate
4. **Document Custom Roles**: Maintain documentation for custom roles to ensure clarity about their permissions

## Examples

### Creating a Data Analyst Role

```bash
nessi rbac roles create --name data-analyst --permissions "table:read,rule:read,report:read,report:write"
```

### Setting Up a Quality Assurance Team

```bash
# Create the team
nessi rbac teams create --name "QA Team" --description "Quality Assurance Team"

# Assign the user role to the team
nessi rbac teams assign-role --team "QA Team" --role user

# Add members to the team
nessi rbac teams add-member --team "QA Team" --user qa1
nessi rbac teams add-member --team "QA Team" --user qa2
```

## Troubleshooting

### Common Issues

1. **Permission Denied Errors**: Ensure the user has the required role and that the role has the necessary permissions
2. **Role Not Found**: Verify that the role exists and is spelled correctly
3. **User Not Found**: Check that the user exists in the system
4. **Team Not Found**: Ensure the team has been created before attempting to add members or assign roles

### Debugging

Use the audit logs to track permission checks and identify why a user might be denied access:

```bash
nessi audit logs --filter "action_type=permission_check,user_id=user1"
```
