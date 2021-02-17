# GoBA - Go Battle Arena

[Try GoBA!](http://goba-env.eba-hiw6diij.ca-central-1.elasticbeanstalk.com/game)

GoBA is a simple multiplayer online battle arena (MOBA) game. The server is written in Go, and the client was built using Typescript and Vue. The game is deployed on AWS using Elastic Beanstalk.

## Performance

The focus of this project is not on gameplay. Instead, it is to see how concurrency features of Go, and websockets can be used to achieve real-time client server communication. GoBA's server is capable of running multiple games at 64 ticks per seconds (TPS) on a single AWS t3nano EC2 instance. This is a great [article](https://technology.riotgames.com/news/valorants-128-tick-servers) by the creators of VALORANT a popular online first person shooter that explains how they were able to achieve 128 TPS and run 3 games on a single CPU core. This is extremely impressive as each frame must be processed in under 2.6 milliseconds.

To analyze the performance of GoBA I used [pprof](https://github.com/google/pprof) a profiling tool for Go built by Google. From a 30 second CPU profile with 1.59 seconds total samples it generated the diagram below. In box 1 we see the go routine that runs a game. In this go routine there's a select statement that chooses between executing a tick, connecting / disconnecting players, and stopping the game. In box 2 we see the execution of a game tick. In box 3 we see the calculation of collisions between circles and rectangles for hit boxes. In box 4 we see the time spent executing a tick for each specific team. This will calculate which enemies members of this team can / can't see and send this information client. This leads to box 5 where we see this data get encoded to JSON.

![pprof](./client/screenshots/profile1-diagram.png)

## Gameplay

The objective of GoBA is to get as many kills as possible until the game automatically ends after 15 minutes. There is an ability to shoot, and an ability to dash. When your character dies you will respawn immediately and the scoreboard will be updated. The scoreboard tracks kills, deaths, and assists. The vision system of GoBA is based on the one in League of Legends. Vision is shared with your teamates meaning anything your team can see you can see. Vision is granted based on line of sight. Walls and bushes obstruct line of sight. If a player is inside of a bush they can see out but other players cannot see in. The projectiles shot by a player also grant vision.

## Next Steps

There were some interesting features I considered adding however, they ended up being outside the current scope of the project:

- Match making system
- Create a more structured game with more characters to play, towers, monst
- Visualize vision on the client with "fog of war". There's an interesting article on how this is done in League of Legends [here](https://technology.riotgames.com/news/story-fog-and-war)
- Implement a path finding algorithm for players to navigate the map automatically. This [article](https://www.researchgate.net/publication/315456384_Applying_Theta_in_Modern_Game) explains a modified version of A\* search algorithm called Theta\* that looks very promising
- Implement various mechanisms to account for latency
- Implement [quad trees](https://en.wikipedia.org/wiki/Quadtree) to improve the efficiency of collision detection
