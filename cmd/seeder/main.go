package main

import (
	"flag"
	"log"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/bootstrap"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/seeder"
)

func main() {
	// Define --force flag
	force := flag.Bool("force", false, "Force reseed: wipes existing data before seeding")
	flag.Parse()

	// Bootstrap application container
	appContainer := bootstrap.Bootstrap()

	// Run seeder with force flag
	if err := seeder.RunSeeder(appContainer, *force); err != nil {
		log.Fatalf("Seeder failed: %v", err)
	}

	log.Println("Seeder completed successfully ✅")
}
