# Pull Request

## Description

Brief description of what this PR accomplishes.

Fixes #(issue number) <!-- If this PR fixes an issue, link it here -->

## Type of Change

Please delete options that are not relevant:

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Code refactoring (no functional changes)
- [ ] Performance improvement
- [ ] CI/CD changes
- [ ] Other (please describe):

## Changes Made

Detailed list of changes made in this PR:

-
-
-

## Testing

### How Has This Been Tested?

Please describe the tests that you ran to verify your changes:

- [ ] Unit tests pass (`make test`)
- [ ] Build completes successfully (`make build`)
- [ ] Manual testing performed
- [ ] Tested with OpenAI provider
- [ ] Tested with Gemini provider
- [ ] Tested error scenarios
- [ ] Tested configuration changes

### Test Configuration

- **OS**: [e.g., macOS, Linux, Windows]
- **Go version**: [e.g., 1.24.7]
- **Git version**: [e.g., 2.42.0]
- **AI Provider(s) tested**: [e.g., OpenAI, Gemini, both]

### Test Scenarios

Describe the test scenarios you executed:

```bash
# Example test commands run
make check
./institutionalized --help
./institutionalized version
./institutionalized commit --dry-run
# etc.
```

## Configuration Changes

If this PR introduces configuration changes:

- [ ] New configuration options are documented
- [ ] Default values are sensible
- [ ] Backwards compatibility is maintained
- [ ] Configuration migration is handled (if needed)

## Breaking Changes

If this PR contains breaking changes, please describe:

- What functionality is affected
- How users can migrate their setup
- Alternative approaches considered

## Documentation

- [ ] I have updated the README.md if needed
- [ ] I have updated CONTRIBUTING.md if needed
- [ ] I have added/updated code comments for complex logic
- [ ] I have updated configuration examples if needed

## Code Quality

- [ ] My code follows the project's style guidelines
- [ ] I have performed a self-review of my code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings or errors
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes

## Dependencies

- [ ] No new dependencies added
- [ ] New dependencies are justified and documented
- [ ] Dependencies are pinned to specific versions
- [ ] License compatibility verified for new dependencies

## Security Considerations

If applicable, describe any security implications:

- [ ] No sensitive information is logged or exposed
- [ ] Input validation is implemented where needed
- [ ] API keys and secrets are handled securely
- [ ] No security vulnerabilities introduced

## Performance Impact

If applicable, describe performance implications:

- [ ] No performance regression
- [ ] Performance improvements measured and documented
- [ ] Memory usage considered
- [ ] API rate limiting respected

## Checklist

- [ ] I have read the [CONTRIBUTING.md](../CONTRIBUTING.md) guidelines
- [ ] I have followed the conventional commit format for my commits
- [ ] I have tested my changes thoroughly
- [ ] I have updated documentation as needed
- [ ] I have considered backwards compatibility
- [ ] I have added appropriate tests
- [ ] All CI checks are passing

## Additional Notes

Any additional information, concerns, or areas where you'd like specific feedback:

<!--
Examples of things to mention:
- Limitations or known issues
- Future improvements planned
- Questions about implementation decisions
- Areas where you'd like code review focus
-->

## Screenshots (if applicable)

If your changes affect the CLI output or user interface, please include before/after screenshots:

### Before

<!-- Screenshot or command output -->

### After

<!-- Screenshot or command output -->

---

**Reviewer Notes**:

<!-- Space for reviewers to add comments during review -->
