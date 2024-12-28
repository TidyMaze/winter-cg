package main

import (
	"fmt"
	"os"
	"sort"
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
	return c.x >= 0 && c.x < globalState.Width && c.y >= 0 && c.y < globalState.Height
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
	PROTEIN_B
	PROTEIN_C
	PROTEIN_D
)

func (t EntityType) isProtein() bool {
	return t == PROTEIN_A || t == PROTEIN_B || t == PROTEIN_C || t == PROTEIN_D
}

func (t EntityType) String() string {
	switch t {
	case WALL:
		return "WALL"
	case ROOT:
		return "ROOT"
	case BASIC:
		return "BASIC"
	case HARVESTER:
		return "HARVESTER"
	case TENTACLE:
		return "TENTACLE"
	case SPORER:
		return "SPORER"
	case PROTEIN_A:
		return "A"
	case PROTEIN_B:
		return "B"
	case PROTEIN_C:
		return "C"
	case PROTEIN_D:
		return "D"
	}
	panic(fmt.Sprintf("Unknown type %d", t))
}

type Dir int

func (d Dir) String() string {
	return showDir(d)
}

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
	if from.x == to.x && from.y == to.y {
		panic(fmt.Sprintf("Same coord %+v", from))
	}

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

func findApproximateDir(from, to Coord) Dir {
	if from.x == to.x && from.y == to.y {
		panic(fmt.Sprintf("Same coord %+v", from))
	}

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
	} else if from.x < to.x {
		return E
	} else {
		return W
	}
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

func showOwner(owner Owner) string {
	switch owner {
	case ME:
		return "ME"
	case OPPONENT:
		return "OPPONENT"
	case NONE:
		return "NONE"
	}
	panic(fmt.Sprintf("Unknown owner %d", owner))
}

func (e Entity) String() string {
	switch e._type {
	case WALL:
		return fmt.Sprintf("Wall at %+v", e.coord)
	case ROOT:
		return fmt.Sprintf("Root at %+v, owner: %s, organId: %d, organParentId: %d, organRootId: %d", e.coord, showOwner(e.owner), e.organId, e.organParentId, e.organRootId)
	case BASIC:
		return fmt.Sprintf("Basic at %+v, owner: %s, organId: %d, organParentId: %d, organRootId: %d", e.coord, showOwner(e.owner), e.organId, e.organParentId, e.organRootId)
	case HARVESTER:
		return fmt.Sprintf("Harvester %s at %+v, owner: %s, organId: %d, organParentId: %d, organRootId: %d", showDir(e.organDir), e.coord, showOwner(e.owner), e.organId, e.organParentId, e.organRootId)
	case TENTACLE:
		return fmt.Sprintf("Tentacle %s at %+v, owner: %s, organId: %d, organParentId: %d, organRootId: %d", showDir(e.organDir), e.coord, showOwner(e.owner), e.organId, e.organParentId, e.organRootId)
	case SPORER:
		return fmt.Sprintf("Sporer %s at %+v, owner: %s, organId: %d, organParentId: %d, organRootId: %d", showDir(e.organDir), e.coord, showOwner(e.owner), e.organId, e.organParentId, e.organRootId)
	case PROTEIN_A:
		return fmt.Sprintf("Protein A at %+v", e.coord)
	case PROTEIN_B:
		return fmt.Sprintf("Protein B at %+v", e.coord)
	case PROTEIN_C:
		return fmt.Sprintf("Protein C at %+v", e.coord)
	case PROTEIN_D:
		return fmt.Sprintf("Protein D at %+v", e.coord)
	default:
		panic(fmt.Sprintf("Unknown entity type %d", e._type))
	}
}

type State struct {
	Height               int
	Width                int
	Entities             []Entity
	Grid                 [][]*Entity
	MyProteins           []int
	OppProteins          []int
	RequiredActionsCount int
}

func (s State) isWalkable(coord Coord, allowOrgans bool) bool {

	walkableEntity := false
	if allowOrgans {
		entity := s.Grid[coord.y][coord.x]

		if entity != nil {
			walkableEntity = entity._type.isProtein() ||
				entity._type == ROOT ||
				entity._type == BASIC ||
				entity._type == HARVESTER ||
				entity._type == TENTACLE ||
				entity._type == SPORER
		}
	} else {
		walkableEntity = s.Grid[coord.y][coord.x] == nil || s.Grid[coord.y][coord.x]._type.isProtein()
	}

	return s.Grid[coord.y][coord.x] == nil || walkableEntity
}

var globalState State

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
	default:
		panic(fmt.Sprintf("Unknown dir %d", dir))
	}
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
	case "B":
		return PROTEIN_B
	case "C":
		return PROTEIN_C
	case "D":
		return PROTEIN_D
	}
	panic(fmt.Sprintf("Unknown type %s", _type))
}

func showOrganType(_type EntityType) string {
	switch _type {
	case BASIC:
		return "BASIC"
	case HARVESTER:
		return "HARVESTER"
	case TENTACLE:
		return "TENTACLE"
	case SPORER:
		return "SPORER"
	default:
		panic(fmt.Sprintf("Unknown organ type %d", _type))
	}
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
	globalState.Grid = make([][]*Entity, globalState.Height)
	for i := 0; i < globalState.Height; i++ {
		globalState.Grid[i] = make([]*Entity, globalState.Width)
		for j := 0; j < globalState.Width; j++ {
			globalState.Grid[i][j] = nil
		}
	}

	var entityCount int
	fmt.Scan(&entityCount)

	globalState.Entities = make([]Entity, entityCount)

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

		if entity._type == ROOT {
			debug("Root: %+v\n", entity)

			if entity.organId == 0 {
				panic(fmt.Sprintf("Root with id 0: %+v", entity))
			}
		}

		globalState.Entities[i] = entity

		globalState.Grid[y][x] = &globalState.Entities[i]
	}

	// debug the entities
	// for _, entity := range globalState.Entities {
	// 	// debug("Entity: %+v\n", entity)
	// }

	// print the grid
	// for i := 0; i < globalState.Height; i++ {
	// 	for j := 0; j < globalState.Width; j++ {
	// 		fmt.Fprintf(os.Stderr, "%d ", globalState.Grid[i][j])
	// 	}
	// 	fmt.Fprintf(os.Stderr, "\n")
	// }

	globalState.MyProteins = make([]int, 4)
	globalState.OppProteins = make([]int, 4)

	// myD: your protein stock
	var myA, myB, myC, myD int
	fmt.Scan(&myA, &myB, &myC, &myD)

	debug("My proteins: A: %d, B: %d, C: %d, D: %d\n", myA, myB, myC, myD)

	globalState.MyProteins[0] = myA
	globalState.MyProteins[1] = myB
	globalState.MyProteins[2] = myC
	globalState.MyProteins[3] = myD

	// oppD: opponent's protein stock
	var oppA, oppB, oppC, oppD int
	fmt.Scan(&oppA, &oppB, &oppC, &oppD)

	debug("Opponent proteins: A: %d, B: %d, C: %d, D: %d\n", oppA, oppB, oppC, oppD)

	globalState.OppProteins[0] = oppA
	globalState.OppProteins[1] = oppB
	globalState.OppProteins[2] = oppC
	globalState.OppProteins[3] = oppD

	// requiredActionsCount: your number of organisms, output an action for each one in any order
	var requiredActionsCount int
	fmt.Scan(&requiredActionsCount)

	debug("Required actions count: %d\n", requiredActionsCount)

	globalState.RequiredActionsCount = requiredActionsCount
}

