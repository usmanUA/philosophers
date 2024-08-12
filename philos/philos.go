package philos

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

type MealSignal struct {
	ID    int
	Eaten bool
}

type Forks struct {
	Fork sync.Mutex
}

type Args struct {
	TotPhilos   int
	TimeToDie   int
	TimeToEat   int
	TimeToSleep int
	TotMeals    int
	DeathStatus bool
	PhilosFull  bool
	PrintLock   sync.Mutex
	DeathLock   sync.Mutex
	StartTime   time.Time
	Forks       []Forks
}

type Philosopher struct {
	Args         *Args
	ID           int
	PrevID       int
	Eaten        bool
	EatLog       sync.Mutex
	LastMealTime time.Duration
}

func NewArgs(clis []string) (*Args, error) {
	var meals int = 0
	phils, err := strconv.Atoi(clis[1])
	if err != nil {
		return nil, errors.New("\033[31mInvalid Input\033[0m")
	}
	die, err := strconv.Atoi(clis[2])
	if err != nil {
		return nil, errors.New("\033[31mInvalid Input\033[0m")
	}
	eat, err := strconv.Atoi(clis[3])
	if err != nil {
		return nil, errors.New("\033[31mInvalid Input\033[0m")
	}
	sleep, err := strconv.Atoi(clis[4])
	if err != nil {
		return nil, errors.New("\033[31mInvalid Input\033[0m")
	}
	if len(clis) == 6 {
		meals, err = strconv.Atoi(clis[5])
		if err != nil {
			return nil, errors.New("\033[31mInvalid Input\033[0m")
		}
	}
	return &Args{
		TotPhilos:   phils,
		TimeToDie:   die,
		TimeToEat:   eat,
		TimeToSleep: sleep,
		TotMeals:    meals,
		DeathStatus: false,
		PhilosFull:  false,
		Forks:       make([]Forks, phils),
		StartTime:   time.Now(),
	}, nil
}

func NewPhilosopher(args *Args) *Philosopher {
	return &Philosopher{
		Args:  args,
		Eaten: false,
	}
}
