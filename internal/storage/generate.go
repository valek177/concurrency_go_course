package storage

//go:generate mockgen -source=./storage.go -destination=./mock/storage.go -package=mock Storage
//go:generate mockgen -source=./engine.go -destination=./mock/engine.go -package=mock Engine
