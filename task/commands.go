package task

type CmdTaskCreate struct {
	ID          string
	Description string
}

type CmdTaskUpdateDescription struct {
	ID             string
	NewDescription string
}

type CmdTaskComplete struct {
	ID string
}
