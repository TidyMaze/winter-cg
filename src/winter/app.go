package main

import (
	"fmt"
	"os"
)

/**
 * Grow and multiply your organisms to end up larger than your opponent.
 	Congratulations
Your organism can fight!

However, your organism has been alone so far. But with the power of a SPORER type organ, you can grow entirely new organisms.


The SPORER organ.
 	SPORER Rules
The SPORER type organ is unique in two ways:

It is the only organ that can create a new ROOT organ.
To create a new ROOT, it shoots out a spore in a straight line, letting you place the new organ in any of the free spaces it is facing.
Note: a ROOT organ never has a parent, even when spawned from a SPORER.



This command will make the SPORER shoot a new ROOT to the South.
When you control multiple organisms, you must output one command for each one. They will perform their actions simultaneously.


The requiredActionsCount variable will keep track of how many organisms you have. You must use the WAIT command for any organism that cannot act.


Note: You can use the organRootId variable to find out which organs belong to the same organism.


To grow a SPORER you need 1 B type protein and 1 D type protein.

To spore a new ROOT you need 1 of each protein.


Here is a table to summarize all organ costs:

Organ	A	B	C	D
BASIC	1	0	0	0
HARVESTER	0	0	1	1
TENTACLE	0	1	1	0
SPORER	0	1	0	1
ROOT	1	1	1	1
In this league, there is one protein source but your starting organism is not close enough to harvest it.

Use a sporer to shoot a new ROOT towards the protein and grow larger than your opponent!


New information added to the Game Protocol section.

 	TENTACLE Rules
On each turn, right after harvesting, any TENTACLE organs facing an opponent organ will attack, causing the target organ to die. Attacks happen simultaneously.

This command will create a new TENTACLE facing E (East), causing the opponent organ to be attacked.
When an organ dies, all of its children also die. This will propagate to the entire organism if the ROOT is hit.


Note: You can use the organParentId variable to keep track of each organ's children.


A tentacle also prevents the opponent from growing onto the tile it is facing.


To grow a TENTACLE you need 1 B type protein and 1 C type protein.


Use them to grow a large organism and attack the opponent's organism!
 	HARVESTER Rules

This command will create new HARVESTER facing N (North).
If a HARVESTER is facing a tile with a protein source, you will receive 1 of that protein on every end of turn.


Note: each player gains only 1 protein from each source per turn, even if multiple harvesters are facing that source.


To grow a HARVESTER you need 1 C type protein and 1 D type protein.

 	Rules
The game is played on a grid.

For the lower leagues, you need only beat the Boss in specific situations.


🔵🔴 The Organisms
Organisms are made up of organs that take up one tile of space on the game grid.


Each player starts with a ROOT type organ. Your organism can GROW a new organ on each turn in order to cover a larger area.


A new organ can grow from any existing organ, onto an empty adjacent location.


In order to GROW, your organism needs proteins. Growing 1 BASIC organ requires 1 protein of type A.


You can obtain more proteins by growing an organ onto a tile of the grid containing a protein source, these are tiles with a letter in them. Doing so will grant you 3 proteins of the corresponding type.


Grow more organs than the Boss to advance to the next league.


You organism can receive the following command:

GROW id x y type: creates a new organ at location x, y from organ with id id. If the target location is not a neighbour of id, the organ will be created on the shortest path to x, y.

This command will create new BASIC organ with the ROOT organ as its parent.

See the Game Protocol section for more information on sending commands to your organism.



⛔ Game end
The game stops when it detects progress can no longer be made or after 100 turns.


Victory Conditions
The winner is the player with the most tiles occupied by one of their organs.
Defeat Conditions
Your program does not provide a command in the alloted time or one of the commands is invalid.

🐞 Debugging tips
Hover over the grid to see extra information on the organ under your mouse.
Append text after any command and that text will appear above your organism.
Press the gear icon on the viewer to access extra display options.
Use the keyboard to control the action: space to play/pause, arrows to step 1 frame at a time.
Click to expand
 	Game Protocol
Initialization Input
First line: two integers width and height for the size of the grid.
Input for One Game Turn
First line: one integer entityCount for the number of entities on the grid.
Next entityCount lines: the following 7 inputs for each entity:
x: X coordinate (0 is leftmost)
y: Y coordinate (0 is topmost)
type:
WALL for a wall
ROOT for a ROOT type organ
BASIC for a BASIC type organ
HARVESTER for a HARVESTER type organ
TENTACLE for a TENTACLE type organ
SPORER for a SPORER type organ
A for an A protein source
B for a B protein source
C for a C protein source
D for a D protein source
owner:
1 if you are the owner of this organ
0 if your opponent owns this organ
-1 if this is not an organ
organId: unique id of this entity if it is an organ, 0 otherwise
organDir: N, W, S, or E for the direction in which this organ is facing
organParentId: if it is an organ, the organId of the organ that this organ grew from (0 for ROOT organs), else 0.
organRootId: if it is an organ, the organId of the ROOT that this organ originally grew from, else 0.
Next line: 4 integers: myA,myB,myC,myD for the amount of each protein type you have.
Next line: 4 integers: oppA,oppB,oppC,oppD for the amount of each protein type your opponent has.
Next line: the integer requiredActionsCount which equals the number of command you have to perform during the turn.
Output
A single line per organism with its action:
GROW id x y type direction : attempt to grow a new organ of type type at location x, y from organ with id id. If the target location is not a neighbour of id, the organ will be created on the shortest path to x, y.
SPORE id x y : attempt to create a new ROOT organ at location x, y from the SPORER with id id.
WAIT : do nothing.
Append text to your command and it will be displayed in the viewer.
Constraints
Response time per turn ≤ 50ms
Response time for the first turn ≤ 1000ms
 **/

