package exp

// Keys returns the keys of the map m.
// The keys will be in an indeterminate order.
//
// Snippet from https://cs.opensource.google/go/x/exp/+/39d4317d:maps/maps.go;l=10
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
