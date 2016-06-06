package f9crew

import "sort"

// Manifest is a crew manifest for a mission. This consists of a slice of
// f9crew.Interface. This type exists so that it can implement the
// sort.Interface interface.
type Manifest []Interface

// Sort is a convenience function around sort.Sort()
func (m Manifest) Sort() { sort.Sort(m) }

// Len is the number of elements in the collection.
func (m Manifest) Len() int { return len(m) }

// Swap swaps the elements with indexes i and j
func (m Manifest) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

// Less reports whether the element with index i should sort before the
// element with index j.
func (m Manifest) Less(i, j int) bool {
	name1, name2 := m[i].Name(), m[j].Name()

	// if the names are the same, decide using HashedKey
	if name1 == name2 {
		return m[i].HashedKey() < m[j].HashedKey()
	}

	return name1 < name2
}
