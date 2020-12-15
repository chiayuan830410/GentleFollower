package main

import (
	"fmt"
	"io"
	"bufio"
	"os"
	"strings"
	"time"
	"math/rand"
	"gopkg.in/ahmdrz/goinsta.v2"
)

func getFollowList(insta* goinsta.Instagram)[]string {
	users := insta.Account.Following()
	if users.Next() {
		follingList := make([]string,len(users.Users))
		for i, user := range users.Users {
			follingList[i] = user.Username
		}
		return follingList
	} else {
		return []string{}
	}
}

func saveIG(insta* goinsta.Instagram) {

	// open output file
	file, err := os.Create("list.txt")
	if err != nil {
		panic(err)
	}
	// close file on exit and check for its returned error
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	foLists:=getFollowList(insta)
	
	for _, fo := range foLists {
		if _, err := file.Write([]byte(fo + "\n")); err != nil {
			panic(err)
		}
	}

}
func restoreIG(insta* goinsta.Instagram) {
	//get following list
	foLists:=getFollowList(insta)
	oldFoLists:=make([]string,0)

	//get history data
	inputFile, err := os.Open("list.txt")
	if err != nil {
		fmt.Println("open error!")
		return
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	
	//get old following list
	for {
		name, Error := inputReader.ReadString('\n')
		if Error == io.EOF {
			break
		}
		if len(strings.TrimSpace(name)) != 0 {
			oldFoLists=append(oldFoLists,strings.TrimSpace(name))
		}
	}

	for i,fo := range foLists {
		for j,oFo := range oldFoLists {
			if fo==oFo {
				foLists[i] = ""
				oldFoLists[j] = ""
				break
			}
		}
	}

	//unfollow 
	fmt.Println("Unfollow :")
	for _,fo := range foLists {	
		if fo != "" {
			user, err := insta.Profiles.ByName(fo)
			if err != nil {
				fmt.Println(err)
				time.Sleep(50*time.Second)
			} else {
				fmt.Println(fo)
				//avoid account be blocked
				time.Sleep(time.Duration(rand.Intn(5000)+500)*time.Millisecond)
				_ = user.Unfollow()
			}
		}
	}

	//follow 
	fmt.Println("Follow :")
	for _,oFo := range oldFoLists {	
		if oFo != "" {
			user, err := insta.Profiles.ByName(oFo)
			if err != nil {
				fmt.Println(err)
				time.Sleep(50*time.Second)
			} else {
				//avoid account be blocked
				time.Sleep(time.Duration(rand.Intn(5000)+500)*time.Millisecond)
				fmt.Println(oFo)
				_ = user.Follow()
			}
		}
	}
}

func main() {
	var yourAccount string
	var yourPasswd string
	var taskType int // 1. save; 2. restore;

	fmt.Println("Please enter Account : ")
	fmt.Scanln(&yourAccount)
	fmt.Println("Please enter Password : ")
	fmt.Scanln(&yourPasswd)
	insta := goinsta.New(yourAccount, yourPasswd)
	if err := insta.Login(); err != nil {
		fmt.Println(err)
		return
	}
	defer insta.Logout()
	fmt.Println("1. Save following\n2. Restore following\nPlease enter a number : ")
	fmt.Scanln(&taskType)
	switch taskType {
	case 1:
		fmt.Println("Save following")
		saveIG(insta)
	case 2:
		fmt.Println("Restore following")
		restoreIG(insta)
	default:
		fmt.Println("Wrong input")
	}
	
}