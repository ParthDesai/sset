package internals

// Dictionary provides efficient lookup of element against its rank
// It is useful for operation that does not require accessing skiplist
type Dictionary map[string]int
