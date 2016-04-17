package live

import (
    "math/rand"
)

type Location struct {
    x uint32
    y uint32
    width uint32
    height uint32
}

type Terrain struct {
    food int32
    individuals map[string]*Individual
    pos Pos
}

type Planet struct {
    Width       uint32
    Height      uint32
    Population  uint32
    grid        [][]Terrain
    Individuals []Individual
}


func CreatePlanet(population uint32) Planet {
    var planet Planet
    planet.Height = population * 10
    planet.Width = population * 10
    planet.grid = make([][]Terrain, planet.Width)
    planet.Individuals = make([]Individual, population)

    var terrain Terrain
    for x := 0; x < planet.Width; x++ {
        planet.grid[x] = make([]Terrain, planet.Height)
        for y := 0; y < planet.Height; y++ {
            terrain.individuals = make(map[string]*Individual)
            terrain.pos = Pos{x, y}
            terrain.food = rand.Int31n(2)
            planet.grid[x][y] = terrain
        }
    }

    var pos Pos
    for i := 0; i < population; i++ {
        pos.x = rand.Intn(planet.Width)
        pos.y = rand.Intn(planet.Height)
        planet.Individuals[i] = createIndivid(i, pos)
        placeIndivid(&planet.Individuals[i], &planet)
        planet.Population++;
    }
    return planet;
}

func takeNearestTerrains(planet *Planet, pos Pos, radius uint32) []*Terrain {
    var minX = pos.x - radius
    var maxX = pos.x + radius
    var minY = pos.y - radius
    var maxY = pos.y + radius
    if (minX < 0) {
        minX = 0
    }
    if (maxX >= planet.Width) {
        maxX = planet.Width - 1
    }
    if (minY < 0) {
        minY = 0
    }
    if (maxY >= planet.Height) {
        maxY = planet.Height - 1
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
