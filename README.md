# OpenCode AI Agents Configuration Repository

This repository serves as a **central configuration hub** for AI agents that conform to the [OpenCode AI configuration schema](https://opencode.ai/config.json). It provides a single source of truth for agent definitions, allowing multiple sites and applications to share consistent agent configurations.

## 📋 What This Repository Contains

- **`opencode.json`**: The main configuration file defining AI agents with their capabilities, models, and permissions
- **Schema Compliance**: All configurations follow the OpenCode AI schema standards
- **Version Control**: Track changes to agent configurations over time
- **Centralized Management**: Single location for maintaining agent definitions

## 🤖 Available Agents

The current configuration defines three specialized AI agents:

### 1. **Code Reviewer** (`code-reviewer`)
- **Purpose**: Reviews code for best practices, security, performance, and maintainability
- **Mode**: Subagent (works under other agents)
- **Model**: Claude Sonnet 4 (2025-05-14)
- **Permissions**: Read-only (cannot write or edit files)
- **Best For**: Code quality assessments, security audits, performance reviews

### 2. **Documentation Router** (`docs`)
- **Purpose**: Routes documentation requests to the appropriate specialized agents
- **Mode**: Primary (can be the main agent handling requests)
- **Model**: Claude Sonnet 4 (2025-05-14)
- **Permissions**: Read-only (analysis and routing only)
- **Best For**: Determining which documentation agent should handle specific requests

### 3. **User Guide Writer** (`user-guide`)
- **Purpose**: Creates end-user documentation like READMEs, tutorials, and usage guides
- **Mode**: Subagent (works under the docs router)
- **Model**: Claude Haiku 4 (2025-10-01)
- **Permissions**: Full write/edit access
- **Best For**: README files, getting started guides, user tutorials, API documentation

## 🚀 Getting Started

### For Users
1. **Reference the Configuration**: Point your OpenCode AI compatible application to this repository
2. **Use Agent Names**: Reference agents by their key names (`code-reviewer`, `docs`, `user-guide`)
3. **Follow Schema**: Ensure your implementation supports the OpenCode AI configuration schema

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
   const codeReviewer = config.agent['code-reviewer'];
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
  - `primary`: Can handle requests directly
  - `subagent`: Works under other agents or routing systems
- **`model`**: The AI model to use (following OpenCode AI model naming conventions)
- **`prompt`**: System-level instructions that define the agent's behavior
- **`tools`**: Permission settings for what the agent can do
  - `write`: Can create new files
  - `edit`: Can modify existing files

## 🔧 Customization

### Adding New Agents
1. Edit `opencode.json`
2. Add your new agent following the schema structure
3. Commit and push your changes
4. Update dependent applications to use the new agent

### Modifying Existing Agents
1. Update the relevant agent configuration
2. Test thoroughly in your development environment
3. Document any breaking changes
4. Deploy updates to consuming applications

## 🌐 Integration Examples

### Web Applications
```html
<!-- Load configuration dynamically -->
<script>
fetch('./opencode.json')
  .then(response => response.json())
  .then(config => {
    const userGuideAgent = config.agent['user-guide'];
    // Initialize agent with configuration
  });
</script>
```

### Python Applications
```python
import json

# Load agent configuration
with open('opencode.json', 'r') as f:
    config = json.load(f)

# Access specific agent
docs_agent = config['agent']['docs']
model = docs_agent['model']
prompt = docs_agent['prompt']
```

### API Integration
```bash
# Use as a remote configuration source
curl -s https://raw.githubusercontent.com/yourusername/opencode-agents/main/opencode.json | jq '.agent["code-reviewer"]'
```

## 📖 Use Cases

This configuration repository is ideal for:

- **Multi-site Deployments**: Share agent configurations across multiple applications
- **Team Collaboration**: Ensure consistent AI behavior across team projects
- **Version Management**: Track evolution of AI agent capabilities
- **A/B Testing**: Easily switch between different agent configurations
- **Standardization**: Maintain consistent AI interactions across platforms

## 🔒 Security Considerations

- **Tool Permissions**: Agents have specific read/write permissions - respect these in your implementation
- **Model Selection**: Different models have different capabilities and cost implications
- **Prompt Engineering**: System prompts define agent behavior - review them carefully
- **Access Control**: Consider who can modify this central configuration

## 📈 Monitoring and Analytics

Consider tracking:
- Agent usage patterns
- Performance metrics per agent type
- Configuration change impact
- User satisfaction with agent responses

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

## 📄 License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). See the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [OpenCode AI Schema](https://opencode.ai/config.json)

---

**Need help?** Open an issue or check the [OpenCode AI documentation](https://opencode.ai) for more information about implementing and using AI agent configurations.