func findOrgansOfOrganism(root Entity) []Entity {
	// find all organs that have the organRootId equal to the root.organId
	var organs []Entity
	for _, entity := range globalState.Entities {
		if entity.organRootId == root.organId {
			organs = append(organs, entity)
		}

	}
	// debug("Organs: %+v\n", organs)
	return organs
}

type Action interface {
	getRootOrganId() int
	getMessage() string
	getStringCommand() string
}

type GrowAction struct {
	rootOrganId int
	organId     int
	coord       Coord
	_type       EntityType
	dir         Dir
	message     string
}

func (a GrowAction) getRootOrganId() int {
	return a.rootOrganId
}

func (a GrowAction) getMessage() string {
	return a.message
}

func (a GrowAction) getStringCommand() string {
	return fmt.Sprintf("GROW %d %d %d %s %s %s", a.organId, a.coord.x, a.coord.y, showOrganType(a._type), showDir(a.dir), a.message)
}

func (a GrowAction) String() string {
	return fmt.Sprintf("Grow %s at %+v, dir: %s, message: %s", showOrganType(a._type), a.coord, showDir(a.dir), a.message)
}

type WaitAction struct {
	rootOrganId int
	message     string
}

func (a WaitAction) getRootOrganId() int {
	return a.rootOrganId
}

func (a WaitAction) getMessage() string {
	return a.message
}

func (a WaitAction) getStringCommand() string {
	return fmt.Sprintf("WAIT %s", a.message)
}

func (a WaitAction) String() string {
	return fmt.Sprintf("Wait, message: %s", a.message)
}

type SporeAction struct {
	rootOrganId int
	sporerId    int
	coord       Coord
	message     string
}

func (a SporeAction) getRootOrganId() int {
	return a.rootOrganId
}

func (a SporeAction) getMessage() string {
	return a.message
}

func (a SporeAction) getStringCommand() string {
	return fmt.Sprintf("SPORE %d %d %d %s", a.sporerId, a.coord.x, a.coord.y, a.message)
}

func (a SporeAction) String() string {
	return fmt.Sprintf("Spore at %+v with sporer %d, message: %s", a.coord, a.sporerId, a.message)
}

func sendActions() {
	// find all roots
	var roots []Entity
	for _, entity := range globalState.Entities {
		if entity._type == ROOT && entity.owner == ME {
			roots = append(roots, entity)
		}
	}

	if len(roots) != globalState.RequiredActionsCount {
		panic(fmt.Sprintf("Expected %d roots, found %d", globalState.RequiredActionsCount, len(roots)))
	}

	enemyTentaclesTargets := findEnemyTentaclesTargets(globalState)

	// find the non-harvested proteins
	//nonHarvestedProteins := findNonHarvestedProteins()

	actions := findBestActions(globalState, roots, enemyTentaclesTargets).actions

	for i := 0; i < globalState.RequiredActionsCount; i++ {
		// get the first root
		var root Entity = roots[i]

		debug("=== Root: %+v ===\n", root)

		for _, a := range actions {
			if a.getRootOrganId() == root.organId {
				fmt.Println(a.getStringCommand())
				break
			}
		}
	}
}

/*
*
  - Get all combinations of slices (ex: [[1,2,3],[a,b,c],[x,y,z]]
    => [[1,a,x],[1,a,y],[1,a,z],[1,b,x],[1,b,y],[1,b,z],[1,c,x],[1,c,y],[1,c,z],[2,a,x],[2,a,y],[2,a,z],[2,b,x],[2,b,y],[2,b,z],[2,c,x],[2,c,y],[2,c,z],[3,a,x],[3,a,y],[3,a,z],[3,b,x],[3,b,y],[3,b,z],[3,c,x],[3,c,y],[3,c,z]]
    )
*/
func allCombinationsOfSlices[T any](slices [][]T) [][]T {
	if len(slices) == 0 {
		return [][]T{{}}
	}

	if len(slices) == 1 {
		perms := make([][]T, 0)
		for _, item := range slices[0] {
			perms = append(perms, []T{item})
		}
		return perms
	}

	perms := make([][]T, 0)
	for _, item := range slices[0] {
		for _, perm := range allCombinationsOfSlices[T](slices[1:]) {
			perms = append(perms, append([]T{item}, perm...))
		}
	}

	return perms
}

type PlayerActions struct {
	actions []Action
	score   float64
	detail  string
}

