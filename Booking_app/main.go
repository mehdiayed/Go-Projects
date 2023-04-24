package main

import (
	"Booking_app/helper"
	"fmt"
	"sync"
	"time"
)

const conferenceTickets = 50

var conferenceName = "Illustrator conference"
var RemainingTickets uint = 50
var bookings = make([]UserData, 0)

type UserData struct {
	firstName       string
	lastName        string
	email           string
	numberOfTickets uint
}

var wg = sync.WaitGroup{}

func main() {

	greetUsers()

	// for {
	firstName, lastName, email, userTickets := getUserInput()

	isValidName, isValidEmail, isValidTickets := helper.ValidateUserInput(firstName, lastName, email, userTickets, RemainingTickets)
	if isValidTickets && isValidEmail && isValidName {

		bookTicket(userTickets, firstName, lastName, email)

		wg.Add(1)

		go sendTicket(userTickets, firstName, lastName, email)

		//printing first names
		// firstNames := GetFirstNames()

		fmt.Printf("These are all the first names of booking: %v\n", GetFirstNames())

		if RemainingTickets == 0 {
			// end programme
			fmt.Println("our conference is booked out.come back next year")
			// break
		}

	} else {
		if !isValidName {
			fmt.Println("your name is short ")
		}
		if !isValidEmail {
			fmt.Println("your name must containes @ ")
		}
		if !isValidTickets {
			fmt.Println("invalid number of tickets ")
		}
		fmt.Println("your input data is invalide please try again ! ")
	}

	// }
	wg.Wait()
}

func greetUsers() {

	fmt.Println("Welcom to our ", conferenceName, " booking application ♥")
	fmt.Printf("We have total %v tickets and %v are still available.\n", conferenceTickets, RemainingTickets)
	fmt.Println("Get your tickets here to attend !")
}

func GetFirstNames() []string {

	firstNames := []string{}

	for _, booking := range bookings {

		firstNames = append(firstNames, booking.firstName)
	}

	return firstNames
}

// func validateUserInput(firstName string, lastName string, email string, userTickets int, RemainingTickets int) (bool, bool, bool) {
// 	isValidName := len(firstName) >= 2 && len(lastName) >= 2
// 	isValidEmail := strings.Contains(email, "@")
// 	isValidTickets := userTickets > 0 && userTickets <= RemainingTickets
// 	return isValidName, isValidEmail, isValidTickets
// }

func getUserInput() (string, string, string, uint) {
	var firstName string
	var lastName string
	var email string
	var userTickets uint

	fmt.Println("Enter your first name: ")
	fmt.Scan(&firstName)

	fmt.Println("Enter your last name: ")
	fmt.Scan(&lastName)

	fmt.Println("Enter your email address: ")
	fmt.Scan(&email)

	fmt.Println("Enter number of tickets: ")
	fmt.Scan(&userTickets)

	return firstName, lastName, email, userTickets

}

func bookTicket(userTickets uint, firstName string, lastName string, email string) {
	RemainingTickets = RemainingTickets - userTickets

	// create map
	var userData = UserData{
		firstName:       firstName,
		lastName:        lastName,
		email:           email,
		numberOfTickets: userTickets,
	}

	bookings = append(bookings, userData)
	fmt.Printf("List of bookings is %v", bookings)

	fmt.Printf("Thank you %v %v for booking %v tickets. You will receive a confirmation email at %v\n", firstName, lastName, userTickets, email)
	fmt.Printf("%v tickets remaining for %v\n", RemainingTickets, conferenceName)
}

func sendTicket(userTickets uint, firstName string, lastName string, email string) {

	time.Sleep(10 * time.Second)
	var ticket = fmt.Sprintf("%v tickets fro %v %v", userTickets, firstName, lastName)
	fmt.Println(" ----------------------------- ")
	fmt.Printf("sending tickets ♥ ♥ ♥ \n %v to email adress %v\n", ticket, email)
	fmt.Println(" ----------------------------- ")
	wg.Done()
}
