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
    actions map[string]int
}

type Individual struct {
    name string
    age  int
    health int
    food int
    pos Pos
    index int

    middle_age int //TODO: move to DNA
    max_food int   //TODO: move to DNA
    desigions []string //TODO: move to DNA
    dna []Gene

    stat Stat
}


func createIndivid(index int, pos Pos) Individual {
    var individual Individual
    individual.name = "i_" + strconv.Itoa(index)
    individual.age = 0
    individual.health = 100
    individual.food = 1
    individual.pos = pos
    individual.index = index
    individual.middle_age = 100
    // individual.dna = TODO: generate simple DNA
    individual.max_food = 3 + rand.Intn(3)
    individual.stat.actions = make(map[string]int)
    return individual;
}

func isYoung(individual *Individual) bool {
    return !isChild(individual) && !isOld(individual)
}

func isOld(individual *Individual) bool {
    return individual.age >= 40
}

func isChild(individual *Individual) bool {
    return individual.age <= 20
}

func calculateDie(individual *Individual) bool {
    if (individual.health == 0) {
        return true
    }
    if (individual.age < individual.middle_age) {
        return false
    }
    return (individual.age - individual.middle_age) > rand.Intn(individual.middle_age)
}

func placeIndivid(individual *Individual, planet *Planet) {
    planet.grid[individual.pos.x][individual.pos.y].
    individuals[individual.name] = individual
}

func moveIndivid(individual *Individual, planet *Planet, new_pos Pos) {
    delete(planet.grid[individual.pos.x][individual.pos.y].individuals,
        individual.name)
    individual.pos = new_pos
    placeIndivid(individual, planet)
}


