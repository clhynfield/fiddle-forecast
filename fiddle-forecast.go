package main

import (
	"fmt"
	"os"

	"github.com/joefitzgerald/forecast"
)

func main() {
	var folks forecast.People

	fmt.Println("vim-go")
	api := forecast.New(
		"https://api.forecastapp.com",
		os.Getenv("FORECAST_ID"),
		os.Getenv("FORECAST_TOKEN"),
	)

	folks, _ = api.People()
	for _, person := range folks {
		if person.LastName == "Hynfield" {
			fmt.Printf("%v, %v\n", person.LastName, person.FirstName)
			fmt.Printf("roles:\n%v", person.Roles)
		}
	}
}
