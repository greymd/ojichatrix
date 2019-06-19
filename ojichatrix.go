package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/docopt/docopt-go"
	"github.com/greymd/ojichat/generator"
	"golang.org/x/crypto/ssh/terminal"
)

var appVersion = `Ojisan Nanchatte Matrix (ojichatrix) version 0.1.0
Copyright (c) 2019 Yamada, Yasuhiro
Released under the MIT License.
https://github.com/greymd/ojichat`

var usage = `Usage:
  ojichatrix [options]

Options:
  -h, --help      Show this help.
  -V, --version   Show version.
  -e <number>     Maximum number of consecutive Emojis.  [default: 4].
  -p <level>      Punctuation frequency level [min:0, max:3] [default: 0].`

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

func createMessage(n int, config generator.Config) []rune {
	s := ""
	for {
		result, _ := generator.Start(config)
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

func screenClear() {
	fmt.Println("\033[2J")
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Trap signals
	go func() {
		<-sigs
		// Clear ASCII Escapes
		fmt.Printf("\033[0;0m")
		screenClear()
		os.Exit(0)
	}()

	parser := &docopt.Parser{}
	args, _ := parser.ParseArgs(usage, nil, appVersion)
	config := generator.Config{}
	err := args.Bind(&config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	width, height, err := terminal.GetSize(0)
	if err != nil {
		os.Exit(1)
	}
	width--
	height--

	// Make bold
	fmt.Println("\033[1;40m")
	screenClear()
	a := map[int]*singleLine{}
	for i := 0; i < width; i++ {
		a[i] = &singleLine{false, 0, darkGreen, createMessage(height, config)}
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
				(*a[x]).message = createMessage(height, config)
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
