package command

type Cmd struct {
	sendFile, getFile []byte
}

func (cmd *Cmd) Init() {
	cmd.sendFile = []byte("a")
	cmd.getFile = []byte("a")
}
