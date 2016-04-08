package main

import (
    "fmt"
    "math/rand"
    "time"
    "live"
)

func main() {
    rand.Seed(int64(time.Now().Nanosecond()))
    var earth = live.CreatePlanet(2)
    var step = 0
    for step <= 60 {
        rand.Seed(int64(time.Now().Nanosecond()))
        live.RunCycle(&earth)
        if (earth.population > 1000) {
            fmt.Println("TOO MUCH")
            return
        }
        step++
    }

    fmt.Println("GAME OVER")
    fmt.Println("Steps: ", step)
    fmt.Println("Individuals")
    var count = len(earth.individuals)
    var individual *Individual
    for i := 0; i < count; i++ {
        individual = &earth.individuals[i]
        fmt.Println(
            " Health: ", individual.health,
            " Position: ", individual.pos,
            " Age: ", individual.age)
        fmt.Println("Actions: ", individual.stat.actions)
    }
    fmt.Println(
        " Width: ", earth.width,
        " Height: ", earth.height,
        " Population: ", earth.population)
}
