//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"flag"
	"fmt"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// beginner limits per request
func beginner(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	status := make(chan bool) 

	go func() {
		process()

		status <-true
	}()

	select {
	case <-status:
		return true
	case <-time.After(10 * time.Second):
		return false
	}
}

// advanced limits per user
func advanced(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	elapsedTime := make(chan time.Duration) 

	remaining := (10 * time.Second) - time.Duration(u.TimeUsed)

	// Returns immediately if the 
	// user has trespassed 10 secs
	if remaining <= 0 {
		fmt.Println("Reached limit before playing")
		return false
	}

	go func() {
		start := time.Now()

		process()

		elapsedTime <-time.Since(start)
	}()

	select {
	case et := <-elapsedTime:
		u.TimeUsed += int64(et)

		return true
	case <-time.After(remaining):
		fmt.Println("Reached limit on playing")
		return false
	}
}

var handler func(func(), *User) bool = beginner

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	return handler(process, u) 
}

func main() {
	advancedHdlr := flag.Bool("advanced", true, "Set advanced handler")

	if *advancedHdlr {
		fmt.Println("Using advanced handler")
		handler = advanced
	}
	

	RunMockServer()
}
