package task

import "fmt"

type CmdRejection interface {
	Reason() string
	Field() string
}

type CmdRejectionRequiredField struct {
	Name string
}

func (r *CmdRejectionRequiredField) Reason() string {
	return fmt.Sprintf("%s: is required", r.Name)
}

func (r *CmdRejectionRequiredField) Field() string {
	return r.Name
}
