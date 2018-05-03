package internals

import (
	"math/rand"
	"strconv"
	"strings"
)

// Node represents individual node in skiplist with unique key
// Each node have at max maxLevels of pointer to next node
// Values is map due to efficiency of checking whether value exists in particular
// node or not.
type Node struct {
	Next   []*Node
	Key    int
	Values map[string]bool
}

// SkipList is one of underlying data structure of sorted set
// It has very high probability O(log n) insert, lookup and delete
// It consists of n number of node, each one with maxLevels pointer
// Pointer at each level point to next node at that level
// Due to power-alike arrnagement of node (from higher level to lower level),
// Each level can act as express way for level below it, thus enabling O(log n)
// probability.
type SkipList struct {
	header               *Node
	maxLevels            int
	currentLevel         int
	levelJumpProbability float32
	minKey               int
}

// DebugPrint returns array of string, that gives you information
// about node position and levels
func (s *SkipList) DebugPrint() []string {
	debugData := make([]string, s.currentLevel+1)
	for i := 0; i <= s.currentLevel; i++ {
		strBuilder := strings.Builder{}
		strBuilder.WriteString(strconv.Itoa(i) + ":")
		current := s.header.Next[i]
		for current != nil {
			strBuilder.WriteString(strconv.Itoa(current.Key) + ",")
			current = current.Next[i]
		}
		debugData[i] = strBuilder.String()
	}
	return debugData
}

func (s *SkipList) createNewNode(key int, initialValues map[string]bool) *Node {
	n := &Node{}
	n.Next = make([]*Node, s.maxLevels)
	n.Key = key
	n.Values = initialValues
	return n
}

// Init Initiates skip list, with maximum number of level supported
// levelJumpProbability indicates the probability by which node in level i, can
// be also present at level i - 1
// minKey restricts key space to [minKey, MAX_INT)
func (s *SkipList) Init(maxLevels int, levelJumpProbability float32, minKey int) {

	if maxLevels < 1 {
		panic("Maximum number of levels should be greater than or equal to 1")
	}

	s.currentLevel = -1
	s.maxLevels = maxLevels
	s.levelJumpProbability = levelJumpProbability
	s.minKey = minKey
	s.header = s.createNewNode(s.minKey-1, nil)
}

func (s *SkipList) generateRandLevel() int {
	level := 0

	// This enables 2^(maxLevel - level) distribution
	for (float32((rand.Int() & 0xFFFF)) < (s.levelJumpProbability * 0xFFFF)) && level < (s.maxLevels-1) {
		level++
	}
	return level
}

// DeleteOrModify either deletes entire node OR deletes one of the value
// This is determined by function reference passed to it.
func (s *SkipList) DeleteOrModify(key int,
	modifyOrDelete func(map[string]bool) (bool, map[string]bool)) bool {

	if key < s.minKey {
		panic("Key must be greater than or equal to minKey")
	}

	current := s.header
	updateArray := make([]*Node, s.maxLevels)
	var nodeToDelete *Node

	// O(log n) time with very high probability
	for i := s.currentLevel; i >= 0; i-- {
		for current.Next[i] != nil && current.Next[i].Key < key {
			current = current.Next[i]
		}

		if current.Next[i] != nil && current.Next[i].Key == key {
			updateArray[i] = current
			nodeToDelete = current.Next[i]
		}
	}

	if nodeToDelete == nil {
		return false
	}

	var modifiedValue map[string]bool
	needToDelete := true

	if modifyOrDelete != nil {
		needToDelete, modifiedValue = modifyOrDelete(nodeToDelete.Values)
	}

	if !needToDelete {
		nodeToDelete.Values = modifiedValue
		return false
	}

	for i := 0; i <= s.currentLevel && updateArray[i] != nil; i++ {
		updateArray[i].Next[i] = nodeToDelete.Next[i]
	}

	newCurrentLevel := s.currentLevel
	for i := s.currentLevel; i >= 0; i-- {
		if s.header.Next[i] == nil {
			newCurrentLevel = i - 1
		}
	}
	s.currentLevel = newCurrentLevel

	return true
}

// SearchRange searches skiplist and finds node which has key less than
// keyMax, and greater than or equal to keyMin
func (s *SkipList) SearchRange(keyMin, keyMax int) []map[string]bool {

	if keyMin < s.minKey || keyMax < s.minKey {
		panic("keyMin and keyMax must be greater than or equal to minKey")
	}

	if keyMin >= keyMax {
		panic("keyMin must be less than keyMax")
	}

	var searchResult []map[string]bool
	current := s.header

	for i := s.currentLevel; i >= 0; i-- {
		for current.Next[i] != nil && current.Next[i].Key < keyMin {
			current = current.Next[i]
		}
	}

	/** Get next node at Level 0, this is the node, which can have one of three values:
	 ** 1) Node with key greater than current key
	 ** 2) Node with key equal to current key
	 ** 3) Nil **/
	current = current.Next[0]

	if current == nil {
		return searchResult
	}

	numberOfNodes := 0
	startNode := current
	for current != nil && current.Key < keyMax {
		numberOfNodes++
		current = current.Next[0]
	}

	searchResult = make([]map[string]bool, numberOfNodes)
	for i, node := 0, startNode; i < numberOfNodes; i, node = i+1, node.Next[0] {
		searchResult[i] = node.Values
	}

	return searchResult
}

// SearchOrModify searches skiplist for element
// You can also supply a function, that can be used to modifiy searched values
func (s *SkipList) SearchOrModify(key int,
	modifier func(map[string]bool) map[string]bool) map[string]bool {

	if key < s.minKey {
		panic("key must be greater than or equal to minKey")
	}

	var searchResult map[string]bool
	current := s.header

	// O(log n) time with very high probability
	for i := s.currentLevel; i >= 0; i-- {
		for current.Next[i] != nil && current.Next[i].Key < key {
			current = current.Next[i]
		}

		if current.Next[i] != nil && current.Next[i].Key == key {
			// Modifier is not nil, then we can invoke it and modify functionality
			if modifier != nil {
				current.Next[i].Values = modifier(current.Next[i].Values)
			}
			searchResult = current.Next[i].Values
			break
		}
	}

	return searchResult
}

// AddOrModify adds new node to skiplist if it does not exists,
// if it is you can supply modifier function that can be used to modify value
// of existing node
func (s *SkipList) AddOrModify(key int, value map[string]bool,
	modifier func(map[string]bool) map[string]bool) {

	if key < s.minKey {
		panic("key must be greater than or equal to minKey")
	}

	current := s.header
	updateArray := make([]*Node, s.maxLevels)
	for i := s.currentLevel; i >= 0; i-- {
		for current.Next[i] != nil && current.Next[i].Key < key {
			current = current.Next[i]
		}
		updateArray[i] = current
	}

	/** Get next node at Level 0, this is the node, which can have one of three values:
	 ** 1) Node with key greater than current key
	 ** 2) Node with key equal to current key
	 ** 3) Nil **/
	current = current.Next[0]

	/** Reached end of list, or current is bigger than key **/
	if current == nil || current.Key != key {
		levels := s.generateRandLevel()
		node := s.createNewNode(key, value)

		if s.currentLevel < levels {
			for i := s.currentLevel + 1; i <= levels; i++ {
				updateArray[i] = s.header
			}

			s.currentLevel = levels
		}

		for i := 0; i <= levels; i++ {
			node.Next[i] = updateArray[i].Next[i]
			updateArray[i].Next[i] = node
		}

	} else {
		if modifier != nil {
			current.Values = modifier(current.Values)
		} else {
			current.Values = value
		}
	}

}
