package live

import (
    "math/rand"
)

type Desigion struct {
    needles int32
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
}

func desigionIndividual(me *Individual, other *Individual, count uint32) Desigion {
    var result Desigion
    // TODO: peace of ...
    result.needles = 0
    result.action = "nothing"

    if (isYoung(me) && isYoung(other) && count < 5 && //TODO: 5 - max terrain // ain population
        me.Food > 2/*TODO: could change*/) {
        result.needles = 10 + rand.Int31n(10)
        result.action = "sex"
    } else if (me.Food < 0) {
        if (isYoung(me)) {
            deltaHealth := int32(me.Health - other.Health)
            if (deltaHealth > 0) {
                deltaHealth = rand.Int31n(deltaHealth)
            } else if (deltaHealth < 0) {
                deltaHealth = -rand.Int31n(-deltaHealth)
            }

            result.needles = 10 +
                calculateDeltaRandom(int32(me.Health), int32(other.Health))
                calculateDeltaRandom(i, i)
            deltaHealth +
                rand.Int31n(int32(other.Age - me.Age))

        }

        if (!isChild(me)) {
            result.needles = 10 + rand.Int31n(int32(me.Health - other.Health)) +
                rand.Int31n(int32(other.Age - me.Age))
            result.action = "kill"
        }
        // TODO: some conditions
    }

    result.withWho = other
    return result
}

func desigionFood(me *Individual, terrain *Terrain) Desigion {
    var result Desigion
    if (terrain.food <= 0) {
        result.needles = 0
    } else {
        result.needles = int32(me.max_food) - me.Food + terrain.food + rand.Int31n(10)
    }
    result.action = "eat"
    return result
}

func calculateDesigion(me *Individual, terrain *Terrain) Desigion {
    var desigion Desigion
    var bestDesigion = Desigion{needles: 0, action: "nothing"}
    var count = len(terrain.individuals)
    if (len(terrain.individuals) > 0) {
        for _, individual := range terrain.individuals {
            if (individual.Health == 0) {
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
    me.Stat.Actions[desigion.action] += 1

    switch desigion.action {
    case "eat":
        var food = int32(me.max_food) - me.Food
        if (food > terrain.food) {
            food = terrain.food
        }
        me.Food += food
        terrain.food -= food
    case "sex":
        me.Food -= 2 //TODO: could change
        desigion.withWho.Food -= 2
        planet.Individuals = append(planet.Individuals,
            createIndivid(len(planet.Individuals), desigion.pos))
        placeIndivid(&planet.Individuals[len(planet.Individuals) - 1], planet)
        planet.Population++;
    case "kill":
        var newHealth = me.Health - desigion.withWho.Health -
                (me.Age - desigion.withWho.Age) / 10 - rand.Intn(10)
        if newHealth > 0 {
            dieIndivid(desigion.withWho, planet)
            me.Health = newHealth
            me.Food -= 3
        } else {
            dieIndivid(me, planet)
            desigion.withWho.Health = newHealth
            me.Food -= 3
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

func calculateDeltaRandom(value1 int32, value2 int32) int32 {
    delta := value1 - value2
    if (delta > 0) {
        delta = rand.Int31n(delta)
    } else if (delta < 0) {
        delta = -rand.Int31n(-delta)
    }
    return delta
}
