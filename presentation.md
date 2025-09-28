# AI-Powered Research Intelligence Platform
## Plyo Hackathon

**Team**: mbvlabs
**Developer**: Morten Bisgaard Vistisen

---

## The Problem

Traditional business research is:
- **Time-consuming** (days to weeks)
- **Fragmented** across multiple sources
- **Incomplete** due to information overload
- **Manual** and error-prone

*Professionals struggle to conduct comprehensive due diligence efficiently*

---

## The Solution

AI-powered research platform with **multi-agent orchestration**:

1. **Input**: Company name or URL
2. **Deploy**: 4 specialized AI agents in parallel
3. **Validate**: Automated fact-checking layer
4. **Generate**: Professional comprehensive reports
5. **Monitor**: Real-time progress tracking

---

## Multi-Agent Architecture

**4 Specialized Research Agents:**

- **Company Intelligence**: Financials, leadership, operations
- **Competitive Intelligence**: Market positioning, competitors
- **Market Dynamics**: Industry trends, growth patterns
- **Trend Analysis**: Emerging trends, future outlook

*All running in parallel for maximum efficiency*

---

## Technical Stack

**Backend:**
- Go 1.25 + Echo framework
- SQLite + SQLC (type-safe queries)
- GoQite job queue

**Frontend:**
- Templ + TailwindCSS
- Datastar for reactivity

**AI & APIs:**
- OpenAI GPT for research agents
- Serper API for web search
- ScrapingBee for web scraping

---

## Key Features

**Parallel Processing** - Multiple agents working simultaneously
**Real-time Tracking** - Live progress dashboard
**Background Jobs** - Non-blocking research execution
**Data Validation** - Automated fact-checking
**Professional Reports** - Comprehensive PDF generation
**Type Safety** - Generated SQL with compile-time safety

---

## Live Demo

**[plyo-hackathon.fly.dev](https://plyo-hackathon.fly.dev/)**

*Experience the platform in action!*

---

## Business Impact

**Time Efficiency**: Research time reduced from **weeks to hours**

**Use Cases:**
- Investment due diligence (VCs, PE firms)
- Competitive analysis for businesses
- Market entry research
- Academic industry studies

**Scalability**: Handle multiple research requests simultaneously

---

## Technical Highlights

**Architecture:**
```
├── agents/        # AI research agents
├── controllers/   # HTTP handlers
├── models/        # Data persistence
├── tools/         # External API integrations
├── views/         # Frontend templates
└── database/      # Schema & migrations
```

**Development Tools:** Air, Goose, SQLC, Golangci-lint, Just

---

## Future Potential

**Enhanced Integration:**
- API access for business tools
- Real-time company monitoring
- More data (better) sources
- Deeper validation

---

## Thank You!

**Questions?**

**Live Demo**: [plyo-hackathon.fly.dev](https://plyo-hackathon.fly.dev/)
**Source Code**: Available in this repository
**LinkedIn**: [linkedin.com/in/mortenvistisen](https://linkedin.com/in/mortenvistisen)

*Transforming business research with AI-powered intelligence*
