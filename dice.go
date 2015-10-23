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

const RollFormat = `\s*[+-]?([0-9]*)[dD]([0-9]+)|([+-]?[0-9]*)`

var splitter = regexp.MustCompile(RollFormat)

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
			r = append(r, NewBonus(v))
		} else {
			s := len(m[0]) == 0 || m[0][0] != '-'
			q, _ := strconv.Atoi(m[1])
			if q == 0 {
				q++
			}
			f, _ := strconv.Atoi(m[2])
			r = append(r, NewDice(s, q, f))
		}
	}
	return &Pouch{s, r}
}

type Pouch struct {
	src   string
	items []Item
}

func (p *Pouch) Roll() {
	for _, i := range p.items {
		i.Roll()
	}
}

func (p *Pouch) Total() int {
	var t = 0
	for _, i := range p.items {
		t += i.Total()
	}
	return t
}

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

type Item interface {
	Roll()
	Total() int
	Partials() []int
	String() string
}

func NewDice(sign bool, number, face int) Item {
	return &Dice{Sign: sign, Face: face, Qty: number}
}

type Dice struct {
	Sign      bool
	Qty, Face int
	Results   []int
}

func (d *Dice) Roll() {
	s := rand.NewSource(time.Now().UnixNano())
	var r = make([]int, d.Qty)
	for i := 0; i < d.Qty; i++ {
		r[i] = 1 + rand.New(s).Intn(d.Face)
	}
	d.Results = r
}

func (d *Dice) Total() int {
	var tot int
	for _, s := range d.Results {
		tot += s
	}
	if !d.Sign {
		return -tot
	}
	return tot
}

func (d *Dice) Partials() []int {
	return d.Results
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

func NewBonus(v int) Item {
	return &Bonus{v}
}

type Bonus struct {
	Value int
}

func (b *Bonus) Roll() {}

func (b *Bonus) Total() int { return b.Value }

func (b *Bonus) Partials() []int { return nil }

func (b *Bonus) String() string {
	var sign = ""
	if b.Value > -1 {
		sign = "+"
	}
	return fmt.Sprintf("%s%v", sign, b.Value)
}
