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

type Rating struct {
	ID            string    `json:"id"`
	Rating        int       `json:"rating"`
	Time          time.Time `json:"time"`
	MotionPicture Movie     `json:"movie"`
	Watcher       User      `json:"user"`
}

type Flags struct {
	numMovies      int
	numGenres      int
	numUsers       int
	numRatings     int
	numMailDomains int
	startDate      time.Time
}

const (
	defaultNumUsers   = 100
	defaultNumMovies  = 100
	defaultNumGenres  = 20
	defaultNumRatings = 1000
	defaultNumDomains = 50
	defaultDate       = "2020-01-01"
	defaultOutputFile = "fake.json"
)

func main() {
	var startTime string
	var outputFile string
	flags := Flags{}
	flag.IntVar(&flags.numUsers, "users", defaultNumUsers, "Number of users to generate, default: "+strconv.Itoa(defaultNumUsers))
	flag.IntVar(&flags.numMovies, "movies", defaultNumMovies, "Number of movies to generate, default: "+strconv.Itoa(defaultNumMovies))
	flag.IntVar(&flags.numGenres, "genres", defaultNumGenres, "Number of movie genres to generate, default: "+strconv.Itoa(defaultNumGenres))
	flag.IntVar(&flags.numRatings, "ratings", defaultNumRatings, "Number of ratings to generate, default: "+strconv.Itoa(defaultNumRatings))
	flag.IntVar(&flags.numMailDomains, "domains", defaultNumDomains, "Number of email domains to generate, default: "+strconv.Itoa(defaultNumDomains))
	flag.StringVar(&startTime, "start", defaultDate, "Start date for all date elements, default: "+defaultDate)
	flag.StringVar(&outputFile, "out", defaultOutputFile, "filename to store results, default: "+defaultOutputFile)
	flag.Parse()
	parsedDate, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		panic(err)
	}
	flags.startDate = parsedDate
	fmt.Println("Starting fake data generation for:")
	fmt.Printf("ratings: %d, users: %d, movies: %d, genres: %d\n", flags.numRatings, flags.numUsers, flags.numMovies, flags.numGenres)
	start := time.Now()
	users := generateUsers(flags)
	movies := generateMovies(flags)
	ratings := generateRatings(flags, movies, users)
	fmt.Printf("generating took: %f sec\n", time.Since(start).Seconds())
	fmt.Println("starting to marshal JSON")
	start = time.Now()
	eventsJson, _ := json.MarshalIndent(ratings, "", "\t")
	fmt.Printf("marshalling JSON took: %f sec\n", time.Since(start).Seconds())
	fmt.Printf("storing in: %s\n", outputFile)
	start = time.Now()
	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(eventsJson)
	fmt.Printf("storing took: %f sec\n", time.Since(start).Seconds())
}

func generateRatings(flags Flags, movies []Movie, users []User) (events []Rating) {
	endDate := time.Now()
	events = make([]Rating, flags.numRatings)
	for i := range events {
		user := users[rand.Intn(len(users))]
		events[i] = Rating{
			ID:            gofakeit.UUID(),
			Rating:        gofakeit.RandomInt([]int{0, 1, 2, 3, 4, 5}),
			Time:          gofakeit.DateRange(user.MemberSince, endDate),
			MotionPicture: movies[rand.Intn(len(movies))],
			Watcher:       user,
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
