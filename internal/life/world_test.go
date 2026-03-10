package life

import "testing"

func TestUnderpopulation(t *testing.T) {
	w := NewWorld(3, 3)
	w.SetAlive(1, 1, true)
	w.SetAlive(2, 1, true)

	w.Step()

	if w.Alive(1, 1) {
		t.Fatal("expected lonely cell to die")
	}
}

func TestSurvival(t *testing.T) {
	w := NewWorld(5, 5)
	w.SetAlive(2, 2, true)
	w.SetAlive(1, 2, true)
	w.SetAlive(3, 2, true)

	w.Step()

	if !w.Alive(2, 2) {
		t.Fatal("expected cell with two neighbors to survive")
	}
}

func TestOverpopulation(t *testing.T) {
	w := NewWorld(5, 5)
	w.SetAlive(2, 2, true)
	w.SetAlive(1, 2, true)
	w.SetAlive(3, 2, true)
	w.SetAlive(2, 1, true)
	w.SetAlive(2, 3, true)

	w.Step()

	if w.Alive(2, 2) {
		t.Fatal("expected crowded cell to die")
	}
}

func TestReproduction(t *testing.T) {
	w := NewWorld(5, 5)
	w.SetAlive(2, 1, true)
	w.SetAlive(1, 2, true)
	w.SetAlive(3, 2, true)

	w.Step()

	if !w.Alive(2, 2) {
		t.Fatal("expected dead cell with three neighbors to become alive")
	}
}

func TestToroidalWrapping(t *testing.T) {
	w := NewWorld(5, 5)
	w.SetAlive(4, 4, true)
	w.SetAlive(0, 4, true)
	w.SetAlive(4, 0, true)

	w.Step()

	if !w.Alive(0, 0) {
		t.Fatal("expected corner cell to be born via wrapped neighbors")
	}
}

func TestDeterministicSeededWorld(t *testing.T) {
	a := NewSeededWorld(8, 8, 42, DefaultDensity)
	b := NewSeededWorld(8, 8, 42, DefaultDensity)

	for i := 0; i < 5; i++ {
		a.Step()
		b.Step()
	}

	for idx := range a.Cells() {
		if a.Cells()[idx] != b.Cells()[idx] {
			t.Fatalf("world state diverged at index %d", idx)
		}
	}
}

func TestApplyPatternAcornSeedsSmallStructure(t *testing.T) {
	w := NewPatternWorld(20, 20, "acorn", 0, DefaultDensity)

	count := 0
	for _, alive := range w.Cells() {
		if alive {
			count++
		}
	}

	if count != 6 {
		t.Fatalf("expected 6 live cells for acorn, got %d", count)
	}
}
