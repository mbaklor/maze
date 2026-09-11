package main

func merge[T string | int](defaultValue T, target *T, sources ...T) {
	var zeroVal T
	for _, src := range sources {
		if src != zeroVal {
			*target = src
			return
		}
	}
	*target = defaultValue
}
