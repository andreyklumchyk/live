package main

import (
    "fmt"
    "math/rand"
    "time"
    "live"
)

func main() {
    rand.Seed(int64(time.Now().Nanosecond()))
    var earth = live.CreatePlanet(10)
    var step = 0
    for step <= 400 {
        rand.Seed(int64(time.Now().Nanosecond()))
        live.RunCycle(&earth)
        if (earth.Population > 1000) {
            fmt.Println("TOO MUCH")
            return
        }
        step++
    }

    fmt.Println("GAME OVER")
    fmt.Println("Steps: ", step)
    fmt.Println("Individuals")

    var count = len(earth.Individuals)
    var individual *live.Individual
    for i := 0; i < count; i++ {
        individual = &earth.Individuals[i]
        fmt.Println(
            " Health: ", individual.Health,
            " Position: ", individual.Pos,
            " Age: ", individual.Age,
            " Food: ", individual.Food)
        fmt.Println("Actions: ", individual.Stat.Actions)
    }

    fmt.Println(
        " Width: ", earth.Width,
        " Height: ", earth.Height,
        " Population: ", earth.Population,
        " Stat: ", earth.Stat)
}
