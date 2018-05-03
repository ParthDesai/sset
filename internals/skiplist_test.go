package internals

import (
	"strconv"
	"testing"
)

func TestSkipListInitialization(t *testing.T) {
	t.Run("PanicForInvalidLevels", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Did not panic, even when level passed was less than 1")
				return
			}
		}()

		s := SkipList{}
		s.Init(0, 0.5, 0)
	})
}

func TestSkipListDebug(t *testing.T) {
	t.Run("DebugPrint", func(t *testing.T) {
		s := SkipList{}
		s.Init(1, 0.5, 0)

		s.AddOrModify(5, map[string]bool{"p": true}, nil)
		s.AddOrModify(6, map[string]bool{"p": true}, nil)
		s.AddOrModify(7, map[string]bool{"p": true}, nil)

		out := s.DebugPrint()
		if len(out) > 1 {
			t.Errorf("Debug Print printed multiple rows")
			return
		}

		if len(out) == 0 {
			t.Errorf("Debug Print did not print any rows")
			return
		}

		if out[0] != "0:5,6,7," {
			t.Errorf("Debug Print output is wrong")
			return
		}
	})
}

func TestSkipListRead(t *testing.T) {
	t.Run("Search", func(t *testing.T) {
		s := SkipList{}
		s.Init(5, 0.5, 0)
		s.AddOrModify(7, map[string]bool{"1": true}, nil)

		result := s.SearchOrModify(7, nil)
		if !result["1"] {
			t.Errorf("Search function is not working properly.")
			return
		}

		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Searching for key less than min key did not panic")
					return
				}
			}()
			s.SearchOrModify(-1, nil)
		}()
	})

	t.Run("ModificationInSearch", func(t *testing.T) {
		s := SkipList{}
		s.Init(5, 0.5, 0)
		s.AddOrModify(5, map[string]bool{"1": true}, nil)

		result := s.SearchOrModify(5, func(existingValue map[string]bool) map[string]bool {
			existingValue["Hello"] = true
			return existingValue
		})

		if !result["1"] {
			t.Errorf("Expected Added key %s to be there", "1")
			return
		}

		if !result["Hello"] {
			t.Errorf("Expected Added key %s to be there", "Hello")
			return
		}
	})

	t.Run("GetRangeSubset", func(t *testing.T) {
		s := SkipList{}
		s.Init(5, 0.5, 0)

		for i := 0; i < 100; i++ {
			s.AddOrModify(i, map[string]bool{strconv.Itoa(i): true}, nil)
		}

		s.DeleteOrModify(90, nil)

		result := s.SearchRange(90, 95)
		if len(result) != 4 {
			t.Errorf("Search range returned incorrect number of result. Expected: %d, Got: %d", 5, len(result))
		}

		if !result[0]["91"] {
			t.Errorf("Wrong element returned from range call")
			return
		}

		if !result[1]["92"] {
			t.Errorf("Wrong element returned from range call")
			return
		}

		if !result[2]["93"] {
			t.Errorf("Wrong element returned from range call")
			return
		}

		if !result[3]["94"] {
			t.Errorf("Wrong element returned from range call")
			return
		}

		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Searching for invalid range (keyMin less than min key) did not panic")
					return
				}
			}()
			s.SearchRange(-1, 5)
		}()

		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Searching for invalid range (keyMin = keyMax) did not panic")
					return
				}
			}()
			s.SearchRange(5, 5)
		}()

		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Searching for invalid range (keyMin > keyMax) did not panic")
					return
				}
			}()
			s.SearchRange(6, 5)
		}()
	})

	t.Run("GetRangeTail", func(t *testing.T) {
		s := SkipList{}
		s.Init(5, 0.5, 0)

		for i := 0; i < 100; i++ {
			s.AddOrModify(i, map[string]bool{strconv.Itoa(i): true}, nil)
		}

		result := s.SearchRange(98, 105)
		if len(result) != 2 {
			t.Errorf("Search range returned incorrect number of result. Expected: %d, Got: %d", 2, len(result))
			return
		}

		if !result[0]["98"] {
			t.Errorf("Wrong element returned from range call")
			return
		}

		if !result[1]["99"] {
			t.Errorf("Wrong element returned from range call")
			return
		}

	})
}

func TestSkipListWrite(t *testing.T) {
	t.Run("Initialization", func(t *testing.T) {

	})

	t.Run("Addition", func(t *testing.T) {
		s := SkipList{}
		s.Init(5, 0.5, 0)
		s.AddOrModify(5, map[string]bool{"1": true}, nil)

		result := s.SearchOrModify(5, nil)
		if !result["1"] {
			t.Errorf("Expected Added key to be there")
			return
		}

		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Addition of key less than minKey did not panic")
					return
				}
			}()
			s.AddOrModify(-1, map[string]bool{"-1": true}, nil)
		}()
	})

	t.Run("ModificationInAddition", func(t *testing.T) {
		s := SkipList{}
		s.Init(5, 0.5, 0)
		s.AddOrModify(5, map[string]bool{"1": true}, nil)
		s.AddOrModify(5, nil, func(existingVal map[string]bool) map[string]bool {
			existingVal["Hello"] = true
			return existingVal
		})

		result := s.SearchOrModify(5, nil)
		if !result["1"] {
			t.Errorf("Expected Added key %s to be there", "1")
			return
		}

		if !result["Hello"] {
			t.Errorf("Expected Added key %s to be there", "Hello")
			return
		}
	})

	t.Run("Deletion", func(t *testing.T) {
		s := SkipList{}
		s.Init(5, 0.5, 0)
		s.AddOrModify(6, map[string]bool{"1": true}, nil)

		result := s.SearchOrModify(6, nil)
		if !result["1"] {
			t.Errorf("Expected Added key to be there")
			return
		}

		if !s.DeleteOrModify(6, nil) {
			t.Errorf("Deletion function did not find the key")
			return
		}

		if len(s.SearchOrModify(6, nil)) != 0 {
			t.Errorf("Even after calling Delete, key is still there")
			return
		}

		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Deletion of key less than minKey did not panic")
					return
				}
			}()

			s.DeleteOrModify(-1, nil)
		}()

	})

	t.Run("ModificationInDeletion", func(t *testing.T) {
		s := SkipList{}
		s.Init(5, 0.5, 0)
		s.AddOrModify(6, map[string]bool{"1": true, "2": true}, nil)

		result := s.SearchOrModify(6, nil)
		if !result["1"] {
			t.Errorf("Expected Added key to be there")
			return
		}

		if s.DeleteOrModify(6, func(existingValue map[string]bool) (bool, map[string]bool) {
			delete(existingValue, "2")
			return false, existingValue
		}) {
			t.Errorf("Deletion function deleted key, where it should have only modified it.")
			return
		}

		if len(s.SearchOrModify(6, nil)) == 0 {
			t.Errorf("Even after calling Delete, key is still there")
			return
		}

		if !s.DeleteOrModify(6, func(existingValue map[string]bool) (bool, map[string]bool) {
			delete(existingValue, "1")
			return true, existingValue
		}) {
			t.Errorf("Deletion function should have deleted key.")
			return
		}

		if len(s.SearchOrModify(6, nil)) != 0 {
			t.Errorf("Search function stil found a key even after it should have been deleted")
			return
		}
	})
}
