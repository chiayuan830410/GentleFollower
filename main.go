package main

import (
	"fmt"
	"io"
	"bufio"
	"os"
	"strings"
	"time"
	"gopkg.in/ahmdrz/goinsta.v2"
)

func saveIG(insta* goinsta.Instagram) {

    // open output file
    fo, err := os.Create("list.txt")
    if err != nil {
        panic(err)
    }
    // close fo on exit and check for its returned error
    defer func() {
        if err := fo.Close(); err != nil {
            panic(err)
        }
    }()

	users := insta.Account.Following()

	for users.Next() {
		fmt.Println("Next:", users.NextID)
		for _, user := range users.Users {
			fmt.Printf("   - %s\n", user.Username)
			if _, err := fo.Write([]byte(user.Username + "\n")); err != nil {
				panic(err)
			}
		}
	}
}
func followIG(insta* goinsta.Instagram) {
	
	inputFile, err := os.Open("list.txt")
	if err != nil {
		fmt.Println("open error!")
		return
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	
	for {
		
		name, Error := inputReader.ReadString('\n')
		
		if Error == io.EOF {
			
			return
		}
		if len(strings.TrimSpace(name)) != 0 {
			for {
				user, err := insta.Profiles.ByName(strings.TrimSpace(name))
				fmt.Println(name)
				if err != nil {
					fmt.Println(err)
					time.Sleep(50*time.Second)
				} else {
					_ = user.Follow()
					break
				}
			}			
		}
	}
}
func main() {  
	insta := goinsta.New("yourAccount", "yourPasswd")
	if err := insta.Login(); err != nil {
		fmt.Println(err)
		return
	}
	defer insta.Logout()
	//saveIG(insta)
	//followIG(insta)
	
	
	
  

  
}