package live

func RunCycle(planet *Planet) {
    // some action here ...
    // posible actions:
    // move
    // find eat
    // create new individual with some mutations
    // die -1 to health
    // fidth for eat
    // adaptation for new location
    // go away from bad
    // want to live
    var count = len(planet.Individuals)
    var individual *Individual
    for i := 0; i < count; i++ {
        individual = &planet.Individuals[i]
        if (individual.IsDed) {
            continue
        }
        if (calculateDie(individual)) {
            dieIndivid(individual, planet)
            planet.Stat.Actions["die"] += 1
            continue
        }
        makeDesigion(individual, planet)
        individual.Age++
    }
}
