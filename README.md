# Simple Wordle Solver

### Helping you select optimal guesses while playing a game of [Wordle](https://www.nytimes.com/games/wordle/index.html).

### Accessible at: [hnasser-dev.github.io/wordle-solver/](https://hnasser-dev.github.io/wordle-solver/)

<img width="500" height="645" alt="image" src="https://github.com/user-attachments/assets/7eeb3dbb-6972-472d-91eb-6743f6f17b22" />

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

### Personal Comments

#### This project isn’t really an ideal use case for WASM.
#### I initially started this project to gain more experience with Go and expected that some functionality would require a backend. However, it soon became apparent that I could make it work with just a static website. Despite this, I chose not to translate the code to JavaScript and instead experimented with WASM.
#### What I learned: WASM introduces unneeded complexity for a small project like this. The business logic would likely have been cleaner and more concise if written in JavaScript, and the page size would have been smaller as well. That being said, it was fun using WASM for the first time, and I can definitely see its merits! It just wasn't ideal for this particular project.
