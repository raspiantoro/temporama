package resp

import "fmt"

type CommandHandler interface {
	Serve(cmd Command) ValueNode
}

type HandlerFunc func(cmd Command) ValueNode

func (h HandlerFunc) Serve(cmd Command) ValueNode {
	return h(cmd)
}

type Mux struct {
	handlers map[string]CommandHandler
}

func NewCommandMux() *Mux {
	return &Mux{
		handlers: make(map[string]CommandHandler),
	}
}

func (m *Mux) Handle(name string, handler CommandHandler) {
	m.handlers[name] = handler
}

func (m *Mux) HandleFunc(name string, handler func(cmd Command) ValueNode) {
	m.handlers[name] = HandlerFunc(handler)
}

func (m *Mux) Serve(cmd Command) ValueNode {
	handler, ok := m.handlers[cmd.name]
	if !ok {
		return ValueNode{
			types: ValueNodeTypeSimpleError,
			val:   fmt.Sprintf("ERROR: \"%s\" command is not supported", cmd.name),
		}
	}

	return handler.Serve(cmd)
}

type Command struct {
	name string
	args []string
}

func NewCommand(name string, args ...string) Command {
	return Command{
		name: name,
		args: args,
	}
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) Key() string {
	var val string

	if len(c.args) == 0 {
		val = ""
	} else {
		val = c.args[0]
	}

	return val
}

func (c *Command) Args() []string {
	var val []string

	if len(c.args) <= 1 {
		val = nil
	} else {
		val = c.args[1:]
	}

	return val
}
