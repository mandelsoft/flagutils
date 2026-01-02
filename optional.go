package flagutils

func AddOptionally[T any](set ExtendableOptionSet, opts ...T) {
	for _, a := range opts {
		if o, ok := any(a).(Options); ok {
			set.Add(o)
		}
	}
}