func findBestActions(s State, roots []Entity, enemyTentaclesTargets [][]bool) PlayerActions {
	allActionsCombinations := make([]PlayerActions, 0)

	debug("Permuting roots %+v\n", roots)

	// find all permutations of the N roots (ex: [[1,2,3], [1,3,2], [2,1,3], [2,3,1], [3,1,2], [3,2,1]])
	rootPermutations := permute(roots)

	debug("Root permutations (%d):\n", len(rootPermutations))

	// for each root, find all possible actions for each organ

	actionsPerRoot := make(map[int][]Action)

	for _, root := range roots {
		// find all organs of the root
		organs := findOrgansOfOrganism(root)

		// find all possible actions for each organ
		actionsPerRoot[root.organId] = findActionsForOrganism(s, root, organs, enemyTentaclesTargets)
	}

	for iPerm, rootPermutation := range rootPermutations {
		debug("Root permutation #%d: %+v\n", iPerm, rootPermutation)

		// once we know the order of roots, we can combine all the actions of each root
		actionsPerRootSorted := make([][]Action, 0)
		for _, root := range rootPermutation {
			actionsPerRootSorted = append(actionsPerRootSorted, actionsPerRoot[root.organId])
		}

		debug("Actions per root sorted (%d):\n", len(actionsPerRootSorted))
		for iRoot, actions := range actionsPerRootSorted {
			debug("Root %d: %d actions\n", iRoot, len(actions))
			//for iAction, action := range actions {
			//	debug("Action %d: %+v\n", iAction, action)
			//}
		}

		combinations := allCombinationsOfSlices(actionsPerRootSorted)

		debug("Combinations (%d)\n", len(combinations))

		// calculate the score of state after applying all the actions
		playerActions := make([]PlayerActions, 0)

		// protein map is built from the current state (stable)
		harvested, nonHarvested := findHarvestedProteins(s)

		proteinMap := buildProteinMap(s, nonHarvested, harvested)

		debug("Protein map:\n%s", showProteinMap(proteinMap))

		disputedCellsMap := findDisputedCells(s)

		debug("Disputed cells map:\n%s", showDisputedCellsMap(disputedCellsMap))

		for _, actions := range combinations {
			//debug("%d actions for comb (%d), ", len(actions), iComb)
			score, detail := scoreActions(globalState, actions, proteinMap, disputedCellsMap)

			playerActions = append(playerActions, PlayerActions{
				actions: actions,
				score:   score,
				detail:  detail,
			})
		}

		allActionsCombinations = append(allActionsCombinations, playerActions...)
	}

	//debug("All actions\n")
	//for i, actions := range allActionsCombinations {
	//	debug("Combination %d\n", i)
	//	for _, action := range actions.actions {
	//		debug("%+v\n", action)
	//	}
	//}

	// find the best combination of actions
	// sorted by score

	sort.Slice(allActionsCombinations, func(i, j int) bool {
		return allActionsCombinations[i].score > allActionsCombinations[j].score
	})

	// print top N combinations with score

	topN := 50
	debug("Top %d combinations\n", topN)
	for i, actions := range allActionsCombinations {
		if i < topN {
			actionsAsStr := ""

			for _, action := range actions.actions {
				actionsAsStr += fmt.Sprintf("%+v, ", action)
			}

			debug("Combination %d, score: %f (%s) %s\n", i, actions.score, actions.detail, actionsAsStr)
		}
	}

	if len(allActionsCombinations) == 0 {
		panic("No actions found")
	}

	return allActionsCombinations[0]
}

func showDisputedCellsMap(cellsMap [][]bool) interface{} {
	str := ""
	for i := 0; i < globalState.Height; i++ {
		for j := 0; j < globalState.Width; j++ {
			if cellsMap[i][j] {
				str += "X "
			} else {
				str += ". "
			}
		}
		str += "\n"
	}
	return str
}

func showProteinMap(proteinMap [][]float64) string {
	str := ""
	for i := 0; i < globalState.Height; i++ {
		for j := 0; j < globalState.Width; j++ {
			str += fmt.Sprintf("%d ", int(proteinMap[i][j]))
		}
		str += "\n"
	}
	return str

}

func scoreActions(s State, actions []Action, proteinMap [][]float64, disputedCellsMap [][]bool) (float64, string) {
	newState := applyActions(s, actions)

	return scoreState(newState, proteinMap, disputedCellsMap)
}

func scoreState(s State, proteinsMap [][]float64, disputedCellsMap [][]bool) (float64, string) {
	// score is the number of harvested proteins plus the number of organs
	harvested, nonHarvested := findHarvestedProteins(s)

	myOrgans := findOrgans(s, ME)
	enemyOrgans := findOrgans(s, OPPONENT)

	//enemyTentaclesTargets := findEnemyTentaclesTargets(s)

	// find the distance from any of my organs to the closest non-harvested protein (malus for being far)
	//path := findShortestPathProt(s, myOrgans, nonHarvested, enemyTentaclesTargets)
	//distanceClosestProtein := len(path)

	// for each of my organs, sum the protein map distance value
	totalDistance := 0.0
	organCount := 0

	for _, organ := range myOrgans {
		totalDistance += proteinsMap[organ.coord.y][organ.coord.x]
		organCount++
	}

	//pathStr := ""
	//for _, coord := range path {
	//	pathStr += fmt.Sprintf("%+v, ", coord)
	//}

	// better to have more proteins left (do not waste them to move)
	proteinScore := s.MyProteins[0] + s.MyProteins[1] + s.MyProteins[2] + s.MyProteins[3]

	for iProt := 0; iProt < 4; iProt++ {
		if s.MyProteins[iProt] <= 0 {
			proteinScore -= 1000
		}
	}

	avgDistance := totalDistance / float64(organCount)

	defendedDisputedCells := findDefendedDisputedCells(s, disputedCellsMap)

	//detailScore := fmt.Sprintf("Score detail: harvested: %d, non-harvested: %d, my organs: %d, enemy organs: %d, distance to closest protein: %d (%s), protein score: %d\n", len(harvested), len(nonHarvested), len(myOrgans), len(enemyOrgans), distanceClosestProtein, pathStr, proteinScore)
	detailScore := fmt.Sprintf("Score detail: harvested: %d, non-harvested: %d, total distance: %f, avgDistance: %f\n, my organs: %d, enemy organs: %d, protein score: %d\n, defended cells: %d", len(harvested), len(nonHarvested), totalDistance, avgDistance, len(myOrgans), len(enemyOrgans), proteinScore, len(defendedDisputedCells))

	//return float64(len(harvested)*100 - len(nonHarvested) + len(myOrgans)*100 - len(enemyOrgans)*100 - distanceClosestProtein*10 + proteinScore), detailScore
	return float64(len(harvested)*1000) - float64(len(nonHarvested)*10) - avgDistance + float64(len(myOrgans)*10000) - float64(len(enemyOrgans)*10000) + float64(proteinScore) + float64(len(defendedDisputedCells))*2000, detailScore
}

func findDefendedDisputedCells(s State, disputedCellsMap [][]bool) []Coord {
	cells := make([]Coord, 0)

	for _, entity := range s.Entities {
		if entity.owner == ME && entity._type == TENTACLE {
			// find the target of the tentacle
			target := entity.coord.add(offsets[entity.organDir])

			if target.isValid() {
				if disputedCellsMap[target.y][target.x] {
					cells = append(cells, target)
				}
			}
		}
	}

	return cells
}

func findProteins(s State) []Entity {
	proteins := make([]Entity, 0)
	for _, entity := range s.Entities {
		if entity._type.isProtein() {
			proteins = append(proteins, entity)
		}
	}
	return proteins
}

func buildDistanceMapForProtein(s State, protein Coord) [][]int {
	// build the grid of distances from the protein
	distances := make([][]int, s.Height)

	for i := 0; i < s.Height; i++ {
		distances[i] = make([]int, s.Width)
		for j := 0; j < s.Width; j++ {
			distances[i][j] = -1
		}
	}

	distances[protein.y][protein.x] = 0

	queue := make([]Coord, 0)
	queue = append(queue, protein)

	for len(queue) > 0 {
		coord := queue[0]
		queue = queue[1:]

		for _, offset := range offsets {
			neighbor := coord.add(offset)

			if neighbor.isValid() && s.isWalkable(neighbor, true) && distances[neighbor.y][neighbor.x] == -1 {
				distances[neighbor.y][neighbor.x] = distances[coord.y][coord.x] + 1
				queue = append(queue, neighbor)
			}
		}
	}

	return distances
}

