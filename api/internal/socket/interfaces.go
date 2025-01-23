package socket

type GameManager interface {
	AddUser(client *Client)
	StartGame(client1, client2 *Client) error
	HandleMove(client *Client, move string) error
}