type Coord struct {
	x, y int
}

func (c Coord) add(offset Coord) Coord {
	return Coord{c.x + offset.x, c.y + offset.y}
}

func (c Coord) isValid() bool {
	return c.x >= 0 && c.x < state.Width && c.y >= 0 && c.y < state.Height
}

type EntityType int

const (
	WALL EntityType = iota
	ROOT
	BASIC
	HARVESTER
	TENTACLE
	SPORER
	PROTEIN_A
)

type Dir int

const (
	N Dir = iota
	S
	W
	E
	NO_DIR
)

var offsets = []Coord{
	// N, S, W, E
	{0, -1},
	{0, 1},
	{-1, 0},
	{1, 0},
}

func findDirRelativeTo(from, to Coord) Dir {
	if from.x == to.x {
		if from.y < to.y {
			return S
		} else {
			return N
		}
	} else if from.y == to.y {
		if from.x < to.x {
			return E
		} else {
			return W
		}
	}
	panic(fmt.Sprintf("Unknown direction from %+v to %+v", from, to))
}

type Owner int

const (
	ME Owner = iota
	OPPONENT
	NONE
)

type Entity struct {
	coord         Coord
	_type         EntityType
	owner         Owner
	organId       int
	organDir      Dir
	organParentId int
	organRootId   int
}

type State struct {
	Height               int
	Width                int
	Entities             []Entity
	Grid                 [][]int
	MyProteins           []int
	OppProteins          []int
	RequiredActionsCount int
}

var state State

func parseDir(dir string) Dir {
	switch dir {
	case "N":
		return N
	case "S":
		return S
	case "E":
		return E
	case "W":
		return W
	default:
		return NO_DIR
	}
}

func showDir(dir Dir) string {
	switch dir {
	case N:
		return "N"
	case S:
		return "S"
	case E:
		return "E"
	case W:
		return "W"
	}
	panic(fmt.Sprintf("Unknown dir %d", dir))
}

func parseType(_type string) EntityType {
	switch _type {
	case "WALL":
		return WALL
	case "ROOT":
		return ROOT
	case "BASIC":
		return BASIC
	case "HARVESTER":
		return HARVESTER
	case "TENTACLE":
		return TENTACLE
	case "SPORER":
		return SPORER
	case "A":
		return PROTEIN_A
	}
	panic(fmt.Sprintf("Unknown type %s", _type))
}

func parseOwner(owner int) Owner {
	switch owner {
	case 1:
		return ME
	case 0:
		return OPPONENT
	case -1:
		return NONE
	}
	panic(fmt.Sprintf("Unknown owner %d", owner))
}

func debug(msg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, v...)
}

func distance(a, b Coord) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type SporePlan struct {
	organ          Entity // the organ from which to grow the sporer
	newSporerCoord Coord  // the coord of the sporer to grow
	sporerDir      Dir    // the direction of the sporer
	target         Coord  // the target coord of the sporer (either a protein or a neighbor of a protein)
}

