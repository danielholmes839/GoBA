# GoBA - Go Battle Arena

[Try GoBA!](http://goba-env.eba-hiw6diij.ca-central-1.elasticbeanstalk.com/game)

GoBA is a simple multiplayer online battle arena (MOBA) game. The server is written in Go, and the client was built using Typescript and Vue. The game is deployed on AWS using Elastic Beanstalk. The focus of this project is not on gameplay. Instead, it is to see how concurrency features of Go, and websockets can be used to achieve real-time client server communication. The server is capable of running multiple games at 64 ticks per seconds (TPS) on a single AWS t3nano EC2 instance. I believe this is a success. This [article](https://technology.riotgames.com/news/valorants-128-tick-servers) by the creators of VALORANT a popular online first person shooter explains how they were able to achieve 128 TPS and run 3 games on a single CPU core. This is extremely impressive as each frame must be processed in under 2.6 milliseconds.

## Features

- Abilities
- Teams
- Vision
- Scoreboard
- Among Us style codes for creating and joining games

## Next Steps

Other interesting features would include:

- Match making system
- More abilties / characters to play.
- Visualize vision on client with "fog of war". There's an interesting article on how this is done in League of Legends [here](https://technology.riotgames.com/news/story-fog-and-war)
- Implement a path finding algorithm for the player to navigate the map automatically. This [article](https://www.researchgate.net/publication/315456384_Applying_Theta_in_Modern_Game) explains a modified version of A\* search algorithm called Theta\* that looks very promising
- Implement various mechanisms to account for latency
