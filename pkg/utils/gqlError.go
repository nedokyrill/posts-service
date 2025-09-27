package utils

type GqlError struct {
	Msg  string
	Type string
}

func (g GqlError) Error() string {
	return g.Msg
}

func (g GqlError) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"msg":  g.Msg,
		"type": g.Type,
	}
}
