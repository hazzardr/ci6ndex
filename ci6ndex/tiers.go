package ci6ndex

import (
	"fmt"
	"math"
)

type Tier struct {
	name  string
	value float64
}

func (t *Tier) Name() string {
	return t.name
}
func (t *Tier) Value() float64 {
	return t.value
}

var (
	S = Tier{
		name:  "S 🏆",
		value: 1.0,
	}
	A = Tier{
		name:  "A ⭐",
		value: 2.0,
	}
	B = Tier{
		name:  "B 👍",
		value: 3.0,
	}
	C = Tier{
		name:  "C 🤔",
		value: 4.0,
	}
	F = Tier{
		name:  "F 💩",
		value: 5.0,
	}
)

func GetTierByName(name string) (*Tier, error) {
	switch name {
	case "S":
		return &S, nil
	case "A":
		return &A, nil
	case "B":
		return &B, nil
	case "C":
		return &C, nil
	case "F":
		return &F, nil
	default:
		return nil, fmt.Errorf("invalid tier name: %s", name)
	}
}

func GetTierByValue(val float64) (*Tier, error) {
	intVal := int(math.Round(val))
	switch intVal {
	case 1:
		return &S, nil
	case 2:
		return &A, nil
	case 3:
		return &B, nil
	case 4:
		return &C, nil
	case 5:
		return &F, nil
	default:
		return nil, fmt.Errorf("invalid tier value: %f", val)
	}
}
