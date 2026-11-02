<p align="center">
    <img alt="Cognit Logo" src="/static/logo-mask.svg?v=0.1.0" width="180" />
    <h3 align="center">Cognit</h3>
    <p align="center">One multi-tenant gateway for your AI agents, APIs and MCP servers.</p>
    <p align="center">
        <a href="https://github.com/clivern/cognit/actions/workflows/ci.yml">
            <img alt="CI" src="https://github.com/clivern/cognit/actions/workflows/ci.yml/badge.svg">
        </a>
        <a href="https://github.com/clivern/cognit/releases">
            <img src="https://img.shields.io/badge/Version-v0.1.0-orange.svg">
        </a>
        <a href="https://github.com/clivern/cognit/blob/main/LICENSE">
            <img src="https://img.shields.io/badge/LICENSE-MIT-blue.svg">
        </a>
    </p>
</p>

Cognit is a multi-tenant gateway and service discovery registry for AI agents, APIs and MCP servers. Register your AI agent, let them discover each other, and route every call through one gateway that handles load balancing, health checks and policy.

> Built for the Meta hackathon

### Setup

Prerequisites: Go 1.27, Node 25 (for the web dashboard), plus Postgres and NATS.

To run locally

```bash
# 1. Start Postgres and NATS
make deps

# 2. Copy and Adjust configs
cp config.dist.yml config.test.yml

# 3. Run the app
make run
```

To run the demo agents

```bash
cd demo && ./run.sh
```

To stop the local dependencies later

```bash
make deps-down
```

This stops Postgres and NATS, then `container system stop`.


### Versioning

For transparency into our release cycle and in striving to maintain backward compatibility, Cognit is maintained under the [Semantic Versioning guidelines](https://semver.org/) and release process is predictable and business-friendly.

See the [Releases section of our GitHub project](https://github.com/clivern/cognit/releases) for changelogs for each release version of Cognit. It contains summaries of the most noteworthy changes made in each release. Also see the [Milestones section](https://github.com/clivern/cognit/milestones) for the future roadmap.


### Bug tracker

If you have any suggestions, bug reports, or annoyances please report them to our issue tracker at https://github.com/clivern/cognit/issues


### Security Issues

If you discover a security vulnerability within Cognit, please send an email to [hello@clivern.com](mailto:hello@clivern.com)


### Contributing

We are an open source, community-driven project so please feel free to join us. see the [contributing guidelines](CONTRIBUTING.md) for more details.


### License

© 2026 Clivern. Released under [MIT License](https://opensource.org/licenses/mit-license.php).

**Cognit** is authored and maintained by [@Clivern](http://github.com/Clivern).
