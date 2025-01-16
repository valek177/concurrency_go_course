package compute

// Query is a struct for query
type Query struct {
	Command string
	Args    []string
}

// NewQuery returns new query object
func NewQuery(cmdName string, arguments []string) Query {
	return Query{
		Command: cmdName,
		Args:    arguments,
	}
}

// CommandName returns command name
func (q *Query) CommandName() string {
	return q.Command
}

// Arguments returns args of command
func (q *Query) Arguments() []string {
	return q.Args
}
