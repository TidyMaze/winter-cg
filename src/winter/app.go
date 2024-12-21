package main

import (
	"fmt"
	"os"
)

/**
 * Grow and multiply your organisms to end up larger than your opponent.
 	Congratulations
Your organism can grow!

However, protein sources on the grid are limited, and once you absorb them, they are gone. This is where the HARVESTER type organ comes in.


The HARVESTER organ.
 	HARVESTER Rules
From this league onwards, organs you place may be given a direction.

This command will create new HARVESTER facing N (North).
If a HARVESTER is facing a tile with a protein source, you will receive 1 of that protein on every end of turn.


Note: each player gains only 1 protein from each source per turn, even if multiple harvesters are facing that source.


To grow a HARVESTER you need 1 C type protein and 1 D type protein.


In this league, you are given an extra 1 C type protein and 1 D type protein, use them to grow a harvester at the correct location to grow your organism indefinitely!


New information added to the Game Protocol section.

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
A for an A protein source
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
Next line: the integer requiredActionsCount which equals 1 in this league.
Output
A single line with your action: GROW id x y type direction : attempt to grow a new organ of type type at location x, y from organ with id id. If the target location is not a neighbour of id, the organ will be created on the shortest path to x, y.

What is in store for me in the higher leagues?

The extra rules available in higher leagues are:
An organ type to attack your opponent
An organ type to spawn more organisms
 **/

type Coord struct {
	x, y int
}

type EntityType int

const (
	WALL EntityType = iota
	ROOT
	BASIC
	HARVESTER
	PROTEINE_A
)

type Dir int

const (
	N Dir = iota
	s
	W
	E
)

var offsets = []Coord{
	{0, -1},
	{0, 1},
	{-1, 0},
	{1, 0},
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

func parseDir(dir string) Dir {
	switch dir {
	case "N":
		return N
	case "S":
		return s
	case "E":
		return E
	case "W":
		return W
	}
	panic(fmt.Sprintf("Unknown dir %s", dir))
}

func showDir(dir Dir) string {
	switch dir {
	case N:
		return "N"
	case s:
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
	case "A":
		return PROTEINE_A
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

func main() {
	// width: columns in the game grid
	// height: rows in the game grid
	var width, height int
	fmt.Scan(&width, &height)

	for {

		grid := make([][]int, height)
		for i := 0; i < height; i++ {
			grid[i] = make([]int, width)
			for j := 0; j < width; j++ {
				grid[i][j] = -1
			}
		}

		var entityCount int
		fmt.Scan(&entityCount)

		entities := make([]Entity, entityCount)

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
				organDir:      N,
				organParentId: organParentId,
				organRootId:   organRootId,
			}

			entities[i] = entity

			grid[y][x] = i
		}

		// debug the entities
		// for _, entity := range entities {
		// 	// debug("Entity: %+v\n", entity)
		// }

		// print the grid
		for i := 0; i < height; i++ {
			for j := 0; j < width; j++ {
				fmt.Fprintf(os.Stderr, "%d ", grid[i][j])
			}
			fmt.Fprintf(os.Stderr, "\n")
		}

		myProteins := make([]int, 4)
		oppProteins := make([]int, 4)

		// myD: your protein stock
		var myA, myB, myC, myD int
		fmt.Scan(&myA, &myB, &myC, &myD)

		debug("My proteins: A: %d, B: %d, C: %d, D: %d\n", myA, myB, myC, myD)

		myProteins[0] = myA
		myProteins[1] = myB
		myProteins[2] = myC
		myProteins[3] = myD

		// oppD: opponent's protein stock
		var oppA, oppB, oppC, oppD int
		fmt.Scan(&oppA, &oppB, &oppC, &oppD)

		debug("Opponent proteins: A: %d, B: %d, C: %d, D: %d\n", oppA, oppB, oppC, oppD)

		oppProteins[0] = oppA
		oppProteins[1] = oppB
		oppProteins[2] = oppC
		oppProteins[3] = oppD

		// requiredActionsCount: your number of organisms, output an action for each one in any order
		var requiredActionsCount int
		fmt.Scan(&requiredActionsCount)

		debug("Required actions count: %d\n", requiredActionsCount)

		for i := 0; i < requiredActionsCount; i++ {

			// get the first root
			var root Entity
			for _, entity := range entities {
				if entity._type == ROOT && entity.owner == ME {
					root = entity
					break
				}
			}

			debug("Root: %+v\n", root)

			// find all organs that have the organRootId equal to the root.organId
			var organs []Entity
			for _, entity := range entities {
				if entity.organRootId == root.organId {
					organs = append(organs, entity)
				}
			}

			debug("Organs: %+v\n", organs)

			// find the prootein that is the closest from any of the organs
			var closestProteine Entity
			var closestOrgan Entity
			minDistance := 1000
			for _, entity := range entities {
				if entity._type == PROTEINE_A {
					for _, organ := range organs {
						dist := distance(entity.coord, organ.coord)
						if dist < minDistance {
							minDistance = dist
							closestProteine = entity
							closestOrgan = organ
						}
					}
				}
			}

			debug("Closest proteine: %+v\n", closestProteine)

			// find the neighbor of the closest organ that is the closest to the closest proteine
			var closestNeighbor Coord
			minDistance = 1000

			for _, offset := range offsets {
				coord := Coord{closestOrgan.coord.x + offset.x, closestOrgan.coord.y + offset.y}
				if coord.x >= 0 && coord.x < width && coord.y >= 0 && coord.y < height {
					if grid[coord.y][coord.x] == -1 || entities[grid[coord.y][coord.x]]._type == PROTEINE_A {
						dist := distance(coord, closestProteine.coord)
						if dist < minDistance {
							minDistance = dist
							closestNeighbor = coord
						}
					}
				}
			}

			debug("Closest neighbor: %+v\n", closestNeighbor)

			if minDistance == 1 {
				// put a harvester facing the proteine
				harvesterDir := N
				if closestNeighbor.x < closestProteine.coord.x {
					harvesterDir = E
				} else if closestNeighbor.x > closestProteine.coord.x {
					harvesterDir = W
				} else if closestNeighbor.y < closestProteine.coord.y {
					harvesterDir = s
				} else if closestNeighbor.y > closestProteine.coord.y {
					harvesterDir = N
				}

				fmt.Printf("GROW %d %d %d HARVESTER %s\n", closestOrgan.organId, closestNeighbor.x, closestNeighbor.y, showDir(harvesterDir))
			} else {
				// grow a basic organ
				fmt.Printf("GROW %d %d %d BASIC\n", closestOrgan.organId, closestNeighbor.x,
					closestNeighbor.y)
			}
		}
	}
}
