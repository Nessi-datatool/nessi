## Description
Fix numpy version requirement in setup.py to use the correct version (1.24.3) instead of the non-existent 2.0.2 version.

## Related Issue
N/A - This is a maintenance fix

## Type of Change
- [x] Bug fix (non-breaking change which fixes an issue)

## How Has This Been Tested?
- Verified numpy 1.24.3 exists and is installable
- Checked compatibility with other dependencies
- Confirmed version matches requirements.txt

## Checklist:
- [x] My code follows the style guidelines of this project
- [x] I have performed a self-review of my own code
- [x] I have commented my code, particularly in hard-to-understand areas
- [x] I have made corresponding changes to the documentation
- [x] My changes generate no new warnings
- [x] I have added tests that prove my fix is effective or that my feature works
- [x] New and existing unit tests pass locally with my changes
- [x] Any dependent changes have been merged and published in downstream modules
- [x] I have checked my code and corrected any misspellings

## Screenshots (if appropriate):

## Additional Notes:
This fix ensures proper installation of the package by using the correct numpy version that exists in PyPI. 