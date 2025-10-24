# Knetz Governance

**Version:** 1.0  
**Last Updated:** October 24, 2025  
**Status:** Active

---

## Overview

This document defines the governance model for the Knetz project. It outlines how decisions are made, how contributions are accepted, and how the project is managed.

---

## Principles

1. **Open**: Decision-making is transparent and documented
2. **Welcoming**: All contributions are valued and respected
3. **Transparent**: All processes are clearly documented
4. **Merit-based**: Contributions determine influence
5. **Collaborative**: Decisions favor consensus

---

## Roles and Responsibilities

### Contributors

**Who:** Anyone who contributes to the project  
**Examples:** Code, documentation, issue reports, testing, community support

**Responsibilities:**
- Follow the Code of Conduct
- Follow contribution guidelines
- Respect maintainer decisions
- Help other community members

**Rights:**
- Submit issues and pull requests
- Participate in discussions
- Vote on community polls

### Committers

**Who:** Regular contributors with write access  
**Requirements:**
- 5+ merged PRs
- Active for 3+ months
- Demonstrated understanding of codebase
- Maintainer nomination and approval

**Responsibilities:**
- Review pull requests
- Triage issues
- Help onboard new contributors
- Maintain code quality
- Follow project standards

**Rights:**
- Merge pull requests
- Create branches
- Label and assign issues
- Vote on technical decisions

### Maintainers

**Who:** Long-term committers who guide the project  
**Requirements:**
- 20+ merged PRs
- Active for 6+ months
- Significant feature contributions
- Community involvement
- Existing maintainer nomination
- 2/3 maintainer approval

**Responsibilities:**
- Set project direction
- Make final decisions on disputes
- Manage releases
- Represent project externally
- Mentor committers
- Enforce governance
- Security response

**Rights:**
- All committer rights
- Vote on governance changes
- Nominate new committers/maintainers
- Approve releases
- Repository admin access

### Lead Maintainer

**Who:** Maintainer elected by other maintainers  
**Term:** 1 year, renewable  
**Election:** Simple majority vote

**Responsibilities:**
- Final decision on deadlocked votes
- Coordinate maintainer activities
- CNCF liaison (if applicable)
- Community representation
- Conflict resolution
- Emergency decisions

**Rights:**
- All maintainer rights
- Casting vote in ties
- Call emergency meetings

---

## Decision Making

### Consensus Building

1. **Discussion**: Issues and proposals discussed openly
2. **Proposal**: Formal proposal created if needed
3. **Feedback**: Community provides input
4. **Revision**: Proposal refined based on feedback
5. **Consensus**: Agreement reached or vote called

### Voting

**When to Vote:**
- Major architectural changes
- Breaking changes
- Governance changes
- Adding/removing maintainers
- CNCF-related decisions

**Voting Process:**
- Proposal posted in GitHub Discussion
- Minimum 7-day comment period
- Vote announced with 3-day notice
- Votes cast via comments or issue
- Results published publicly

**Voting Rights:**
- Contributors: Advisory votes on features
- Committers: Binding votes on technical decisions
- Maintainers: Binding votes on all decisions
- Lead Maintainer: Casting vote in ties

**Thresholds:**
- Technical decisions: Simple majority (>50%)
- Maintainer changes: Supermajority (>66%)
- Governance changes: Supermajority (>66%)
- Emergency decisions: Lead maintainer discretion

### Lazy Consensus

For minor changes:
- Propose change in PR or issue
- If no objections in 72 hours, proceed
- Any maintainer can request full review
- Applies to: bug fixes, docs, minor features

---

## Contribution Process

### Code Contributions

1. **Issue**: Create issue describing problem/feature
2. **Discussion**: Discuss approach with maintainers
3. **Implementation**: Develop solution following guidelines
4. **Pull Request**: Submit PR with tests and docs
5. **Review**: Address feedback from committers/maintainers
6. **Approval**: Get required approvals
7. **Merge**: Committer merges approved PR

### Required Approvals

- Bug fixes: 1 committer
- New features: 1 maintainer
- Breaking changes: 2 maintainers
- Governance changes: 2/3 maintainers

### Review Guidelines

