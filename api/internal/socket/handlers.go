package socket

import "fmt"

func SendMessage(event Event, c *Client) error {
	fmt.Println(event)
	return nil
}
