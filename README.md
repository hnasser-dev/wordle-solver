# Simple Wordle Solver

### Helping you select optimal guesses while playing a game of [Wordle](https://www.nytimes.com/games/wordle/index.html).

### Accessible at: [hnasser-dev.github.io/wordle-solver/](https://hnasser-dev.github.io/wordle-solver/)

<img width="699" height="887" alt="image" src="https://github.com/user-attachments/assets/bee083be-9d56-4847-9c61-5d850643f23d" />

### Instructions

- Make a guess in Wordle, and note down the colour of the tiles.
- In the solver, populate the row with the guess you just made.
    - If you made the optimal guess, the solver will already be populated with that guess.
    - If you chose your own guess, use the dropdown menu or keyboard (on-screen or physical) to input the guess you made.
- Click on the tiles to toggle their colours - ensure the colours match what you observed in the Wordle game.
- Click `Submit`, and the solver will auto-populate the next row with the optimal* next guess.
- Repeat the process outlined above (make a guess, input information about your guess into the solver, then follow the optimal guess suggested at each step).
- If you make a mistake at any time, you may `Undo` your guess, `Restart` the game or refresh the page.

#### \* The guess that will, on average, reduce the number of remaining possible answers by the largest amount. This does not necessarily result in an optimal game.

### Additional Info

- This repository contains some additional code which doesn't run in the web application.
- Most notably, some Go code that allows you to simulate games, including with user-specified first guesses, and a "dumb" mode that will select the *least* optimal guesses.
- Feel free to play around with these if you like.
