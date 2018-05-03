package sset

import (
	"testing"
)

func TestSortedSetWrite(t *testing.T) {
	t.Run("SortedSetAdd", func(t *testing.T) {
		s := SortedSet{}
		s.Init()

		result := s.Get(5)
		if len(result) != 0 {
			t.Errorf("Result length should not be greater than zero for Get call")
			return
		}

		result = s.GetRange(5, 10)
		if len(result) != 0 {
			t.Errorf("Result length should not be greater than zero for GetRange call")
			return
		}

		s.Add("Hello", 5)
		s.Add("World", 5)
		s.Add("World2", 6)

		result = s.Get(5)
		if result[0] != "Hello" && result[1] != "Hello" {
			t.Errorf("Hello not present in result.")
			return
		}

		if result[0] != "World" && result[1] != "World" {
			t.Errorf("World not present in result")
			return
		}

		if s.Add("World2", 7) {
			t.Errorf("World2 should not have been added")
			return
		}

		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Negative Rank addition did not panic")
					return
				}
			}()
			s.Add("NegativeRank", -1)
		}()
	})

	t.Run("SortedSetRemove", func(t *testing.T) {
		s := SortedSet{}
		s.Init()

		if s.Remove("Hello") {
			t.Errorf("Removal of non existant key should return false")
			return
		}

		s.Add("Hello", 5)
		s.Add("World", 5)

		if !s.Remove("Hello") {
			t.Errorf("Removal of existing key should return true")
			return
		}

		result := s.Get(5)
		if len(result) != 1 {
			t.Errorf("Length of result should be: 1, got: %d", len(result))
			return
		}

		if result[0] != "World" {
			t.Errorf("World not pesent in the result")
			return
		}
	})

}

func TestSortedSetRead(t *testing.T) {
	t.Run("SortedSetGet", func(t *testing.T) {
		s := SortedSet{}
		s.Init()

		s.Add("Member1", 5)
		s.Add("Member2", 5)
		s.Add("Member3", 6)

		result := s.Get(5)

		if len(result) != 2 {
			t.Errorf("Expected length of result to be 2, actual is: %d", len(result))
			return
		}

		if result[0] != "Member1" && result[1] != "Member1" {
			t.Errorf("Member1 is not present in the result")
			return
		}

		if result[0] != "Member2" && result[1] != "Member2" {
			t.Errorf("Member2 is not present in the result")
			return
		}

		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Negative Rank searching did not panic")
					return
				}
			}()
			s.Get(-1)
		}()

	})

	t.Run("SortedSetGetRange", func(t *testing.T) {
		s := SortedSet{}
		s.Init()

		result := s.GetRange(5, 10)
		if len(result) != 0 {
			t.Errorf("Result length should not be greater than zero for GetRange call")
			return
		}

		s.Add("Hello", 5)
		s.Add("World", 5)
		s.Add("World2", 6)

		result = s.GetRange(4, 6)
		if len(result) != 2 {
			t.Errorf("Result length for GetRange expected: %d, got: %d", 3, len(result))
			return
		}

		if result[0] != "Hello" && result[1] != "Hello" {
			t.Errorf("Hello not present in result.")
			return
		}

		if result[0] != "World" && result[1] != "World" {
			t.Errorf("World not present in result")
			return
		}

		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Negative Rank searching did not panic")
					return
				}
			}()
			s.GetRange(-1, 5)
		}()

		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Invalid Range Rank searching did not panic")
					return
				}
			}()
			s.GetRange(5, 5)
		}()

		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Invalid Range Rank searching did not panic")
					return
				}
			}()
			s.GetRange(5, 1)
		}()
	})

	t.Run("SortedSetExists", func(t *testing.T) {
		s := SortedSet{}
		s.Init()

		if s.Exists("Hello") {
			t.Errorf("Returned True for non existant member")
			return
		}

		s.Add("World", 5)
		if !s.Exists("World") {
			t.Errorf("Returned False for existing key")
			return
		}
	})

	t.Run("SortedSetGetRank", func(t *testing.T) {
		s := SortedSet{}
		s.Init()

		if s.GetRank("Hello") != -1 {
			t.Errorf("Returned rank for non existant member")
			return
		}

		s.Add("World", 5)
		if s.GetRank("World") != 5 {
			t.Errorf("Returned Wrong rank for member")
			return
		}
	})
}
