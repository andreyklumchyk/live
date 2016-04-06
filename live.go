package main

import (
    "fmt"
    "math/rand"
    "strconv"
    "time"
    live "live"
)

type Pos struct {
    x, y int
}

type Gene struct {
    id int
    value []float32
    ability string
}

type Location struct {
    x int
    y int
    width int
    height int
}

type Desigion struct {
    needles int
    action string
    pos Pos
    withWho *Individual
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

type Terrain struct {
    food int
    individuals map[string]*Individual
    pos Pos
}

type Planet struct {
    width int
    height int
    population int
    grid [][]Terrain
    individuals []Individual
}



func main() {
    rand.Seed(int64(time.Now().Nanosecond()))
    var earth = createPlanet(2)
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

func createPlanet(population int) Planet {
    var planet Planet
    planet.height = population * 10
    planet.width = population * 10
    planet.grid = make([][]Terrain, planet.width)
    planet.individuals = make([]Individual, population)

    var terrain Terrain
    for x := 0; x < planet.width; x++ {
        planet.grid[x] = make([]Terrain, planet.height)
        for y := 0; y < planet.height; y++ {
            terrain.individuals = make(map[string]*Individual)
            terrain.pos = Pos{x, y}
            terrain.food = rand.Intn(2)
            planet.grid[x][y] = terrain
        }
    }

    var pos Pos
    for i := 0; i < population; i++ {
        pos.x = rand.Intn(planet.width)
        pos.y = rand.Intn(planet.height)
        planet.individuals[i] = createIndivid(i, pos)
        placeIndivid(&planet.individuals[i], &planet)
        planet.population++;
    }
    return planet;
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

func calculateDie(individual *Individual) bool {
    if (individual.health == 0) {
        return true
    }
    if (individual.age < individual.middle_age) {
        return false
    }
    return (individual.age - individual.middle_age) > rand.Intn(individual.middle_age)
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

func takeNearestTerrains(planet *Planet, pos Pos, radius int) []*Terrain {
    var minX = pos.x - radius
    var maxX = pos.x + radius
    var minY = pos.y - radius
    var maxY = pos.y + radius
    if (minX < 0) {
        minX = 0
    }
    if (maxX >= planet.width) {
        maxX = planet.width - 1
    }
    if (minY < 0) {
        minY = 0
    }
    if (maxY >= planet.height) {
        maxY = planet.height - 1
    }

    var count = (maxX - minX + 1) * (maxY - minY + 1)
    var terrains = make([]*Terrain, count)
    var i = 0
    for x := minX; x <= maxX; x++ {
        for y := minY; y <= maxY; y++ {
            terrains[i] = &planet.grid[x][y]
            i++
        }
    }
    return terrains
}

func makeDesigion(me *Individual, planet *Planet) {
    var terrains = takeNearestTerrains(planet, me.pos, 1)
    var count = len(terrains)
    var bestDesigion Desigion
    var desigion Desigion
    var terrain *Terrain

    for i := 0; i < count; i++ {
        terrain = terrains[i]
        desigion = calculateDesigion(me, terrain)
        if (desigion.needles > bestDesigion.needles) {
            bestDesigion = desigion
        }
    }

    processDesigion(me, planet, bestDesigion)
}

func desigionIndividual(me *Individual, other *Individual, count int) Desigion {
    var result Desigion
    // TODO: peace of ...
    if (isYoung(me) && isYoung(other) && count < 5 && //TODO: 5 - max terrain // ain population
        me.food > 2/*TODO: could change*/) {
        result.needles = 10 + rand.Intn(10)
        result.action = "sex"
    } else if (me.food < 0) {
        if (isChild(me)) {
            result.needles = 0
            result.action = "nothing"
        } else {
            result.needles = 10 + rand.Intn(me.health - other.health) +
                rand.Intn(other.age - me.age)
            result.action = "kill"
        }
    }
    result.withWho = other
    return result
}

func desigionFood(me *Individual, terrain *Terrain) Desigion {
    var result Desigion
    if (terrain.food == 0) {
        result.needles = 0
    } else {
        result.needles = me.max_food - me.food + terrain.food + rand.Intn(10)
    }
    result.action = "eat"
    return result
}

func calculateDesigion(me *Individual, terrain *Terrain) Desigion {
    var desigion Desigion
    var bestDesigion Desigion
    var count = len(terrain.individuals)
    if (len(terrain.individuals) > 0) {
        for _, individual := range terrain.individuals {
            desigion = desigionIndividual(me, individual, count)
            if (desigion.needles > bestDesigion.needles) {
                bestDesigion = desigion
            }
        }
    }
    desigion = desigionFood(me, terrain)
    if (desigion.needles > bestDesigion.needles) {
        bestDesigion = desigion
    }

    bestDesigion.pos = terrain.pos
    return bestDesigion
}

func processDesigion(me *Individual, planet *Planet, desigion Desigion) {
    var terrain = &planet.grid[desigion.pos.x][desigion.pos.y]
    moveIndivid(me, planet, desigion.pos)

    switch desigion.action {
    case "eat":
        var food = me.max_food - me.food
        if (food > terrain.food) {
            food = terrain.food
        }
        me.food += food
        terrain.food -= food
    case "sex":
        me.food -= 2 //TODO: could change
        planet.individuals = append(planet.individuals,
            createIndivid(len(planet.individuals), desigion.pos))
        placeIndivid(&planet.individuals[len(planet.individuals) - 1], planet)
        planet.population++;
    case "kill":
        var newHealth = me.health - desigion.withWho.health -
            (me.age - desigion.withWho.age) % 10 - rand.Intn(10)
        if newHealth > 0 {
            if (desigion.withWho.food > 0) {
                terrain.food += desigion.withWho.food
            }
            terrain.food += desigion.withWho.health % 10
            me.health = newHealth
            desigion.withWho.health = 0
        } else {
            if (me.food > 0) {
                terrain.food += me.food
            }
            terrain.food += me.health % 10
            me.health = 0
            desigion.withWho.health = newHealth
        }
        me.food -= 3
        planet.population--;
    }

    if (me.health < 1) {
        me.health = 0
    } else if (me.food <= 0) {
        me.health -= 5
    } else if (me.health < 100) {
        me.health += 1
    }
    me.stat.actions[desigion.action] += 1
}



