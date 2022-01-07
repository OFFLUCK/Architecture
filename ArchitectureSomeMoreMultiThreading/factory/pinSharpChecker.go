package factory

import (
	"fmt"
	"go.uber.org/atomic"
	"sync"
)

type PinSharpChecker struct {
	pins     []*Pin
	lock     sync.Mutex
	pinsLeft *atomic.Uint32
	ans      []*Pin
}

func NewCheck(pinsLeft *atomic.Uint32, ans []*Pin) *PinSharpChecker {
	checker := new(PinSharpChecker)
	checker.pinsLeft = pinsLeft
	checker.ans = ans
	return checker
}

func (checker *PinSharpChecker) Run(wg *sync.WaitGroup) {
	checker.lock.Lock()
	defer wg.Done()

	if len(checker.pins) == 0 {
		checker.lock.Unlock()
		return
	}

	//fmt.Printf("\n%f\n\n", checker.pins[len(checker.pins)-1].sharpness)

	if checker.pins[len(checker.pins)-1].sharpness > 0.75 {
		checker.lock.Unlock()
		checker.returnPin()
	} else {
		fmt.Printf("Sharp checker disapproved and threw away a pin with curvature %f and sharpness %f\n",
			checker.pins[len(checker.pins)-1].curvature,
			checker.pins[len(checker.pins)-1].sharpness)
		checker.pinsLeft.Dec()
	}

	checker.lock.Lock()
	fmt.Println(len(checker.pins))
	checker.pins = checker.pins[:len(checker.pins)-1]
	fmt.Println(len(checker.pins))
	checker.lock.Unlock()
}

func (checker *PinSharpChecker) receivePin(pin *Pin) {
	checker.lock.Lock()
	// Adding a new pin to our current collection.
	checker.pins = append(checker.pins, pin)
	checker.lock.Unlock()
}

func (checker *PinSharpChecker) returnPin() {
	checker.lock.Lock()
	fmt.Printf("Curvature checker approved and gave grinder man a pin with curvature %f and sharpness %f\n",
		checker.pins[len(checker.pins)-1].curvature,
		checker.pins[len(checker.pins)-1].sharpness)
	checker.pinsLeft.Dec()
	checker.ans = append(checker.ans, checker.pins[len(checker.pins)-1])
	checker.lock.Unlock()
}