func parseTurnState() {
	state.Grid = make([][]int, state.Height)
	for i := 0; i < state.Height; i++ {
		state.Grid[i] = make([]int, state.Width)
		for j := 0; j < state.Width; j++ {
			state.Grid[i][j] = -1
		}
	}

	var entityCount int
	fmt.Scan(&entityCount)

	state.Entities = make([]Entity, entityCount)

	for i := 0; i < entityCount; i++ {
		// y: grid coordinate
		// _type: WALL, ROOT, BASIC, TENTACLE, HARVESTER, SPORER, A, B, C, D
		// owner: 1 if your organ, 0 if enemy organ, -1 if neither
		// organId: id of this entity if it's an organ, 0 otherwise
		// organDir: N,E,S,W or X if not an organ
		var x, y int
		var _type string
		var owner, organId int
		var organDir string
		var organParentId, organRootId int
		fmt.Scan(&x, &y, &_type, &owner, &organId, &organDir, &organParentId, &organRootId)

		// debug("x: %d, y: %d, type: %s, owner: %d, organId: %d, organDir: %s, organParentId: %d, organRootId: %d\n", x, y, _type, owner, organId, organDir, organParentId, organRootId)

		entity := Entity{
			coord:         Coord{x, y},
			_type:         parseType(_type),
			owner:         parseOwner(owner),
			organId:       organId,
			organDir:      parseDir(organDir),
			organParentId: organParentId,
			organRootId:   organRootId,
		}

		state.Entities[i] = entity

		state.Grid[y][x] = i
	}

	// debug the entities
	// for _, entity := range state.Entities {
	// 	// debug("Entity: %+v\n", entity)
	// }

	// print the grid
	for i := 0; i < state.Height; i++ {
		for j := 0; j < state.Width; j++ {
			fmt.Fprintf(os.Stderr, "%d ", state.Grid[i][j])
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	state.MyProteins = make([]int, 4)
	state.OppProteins = make([]int, 4)

	// myD: your protein stock
	var myA, myB, myC, myD int
	fmt.Scan(&myA, &myB, &myC, &myD)

	debug("My proteins: A: %d, B: %d, C: %d, D: %d\n", myA, myB, myC, myD)

	state.MyProteins[0] = myA
	state.MyProteins[1] = myB
	state.MyProteins[2] = myC
	state.MyProteins[3] = myD

	// oppD: opponent's protein stock
	var oppA, oppB, oppC, oppD int
	fmt.Scan(&oppA, &oppB, &oppC, &oppD)

	debug("Opponent proteins: A: %d, B: %d, C: %d, D: %d\n", oppA, oppB, oppC, oppD)

	state.OppProteins[0] = oppA
	state.OppProteins[1] = oppB
	state.OppProteins[2] = oppC
	state.OppProteins[3] = oppD

	// requiredActionsCount: your number of organisms, output an action for each one in any order
	var requiredActionsCount int
	fmt.Scan(&requiredActionsCount)

	debug("Required actions count: %d\n", requiredActionsCount)

	state.RequiredActionsCount = requiredActionsCount
}

func sendActions() {
	for i := 0; i < state.RequiredActionsCount; i++ {
		// get the first root
		var root Entity
		for _, entity := range state.Entities {
			if entity._type == ROOT && entity.owner == ME {
				root = entity
				break
			}
		}

		debug("Root: %+v\n", root)

		// find all organs that have the organRootId equal to the root.organId
		var organs []Entity
		for _, entity := range state.Entities {
			if entity.organRootId == root.organId {
				organs = append(organs, entity)
			}

		}
		debug("Organs: %+v\n", organs)

		// find the non-harvested proteins
		var nonHarvestedProteins []Entity

		for _, entity := range state.Entities {
			if entity._type == PROTEIN_A {
				// find my neighbor harvesters of this protein (must be facing the protein)
				myHarvesters := make([]Entity, 0)
				for _, offset := range offsets {
					coord := entity.coord.add(offset)
					if coord.isValid() {
						if state.Grid[coord.y][coord.x] != -1 {
							neighbor := state.Entities[state.Grid[coord.y][coord.x]]
							if neighbor._type == HARVESTER && neighbor.owner == ME {
								if findDirRelativeTo(neighbor.coord, entity.coord) == neighbor.organDir {
									myHarvesters = append(myHarvesters, neighbor)
								} else {
									debug("Neighbor harvester %+v is not facing the protein %+v\n", neighbor, entity)
								}
							}
						}
					}
				}

				if len(myHarvesters) > 0 {
					debug("My harvesters for protein: %+v: %+v\n", entity, myHarvesters)
				} else {
					debug("No harvesters for protein: %+v\n", entity)
					nonHarvestedProteins = append(nonHarvestedProteins, entity)
				}
			}
		}

		debug("Non-harvested proteins: %+v\n", nonHarvestedProteins)

		if len(nonHarvestedProteins) > 0 {

			// build a map of the intersting spore cells (cells that are at distance max 1 from a non harvested protein)
			sporeCells := make([][]bool, state.Height)
			for i := 0; i < state.Height; i++ {
				sporeCells[i] = make([]bool, state.Width)
			}

			for _, protein := range nonHarvestedProteins {
				for _, offset := range offsets {
					coord := protein.coord.add(offset)
					if coord.isValid() && state.Grid[coord.y][coord.x] == -1 {
						sporeCells[coord.y][coord.x] = true
					}
				}
			}

			debug("Spore cells:\n")
			for i := 0; i < state.Height; i++ {
				for j := 0; j < state.Width; j++ {
					if sporeCells[i][j] {
						fmt.Fprintf(os.Stderr, "X ")
					} else {
						fmt.Fprintf(os.Stderr, "  ")
					}
				}
				fmt.Fprintf(os.Stderr, "\n")
			}

			// check if I have a sporer that can spore a new root into a spore cell
			sporer := Entity{}
			sporeCoord := Coord{-1, -1}
			for _, entity := range state.Entities {
				if entity._type == SPORER && entity.owner == ME {
					sporeCooord := findSporeCellInDirection(entity.coord, entity.organDir, sporeCells)
					if sporeCooord.isValid() {
						sporer = entity
						sporeCoord = sporeCooord
						break
					}
				}
			}

			if sporeCoord.isValid() {
				debug("Found a spore cell: %+v for sporer: %+v\n", sporeCoord, sporer)
				fmt.Printf("SPORE %d %d %d\n", sporer.organId, sporeCoord.x, sporeCoord.y)
			} else {

				if canGrow(state.MyProteins, SPORER) {
					// check if the closest organ can reach the closest protein using sporers

					debug("Can grow a sporer\n")

					// find any neighbor of my organs that can reach a spore cell in any direction

					sporerPlans := make([]SporePlan, 0)

					for _, organ := range organs {
						for _, offset := range offsets {
							coord := organ.coord.add(offset)
							if coord.isValid() && state.Grid[coord.y][coord.x] == -1 {
								// simulate the spore in all directions until it reaches a spore cell
								for _, dir := range []Dir{N, S, W, E} {
									sporeCoord := findSporeCellInDirection(coord, dir, sporeCells)

									if sporeCoord.isValid() {
										debug("Organ: %+v can reach spore cell: %+v after sporing in direction: %s from cell: %+v\n", organ, sporeCoord, showDir(dir), coord)
										sporerPlans = append(sporerPlans, SporePlan{
											organ:          organ,
											newSporerCoord: coord,
											sporerDir:      dir,
											target:         sporeCoord,
										})
									}
								}
							}
						}
					}

					if len(sporerPlans) > 0 {
						debug("Spore plans:\n")
						for _, plan := range sporerPlans {
							debug("Organ: %+v, new sporer coord: %+v, sporer dir: %s, target: %+v\n", plan.organ, plan.newSporerCoord, showDir(plan.sporerDir), plan.target)
						}

						// choose the best spore plan
						bestPlan := sporerPlans[0]

						// grow the sporer
						fmt.Printf("GROW %d %d %d SPORER %s\n", bestPlan.organ.organId, bestPlan.newSporerCoord.x, bestPlan.newSporerCoord.y, showDir(bestPlan.sporerDir))
					} else {
						debug("No spore plans\n")
					}
				} else {
					// find the protein that is the closest from any of the organs and that is not harvested
					var closestProtein Entity
					var closestOrgan Entity
					minDistance := 1000
					for _, entity := range state.Entities {
						if entity._type == PROTEIN_A && !isAlreadyHarvested(entity, nonHarvestedProteins) {
							for _, organ := range organs {
								dist := distance(entity.coord, organ.coord)
								if dist < minDistance {
									minDistance = dist
									closestProtein = entity
									closestOrgan = organ
								}
							}
						}
					}

					debug("Closest protein: %+v\n from organ: %+v\n", closestProtein, closestOrgan)

					// find the neighbor of the closest organ that is the closest to the closest protein
					var closestNeighbor Coord
					minDistanceFromNeighbor := 1000

					for _, offset := range offsets {
						coord := closestOrgan.coord.add(offset)
						if coord.isValid() {
							// never grow on a protein
							if state.Grid[coord.y][coord.x] == -1 {
								dist := distance(coord, closestProtein.coord)
								if dist < minDistanceFromNeighbor {
									minDistanceFromNeighbor = dist
									closestNeighbor = coord
								}
							}
						}
					}

					debug("Closest neighbor: %+v\n", closestNeighbor)

					if minDistanceFromNeighbor == 1 {
						// put a harvester facing the protein
						harvesterDir := N

						harvesterDir = findDirRelativeTo(closestNeighbor, closestProtein.coord)

						fmt.Printf("GROW %d %d %d HARVESTER %s\n", closestOrgan.organId, closestNeighbor.x, closestNeighbor.y, showDir(harvesterDir))
					} else {
						// grow a basic organ to get closer to the protein
						fmt.Printf("GROW %d %d %d BASIC\n", closestOrgan.organId, closestNeighbor.x,
							closestNeighbor.y)
					}
				}
			}
		} else {
			// there is no protein on the grid, find a cell that is at the frontier of players' organisms

			var enemyOrgans []Entity
			for _, entity := range state.Entities {
				if entity.owner == OPPONENT && entity._type == BASIC {
					enemyOrgans = append(enemyOrgans, entity)
				}
			}

			debug("Enemy organs: %+v\n", enemyOrgans)

			var bestCell Coord
			var bestOfMyOrgans Entity
			var bestOfEnemyOrgans Entity
			bestScore := -1000

			for i := 0; i < state.Height; i++ {
				for j := 0; j < state.Width; j++ {
					if state.Grid[i][j] == -1 {
						cell := Coord{j, i}

						// find the closest of my organs
						minDistance := 1000
						closestOfMyOrgans := Entity{}

						for _, organ := range organs {
							dist := distance(cell, organ.coord)
							if dist < minDistance {
								minDistance = dist
								closestOfMyOrgans = organ
							}
						}

						// find the closest of enemy organs
						minDistance = 1000
						closestOfEnemyOrgans := Entity{}

						for _, organ := range enemyOrgans {
							dist := distance(cell, organ.coord)
							if dist < minDistance {
								minDistance = dist
								closestOfEnemyOrgans = organ
							}
						}

						diffDist := distance(closestOfMyOrgans.coord, cell) - distance(closestOfEnemyOrgans.coord, cell)

						// only keep if the cell is closest to one of my organs than to any of the enemy organs (score < 0) but with a score that is the closest to 0 (frontier)
						if diffDist <= 0 {
							if diffDist > bestScore {
								bestScore = diffDist
								bestCell = cell
								bestOfMyOrgans = closestOfMyOrgans
								bestOfEnemyOrgans = closestOfEnemyOrgans
							}
						}
					}
				}
			}

			debug("Grow target cell: %+v from organ: %+v and enemy organ: %+v\n", bestCell, bestOfMyOrgans, bestOfEnemyOrgans)

			fmt.Printf("GROW %d %d %d BASIC\n", bestOfMyOrgans.organId, bestCell.x, bestCell.y)
		}
	}
}

func findSporeCellInDirection(coord Coord, dir Dir, sporeCells [][]bool) Coord {
	sporeCoord := coord
	for {
		sporeCoord = sporeCoord.add(offsets[dir])
		if !sporeCoord.isValid() {
			break
		}

		if state.Grid[sporeCoord.y][sporeCoord.x] != -1 {
			break
		}

		if sporeCells[sporeCoord.y][sporeCoord.x] {
			return sporeCoord
		}
	}

	return Coord{-1, -1}
}

func main() {
	// width: columns in the game grid
	// height: rows in the game grid
	fmt.Scan(&state.Width, &state.Height)

	for {
		parseTurnState()
		sendActions()
	}
}

func canGrow(proteinCounts []int, _type EntityType) bool {
	switch _type {
	case BASIC:
		return proteinCounts[0] >= 1
	case HARVESTER:
		return proteinCounts[2] >= 1 && proteinCounts[3] >= 1
	case TENTACLE:
		return proteinCounts[1] >= 1 && proteinCounts[2] >= 1
	case SPORER:
		return proteinCounts[1] >= 1 && proteinCounts[3] >= 1
	case ROOT:
		return proteinCounts[0] >= 1 && proteinCounts[1] >= 1 && proteinCounts[2] >= 1 && proteinCounts[3] >= 1
	}
	return false
}

func isAlreadyHarvested(entity Entity, nonHarvestedProteins []Entity) bool {
	for _, protein := range nonHarvestedProteins {
		if protein.coord.x == entity.coord.x && protein.coord.y == entity.coord.y {
			return false
		}
	}
	return true
}
