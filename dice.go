// The package dice creates rollable dice and pouch (group of dice)
package dice

import (
	"bytes"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// The regexp used to split dice and bonus
const RollFormat = `\s*[+-]?([0-9]*)[dD]([0-9]+)|([+-]?[0-9]*)`

var splitter = regexp.MustCompile(RollFormat)

// Creates a new pouch by parsing dice and bonus from the string
func NewPouch(s string) *Pouch {
	matches := splitter.FindAllStringSubmatch(strings.Replace(s, " ", "", -1), -1)
	var r = make([]Item, 0, len(matches))
	for i := 0; i < len(matches); i++ {
		m := matches[i]
		if m[0] == "" {
			continue
		}
		if m[3] != "" {
			v, _ := strconv.Atoi(m[0])
			r = append(r, Bonus(v))
		} else {
			s := len(m[0]) == 0 || m[0][0] != '-'
			q, _ := strconv.Atoi(m[1])
			if q == 0 {
				q++
			}
			f, _ := strconv.Atoi(m[2])
			r = append(r, &Dice{s, q, f, nil})
		}
	}
	return &Pouch{s, r}
}

// A group of dice
type Pouch struct {
	src   string
	items []Item
}

// Rolls all the dices in the pouch
func (p *Pouch) Roll() {
	for _, i := range p.items {
		i.Roll()
	}
}

// Gets the result from thelast roll
func (p *Pouch) Total() int {
	var t = 0
	for _, i := range p.items {
		t += i.Total()
	}
	return t
}

// Pretty print of the pouch with its result
func (p *Pouch) String() string {
	var b = bytes.NewBuffer(nil)
	for _, i := range p.items {
		var ps string
		if p := i.Partials(); p != nil {
			ps = fmt.Sprintf("%v", p)
		}
		fmt.Fprintf(b, "%-6s = %+4d %v\n", i, i.Total(), ps)
	}
	fmt.Fprintf(b, "-------------------------\nTotal\t%+4d\n", p.Total())
	return b.String()
}

// An interface for Bonus and Dice
type Item interface {
	Roll()
	Total() int
	Partials() []int
	String() string
}

// A group of dice of the same type (number of faces)
type Dice struct {
	Sign      bool
	Qty, Face int
	results   []int
}

// Rolls each dice
func (d *Dice) Roll() {
	s := rand.NewSource(time.Now().UnixNano())
	var r = make([]int, d.Qty)
	for i := 0; i < d.Qty; i++ {
		r[i] = 1 + rand.New(s).Intn(d.Face)
	}
	d.results = r
}

// Return the total form thge last roll
func (d *Dice) Total() int {
	var tot int
	for _, s := range d.results {
		tot += s
	}
	if !d.Sign {
		return -tot
	}
	return tot
}

// Return the result fo the single rolls
func (d *Dice) Partials() []int {
	return d.results
}

func (d *Dice) String() string {
	var b = bytes.NewBuffer(nil)
	if !d.Sign {
		b.WriteRune('-')
	} else {
		b.WriteRune('+')
	}
	fmt.Fprintf(b, "%dd%d", d.Qty, d.Face)
	return b.String()
}

// An integer modifier to a roll (positive or negative)
type Bonus int

// Does nothing
func (b Bonus) Roll() {}

// Returns the modifier
func (b Bonus) Total() int { return int(b) }

// Returns nothing
func (b Bonus) Partials() []int { return nil }

func (b Bonus) String() string {
	var sign = ""
	if b > -1 {
		sign = "+"
	}
	return fmt.Sprintf("%s%d", sign, b)
}
