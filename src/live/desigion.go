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
    me.Food > 2/*TODO: could change*/) {
        result.needles = 10 + rand.Intn(10)
        result.action = "sex"
    } else if (me.Food < 0) {
        if (isChild(me)) {
            result.needles = 0
            result.action = "nothing"
        } else {
            result.needles = 10 + rand.Intn(me.Health - other.Health) +
            rand.Intn(other.Age - me.Age)
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
        result.needles = me.max_food - me.Food + terrain.food + rand.Intn(10)
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
        var food = me.max_food - me.Food
        if (food > terrain.food) {
            food = terrain.food
        }
        me.Food += food
        terrain.food -= food
    case "sex":
        me.Food -= 2 //TODO: could change
        planet.Individuals = append(planet.Individuals,
            createIndivid(len(planet.Individuals), desigion.pos))
        placeIndivid(&planet.Individuals[len(planet.Individuals) - 1], planet)
        planet.Population++;
    case "kill":
        var newHealth = me.Health - desigion.withWho.Health -
        (me.Age - desigion.withWho.Age) % 10 - rand.Intn(10)
        if newHealth > 0 {
            if (desigion.withWho.Food > 0) {
                terrain.food += desigion.withWho.Food
            }
            terrain.food += desigion.withWho.Health % 10
            me.Health = newHealth
            desigion.withWho.Health = 0
        } else {
            if (me.Food > 0) {
                terrain.food += me.Food
            }
            terrain.food += me.Health % 10
            me.Health = 0
            desigion.withWho.Health = newHealth
        }
        me.Food -= 3
        planet.Population--;
    }

    if (me.Health < 1) {
        me.Health = 0
    } else if (me.Food <= 0) {
        me.Health -= 5
    } else if (me.Health < 100) {
        me.Health += 1
    }
    me.Stat.Actions[desigion.action] += 1
}
