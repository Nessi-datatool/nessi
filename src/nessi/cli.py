import click
from .license import LicenseValidator

@click.group()
def cli():
    """Nessi Data Tool CLI"""
    pass

@cli.command()
def version():
    """Show the current version"""
    from . import __version__
    click.echo(f"Nessi version {__version__}")

@cli.command()
def status():
    """Check license status"""
    validator = LicenseValidator()
    status = validator.get_status()
    click.echo(status["message"])

@cli.command()
@click.argument('license_key')
def activate(license_key):
    """Activate Nessi with a license key"""
    validator = LicenseValidator()
    if validator.validate_license_key(license_key):
        click.echo("License activated successfully")
    else:
        click.echo("Invalid license key")

if __name__ == '__main__':
    cli() 