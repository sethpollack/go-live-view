package channel

type Channel interface {
	Join(Socket, any) error
	Leave(Socket) error
	Message(Socket, string, any) error
	Broadcast(Socket, string, any) error
}
