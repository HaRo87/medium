package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"user"`
	Name        string    `json:"name"`
	Email       string    `json:"mail"`
	Country     string    `json:"country"`
	MemberSince time.Time `json:"date"`
}

type Movie struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Genre  string `json:"genre"`
	Length int    `json:"length"`
}

type WatchEvent struct {
	ID            string    `json:"id"`
	EventType     string    `json:"type"`
	Time          time.Time `json:"time"`
	MotionPicture Movie     `json:"movie"`
	Watcher       User      `json:"user"`
}

type Flags struct {
	numMovies      int
	numGenres      int
	numUsers       int
	numEvents      int
	numMailDomains int
	startDate      time.Time
}

const (
	defaultNumUsers   = 100
	defaultNumMovies  = 400
	defaultNumGenres  = 50
	defaultNumEvents  = 1000
	defaultNumDomains = 300
	defaultDate       = "2020-01-01"
)

func main() {
	var startTime string
	flags := Flags{}
	flag.IntVar(&flags.numUsers, "users", defaultNumUsers, "Number of users to generate, default: "+strconv.Itoa(defaultNumUsers))
	flag.IntVar(&flags.numMovies, "movies", defaultNumMovies, "Number of movies to generate, default: "+strconv.Itoa(defaultNumMovies))
	flag.IntVar(&flags.numGenres, "genres", defaultNumGenres, "Number of movie genres to generate, default: "+strconv.Itoa(defaultNumGenres))
	flag.IntVar(&flags.numEvents, "events", defaultNumEvents, "Number of events to generate, default: "+strconv.Itoa(defaultNumEvents))
	flag.IntVar(&flags.numMailDomains, "domains", defaultNumDomains, "Number of email domains to generate, default: "+strconv.Itoa(defaultNumDomains))
	flag.StringVar(&startTime, "start", defaultDate, "Start date for all date elements, default: "+defaultDate)
	flag.Parse()
	parsedDate, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		panic(err)
	}
	flags.startDate = parsedDate
	users := generateUsers(flags)
	for i, user := range users {
		if i == 5 {
			break
		}
		fmt.Printf("user %d: %v", i, user)
	}
	movies := generateMovies(flags)
	for i, movie := range movies {
		if i == 5 {
			break
		}
		fmt.Printf("movie %d: %v", i, movie)
	}
	events := generateEvents(flags, movies, users)
	for i, event := range events {
		if i == 5 {
			break
		}
		fmt.Printf("event %d: %v", i, event)
	}
	eventsJson, _ := json.Marshal(events)
	f, err := os.Create("output.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(eventsJson)
}

func generateEvents(flags Flags, movies []Movie, users []User) (events []WatchEvent) {
	endDate := time.Now()
	events = make([]WatchEvent, flags.numEvents)
	for i := range events {
		events[i] = WatchEvent{
			ID:            gofakeit.UUID(),
			EventType:     gofakeit.RandomString([]string{"started", "stopped"}),
			Time:          gofakeit.DateRange(flags.startDate, endDate),
			MotionPicture: movies[rand.Intn(len(movies))],
			Watcher:       users[rand.Intn(len(users))],
		}
	}
	return events
}

func generateMovies(flags Flags) (movies []Movie) {
	genres := generateGenres(flags.numGenres)
	movies = make([]Movie, flags.numMovies)
	for i := range movies {
		movies[i] = Movie{
			ID:     gofakeit.UUID(),
			Title:  gofakeit.MovieName(),
			Genre:  gofakeit.RandomString(genres),
			Length: gofakeit.Number(60, 200),
		}
	}
	return movies
}

func generateGenres(amount int) (genres []string) {
	genres = make([]string, amount)
	for i := range genres {
		genres[i] = gofakeit.MovieGenre()
	}
	return genres
}

func generateUsers(flags Flags) (users []User) {
	domains := generateDomains(flags.numMailDomains)
	users = make([]User, flags.numUsers)
	endDate := time.Now()
	for i := range users {
		name := gofakeit.Name()
		mail := strings.Replace(name, " ", ".", -1) + "@" + gofakeit.RandomString(domains)
		users[i] = User{
			ID:          gofakeit.UUID(),
			Username:    gofakeit.Username(),
			Name:        name,
			Email:       mail,
			Country:     gofakeit.Country(),
			MemberSince: gofakeit.DateRange(flags.startDate, endDate),
		}
	}
	return users
}

func generateDomains(amount int) (domains []string) {
	domains = make([]string, amount)
	for i := range domains {
		domains[i] = gofakeit.DomainName() + gofakeit.DomainSuffix()
	}
	return domains
}
