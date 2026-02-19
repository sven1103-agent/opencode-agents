# OpenCode AI Agents Configuration Repository

This repository serves as a **central configuration hub** for AI agents that conform to the [OpenCode AI configuration schema](https://opencode.ai/config.json). It provides a single source of truth for agent definitions, allowing multiple sites and applications to share consistent agent configurations.

## 📋 What This Repository Contains

- **`opencode.json`**: The main configuration file defining AI agents with their capabilities, models, and permissions
- **Schema Compliance**: All configurations follow the OpenCode AI schema standards
- **Version Control**: Track changes to agent configurations over time
- **Centralized Management**: Single location for maintaining agent definitions

## 🤖 Available Agents

The current configuration defines **eight specialized AI agents** organized into primary agents and subagents across three functional domains.

---

## 🎯 Primary Agents (Entry Points)

These agents serve as main entry points for different types of work:

### 1. **Coding Boss** (`coding-boss`)
- **Purpose**: Routes coding tasks to junior/senior/architect engineers based on complexity and risk
- **Mode**: Primary (main entry point for coding tasks)
- **Model**: Claude Haiku 4.5 (2025-10-01)
- **Permissions**: Read-only (routing and analysis only)
- **Classification Logic**:
  - **JUNIOR**: Local changes, low risk, ≤2 files, straightforward work
  - **SENIOR**: Debugging, non-trivial refactors, multi-module tests, medium risk
  - **ARCHITECT**: Architecture/API changes, security-sensitive, multi-module, migrations, high risk
- **Best For**: Automatically routing coding requests to the appropriate engineering level

### 2. **Documentation Router** (`docs`)
- **Purpose**: Routes documentation requests to the appropriate specialized documentation subagent
- **Mode**: Primary (main entry point for documentation tasks)
- **Model**: Claude Haiku 4.5 (2025-10-01)
- **Permissions**: Read-only (routing and analysis only)
- **Best For**: Determining which documentation agent should handle specific requests

---

## 👨‍💼 Engineering Subagents

These agents handle different levels of coding complexity and implementation:

### 3. **Junior Engineer** (`junior`)
- **Purpose**: Quick, cheap implementation of small scoped tasks with minimal risk
- **Mode**: Subagent (works under coding-boss routing)
- **Model**: Claude Haiku 4.5 (2024-01-01)
- **Permissions**: Full write/edit access
- **Guidelines**: Makes minimal, safe changes; prefers small diffs; adds tests when relevant
- **Escalation**: Stops and recommends escalation to @senior if task expands beyond 2 files or discovers architectural risk
- **Best For**: Bug fixes, small features, straightforward implementations

### 4. **Senior Engineer** (`senior`)
- **Purpose**: Robust solutions, refactors, debugging, and code reviews with attention to maintainability
- **Mode**: Subagent (works under coding-boss routing)
- **Model**: Claude Sonnet 4.6
- **Permissions**: Full write/edit access
- **Focus**: Correctness, maintainability, comprehensive tests, tradeoff explanations
- **Escalation**: Recommends escalation to @architect for public API changes, security boundaries, or cross-service contracts
- **Best For**: Complex bug fixes, refactors, multi-module changes, architectural reviews

### 5. **Architect** (`architect`)
- **Purpose**: High-stakes design and cross-cutting system changes
- **Mode**: Subagent (works under coding-boss routing)
- **Model**: Claude Sonnet 4.6
- **Permissions**: Read-only (design and planning focus)
- **Approach**: Proposes plan first, then implements in safe steps with tests and rollback notes
- **Focus**: System design, contracts, migration strategy, security, long-term maintainability
- **Best For**: Architecture redesigns, security-sensitive changes, multi-service coordination, major migrations

### 6. **Code Reviewer** (`code-reviewer`)
- **Purpose**: Reviews code for best practices, security, performance, and maintainability
- **Mode**: Subagent (works independently or under routing agents)
- **Model**: Claude Sonnet 4.6
- **Permissions**: Read-only (cannot write or edit files)
- **Best For**: Code quality assessments, security audits, performance reviews, quality gates

---

## 📚 Documentation Subagents

These agents specialize in creating different types of documentation:

### 7. **User Guide Writer** (`user-guide`)
- **Purpose**: Creates end-user documentation like READMEs, tutorials, and usage guides
- **Mode**: Subagent (works under docs routing agent)
- **Model**: Claude Haiku 4.5 (2025-10-01)
- **Permissions**: Full write/edit access
- **Focus**: Clear, practical documentation avoiding implementation details
- **Best For**: README files, getting started guides, user tutorials, API documentation, usage examples

### 8. **Agent Architect** (`agent-architect`)
- **Purpose**: Analyzes codebases and designs AGENTS.md files for multi-agent workflows
- **Mode**: Subagent (works under docs routing agent)
- **Model**: Claude Sonnet 4.6
- **Permissions**: Read-only (design and analysis focus)
- **Focus**: Agent role definition, delegation patterns, interaction guidelines, system design
- **Best For**: Designing AGENTS.md documentation, multi-agent system architecture, workflow planning

---

## 🚀 Getting Started

### For Users
1. **Reference the Configuration**: Point your OpenCode AI compatible application to this repository
2. **Use Agent Names**: Reference agents by their key names (e.g., `coding-boss`, `docs`, `code-reviewer`)
3. **Follow Schema**: Ensure your implementation supports the OpenCode AI configuration schema
4. **Select Appropriate Agent**:
   - Use `coding-boss` for coding tasks (it will route to the right engineering level)
   - Use `docs` for documentation tasks (it will route to appropriate doc writers)
   - Use `code-reviewer` directly for standalone code reviews

### For Developers
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/sven1103-agent/opencode-agents.git
   cd opencode-agents
   ```

2. **Validate Configuration**:
   ```bash
   # Ensure your JSON is valid
   cat opencode.json | jq .
   ```

3. **Use in Your Application**:
   ```javascript
   // Example usage in a Node.js application
   const config = require('./opencode.json');
   const codingBoss = config.agent['coding-boss'];
   const docsRouter = config.agent['docs'];
   ```

## 📁 Configuration Schema

Each agent in the configuration follows this structure:

```json
{
  "agent-name": {
    "description": "What this agent does",
    "mode": "primary|subagent",
    "model": "anthropic/claude-model-version",
    "prompt": "System prompt for the agent",
    "tools": {
      "write": true|false,
      "edit": true|false
    }
  }
}
```

### Key Properties Explained

- **`description`**: Human-readable explanation of the agent's purpose
- **`mode`**: 
  - `primary`: Can handle requests directly as entry points
  - `subagent`: Works under other agents or routing systems
- **`model`**: The AI model to use (following OpenCode AI model naming conventions)
  - Claude Haiku 4.5: Lightweight model for routing, analysis, and simple tasks
  - Claude Haiku 4: Lightweight model for junior-level implementations
  - Claude Sonnet 4.6: Powerful model for complex tasks and architecture
- **`prompt`**: System-level instructions that define the agent's behavior and decision-making
- **`tools`**: Permission settings for what the agent can do
  - `write`: Can create new files
  - `edit`: Can modify existing files

## 🔧 Customization

### Adding New Agents
1. Edit `opencode.json`
2. Add your new agent following the schema structure
3. Decide if it should be a `primary` (entry point) or `subagent` (called by others)
4. Commit and push your changes
5. Update dependent applications to use the new agent

### Modifying Existing Agents
1. Update the relevant agent configuration
2. Test thoroughly in your development environment
3. Document any breaking changes in commit messages
4. Deploy updates to consuming applications

## 🌐 Integration Examples

### Web Applications
```html
<!-- Load configuration dynamically -->
<script>
fetch('./opencode.json')
  .then(response => response.json())
  .then(config => {
    const codingBoss = config.agent['coding-boss'];
    const docsRouter = config.agent['docs'];
    // Initialize agents with configuration
  });
</script>
```

### Python Applications
```python
import json

# Load agent configuration
with open('opencode.json', 'r') as f:
    config = json.load(f)

# Access specific agents
coding_boss = config['agent']['coding-boss']
docs_router = config['agent']['docs']
junior_eng = config['agent']['junior']
```

### API Integration
```bash
# Use as a remote configuration source
curl -s https://raw.githubusercontent.com/yourusername/opencode-agents/main/opencode.json | jq '.agent["coding-boss"]'
```

## 📖 Use Cases

This configuration repository is ideal for:

- **Multi-site Deployments**: Share agent configurations across multiple applications
- **Team Collaboration**: Ensure consistent AI-assisted development workflows across team projects
- **Automatic Task Routing**: Let the coding-boss automatically assign work to appropriate engineering levels
- **Version Management**: Track evolution of AI agent capabilities and routing logic
- **A/B Testing**: Easily switch between different agent models or prompts
- **Standardization**: Maintain consistent AI interactions and quality standards across platforms

## 🔒 Security Considerations

- **Tool Permissions**: Agents have specific read/write permissions - respect these in your implementation
- **Model Selection**: Different models have different capabilities and cost implications
  - Haiku models: Cost-effective for routing and simple tasks
  - Sonnet models: Recommended for complex implementations and architectural decisions
- **Prompt Engineering**: System prompts define agent behavior - review them carefully before deploying
- **Access Control**: Consider who can modify this central configuration
- **Escalation Chains**: Respect the built-in escalation logic (junior → senior → architect)

## 📈 Monitoring and Analytics

Consider tracking:
- Agent usage patterns and routing decisions
- Performance metrics per agent type and model
- Configuration change impact on team workflows
- User satisfaction with agent responses
- Cost implications of model selection

## 🤝 Contributing

1. **Fork** this repository
2. **Create** a feature branch for your changes
3. **Test** your configuration thoroughly
4. **Submit** a pull request with clear documentation
5. **Collaborate** on reviews and improvements

### Contribution Guidelines
- Follow the OpenCode AI schema strictly
- Provide clear descriptions for new agents
- Test configurations before submitting
- Document any breaking changes
- Consider backward compatibility
- Update this README if you add or modify agents

## 📄 License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). See the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [OpenCode AI Schema](https://opencode.ai/config.json)

---

**Need help?** Open an issue or check the [OpenCode AI documentation](https://opencode.ai) for more information about implementing and using AI agent configurations.
