# Game "Winter is coming"

Test task for **Mysterium Network**

## Run server
```
go run main.go
```

## Run tests
```
env GOCACHE=off go test ./tests
```

## Connect via telnet
Default port is 3000
```
telnet {IP} 3000
telnet 127.0.0.1 3000
```

## Command list
```
START {name}   - starts a new game (ex. START John)
SHOOT {x} {y}  - shoots at given coordinates (ex. SHOOT 0 1)
SHARE          - shares your current game to be accessible for friends
JOIN {GAME_ID} - joins the provided game
```
