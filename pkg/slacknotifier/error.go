package slacknotifier

type NoChannelFoundErr struct {
	Err error
}

type UserLookupByEmailErr struct {
	Err error
}

type JoinedChannelFetchErr struct {
	Err error
}

type MsgSendErr struct {
	Err error
}

type SlackNotifierErr struct {
	Err error
}

func (n *NoChannelFoundErr) Error() string {
	return n.Err.Error()
}

func (n *UserLookupByEmailErr) Error() string {
	return n.Err.Error()
}

func (n *JoinedChannelFetchErr) Error() string {
	return n.Err.Error()
}

func (n *MsgSendErr) Error() string {
	return n.Err.Error()
}

func (n *SlackNotifierErr) Error() string {
	return n.Err.Error()
}
