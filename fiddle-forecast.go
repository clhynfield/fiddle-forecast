package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joefitzgerald/forecast"
)

var api *forecast.API
var debug bool

func init() {
	api = forecast.New(
		"https://api.forecastapp.com",
		os.Getenv("FORECAST_ID"),
		os.Getenv("FORECAST_TOKEN"),
	)
}

func debugMsg(message string) {
	if debug {
		fmt.Println("DEBUG: ", message)
	}
}

func rolesInclude(roles []string, rolesRequired []string) bool {
	for _, roleRequired := range rolesRequired {
		found := false
		for _, role := range roles {
			if role == roleRequired {
				found = true
			}
		}
		if found == false {
			return false
		}
	}
	return true
}

func platformPeopleIDs() []int {
	people, _ := api.People()
	var ids []int
	platformAmerRoles := []string{"PCFS", "amer", "billable"}

	for _, person := range people {
		if person.Archived {
			continue
		}
		personRoles := person.Roles
		if rolesInclude(personRoles, platformAmerRoles) {
			ids = append(ids, person.ID)
		}
	}

	return ids
}

func printUnassignedPlatformPeople() {
	ids := unassignedPlatformPeopleIDs()

	for _, id := range ids {
		person, _ := api.Person(id)
		fmt.Println(person.FirstName, person.LastName)
	}
}

func remove(s []int, i int) []int {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func unassignedPlatformPeopleIDs() []int {
	peopleIDs := platformPeopleIDs()
	assignedPeopleIDs := pcfsPeopleAssignedNowIDs()
	var unassignedPeopleIDs []int

	for _, id := range peopleIDs {
		assigned := false
		for _, aid := range assignedPeopleIDs {
			if id == aid {
				assigned = true
			}
		}
		if !assigned {
			unassignedPeopleIDs = append(unassignedPeopleIDs, id)
		}
	}
	return unassignedPeopleIDs
}

func isoDateToday() string {
	return "2019-08-27"
}

func isoDateAYearAgo() string {
	// now := time.Now()
	return "2018-08-20"
}

func isoDateAYearFromNow() string {
	// now := time.Now()
	return "2020-08-20"
}

func dateFallsBetween(date, start, end string) bool {
	d, _ := strconv.Atoi(strings.ReplaceAll(date, "-", ""))
	s, _ := strconv.Atoi(strings.ReplaceAll(start, "-", ""))
	e, _ := strconv.Atoi(strings.ReplaceAll(end, "-", ""))

	return d <= e && d >= s
}

func pcfsPeopleAssignedNowIDs() []int {
	idMap := make(map[int]struct{})
	var ids []int
	var filter forecast.AssignmentFilter

	assignmentDate := isoDateToday()
	filterStartDate := isoDateAYearAgo()
	filterEndDate := isoDateAYearFromNow()

	filter.StartDate = filterStartDate
	filter.EndDate = filterEndDate

	assignments, _ := api.AssignmentsWithFilter(filter)

	for _, assignment := range assignments {
		startDate := assignment.StartDate
		endDate := assignment.EndDate
		if dateFallsBetween(assignmentDate, startDate, endDate) {
			idAssigned := assignment.PersonID
			if _, ok := idMap[idAssigned]; !ok {
				idMap[idAssigned] = struct{}{}
				ids = append(ids, idAssigned)
			}
		}
	}

	return ids
}

func main() {
	debug = true
	printUnassignedPlatformPeople()
}
