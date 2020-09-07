# Repo Structure
1. Go webservice
2. react front-end (front-end)

eventually, serve the frontend from the webservice so you don't have to deploy twice.

## Getting set up

install golang and npm

get all the go dependencies 
```
go build ./...
```
get all the npm dependencies

run:
```
npm install
```
inside `front-end`

## Test it
```
go run .* 
```
that's going to launch a webservice at localhost:8080

cd in to `front-end`

```
npm start
```



### Possible Architecture?
There's so many things that things can happen
Judge happens at the end of the day of no one voted -> specific conditional.

Checks at the start of every role.

## TODO
move players over to a map with people actually in the game.

# "Stages"
Day/Night_Stage_Role

Day_vote

Day_judge_judge

night

night_gravedigger_gravedigger -> show the gravedigger screen with the "read" button<br>
night_gravedigger_dead -> show the gravedigger screen without the "read" button

night_gravedigger_ready -> check whether gravedigger has pressed the "read" button

night_angel_angel<br>
night_angle_demon

night_angel_ready -> check whether the angels have voted for a person to protect

night_demon_angel<br>
night_demon_demon
- json : first vote, second vote




## Planning
- Get 5 man game working
- Add roles one at a time




