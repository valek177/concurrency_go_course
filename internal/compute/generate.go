package compute

//go:generate mockgen -source=./compute.go -destination=./mock/handler.go -package=mock Handler
//go:generate mockgen -source=./parser.go -destination=./mock/parser.go -package=mock Parser
