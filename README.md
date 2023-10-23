# Omokogo

online gomoku 

## user senarios

1. auth (db)
    - create account
    - login, logout
3. queue
    - user enqueue: push user(Id, Score) to queue
    - match made: match maker
    - create game: db
2. game
    - connects to game
        - handle connection
        - handle disconnection
    - handle messages: messageHandler
        1. places stone
            - handle auth (check if player's turn to place)
            - check if stone location is allowed
            - broadcast game state if board changed
            - handle win, lose, draw

## data structure

```
there are games 
games are played by 2 users 1 vs 1
game (
    id
    userId1
    userId2
    stones [](userNum, x, y)
    status (playing, done)
)

there are users
users can have mulit websockets
user (
    id
    name
    score
    status (playing, idle, queue) 
    gameId
    sockets []*websocket.Conn
)
```

## Dev order

- [ ] new user
- [ ] some rest api stuff
- [ ] queue
- [ ] new game
- [ ] think about websocket handling
- [ ] handle game logic

## some things to think about

1. database

first, thought about lightweight sqlite but it doesn't support row level locking which is important for write heavy programs like games. but I havent tested yet...
postgres or mysql sounds heavy I want something light
never tried redis but by the nature of key value store I think it won't handle structured data well but I guess I'll have to use this? think I have to do some more searching
