package live

import (
    "math/rand"
)

type Location struct {
    x int
    y int
    width int
    height int
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


func CreatePlanet(population int) Planet {
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