func computeTurnIncome(harvestedProteins []Entity) []float64 {
	turnIncome := make([]float64, 4)

	// calculate the income of the turn
	for _, protein := range harvestedProteins {
		turnIncome[protein._type-PROTEIN_A]++
	}

	debug("Turn income: %+v\n", turnIncome)

	return turnIncome
}

/*
*
Normalize the array to have values between 0 and 1
*/
func normalizeArray(arr []float64) []float64 {
	normalized := make([]float64, len(arr))

	max := 0.0
	min := 1000000000.0

	for _, v := range arr {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}

	if max == min {
		return arr
	}

	for i, v := range arr {
		normalized[i] = (v - min) / (max - min)
	}

	return normalized
}

/**
 * For each protein, build the grid of distances from each cell to the protein.
 * Then sum the distances for each cell (merge the grids).
 */
func buildProteinMap(s State, nonHarvestedProteins []Entity, harvestedProteins []Entity) [][]float64 {

	turnIncome := computeTurnIncome(harvestedProteins)

	normalizedTurnIncome := normalizeArray(turnIncome)

	debug("Normalized turn income: %+v\n", normalizedTurnIncome)

	finalMap := make([][]float64, s.Height)

	for i := 0; i < s.Height; i++ {
		finalMap[i] = make([]float64, s.Width)
		for j := 0; j < s.Width; j++ {
			finalMap[i][j] = 0
		}
	}

	for _, protein := range nonHarvestedProteins {
		singleProteinMap := buildDistanceMapForProtein(s, protein.coord)

		// add the distances to the final map
		for i := 0; i < s.Height; i++ {
			for j := 0; j < s.Width; j++ {
				finalMap[i][j] += float64(singleProteinMap[i][j]) * (1 + normalizedTurnIncome[protein._type-PROTEIN_A])
			}
		}
	}

	return finalMap
}

// my organs (any root)
func findOrgans(s State, o Owner) []Entity {
	organs := make([]Entity, 0)
	for _, entity := range s.Entities {
		if entity.owner == o {
			organs = append(organs, entity)
		}
	}
	return organs
}

func applyActions(s State, actions []Action) State {
	// copy the state
	newState := copyState(s)

	growCoords := make([][]bool, globalState.Height)

	for i := 0; i < globalState.Height; i++ {
		growCoords[i] = make([]bool, globalState.Width)
		for j := 0; j < globalState.Width; j++ {
			growCoords[i][j] = false
		}
	}

	// apply each action in order
	for _, action := range actions {
		switch a := action.(type) {
		case GrowAction:
			// grow the organ
			newOrganId := maxOrganId(newState.Entities) + 1

			newEntity := Entity{
				coord:         a.coord,
				_type:         a._type,
				owner:         ME,
				organId:       newOrganId,
				organDir:      a.dir,
				organParentId: a.organId,
				organRootId:   a.rootOrganId,
			}

			oldEntityAtCoord := newState.Grid[a.coord.y][a.coord.x]
			if oldEntityAtCoord != nil {

				if oldEntityAtCoord.coord != a.coord {
					panic(fmt.Sprintf("Old entity coord %+v different from new coord %+v", oldEntityAtCoord.coord, a.coord))
				}

				if oldEntityAtCoord._type != PROTEIN_A && oldEntityAtCoord._type != PROTEIN_B && oldEntityAtCoord._type != PROTEIN_C && oldEntityAtCoord._type != PROTEIN_D {
					panic(fmt.Sprintf("Entity at %+v is not a protein", a.coord))
				}

				newState.Grid[a.coord.y][a.coord.x] = nil

				// remove the old entity from the entities list
				newState.Entities = removeEntity(newState.Entities, *oldEntityAtCoord)
			}

			newState.Entities = append(newState.Entities, newEntity)
			newState.Grid[a.coord.y][a.coord.x] = &newEntity

			//debug("Grew organ %+v\n", newEntity)

			if growCoords[a.coord.y][a.coord.x] {
				panic(fmt.Sprintf("Already grew organ at %+v", a.coord))
			}

			// kill neighbors of tentacles
			if a._type == TENTACLE {
				neighborCoord := a.coord.add(offsets[a.dir])
				if neighborCoord.isValid() {
					neighborEntity := newState.Grid[neighborCoord.y][neighborCoord.x]
					if neighborEntity != nil && neighborEntity.owner == OPPONENT {

						// get all children of the killed entity
						destroyedChildren := findDestroyed(newState, *neighborEntity)

						if len(destroyedChildren) == 0 {
							panic(fmt.Sprintf("No destroyed children by playing action %+v: %+v", a, destroyedChildren))
						}

						for _, destroyedChild := range destroyedChildren {
							newState.Grid[destroyedChild.coord.y][destroyedChild.coord.x] = nil
							newState.Entities = removeEntity(newState.Entities, destroyedChild)
						}
					}
				}
			}

			// apply the grow cost to my proteins
			growCost := growCost(a._type)
			newState.MyProteins[0] -= growCost.costA
			newState.MyProteins[1] -= growCost.costB
			newState.MyProteins[2] -= growCost.costC
			newState.MyProteins[3] -= growCost.costD
		case WaitAction:
			// do nothing
		case SporeAction:
			panic(fmt.Sprintf("Spore action not implemented for action %+v", a))
		}
	}

	return newState
}

func removeEntity(entities []Entity, entity Entity) []Entity {
	// find index of the entity
	index := -1

	for i, e := range entities {
		if e.coord == entity.coord {
			index = i
			break
		}
	}

	if index == -1 {
		panic(fmt.Sprintf("Entity %+v not found in entities", entity))
	}

	// remove the entity
	entities[index] = entities[len(entities)-1]
	return entities[:len(entities)-1]
}

func maxOrganId(entities []Entity) int {
	maxId := -1
	for _, entity := range entities {
		if entity.organId > maxId {
			maxId = entity.organId
		}
	}
	return maxId
}

func copyState(s State) State {
	newState := State{
		Height:               s.Height,
		Width:                s.Width,
		Entities:             make([]Entity, len(s.Entities)),
		Grid:                 make([][]*Entity, s.Height),
		MyProteins:           make([]int, 4),
		OppProteins:          make([]int, 4),
		RequiredActionsCount: s.RequiredActionsCount,
	}

	// copy entities
	for i, entity := range s.Entities {
		// entity is copied
		newState.Entities[i] = entity
	}

	// copy grid
	for i := 0; i < s.Height; i++ {
		newState.Grid[i] = make([]*Entity, s.Width)
		for j := 0; j < s.Width; j++ {
			// pointer is copied (same ref)
			newState.Grid[i][j] = s.Grid[i][j]
		}
	}

	// copy proteins
	newState.MyProteins[0] = s.MyProteins[0]
	newState.MyProteins[1] = s.MyProteins[1]
	newState.MyProteins[2] = s.MyProteins[2]
	newState.MyProteins[3] = s.MyProteins[3]

	newState.OppProteins[0] = s.OppProteins[0]
	newState.OppProteins[1] = s.OppProteins[1]
	newState.OppProteins[2] = s.OppProteins[2]
	newState.OppProteins[3] = s.OppProteins[3]

	return newState
}

