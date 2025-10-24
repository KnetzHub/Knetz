# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1   | :x:                |

## Reporting a Vulnerability

The Knetz team and community take security bugs seriously. We appreciate your efforts to responsibly disclose your findings, and will make every effort to acknowledge your contributions.

### Where to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to:

**security@knetz.io**

If you prefer encrypted communication, you can use our PGP key (available upon request).

### What to Include

To help us better understand the nature and scope of the possible issue, please include as much of the following information as possible:

* Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
* Full paths of source file(s) related to the manifestation of the issue
* The location of the affected source code (tag/branch/commit or direct URL)
* Any special configuration required to reproduce the issue
* Step-by-step instructions to reproduce the issue
* Proof-of-concept or exploit code (if possible)
* Impact of the issue, including how an attacker might exploit the issue

### Response Timeline

* **Acknowledgment**: We will acknowledge receipt of your vulnerability report within 48 hours
* **Initial Assessment**: We will send an initial assessment within 7 days
* **Status Updates**: We will keep you informed of the progress towards a fix
* **Resolution**: We aim to resolve critical vulnerabilities within 90 days

### What to Expect

After you submit a report:

1. **Confirmation**: We'll confirm receipt and begin investigation
2. **Assessment**: We'll assess the severity and impact
3. **Fix Development**: We'll develop a fix in a private repository
4. **Notification**: We'll notify you when the fix is ready
5. **Public Disclosure**: We'll coordinate disclosure timing with you
6. **Credit**: We'll credit you in the security advisory (if you wish)

## Security Update Process

### For Critical Vulnerabilities (CVSS >= 7.0)

1. Private fix development
2. Security advisory created (GitHub Security Advisory)
3. CVE requested if applicable
4. Coordinated disclosure with reporter
5. Patch release published
6. Public advisory released
7. Community notification

### For Non-Critical Vulnerabilities (CVSS < 7.0)

1. Fix developed in private or public depending on severity
2. Regular release cycle
3. Security fix noted in changelog
4. Optional security advisory

## Security Best Practices for Users

### Running Knetz Securely

* Always use the latest version
* Run with minimal required permissions
* Use RBAC to limit cluster access
* Store credentials securely (never in code)
* Enable audit logging
* Monitor for suspicious activity
* Review security advisories regularly

### Cluster Access

* Use service accounts with limited permissions
* Implement network policies
* Enable pod security policies
* Use secure kubeconfig management
* Rotate credentials regularly
* Monitor API access logs

### Configuration Security

* Protect configuration files (600 permissions)
* Use secrets management (Vault, Sealed Secrets)
* Avoid hardcoding credentials
* Use environment variables for sensitive data
* Implement least privilege access
* Enable TLS for all connections

## Security Features

### Current Security Features

* **Authentication**: Kubernetes RBAC integration
* **Authorization**: Namespace and cluster scoping
* **Audit Logging**: Coming soon
* **Encryption**: TLS for API communication
* **Secrets Management**: Integration with K8s secrets

### Planned Security Features

* Enhanced audit logging
* SPIFFE/SPIRE integration
* OPA policy enforcement
* Sigstore integration for supply chain security
* mTLS support
* Key rotation capabilities

## Security Scanning

### Our Scanning Process

We employ the following security scanning tools:

* **Gosec**: Static analysis for Go code
* **Trivy** (planned): Container and dependency scanning
* **Snyk** (planned): Dependency vulnerability scanning
* **CodeQL** (planned): Semantic code analysis
* **Dependabot**: Automated dependency updates

### Dependency Management

* Regular dependency updates
* Automated vulnerability scanning
* SBOM (Software Bill of Materials) generation
* License compliance checking

## Vulnerability Disclosure Policy

### Coordinated Disclosure

We follow coordinated disclosure principles:

* 90-day disclosure deadline after initial report
* Earlier disclosure if fix is publicly available
* Earlier disclosure if actively exploited
* Credit to reporter (if desired)
* CVE assignment for qualifying vulnerabilities

### Public Communication

* Security advisories published in GitHub Security Advisories
* Critical vulnerabilities announced via multiple channels
* Patch releases clearly marked as security releases
* Detailed fix information in changelogs

## Bug Bounty Program

We currently do not have a formal bug bounty program, but we greatly appreciate responsible disclosure and will:

* Acknowledge your contribution publicly (if you wish)
* List you in our security hall of fame
* Provide swag/merchandise when available
* Consider a bounty program as the project matures

## Security Team

### Current Security Team

See [MAINTAINERS.md](MAINTAINERS.md) for current security team members.

### Responsibilities

* Respond to security reports
* Coordinate vulnerability remediation
* Manage security advisories
* Conduct security reviews
* Update security documentation

## Compliance

### Standards

We strive to comply with:

* **OWASP Top 10**: Web application security risks
* **CIS Benchmarks**: Configuration security
* **NIST Cybersecurity Framework**: Security practices
* **Cloud Native Security**: CNCF security best practices

### Audits

* Internal security reviews quarterly
* External security audit (planned for Sandbox)
* Penetration testing (planned for Incubation)

## Security Resources

### For Developers

* [Secure Coding Guidelines](docs/DEVELOPER_GUIDE.md)
* [Security Testing Guide](docs/TESTING.md)
* [Dependency Management](docs/DEPENDENCIES.md)

### For Users

* [Deployment Security](docs/DEPLOYMENT.md)
* [Configuration Security](docs/CONFIGURATION.md)
* [Troubleshooting Security Issues](docs/TROUBLESHOOTING.md)

## Past Security Advisories

No security advisories have been published yet.

Security advisories will be published at:
https://github.com/knetz-io/knetz/security/advisories

## Contact

* **Security Issues**: security@knetz.io
* **General Questions**: maintainers@knetz.io
* **Community**: See [README.md](README.md) for community channels

## Acknowledgments

We thank the following security researchers for responsibly disclosing vulnerabilities:

(No disclosures yet)

---

**Last Updated**: October 24, 2025  
**Version**: 1.0

