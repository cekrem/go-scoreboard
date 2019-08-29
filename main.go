package main

import (
	"bufio"
	cRand "crypto/rand"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	figure "github.com/common-nighthawk/go-figure"
)

const zeldaASCII = `
__________________________††¥¥
_________________________††††¥¥
________________________††††††¥¥
_______________________††††††††¥¥
______________________††††††††††¥¥
_____________________††††††††††††¥¥
____________________††††††††††††††¥¥
___________________††††††††††††††††¥¥
__________________††††††††††††††††††¥¥
_________________††††††††††††††††††††¥¥
________________††††††††††††††††††††††¥¥
_______________††††††††††††††††††††††††¥¥
______________††††††††††††††††††††††††††¥¥
_____________††††††††††††††††††††††††††††¥¥
____________††¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥††¥¥
___________††††¥¥¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯††††¥¥
__________††††††¥¥¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯††††††¥¥
_________††††††††¥¥¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯††††††††¥¥
________††††††††††¥¥¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯††††††††††¥¥
_______††††††††††††¥¥¯¯¯¯¯¯¯¯¯¯¯¯¯¯††††††††††††¥¥
______††††††††††††††¥¥¯¯¯¯¯¯¯¯¯¯¯¯††††††††††††††¥¥
_____††††††††††††††††¥¥¯¯¯¯¯¯¯¯¯¯††††††††††††††††¥¥
____††††††††††††††††††¥¥¯¯¯¯¯¯¯¯††††††††††††††††††¥¥
___††††††††††††††††††††¥¥¯¯¯¯¯¯††††††††††††††††††††¥¥
__††††††††††††††††††††††¥¥¯¯¯¯††††††††††††††††††††††¥¥
_††††††††††††††††††††††††¥¥¯¯††††††††††††††††††††††††¥¥
††††††††††††††††††††††††††††††††††††††††††††††††††††††¥¥
¯¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥¥
____________¶¶3333¶___¶¶¶¶ÿÿÿÿ¶¶¶¶___¶3333¶¶
____________¶¶33333¶_¶ÿÿÿÿÿÿÿÿÿÿÿÿ¶_¶33333¶¶
____________¶¶¶¶¶¶¶¶ÿÿ¶¶¶¶¶¶¶¶¶¶¶¶ÿÿ¶¶¶¶¶¶¶¶
____________¶¶0000¶¶¶77777777777777¶¶¶0000¶¶
____________¶¶0000¶¶7¶¶¶¶¶¶¶¶¶¶¶¶¶¶7¶¶0000¶¶
_____________¶¶000¶¶¶a¯¯¶¶aaaa¶¶¯¯a¶¶¶000¶¶
_____________¶¶000¶¶aa¯¯¶¶aaaa¶¶¯¯aa¶¶000¶¶
______________¶¶00¶¶aaa¯¯aaaaaa¯¯aaa¶¶00¶¶
_______________¶¶000¶¶aaaaaaaaaaaa¶¶000¶¶
________________¶¶00¶¶¶¶aaaaaaaa¶¶¶¶00¶¶
_________________¶¶88888¶¶¶¶¶¶¶¶88888¶¶
__________________¶¶8855888888885588¶¶
__________________¶¶8855555555555588¶¶
________________¶¶11¶¶888855558888¶¶11¶¶
______________¶¶88881111¶¶1111¶¶11118888¶¶
______________¶¶¶¶¶¶8888111111118888¶¶¶¶¶¶
____________¶¶ƒƒ§§¶¶¶¶¶¶88888888¶¶¶¶¶¶§§ƒƒ¶¶
____________¶¶ƒƒƒƒ§§¶¶¯¶¶¶¶¶¶¶¶¶¶¯¶¶§§ƒƒƒƒ¶¶
____________¶¶¶¶¶¶¶¶¶______________¶¶¶¶¶¶¶¶¶
`

var winnerASCII = figure.NewFigure("Winner:", "", true)
var re = regexp.MustCompile(`(?:@)(?P<name>.*)\swith\s(?P<points>\d+)`)
var score []string

func main() {
	// parse input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Paste scoreboard results from `/scoreboard list`, hit enter twice when you're ready to draw Champion™.")
	for {
		entry, err := reader.ReadString('\n')
		failOnErr(err)

		if entry == "\n" {
			break
		}

		matches := re.FindStringSubmatch(entry)
		if len(matches) < 3 {
			log.Println("could not parse line!")
			continue
		}

		points, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Println("could not parse line!")
			continue
		}
		name := matches[1]

		for i := 0; i < points; i++ {
			score = append(score, name)
		}
	}

	if len(score) < 1 {
		log.Fatal("no valid score entries!")
	}

	// shuffle score
	seed, err := cRand.Int(cRand.Reader, big.NewInt(100000000))
	failOnErr(err)

	rand.Seed(seed.Int64() + time.Now().UnixNano())
	for countdown := len(score); countdown > 0; countdown-- {
		rand.Shuffle(len(score), func(i, j int) {
			score[i], score[j] = score[j], score[i]
		})

		clearConsole()
		fmt.Printf("%s\n%s\n", "Shuffling names like Øyvind would:", strings.Join(score[:countdown], "  "))
		time.Sleep(time.Millisecond * 50)
	}

	nameASCII := figure.NewFigure(score[0], "", true)
	fmt.Println("", winnerASCII, nameASCII, zeldaASCII)

	// In some browsers, connection will close before all is printed
	// - this avoids that :)
	time.Sleep(5 * time.Second)
}

func clearConsole() {
	fmt.Print("\033[H\033[2J")
}

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
