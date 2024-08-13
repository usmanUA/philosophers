package main

import (
	"fmt"
	"os"
	"philoGolang/philos"
	"sync"
	"time"
)

func eat(p *philos.Philosopher) {
	p.Args.DeathLock.Lock()
	if p.Args.DeathStatus == true {
		p.Args.DeathLock.Unlock()
		return
	}
	p.Args.DeathLock.Unlock()
	p.Args.PrintLock.Lock()
	current := time.Now()
	change := current.Sub(p.Args.StartTime)
	fmt.Printf("\033[33m%d %d is eating\033[0m\n", change.Milliseconds(), p.ID+1)
	p.Args.PrintLock.Unlock()
	time.Sleep(time.Duration(p.Args.TimeToEat) * time.Millisecond)
}

func sleep(p *philos.Philosopher) {
	p.Args.DeathLock.Lock()
	if p.Args.DeathStatus == true {
		p.Args.DeathLock.Unlock()
		return
	}
	p.Args.DeathLock.Unlock()
	p.Args.PrintLock.Lock()
	current := time.Now()
	change := current.Sub(p.Args.StartTime)
	fmt.Printf("\033[33m%d %d is sleeping\033[0m\n", change.Milliseconds(), p.ID+1)
	p.Args.PrintLock.Unlock()
	time.Sleep(time.Duration(p.Args.TimeToSleep) * time.Millisecond)
}

func think(p *philos.Philosopher) {
	p.Args.DeathLock.Lock()
	if p.Args.DeathStatus == true {
		p.Args.DeathLock.Unlock()
		return
	}
	p.Args.DeathLock.Unlock()
	p.Args.PrintLock.Lock()
	current := time.Now()
	change := current.Sub(p.Args.StartTime)
	fmt.Printf("\033[33m%d %d is thinking\033[0m\n", change.Milliseconds(), p.ID+1)
	p.Args.PrintLock.Unlock()
}

func forkTaken(p *philos.Philosopher) {
	p.Args.DeathLock.Lock()
	if p.Args.DeathStatus == true {
		p.Args.DeathLock.Unlock()
		return
	}
	p.Args.DeathLock.Unlock()
	p.Args.PrintLock.Lock()
	current := time.Now()
	change := current.Sub(p.Args.StartTime)
	fmt.Printf("\033[33m%d %d has taken a fork\033[0m\n", change.Milliseconds(), p.ID+1)
	p.Args.PrintLock.Unlock()
}

func checkStatus(p *philos.Philosopher) bool {
	p.Args.DeathLock.Lock()
	if p.Args.DeathStatus == true {
		p.Args.DeathLock.Unlock()
		return true
	}
	p.Args.DeathLock.Unlock()
	return false
}

func stopSimulation(p *philos.Philosopher, forks int) bool {
	if checkStatus(p) == true {
		if forks == 1 {
			p.Args.Forks[p.ID].Fork.Unlock()
		}
		return true
	}
	return false
}

func eatSleepRepeat(p *philos.Philosopher) bool {
	p.Args.Forks[p.ID].Fork.Lock()
	forkTaken(p)
	if stopSimulation(p, 1) == true {
		return false
	}
	p.Args.Forks[p.PrevID].Fork.Lock()
	forkTaken(p)
	p.EatLog.Lock()
	p.Eaten = true
	p.EatLog.Unlock()
	eat(p)
	p.Args.Forks[p.ID].Fork.Unlock()
	p.Args.Forks[p.PrevID].Fork.Unlock()
	sleep(p)
	if stopSimulation(p, 0) == true {
		return false
	}
	think(p)
	return true
}

func philosAtWork(p *philos.Philosopher) {
	if p.Args.TotPhilos == 1 {
		p.Args.Forks[p.ID].Fork.Lock()
		p.Args.PrintLock.Lock()
		fmt.Printf("\033[33m%d %d has taken a fork\033[0m\n", time.Since(p.Args.StartTime)/time.Millisecond, p.ID+1)
		p.Args.PrintLock.Unlock()
		time.Sleep(time.Duration(p.Args.TimeToDie) * time.Millisecond)
		p.Args.PrintLock.Lock()
		duration := time.Since(p.Args.StartTime)
		fmt.Printf("\033[31m%d %d died\033[0m\n", duration/time.Millisecond, p.ID+1)
		p.Args.PrintLock.Unlock()
		return
	}
	if p.ID%2 == 1 {
		time.Sleep(7 * time.Millisecond)
	}
	// if p.Args.TotPhilos%2 == 1 && p.ID+1 == p.Args.TotPhilos {
	// 	time.Sleep(5 * time.Millisecond)
	// }
	for {
		if eatSleepRepeat(p) == false {
			return
		}
	}
}

func monitorPhilos(p []*philos.Philosopher) {
	start := p[0].Args.StartTime
	dieTime := p[0].Args.TimeToDie
	total := p[0].Args.TotPhilos
	for {
		for i := range total {
			p[i].EatLog.Lock()
			if p[i].Eaten {
				p[i].Eaten = false
				p[i].LastMealTime = time.Since(start)
			}
			p[i].EatLog.Unlock()
			if time.Since(start)-p[i].LastMealTime >= time.Duration(dieTime*int(time.Millisecond)) {
				p[i].Args.DeathLock.Lock()
				p[i].Args.DeathStatus = true
				p[i].Args.DeathLock.Unlock()
				p[i].Args.PrintLock.Lock()
				fmt.Printf("\033[31m%d %d died\033[0m\n", time.Since(start)/time.Millisecond, p[i].ID+1)
				p[i].Args.PrintLock.Unlock()
				return
			}

		}
	}
}

func main() {
	argsLen := len(os.Args)
	if argsLen == 5 || argsLen == 6 {
		arg, err := philos.NewArgs(os.Args)
		if err != nil {
			fmt.Println(err)
		}
		philosophers := make([]*philos.Philosopher, arg.TotPhilos)
		for i := 0; i < arg.TotPhilos; i++ {
			philosopher := philos.NewPhilosopher(arg)
			philosopher.ID = i
			philosopher.PrevID = (i - 1 + arg.TotPhilos) % arg.TotPhilos
			philosophers[i] = philosopher
		}
		var wg sync.WaitGroup
		wg.Add(arg.TotPhilos)
		for i := 0; i < arg.TotPhilos; i++ {
			go func(p *philos.Philosopher) {
				defer wg.Done()
				philosAtWork(p)
			}(philosophers[i])
		}
		if arg.TotPhilos > 1 {
			monitorPhilos(philosophers)
		}
		wg.Wait()
		return
	}
	fmt.Println("\033[31mUSAGE:\n\t./philosophers <tot_philos> <time_to_die> <time_to_eat> <time_to_sleep>")
}
