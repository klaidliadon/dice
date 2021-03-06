package dice

import "testing"

func TestBonus(t *testing.T) {
	b := Bonus(2)
	if e, o := 2, b.Total(); e != o {
		t.Errorf("Invalid value %d, expected %d", o, e)
	}
	b.Roll()
	if e, o := 2, b.Total(); e != o {
		t.Errorf("Invalid value %d, expected %d", o, e)
	}
}

func TestDice(t *testing.T) {
	d := Dice{true, 1, 8, nil}
	if d.Total() != 0 {
		t.Error("Total != 0 on creation")
	}
	d.Roll()
	if d.Total() <= 0 {
		t.Error("Total <= 0 after roll")
	}
}

func TestDiceNeg(t *testing.T) {
	d := Dice{false, 1, 8, nil}
	if d.Total() != 0 {
		t.Error("Total != 0 on creation")
	}
	d.Roll()
	if d.Total() >= 0 {
		t.Error("Total >= 0 after roll")
	}

}

func TestPouch(t *testing.T) {
	p := NewPouch("d8-1-2d6+2")
	if e, o := 4, len(p.items); e != o {
		t.Errorf("Unexpected length %s, expected %s", o, e)
	}
	p.Roll()
	var length = []int{1, 0, 2, 0}
	for i, d := range p.items {
		if e, o := length[i], len(d.Partials()); e != o {
			t.Errorf("Unexpected partials length %s for %d, expected %s", o, i, e)
		}
	}
	t.Log(p)
}

func TestEmptyPouch(t *testing.T) {
	p := NewPouch("")
	if e, o := 0, len(p.items); e != o {
		t.Errorf("Unexpected length %s, expected %s", o, e)
	}
}