func emptyMapBool(height, width int) [][]bool {
	m := make([][]bool, height)
	for i := 0; i < height; i++ {
		m[i] = make([]bool, width)
		for j := 0; j < width; j++ {
			m[i][j] = false
		}
	}
	return m
}

/*
*
Find all the empty cells that are:
- less than 2 cells away from player (path length <= 3)
- less than 2 cells away from enemy (path length <= 3)
*/
func findDisputedCells(s State) [][]bool {
	disputedCells := make([][]bool, s.Height)

	for i := 0; i < s.Height; i++ {
		disputedCells[i] = make([]bool, s.Width)
		for j := 0; j < s.Width; j++ {
			disputedCells[i][j] = false
		}
	}

	myOrgans := findOrgans(s, ME)
	enemyOrgans := findOrgans(s, OPPONENT)

	myOrgansCoords := make([]Coord, 0)
	for _, organ := range myOrgans {
		myOrgansCoords = append(myOrgansCoords, organ.coord)
	}

	enemyOrgansCoords := make([]Coord, 0)
	for _, organ := range enemyOrgans {
		enemyOrgansCoords = append(enemyOrgansCoords, organ.coord)
	}

	forbiddenCells := make([][]bool, s.Height)
	for i := 0; i < s.Height; i++ {
		forbiddenCells[i] = make([]bool, s.Width)
		for j := 0; j < s.Width; j++ {
			forbiddenCells[i][j] = false
		}
	}

	// mark walls as forbidden
	for _, entity := range s.Entities {
		if entity._type == WALL {
			forbiddenCells[entity.coord.y][entity.coord.x] = true
		}
	}

	for i := 0; i < s.Height; i++ {
		for j := 0; j < s.Width; j++ {
			coord := Coord{j, i}

			if s.isWalkable(coord, false) {
				// find the distance from any of my organs to the cell
				pathMy := findShortestPath(s, myOrgansCoords, []Coord{coord}, forbiddenCells)
				pathEnemy := findShortestPath(s, enemyOrgansCoords, []Coord{coord}, forbiddenCells)

				if len(pathMy) <= 3 && len(pathEnemy) <= 3 {
					disputedCells[i][j] = true
				}
			}
		}
	}

	return disputedCells
}

func findActionsForOrganism(s State, root Entity, organs []Entity, enemyTentaclesTargets [][]bool) []Action {
	actions := make([]Action, 0)

	for _, organ := range organs {
		// find all possible actions for the organ
		actionsForOrgan := findActionsForOrgan(s, root, organ, enemyTentaclesTargets)
		actions = append(actions, actionsForOrgan...)
		//debug("Actions for organ %+v: %d\n", organ.organId, len(actions))
	}

	actions = append(actions, WaitAction{root.organId, ""})

	return actions
}

func findActionsForOrgan(s State, root, organ Entity, enemyTentaclesTargets [][]bool) []Action {
	actions := make([]Action, 0)

	// find the grow actions
	growActions := findGrowActions(s, root, organ, enemyTentaclesTargets)
	actions = append(actions, growActions...)

	// find the spore actions
	sporeActions := findSporeActions(root, organ, enemyTentaclesTargets)
	actions = append(actions, sporeActions...)

	return actions
}

func findSporeActions(root Entity, organ Entity, enemyTentaclesTargets [][]bool) []Action {
	return make([]Action, 0)
}

func findWaitActions(root Entity, organ Entity) []Action {
	return []Action{WaitAction{
		rootOrganId: root.organId,
		message:     "",
	}}
}

func findGrowActions(s State, root, organ Entity, enemyTentaclesTargets [][]bool) []Action {
	actions := make([]Action, 0)

	// find all the possible grow actions for the organ
	for _, offset := range offsets {
		coord := organ.coord.add(offset)
		if coord.isValid() && s.isWalkable(coord, false) && !enemyTentaclesTargets[coord.y][coord.x] {
			for _, _type := range []EntityType{BASIC, HARVESTER, TENTACLE, SPORER} {

				if _type == BASIC && canGrow(s.MyProteins, BASIC) {
					// for basic, direction doesn't matter
					actions = append(actions, GrowAction{
						rootOrganId: root.organId,
						organId:     organ.organId,
						coord:       coord,
						_type:       _type,
						dir:         N,
						message:     "",
					})
				}

				if (_type == HARVESTER && canGrow(s.MyProteins, HARVESTER)) ||
					(_type == TENTACLE && canGrow(s.MyProteins, TENTACLE)) ||
					(_type == SPORER && canGrow(s.MyProteins, SPORER)) {
					for _, dir := range []Dir{N, S, W, E} {
						actions = append(actions, GrowAction{
							rootOrganId: root.organId,
							organId:     organ.organId,
							coord:       coord,
							_type:       _type,
							dir:         dir,
							message:     "",
						})
					}
				}
			}
		}
	}

	return actions
}

func permute[T any](items []T) [][]T {
	if len(items) == 0 {
		return [][]T{{}}
	}

	if len(items) == 1 {
		return [][]T{items}
	}

	perms := make([][]T, 0)

	for i, item := range items {
		remaining := make([]T, 0)
		remaining = append(remaining, items[:i]...)
		remaining = append(remaining, items[i+1:]...)

		for _, perm := range permute[T](remaining) {
			perms = append(perms, append([]T{item}, perm...))
		}
	}

	return perms
}

