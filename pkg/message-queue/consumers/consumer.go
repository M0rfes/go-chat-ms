package consumers

type Consumer[T any] func(*T, error)
