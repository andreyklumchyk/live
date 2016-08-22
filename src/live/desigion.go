package live

import (
    "math/rand"
)

type Desigion struct {
    needles int
    action string
    pos Pos
    withWho *Individual
}


func makeDesigion(me *Individual, planet *Planet) {
    var terrains = takeNearestTerrains(planet, me.Pos, 1)
    var count = len(terrains)
    var bestDesigion = Desigion{needles: 0, action: "nothing"}
    var desigion Desigion
    var terrain *Terrain

    for i := 0; i < count; i++ {
        terrain = terrains[i]
        desigion = calculateDesigion(me, terrain)
        if (desigion.needles >= bestDesigion.needles) {
            bestDesigion = desigion
        }
    }

    processDesigion(me, planet, bestDesigion)
    me.Stat.Actions[desigion.action] += 1
    planet.Stat.Actions[bestDesigion.action] += 1
}

func desigionIndividual(me *Individual, other *Individual, count uint) Desigion {
    var result Desigion
    // TODO: peace of ...
    result.needles = 0
    result.action = "nothing"

    if (isYoung(me) && isYoung(other) && count < 5 && //TODO: 5 - max terrain population
        me.Food > 2/*TODO: could change*/) {
        result.needles = 10 + rand.Intn(10)
        result.action = "sex"
    } else if (me.Food < 0) {
        if (isYoung(me) ||
            (isChild(me) && isChild(other)) ||
            (isOld(me) && !isYoung(other))) {
            result.needles = 10 - me.Food +
                calculateDeltaRandom(int(me.Health), int(other.Health)) +
                calculateDeltaRandom(int(other.Age), int(me.Age))
            result.action = "kill"
        }
    }

    result.withWho = other
    return result
}

func desigionFood(me *Individual, terrain *Terrain) Desigion {
    var result Desigion
    if (terrain.food <= 0) {
        result.needles = 0
    } else {
        result.needles = int(me.max_food) - me.Food +
                         int(terrain.food / 10) + rand.Intn(10)
    }
    result.action = "eat"
    return result
}

func calculateDesigion(me *Individual, terrain *Terrain) Desigion {
    var desigion Desigion
    var bestDesigion = Desigion{needles: 0, action: "nothing"}
    var count = uint(len(terrain.individuals))
    if (len(terrain.individuals) > 0) {
        for _, individual := range terrain.individuals {
            if (individual.index == me.index || individual.IsDed) {
                continue
            }
            desigion = desigionIndividual(me, individual, count)
            if (desigion.needles > bestDesigion.needles) {
                bestDesigion = desigion
            }
        }
    }
    desigion = desigionFood(me, terrain)
    if (desigion.needles >= bestDesigion.needles) {
        bestDesigion = desigion
    }

    bestDesigion.pos = terrain.pos
    return bestDesigion
}

func processDesigion(me *Individual, planet *Planet, desigion Desigion) {
    var terrain = &planet.grid[desigion.pos.x][desigion.pos.y]
    if (me.Pos != desigion.pos) {
        moveIndivid(me, planet, desigion.pos)
        me.Stat.Actions["move"] += 1
    }

    switch desigion.action {
    case "eat":
        var food = int(me.max_food) - me.Food
        if (food > terrain.food) {
            food = terrain.food
        }
        me.Food += food
        terrainIncFood(planet, terrain, -food)
    case "sex":
        me.Food -= 2 //TODO: could change
        desigion.withWho.Food -= 2
        planet.Individuals = append(planet.Individuals,
            createIndivid(len(planet.Individuals), desigion.pos))
        placeIndivid(&planet.Individuals[len(planet.Individuals) - 1], planet)
        planet.Population++;
    case "kill":
        var newHealth = me.Health - desigion.withWho.Health -
                uint(calculateDeltaDevided(
                    int(me.Age), int(desigion.withWho.Age), 10)) -
                uint(rand.Intn(10))

        if newHealth >= 0 {
            dieIndivid(desigion.withWho, planet)
            me.Health = newHealth
            me.Food -= 3
            me.Stat.Actions["kill"] += 1
            desigion.withWho.Stat.Actions["killed_by"] = me.index
        } else {
            dieIndivid(me, planet)
            desigion.withWho.Health = -newHealth
            desigion.withWho.Food -= 3
            me.Stat.Actions["killed_by"] = desigion.withWho.index
            desigion.withWho.Stat.Actions["kill"] += 1
            return
        }
    }

    if (me.Health < 1) {
        me.Health = 0
    } else if (me.Food <= 0) {
        me.Health -= 5
    } else if (me.Health < 100) {
        me.Health += 1
    }
}

func calculateDeltaRandom(value1 int, value2 int) int {
    delta := value1 - value2
    if (delta > 0) {
        delta = rand.Intn(delta)
    } else if (delta < 0) {
        delta = -rand.Intn(-delta)
    }
    return delta
}

func calculateDeltaDevided(value1 int, value2 int, devider int) int {
    return (value1 - value2) / devider
}