- Focus on code quality and correctness
- Be respectful and constructive
- Suggest alternatives when rejecting
- Approve when requirements met
- Anyone can review, but approval rights differ

---

## Nomination Process

### Nominating a Committer

1. Maintainer nominates contributor via GitHub Discussion
2. Include evidence of contributions
3. 7-day discussion period
4. Simple majority vote by maintainers
5. Nominee notified and invited
6. Access granted upon acceptance

### Nominating a Maintainer

1. Maintainer nominates committer via GitHub Discussion
2. Include evidence of sustained contributions
3. 14-day discussion period
4. Supermajority (>66%) vote by maintainers
5. Nominee notified and invited
6. Access granted upon acceptance

### Emeritus Status

Maintainers/committers may become emeritus:
- Request emeritus status
- Inactive for 12+ months
- No longer have time to contribute
- Moving to different projects

Emeritus members:
- Listed in MAINTAINERS.md
- Can return with maintainer approval
- Retain honor and recognition

---

## Meetings

### Community Meetings

**Frequency:** Bi-weekly  
**Format:** Video call (Zoom/Google Meet)  
**Agenda:** Published 48 hours before  
**Minutes:** Published within 24 hours  
**Open to:** Everyone

**Purpose:**
- Project updates
- Community discussion
- Feature proposals
- Q&A

### Maintainer Meetings

**Frequency:** Monthly  
**Format:** Video call  
**Agenda:** Circulated 48 hours before  
**Minutes:** Published within 48 hours  
**Open to:** Maintainers (others by invitation)

**Purpose:**
- Strategic planning
- Governance decisions
- Security issues
- Sensitive topics

---

## Conflict Resolution

### Process

1. **Direct Communication**: Parties discuss directly
2. **Mediation**: Maintainer mediates if needed
3. **Maintainer Review**: Maintainers review and decide
4. **Lead Decision**: Lead maintainer makes final call
5. **Code of Conduct**: CoC violations handled separately

### Escalation

- Technical disputes: Maintainer vote
- Behavioral issues: Code of Conduct process
- Governance disputes: Lead maintainer decision
- External conflicts: CNCF mediation (if applicable)

---

## Changes to Governance

### Amendment Process

1. Proposal created in GitHub Discussion
2. Labeled "governance"
3. 14-day comment period
4. Refined based on feedback
5. Vote called with 7-day notice
6. Supermajority (>66%) required
7. Announcement and implementation

### Historical Record

- All governance changes tracked in git
- Discussion links preserved
- Vote results recorded

---

## Communication

### Channels

- **GitHub Issues**: Bug reports, feature requests
- **GitHub Discussions**: General discussion, proposals
- **Pull Requests**: Code review and discussion
- **Slack/Discord** (future): Real-time chat
- **Mailing List** (future): Announcements
- **Monthly Newsletter** (future): Project updates

### Transparency

- All decisions documented publicly
- Meeting minutes published
- Vote results published
- Roadmap publicly visible
- Security issues handled privately initially

---

## Project Assets

### Repository Access

- Public GitHub repository
- Maintainers: Admin access
- Committers: Write access
- Contributors: Read access

### Release Management

- Maintainers create releases
- Automated via CI/CD
- Follows semantic versioning
- Changelog auto-generated
- Release notes reviewed by maintainers

### Security

- Security issues reported privately
- Security team responds within 48 hours
- Fixes developed privately
- Public disclosure after patch
- CVEs assigned as needed

---

## CNCF Considerations

If Knetz joins CNCF:
- Governance must align with CNCF principles
- CNCF Code of Conduct adopted
- Lead maintainer serves as CNCF liaison
- Regular reports to CNCF TOC
- Participate in CNCF community

---

## Current Maintainers

See [MAINTAINERS.md](MAINTAINERS.md) for current list of maintainers and committers.

---

## Credits

This governance model is inspired by:
- CNCF Projects (Kubernetes, Prometheus, etc.)
- Apache Software Foundation
- Cloud Native Computing Foundation guidelines
- Open Source best practices

---

## Questions

Questions about governance? Open a GitHub Discussion with the "governance" label.

---

**Adopted:** October 24, 2025  
**Next Review:** October 24, 2026  
**Contact:** maintainers@knetz.io

