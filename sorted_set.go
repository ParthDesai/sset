package sset

import (
	"sync"

	"github.com/parthdesai/sset/internals"
)

const levelJumpProbability = 0.5 // 1/2 probability of level jump
const maxLevels = 32             // log n distribution, so 2^32
const minKey = 0

// SortedSet struct represent sorted set abstract data structure
// Under the hood, it uses skiplist and dict for sorted set functionality and
// Read write mutex for thread safe operation
type SortedSet struct {
	dict     internals.Dictionary
	skiplist internals.SkipList
	rwMutex  *sync.RWMutex
}

// Init Initiates sorted set.
func (s *SortedSet) Init() {
	s.dict = internals.Dictionary{}
	s.skiplist = internals.SkipList{}
	s.rwMutex = &sync.RWMutex{}
	s.skiplist.Init(maxLevels, levelJumpProbability, minKey)
}

// Add adds a string element to sorted set, with rank indicated by
// rank parameter
func (s *SortedSet) Add(member string, rank int) bool {

	if rank < 0 {
		panic("Rank must be greater than or equal to zero")
	}

	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if _, ok := s.dict[member]; ok {
		return false
	}

	s.dict[member] = rank
	s.skiplist.AddOrModify(rank, map[string]bool{member: true}, func(currentVal map[string]bool) map[string]bool {
		currentVal[member] = true
		return currentVal
	})

	return true
}

// Remove Removes member from sorted set
func (s *SortedSet) Remove(member string) bool {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	val, ok := s.dict[member]
	if !ok {
		return false
	}

	deleted := false

	delete(s.dict, member)
	s.skiplist.DeleteOrModify(val, func(memberMap map[string]bool) (bool, map[string]bool) {
		deleted = true
		delete(memberMap, member)
		if len(memberMap) == 0 {
			return true, nil
		}
		return false, memberMap
	})

	return deleted
}

// Get gets member(s) with given rank
func (s *SortedSet) Get(rank int) []string {
	if rank < 0 {
		panic("Rank must be greater than or equal to zero")
	}

	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	memberMap := s.skiplist.SearchOrModify(rank, nil)
	members := make([]string, len(memberMap))

	i := 0
	for value := range memberMap {
		members[i] = value
		i++
	}

	return members
}

// GetRange returns all members with rank in between rankMin and rankMax
// rankMin is inclusive
func (s *SortedSet) GetRange(rankMin, rankMax int) []string {
	if rankMin < 0 || rankMax < 0 {
		panic("rankMin and rankMax must be greater than equal to zero")
	}

	if rankMin >= rankMax {
		panic("rankMin must be less than rankMax")
	}

	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	searchResult := s.skiplist.SearchRange(rankMin, rankMax)

	numberOfMembers := 0
	for _, memberMap := range searchResult {
		numberOfMembers += len(memberMap)
	}

	members := make([]string, numberOfMembers)
	fillIndex := 0
	for _, memberMap := range searchResult {
		for value := range memberMap {
			members[fillIndex] = value
			fillIndex++
		}
	}

	return members
}

// Exists check for membership of member in sorted set
func (s *SortedSet) Exists(member string) bool {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	_, ok := s.dict[member]
	return ok
}

// GetRank gives rank of member
func (s *SortedSet) GetRank(member string) int {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if val, ok := s.dict[member]; ok {
		return val
	}
	return -1
}
