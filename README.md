# Plyo Hackathon - AI-Powered Research Intelligence Platform

A comprehensive research intelligence platform that leverages AI agents to generate detailed company analysis reports. The platform automatically researches companies across multiple dimensions including company intelligence, competitive analysis, market dynamics, and trend analysis.

Check the live app [here](https://plyo-hackathon.fly.dev/) - assuming the tokens are still valid and have credits!

## Project Write-Up

### The Problem I Tackled

Traditional business research and due diligence processes are time-consuming, fragmented, and often incomplete. Professionals need to manually gather information from multiple sources, analyze competitors, understand market dynamics, and synthesize findings into coherent reports. This process can take days or weeks and often misses critical insights due to the overwhelming amount of available data and the difficulty of conducting comprehensive analysis across multiple domains simultaneously.

### The Solution and How It Works

Our AI-powered research intelligence platform automates and orchestrates comprehensive business research through a multi-agent system. Here's how it works:

1. **Input Processing**: Users provide a candidate name and company URL
2. **Agent Orchestration**: The system deploys four specialized AI research agents in parallel:
   - **Company Intelligence Agent**: Researches company background, financials, leadership, and operations
   - **Competitive Intelligence Agent**: Analyzes direct and indirect competitors, market positioning
   - **Market Dynamics Agent**: Investigates industry trends, market size, growth patterns, and regulatory environment
   - **Trend Analysis Agent**: Identifies emerging trends, technological disruptions, and future outlook
3. **Data Validation**: Each agent's findings are processed through a validation agent to ensure accuracy and relevance
4. **Report Generation**: A final agent synthesizes all research into a comprehensive, professional report
5. **Real-time Tracking**: Users monitor progress through a live dashboard showing completion status of each research domain

### Core Features and Technical Choices

**Core Features:**
- **Parallel Multi-Agent Processing**: Simultaneous research across multiple domains for faster results
- **Real-time Progress Tracking**: Live updates on research completion status
- **Comprehensive Report Generation**: Professional PDF reports combining all research findings
- **Background Job Processing**: Asynchronous execution prevents UI blocking during long research tasks
- **Data Validation Layer**: Automated fact-checking and relevance filtering of research findings

**Technical Choices:**
- **Go + Echo Framework**: High-performance backend with excellent concurrency for handling multiple AI agent requests
- **SQLite + SQLC**: Lightweight database with type-safe query generation for rapid development
- **Templ + TailwindCSS**: Modern frontend stack with server-side rendering for fast, responsive UI
- **GoQite Job Queue**: Reliable background job processing for long-running research tasks
- **Multi-API Integration**: Combines web search (Serper), web scraping (ScrapingBee), and LLM capabilities (OpenAI)

### Why It Matters and Possible Impact

**Business Impact:**
- **Time Efficiency**: Reduces research time from days/weeks to hours
- **Comprehensive Coverage**: Ensures no critical research domain is overlooked
- **Consistency**: Standardized research methodology across all analyses
- **Scalability**: Can handle multiple research requests simultaneously

**Market Applications:**
- **Investment Due Diligence**: VCs and PE firms conducting company evaluations
- **Competitive Analysis**: Businesses understanding their competitive landscape
- **Market Entry Research**: Companies exploring new markets or partnerships
- **Academic Research**: Researchers conducting industry or company studies

**Future Potential:**
- Integration with financial databases for deeper quantitative analysis
- Customizable research templates for different industries
- API access for integration with existing business tools
- Real-time monitoring of tracked companies for ongoing intelligence

### External APIs, Datasets, and Tools

**APIs and External Services:**
- **OpenAI GPT API**: Powers all AI research agents with advanced language understanding and generation
- **Serper API**: Provides web search capabilities for finding current information and news
- **ScrapingBee API**: Enables structured data extraction from company websites and online sources

**Development Tools:**
- **Air**: Live reload development server for efficient development cycles
- **Goose**: Database migration management for schema evolution
- **SQLC**: Generates type-safe Go code from SQL queries
- **Golangci-lint**: Code quality assurance and style consistency
- **Just**: Modern command runner for development workflow automation

**Data Sources** (accessed through APIs):
- Public company websites and press releases
- Industry news and analysis articles
- Market research reports and publications
- Social media and professional networks
- Government and regulatory filings
- Academic and research publications

## Project Description

This application is an AI-powered research platform that automates the generation of comprehensive business intelligence reports. Users input a candidate name and company URL, and the system deploys multiple specialized AI agents to conduct research across different domains:

- **Company Intelligence**: Deep dive into company background, financials, and operations
- **Competitive Intelligence**: Analysis of competitors and market positioning
- **Market Dynamics**: Industry trends, market size, and growth patterns
- **Trend Analysis**: Emerging trends and future outlook

The platform provides real-time progress tracking and generates professional reports combining insights from all research agents.

## How to Use

1. **Setup Environment**:
   ```bash
   # Copy environment file and configure API keys
   cp .env.example .env
   # Edit .env with your API keys (OpenAI, Serper, ScrapingBee)
   ```

2. **Database Setup**:
   ```bash
   # Run database migrations
   just up-migrations
   ```

3. **Development**:
   ```bash
   # Install dependencies and start development server
   go mod tidy
   just run
   ```

4. **Usage**:
   - Navigate to the web interface
   - Enter candidate name and company URL
   - Monitor real-time research progress
   - Download comprehensive PDF reports

## Team Details

**Team Name**: mbvlabs

**Team Members**:
- Morten Bisgaard Vistisen - [LinkedIn](https://linkedin.com/in/mortenvistisen)

## Tech Stack

### Backend
- **Go 1.25** - Primary backend language
- **Echo v4** - Web framework
- **SQLite** - Database with migrations via Goose
- **SQLC** - Type-safe SQL code generation
- **GoQite** - Job queue for background processing

### Frontend
- **Templ** - Go templating engine
- **TailwindCSS** - Utility-first CSS framework
- **Datastar** - Frontend reactivity library
- **HTMX-style** interactions

### AI & External Services
- **OpenAI GPT** - LLM for research agents
- **Serper API** - Web search capabilities
- **ScrapingBee** - Web scraping service

### Development Tools
- **Air** - Live reload development server
- **Just** - Command runner (alternative to Make)
- **Golangci-lint** - Code linting
- **Golines** - Code formatting

## Project Structure

```
├── agents/                 # AI research agents
│   ├── company_intelligence_agent.go
│   ├── competitive_intelligence_agent.go
│   ├── market_dynamics.go
│   ├── trend_analysis.go
│   └── report_generator.go
├── cmd/app/               # Application entry point
├── controllers/           # HTTP request handlers
├── database/             # Database schema and migrations
├── models/               # Data models and database queries
├── providers/            # External service providers (OpenAI)
├── router/               # HTTP routing and middleware
├── tools/                # External API integrations
├── views/                # HTML templates
└── assets/               # Static assets (CSS, JS)
```

## Available Commands

```bash
# Development
just run                   # Start development server with live reload
just live-server          # Backend server only
just live-templ           # Template generation watcher
just live-tailwind        # TailwindCSS watcher

# Database
just new-migration <name>  # Create new migration
just up-migrations        # Run pending migrations

# Code Quality
just vet                  # Run go vet
just golangci             # Run linter
just golines              # Format code
```

## Features

- **Multi-Agent Research System**: Specialized AI agents for different research domains
- **Real-time Progress Tracking**: Live updates on research completion status
- **Background Job Processing**: Asynchronous research execution with job queues
- **Professional Report Generation**: Comprehensive PDF reports with structured insights
- **Web Scraping Integration**: Automated data collection from multiple sources
- **Type-safe Database Operations**: Generated SQL queries with compile-time safety

## Environment Configuration

Required environment variables:
- `OPENAI_API_KEY` - OpenAI API access
- `SERPER_API_KEY` - Serper search API
- `SCRAPINGBEE_API_KEY` - ScrapingBee web scraping
- `SERVER_HOST` - Server host (default: localhost)
- `SERVER_PORT` - Server port (default: 8080)

## Assets and Documentation

- All source code is available in this repository
- Database schema in `database/migrations/`
- API documentation can be found in the route handlers
- Frontend templates in `views/` directory

## Architecture

The application follows a clean architecture pattern with:
- **Controllers** handling HTTP requests/responses
- **Services** containing business logic
- **Models** managing data persistence
- **Agents** implementing AI research capabilities
- **Tools** providing external API integrations

Background jobs enable scalable research processing, while the web interface provides real-time feedback to users.
