package live

import (
    "strconv"
    "math/rand"
)

type Pos struct {
    x, y int
}

type Gene struct {
    id int
    value []float32
    ability string
}

type Stat struct {
    Actions map[string]int
}

type Individual struct {
    name       string
    Age        int
    Health     int
    Food       int
    Pos        Pos
    index      int

    middle_age int //TODO: move to DNA
    max_food   int   //TODO: move to DNA
    desigions  []string //TODO: move to DNA
    dna        []Gene

    Stat       Stat
}


func createIndivid(index int, pos Pos) Individual {
    var individual Individual
    individual.name = "i_" + strconv.Itoa(index)
    individual.Age = 0
    individual.Health = 100
    individual.Food = 1
    individual.Pos = pos
    individual.index = index
    individual.middle_age = 100
    // individual.dna = TODO: generate simple DNA
    individual.max_food = 3 + rand.Intn(3)
    individual.Stat.Actions = make(map[string]int)
    return individual;
}

func isYoung(individual *Individual) bool {
    return !isChild(individual) && !isOld(individual)
}

func isOld(individual *Individual) bool {
    return individual.Age >= 40
}

func isChild(individual *Individual) bool {
    return individual.Age <= 20
}

func calculateDie(individual *Individual) bool {
    if (individual.Health == 0) {
        return true
    }
    if (individual.Age < individual.middle_age) {
        return false
    }
    return (individual.Age - individual.middle_age) > rand.Intn(individual.middle_age)
}

func placeIndivid(individual *Individual, planet *Planet) {
    planet.grid[individual.Pos.x][individual.Pos.y].
    individuals[individual.name] = individual
}

func moveIndivid(individual *Individual, planet *Planet, new_pos Pos) {
    delete(planet.grid[individual.Pos.x][individual.Pos.y].individuals,
        individual.name)
    individual.Pos = new_pos
    placeIndivid(individual, planet)
}
