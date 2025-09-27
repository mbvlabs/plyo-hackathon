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
	serperScrape := tools.NewSerperScrape(config.App.SerperAPIkey)
	scrapingBee := tools.NewScrapingBee(config.App.ScrapingBeeAPIKey)
	openai := providers.NewClient(config.App.OpenAPIKey)

	toolsMap := map[string]tools.Tooler{
		serper.GetName():       &serper,
		serperScrape.GetName(): &serperScrape,
		scrapingBee.GetName():  &scrapingBee,
	}

	// Create research agent
	agent := agents.NewCompetitiveIntelligence(openai, toolsMap)
	validator := agents.NewDataValidation(openai, toolsMap)

	// Test the agent
	companyName := "kfund"
	companyURL := "https://www.kfund.vc/"

	fmt.Printf("Researching company: %s\n\n", companyName)

	result, err := agent.Research(ctx, companyName, companyURL)
	if err != nil {
		log.Fatalf("Research failed: %v", err)
	}

	valiResult, err := validator.Research(ctx, companyName, companyURL, result)
	if err != nil {
		log.Fatalf("Research failed: %v", err)
	}

	fmt.Println("Research Results:")
	fmt.Println("=================")
	fmt.Println(result)
	fmt.Println("=================")
	fmt.Println(valiResult)
}
