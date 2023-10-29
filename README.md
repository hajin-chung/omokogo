# Omokogo

online gomoku 

## notes

every realtime communications happen through websocket on `/ws`

#### input messages

- ENQ: enqueue
- DEQ: dequeue, cancel queue
- PLC y x: place stone on y, x

#### output messages

- MTC gameId: match made on gameId
- GST gameId userId1 userId2 (y1, x1), (y2, x2), ... : game state
- ERR msg: error messages

## Design

1. auth
    - create account
    - login, logout
2. hub
    - handle connection
    - handle disconnection
3. core
    - queue
        - user enqueue: push user(Id, Score) to queue
        - match made: match maker
        - create game
    - handle commands: commandHandler
        1. places stone
            - check if player's turn to place
            - check if stone location is allowed
            - broadcast game state if board changed
            - handle win, lose, draw

## data structure

```
user (
    id
    status (playing, queue, idle)
    gameId
)

game (
    id
    userId1
    userId2
    stones [](userNum, x, y)
    status (playing, done)
)
```
