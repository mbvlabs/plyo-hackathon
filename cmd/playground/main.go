package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mbvlabs/plyo-hackathon/agents"
	"github.com/mbvlabs/plyo-hackathon/config"
	"github.com/mbvlabs/plyo-hackathon/providers"
	"github.com/mbvlabs/plyo-hackathon/tools"
)

func main() {
	ctx := context.Background()
	serper := tools.NewSerper(config.App.SerperAPIkey)
	openai := providers.NewClient(config.App.OpenAPIKey)

	// Create research agent
	agent := agents.NewPreliminaryResearch(
		openai,
		map[string]tools.Tooler{serper.GetName(): &serper},
	)

	// Test the agent
	companyName := "hercule"

	fmt.Printf("Researching company: %s\n\n", companyName)

	result, err := agent.Research(ctx, companyName)
	if err != nil {
		log.Fatalf("Research failed: %v", err)
	}

	fmt.Println("Research Results:")
	fmt.Println("=================")
	fmt.Println(result)
}
