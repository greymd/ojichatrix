package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"unicode/utf8"

	"github.com/greymd/ojichat/generator"
	"golang.org/x/crypto/ssh/terminal"
)

type printColor int

const (
	lightGreen printColor = iota
	darkGreen
)

type singleLine struct {
	flag    bool
	cursor  int
	pcolor  printColor
	message []rune
}

func createMessage(n int) []rune {
	s := ""
	for {
		result, _ := generator.Start(generator.Config{EmojiNum: 3})
		s += result
		if utf8.RuneCountInString(s) > n {
			return []rune(s)
		}
	}
}

func printLightGreen(y, x int, ch string) {
	fmt.Printf("\033[%d;%dH\033[1;32m%s\033[0;0m", y, x, ch)
}

func printDarkGreen(y, x int, ch string) {
	fmt.Printf("\033[%d;%dH\033[2;32m%s\033[0;0m", y, x, ch)
}

func printWhite(y, x int, ch string) {
	fmt.Printf("\033[%d;%dH\033[40;37m%s\033[0;0H", y, x, ch)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	width, height, err := terminal.GetSize(0)
	if err != nil {
		os.Exit(1)
	}
	width--
	height--

	// Make bold
	fmt.Println("\033[1;40m")
	// clear
	fmt.Println("\033[2J")
	a := map[int]*singleLine{}
	for i := 0; i < width; i++ {
		a[i] = &singleLine{false, 0, darkGreen, createMessage(height)}
	}
	for {
		(*a[rand.Intn(width)]).flag = true
		for x := range a {
			if !(*a[x]).flag {
				continue
			}
			y := (*a[x]).cursor
			(*a[x]).cursor++
			ny := (*a[x]).cursor
			ch := string((*a[x]).message[(*a[x]).cursor])
			if (*a[x]).pcolor == lightGreen {
				printLightGreen(y, x, ch)
			} else if (*a[x]).pcolor == darkGreen {
				printDarkGreen(y, x, ch)
			}
			printWhite(ny, x, ch)
			if (*a[x]).cursor >= height {
				(*a[x]).cursor = 0
				(*a[x]).message = createMessage(height)
				if rand.Intn(10) <= 8 {
					(*a[x]).pcolor = darkGreen
				} else {
					(*a[x]).pcolor = lightGreen
				}
			}
		}
		time.Sleep(50000000)
	}
}
