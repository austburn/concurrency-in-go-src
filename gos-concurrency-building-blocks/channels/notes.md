Send and receive channel example: https://play.golang.org/p/bc-yPpQ27Nk
On this, the channel is unbuffered, which means that a send and receive goroutine need to be ready for communication to occur. That means that a minimum of 2 goroutines is required (including the main goroutine). Send or receive, doesn't matter, but one side has to live in a goroutine.

- https://play.golang.org/p/nlzqAeVInSB
A channel can be written to, closed, and then return (1, true) in this case, as 1 was indeed written to the channel.
