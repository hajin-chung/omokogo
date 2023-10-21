package main

/*
Match Making Algorithm

1. Fast simple one --- X

wait until two players
match make those two players

2. Waiting Matcher --- X

each `n` seconds look at the queue
sort queue by score and make 0, 1 / 2, 3 / 4, 5 and so on...

3. Exponential Matcher

each players  in the queue starts with `a` diff points
each `n` seconds look at the queue
if playerN's and playerM's score diff <= 
	playerN's diff point + playerM's diff point
then match those two, for players missed match multiply diff point by `r`
(1 <= `r`)
lets think of reasonable numbers
*/