func findClosestOrganTo(to []Coord, from Coord, tentacleTargets [][]bool) Coord {
	// use BFS to find the closest organ from the root to the target

	//debug("Finding closest organ to target: %+v\n", to)
	//debug("From: %+v\n", from)

	visited := make([][]bool, globalState.Height)
	for i := 0; i < globalState.Height; i++ {
		visited[i] = make([]bool, globalState.Width)
	}

	queue := make([]Coord, 0)
	queue = append(queue, from)
	visited[from.y][from.x] = true

	toMap := make([][]bool, globalState.Height)
	for i := 0; i < globalState.Height; i++ {
		toMap[i] = make([]bool, globalState.Width)
	}

	for _, coord := range to {
		toMap[coord.y][coord.x] = true
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if toMap[current.y][current.x] {
			// found the target
			return current
		}

		for _, offset := range offsets {
			neighbor := current.add(offset)
			if neighbor.isValid() &&
				!visited[neighbor.y][neighbor.x] &&
				!tentacleTargets[neighbor.y][neighbor.x] &&
				(globalState.Grid[neighbor.y][neighbor.x] == nil ||
					globalState.Grid[neighbor.y][neighbor.x]._type.isProtein() ||
					globalState.Grid[neighbor.y][neighbor.x].owner != NONE) {
				visited[neighbor.y][neighbor.x] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return Coord{-1, -1}
}

func findClosestEnemyOrgan(root Entity) Coord {
	debug("Finding closest enemy organ from root: %+v\n", root)

	//for _, entity := range globalState.Entities {
	//	if entity.owner == OPPONENT {
	//		debug("Possible enemy organ: %+v\n", entity)
	//	}
	//}

	// use BFS to find the closest enemy organ from the root
	visited := make([][]bool, globalState.Height)
	for i := 0; i < globalState.Height; i++ {
		visited[i] = make([]bool, globalState.Width)
	}

	queue := make([]Coord, 0)
	queue = append(queue, root.coord)
	visited[root.coord.y][root.coord.x] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if globalState.Grid[current.y][current.x] != nil {
			entity := globalState.Grid[current.y][current.x]
			if entity.owner == OPPONENT {
				debug("Found enemy organ at %+v\n", current)
				return current
			}
		}

		for _, offset := range offsets {
			neighbor := current.add(offset)
			if neighbor.isValid() && !visited[neighbor.y][neighbor.x] && (globalState.Grid[neighbor.y][neighbor.x] == nil ||
				globalState.Grid[neighbor.y][neighbor.x]._type.isProtein() ||
				globalState.Grid[neighbor.y][neighbor.x].owner != NONE) {
				visited[neighbor.y][neighbor.x] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return Coord{-1, -1}
}

func findShortestPathProt(s State, organs []Entity, nonHarvestedProteins []Entity, enemyTentaclesTargets [][]bool) []Coord {
	from := make([]Coord, 0)
	for _, organ := range organs {
		from = append(from, organ.coord)
	}

	to := make([]Coord, 0)
	for _, protein := range nonHarvestedProteins {
		to = append(to, protein.coord)
	}

	blockedCoords := make([][]bool, s.Height)
	for i := 0; i < s.Height; i++ {
		blockedCoords[i] = make([]bool, s.Width)
	}

	// block the cells that are targeted by the enemy tentacles
	for i := 0; i < s.Height; i++ {
		for j := 0; j < s.Width; j++ {
			if enemyTentaclesTargets[i][j] {
				blockedCoords[i][j] = true
			}
		}
	}

	// block the cells that are not walkable
	for i := 0; i < s.Height; i++ {
		for j := 0; j < s.Width; j++ {
			if !s.isWalkable(Coord{j, i}, false) {
				blockedCoords[i][j] = true
			}
		}
	}

	return findShortestPath(s, from, to, blockedCoords)
}

func findEnemyTentaclesTargets(s State) [][]bool {
	// find all the cells that are targeted by the enemy tentacles (cannot grow there)
	tentacleTargets := make([][]bool, s.Height)
	for i := 0; i < s.Height; i++ {
		tentacleTargets[i] = make([]bool, s.Width)
	}

	for _, entity := range s.Entities {
		if entity._type == TENTACLE && entity.owner == OPPONENT {
			coord := entity.coord.add(offsets[entity.organDir])
			if coord.isValid() {
				tentacleTargets[coord.y][coord.x] = true
			}
		}
	}

	return tentacleTargets
}

type TentacleGrowPlan struct {
	organ          Entity
	growCoord      Coord
	dir            Dir
	attacked       Entity
	destroyedCount int
}

func findTentacleAttacks(s State, organs []Entity, enemyTentaclesTargets [][]bool) []TentacleGrowPlan {
	// find all the tentacles that I can grow to instantly kill an opponent organ
	attacks := make([]TentacleGrowPlan, 0)

	if !canGrow(s.MyProteins, TENTACLE) {
		return attacks
	}

	for _, organ := range organs {
		for _, offset := range offsets {
			coord := organ.coord.add(offset)
			if coord.isValid() && s.isWalkable(coord, false) && !enemyTentaclesTargets[coord.y][coord.x] {
				// check if there is an opponent organ in the direction of the offset
				for _, dir := range []Dir{N, S, W, E} {
					attackedCoord := coord.add(offsets[dir])
					if attackedCoord.isValid() && s.Grid[attackedCoord.y][attackedCoord.x] != nil {
						attacked := s.Grid[attackedCoord.y][attackedCoord.x]
						if attacked.owner == OPPONENT {

							destroyed := findDestroyed(s, *attacked)

							attacks = append(attacks, TentacleGrowPlan{
								organ:          organ,
								growCoord:      coord,
								dir:            dir,
								attacked:       *attacked,
								destroyedCount: len(destroyed),
							})
						}
					}
				}
			}
		}
	}

	return attacks
}

func findDestroyed(s State, attacked Entity) []Entity {
	// find the enemy organs that will be destroyed if the attacked organ is destroyed
	// use BFS to find all the children of the attacked organ

	visited := make([][]bool, s.Height)
	for i := 0; i < s.Height; i++ {
		visited[i] = make([]bool, s.Width)
	}

	queue := make([]Entity, 0)
	queue = append(queue, attacked)
	visited[attacked.coord.y][attacked.coord.x] = true
	destroyed := make([]Entity, 0)
	destroyed = append(destroyed, attacked)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, offset := range offsets {
			neighbor := current.coord.add(offset)
			if neighbor.isValid() && !visited[neighbor.y][neighbor.x] && s.Grid[neighbor.y][neighbor.x] != nil {
				entity := s.Grid[neighbor.y][neighbor.x]
				if entity.organParentId == current.organId {
					visited[neighbor.y][neighbor.x] = true
					queue = append(queue, *entity)
					destroyed = append(destroyed, *entity)
				}
			}
		}
	}

	return destroyed
}

func growToFrontier(s State, organs []Entity, enemyTentaclesTargets [][]bool) {
	// there is no protein on the grid, find a cell that is at the frontier of players' organisms

	var enemyOrgans []Entity
	for _, entity := range s.Entities {
		if entity.owner == OPPONENT {
			enemyOrgans = append(enemyOrgans, entity)
		}
	}

	// debug("Enemy organs: %+v\n", enemyOrgans)

	var bestCell Coord = Coord{-1, -1}
	var bestOfMyOrgans Entity
	var bestOfEnemyOrgans Entity
	bestScore := -1000

	for i := 0; i < s.Height; i++ {
		for j := 0; j < s.Width; j++ {
			if s.Grid[i][j] == nil && !enemyTentaclesTargets[i][j] {
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

	if bestCell == (Coord{-1, -1}) {
		fmt.Println("WAIT no cell")
	} else {
		growType := findGrowType(s)

		debug("Grow target cell: %+v from organ: %+v and enemy organ: %+v\n", bestCell, bestOfMyOrgans, bestOfEnemyOrgans)

		growDir := findApproximateDir(bestOfMyOrgans.coord, bestCell)

		if growType == -1 {
			fmt.Println("WAIT cannot grow frontier")
		} else {
			fmt.Printf("GROW %d %d %d %s %s frontier\n", bestOfMyOrgans.organId, bestCell.x, bestCell.y, showOrganType(growType), showDir(growDir))
		}
	}
}

func buildSporeCellsMap(s State, nonHarvestedProteins []Entity) [][]bool {
	sporeCells := make([][]bool, s.Height)
	for i := 0; i < s.Height; i++ {
		sporeCells[i] = make([]bool, s.Width)
	}

	for _, protein := range nonHarvestedProteins {
		for _, offset := range offsets {
			coord := protein.coord.add(offset)
			if coord.isValid() && s.Grid[coord.y][coord.x] == nil {
				sporeCells[coord.y][coord.x] = true
			}
		}
	}

	// debug("Spore cells:\n")
	// for i := 0; i < globalState.Height; i++ {
	// 	for j := 0; j < globalState.Width; j++ {
	// 		if sporeCells[i][j] {
	// 			fmt.Fprintf(os.Stderr, "X ")
	// 		} else {
	// 			fmt.Fprintf(os.Stderr, ". ")
	// 		}
	// 	}
	// 	fmt.Fprintf(os.Stderr, "\n")
	// }

	return sporeCells
}

func findHarvestedProteins(s State) ([]Entity, []Entity) {
	var nonHarvestedProteins []Entity
	var harvestedProteins []Entity

	for _, entity := range s.Entities {
		if entity._type.isProtein() {
			// find my neighbor harvesters of this protein (must be facing the protein)
			myHarvesters := make([]Entity, 0)
			for _, offset := range offsets {
				coord := entity.coord.add(offset)
				if coord.isValid() {
					if s.Grid[coord.y][coord.x] != nil {
						neighbor := s.Grid[coord.y][coord.x]
						if neighbor._type == HARVESTER && neighbor.owner == ME {
							if findDirRelativeTo(neighbor.coord, entity.coord) == neighbor.organDir {
								myHarvesters = append(myHarvesters, *neighbor)
							} else {
								//debug("Neighbor harvester %+v is not facing the protein %+v\n", neighbor, entity)
							}
						}
					}
				}
			}

			if len(myHarvesters) > 0 {
				harvestedProteins = append(harvestedProteins, entity)
			} else {
				nonHarvestedProteins = append(nonHarvestedProteins, entity)
			}
		}
	}

	//nonHarvestedProteinsCoords := make([]Coord, 0)
	//harvestedProteinsCoords := make([]Coord, 0)
	//
	//for _, protein := range nonHarvestedProteins {
	//	nonHarvestedProteinsCoords = append(nonHarvestedProteinsCoords, protein.coord)
	//}
	//
	//for _, protein := range harvestedProteins {
	//	harvestedProteinsCoords = append(harvestedProteinsCoords, protein.coord)
	//}

	//debug("%d Non-harvested proteins, %d Harvested proteins\n", len(nonHarvestedProteins), len(harvestedProteins))

	return harvestedProteins, nonHarvestedProteins
}

func findGrowType(s State) EntityType {
	// grow a tentacle if I have enough proteins, better for attack and defense
	if s.MyProteins[1] >= 10 && s.MyProteins[2] >= 10 {
		return TENTACLE
	}

	growType := EntityType(-1)

	for _, _type := range []EntityType{BASIC, HARVESTER, TENTACLE, SPORER} {
		if canGrow(s.MyProteins, _type) {
			return _type
		}
	}

	return growType
}

/*
N to M pathfinding. Where N are my organs and M are the non=harvested proteins.
Finds the shortest path from any of my organs to any of the non-harvested proteins.
It must avoid the enemy tentacles.
Cannot go through existing organs.
*/
func findShortestPath(s State, from, to []Coord, forbiddenCells [][]bool) []Coord {
	// the chosen algorithm is BFS

	toMap := make([][]bool, s.Height)
	for i := 0; i < s.Height; i++ {
		toMap[i] = make([]bool, s.Width)
	}

	for _, coord := range to {
		toMap[coord.y][coord.x] = true
	}

	previous := make([][]Coord, s.Height)
	for i := 0; i < s.Height; i++ {
		previous[i] = make([]Coord, s.Width)
		for j := 0; j < s.Width; j++ {
			previous[i][j] = Coord{-1, -1}
		}
	}

	visited := make([][]bool, s.Height)
	for i := 0; i < s.Height; i++ {
		visited[i] = make([]bool, s.Width)
	}

	queue := make([]Coord, 0)

	for _, coord := range from {
		queue = append(queue, coord)
		visited[coord.y][coord.x] = true
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if toMap[current.y][current.x] {
			// found the target
			path := make([]Coord, 0)

			for current != (Coord{-1, -1}) {
				path = append(path, current)
				current = previous[current.y][current.x]
			}

			// reverse the path
			for i := 0; i < len(path)/2; i++ {
				path[i], path[len(path)-1-i] = path[len(path)-1-i], path[i]
			}

			return path
		}

		for _, offset := range offsets {
			neighbor := current.add(offset)
			if neighbor.isValid() && !visited[neighbor.y][neighbor.x] && !forbiddenCells[neighbor.y][neighbor.x] {
				visited[neighbor.y][neighbor.x] = true
				previous[neighbor.y][neighbor.x] = current
				queue = append(queue, neighbor)
			}
		}
	}

	debug("No path found\n")

	return nil
}

func growTowardsProtein(s State, nonHarvestedProteins []Entity, organs []Entity, enemyTentaclesTargets [][]bool, shortestPath []Coord) {
	// find the closest protein and organ

	if len(shortestPath) > 1 {
		// use the shortest path to grow towards the closest protein
		fromCell := shortestPath[0]
		fromEntity := s.Grid[fromCell.y][fromCell.x]

		stepCell := shortestPath[1]

		toCell := shortestPath[len(shortestPath)-1]
		// toEntity := globalState.Entities[s.Grid[toCell.y][toCell.x]]

		debug("Path from cell: %+v, to cell: %+v via step cell: %+v\n", fromCell, toCell, stepCell)

		if distance(stepCell, toCell) == 1 && canGrow(s.MyProteins, HARVESTER) {
			harvesterDir := findDirRelativeTo(stepCell, toCell)
			fmt.Printf("GROW %d %d %d HARVESTER %s path_harv_prot\n", fromEntity.organId, stepCell.x, stepCell.y, showDir(harvesterDir))
		} else {
			growType := findGrowType(s)

			growDir := N

			if len(shortestPath) >= 3 {
				growDir = findDirRelativeTo(stepCell, shortestPath[2])
			} else {
				growDir = findDirRelativeTo(fromCell, stepCell)
			}

			if growType == -1 {
				fmt.Println("WAIT cannot grow path")
			} else {
				fmt.Printf("GROW %d %d %d %s %s path_closer_prot\n", fromEntity.organId, stepCell.x, stepCell.y, showOrganType(growType), showDir(growDir))
			}
		}
	} else {
		closestProtein, closestOrgan := findClosestProteinAndOrgan(nonHarvestedProteins, organs)
		debug("Closest protein: %+v\n from organ: %+v\n", closestProtein, closestOrgan)

		// find the closest neighbor of the closest protein that can be reached by the closest organ
		closestNeighbor, closestOrgan := findClosestNeighborToProtein(s, closestProtein, organs, enemyTentaclesTargets)

		if closestNeighbor == (Coord{-1, -1}) {
			debug("No neighbor found for protein: %+v\n", closestProtein)
			fmt.Println("WAIT no neighbor")
		} else {
			debug("Closest neighbor: %+v\n", closestNeighbor)

			if distance(closestNeighbor, closestProtein.coord) == 1 && canGrow(s.MyProteins, HARVESTER) {
				harvesterDir := findDirRelativeTo(closestNeighbor, closestProtein.coord)
				fmt.Printf("GROW %d %d %d HARVESTER %s harv_prot\n", closestOrgan.organId, closestNeighbor.x, closestNeighbor.y, showDir(harvesterDir))
			} else {
				growType := findGrowType(s)

				growDir := findApproximateDir(closestNeighbor, closestProtein.coord)

				if growType == -1 {
					fmt.Println("WAIT cannot grow")
				} else {
					fmt.Printf("GROW %d %d %d %s %s closer_prot\n", closestOrgan.organId, closestNeighbor.x, closestNeighbor.y, showOrganType(growType), showDir(growDir))
				}
			}
		}
	}
}

func findClosestProteinAndOrgan(nonHarvestedProteins []Entity, organs []Entity) (Entity, Entity) {
	var closestProtein Entity
	var closestOrgan Entity
	minDistance := 1000

	for _, protein := range nonHarvestedProteins {
		for _, organ := range organs {
			dist := distance(protein.coord, organ.coord)
			if dist < minDistance {
				minDistance = dist
				closestProtein = protein
				closestOrgan = organ
			}
		}
	}

	return closestProtein, closestOrgan
}

func findClosestNeighborToProtein(s State, protein Entity, organs []Entity, enemyTentaclesTargets [][]bool) (Coord, Entity) {
	var closestNeighbor Coord = Coord{-1, -1}
	var closestOrgan Entity
	minDistance := 1000

	for _, organ := range organs {
		for _, offset := range offsets {
			neighbor := organ.coord.add(offset)
			if neighbor.isValid() &&
				s.isWalkable(neighbor, false) &&
				!enemyTentaclesTargets[neighbor.y][neighbor.x] {
				dist := distance(neighbor, protein.coord)
				if dist < minDistance {
					minDistance = dist
					closestNeighbor = neighbor
					closestOrgan = organ
				}
			}
		}
	}

	return closestNeighbor, closestOrgan
}

func growSporerIfPossible(s State, sporeCells [][]bool, organs []Entity) bool {
	if canGrow(s.MyProteins, SPORER) {
		// check if the closest organ can reach the closest protein using sporers

		debug("Can grow a sporer\n")

		// find any neighbor of my organs that can reach a spore cell in any direction

		sporerPlans := make([]SporePlan, 0)

		for _, organ := range organs {
			for _, offset := range offsets {
				sporerCoord := organ.coord.add(offset)
				if sporerCoord.isValid() && s.Grid[sporerCoord.y][sporerCoord.x] == nil {
					// simulate the spore in all directions until it reaches a spore cell
					for _, dir := range []Dir{N, S, W, E} {
						sporeCoord := findSporeCellInDirection(s, sporerCoord, dir, sporeCells)

						if sporeCoord.isValid() &&
							distance(sporerCoord, sporeCoord) > 5 {
							debug("Organ: %+v can reach spore cell: %+v after sporing in direction: %s from cell: %+v\n", organ, sporeCoord, showDir(dir), sporerCoord)
							sporerPlans = append(sporerPlans, SporePlan{
								organ:          organ,
								newSporerCoord: sporerCoord,
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
			return true
		} else {
			debug("No spore plans\n")
		}
	}

	return false
}

func sporeIfPossible(s State, sporeCells [][]bool) bool {
	if canSpore(s.MyProteins) {
		// check if I have a sporer that can spore a new root into a spore cell
		sporer := Entity{}
		sporeCoord := Coord{-1, -1}
		for _, entity := range s.Entities {
			if entity._type == SPORER && entity.owner == ME {
				sporeCooord := findSporeCellInDirection(s, entity.coord, entity.organDir, sporeCells)
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
			return true
		}
	}

	return false
}

func canSpore(proteinCounts []int) bool {
	return proteinCounts[0] >= 1 && proteinCounts[1] >= 1 && proteinCounts[2] >= 1 && proteinCounts[3] >= 1
}

func findSporeCellInDirection(s State, coord Coord, dir Dir, sporeCells [][]bool) Coord {
	sporeCoord := coord
	for {
		sporeCoord = sporeCoord.add(offsets[dir])
		if !sporeCoord.isValid() {
			break
		}

		if s.Grid[sporeCoord.y][sporeCoord.x] != nil &&
			!(s.Grid[sporeCoord.y][sporeCoord.x]._type.isProtein()) {
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
	fmt.Scan(&globalState.Width, &globalState.Height)

	for {
		parseTurnState()
		sendActions()
	}
}

type OrganGrowCost struct {
	costA, costB, costC, costD int
}

func growCost(_type EntityType) OrganGrowCost {
	switch _type {
	case BASIC:
		return OrganGrowCost{1, 0, 0, 0}
	case HARVESTER:
		return OrganGrowCost{0, 0, 1, 1}
	case TENTACLE:
		return OrganGrowCost{0, 1, 1, 0}
	case SPORER:
		return OrganGrowCost{0, 1, 0, 1}
	default:
		panic(fmt.Sprintf("Unknown type %d", _type))
	}
}

func canGrow(proteinCounts []int, _type EntityType) bool {
	cost := growCost(_type)
	return proteinCounts[0] >= cost.costA && proteinCounts[1] >= cost.costB && proteinCounts[2] >= cost.costC && proteinCounts[3] >= cost.costD
}
