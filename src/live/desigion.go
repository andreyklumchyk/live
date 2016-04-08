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
