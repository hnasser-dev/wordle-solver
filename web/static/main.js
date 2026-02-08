const colourClasses = [
    ["bg-gray-50", "grey"],
    ["bg-orange-200", "yellow"],
    ["bg-green-500", "green"],
];
const disabledColour = "bg-gray-600";

let guesses = [];
let colourPatterns = [];
let currentGuessArr = [];

const executeWithLoadingSpinner = (callback, ...args) => {
    const loadingSpinner = document.querySelector("#loading-spinner");
    loadingSpinner.classList.remove("hidden");
    setTimeout(() => {
        callback(...args);
        loadingSpinner.classList.add("hidden");
    }, 10);
};

const showErrorPopup = (msg) => {
    document.querySelector("#game-error-inner").innerHTML = msg;
    document.querySelector("#game-error-outer").classList.remove("hidden");
};

const hideErrorPopup = () => {
    document.querySelector("#game-error-inner").innerHTML = "";
    document.querySelector("#game-error-outer").classList.add("hidden");
};

const showGameCompletePopup = (msg) => {
    document.querySelector("#game-complete-inner").innerHTML = msg;
    document.querySelector("#game-complete-outer").classList.remove("hidden");
};

const hideGameCompletePopup = () => {
    document.querySelector("#game-complete-inner").innerHTML = "";
    document.querySelector("#game-complete-outer").classList.add("hidden");
};

const showGameHelpPopup = () => {
    const popup = document.querySelector("#game-help-outer");
    popup.classList.remove("pointer-events-none");
    requestAnimationFrame(() => {
        popup.classList.remove("opacity-0");
    });
};

const hideGameHelpPopup = () => {
    const popup = document.querySelector("#game-help-outer");
    popup.classList.add("pointer-events-none", "opacity-0");
    localStorage.setItem("seenHelpPopup", Date.now());
};

const shakeActiveLetterPanels = () => {
    const activeLetterPanels = document
        .querySelector(`#game-row-${guesses.length}`)
        .querySelectorAll(".letter-panel");
    activeLetterPanels.forEach((panel) => {
        panel.classList.add("shake");
        setTimeout(() => {
            panel.classList.remove("shake");
        }, 200);
    });
};

const charIsLetter = (char) => {
    return /^[a-z]$/i.test(char);
};

const removeBgColours = (elem) => {
    [...elem.classList].forEach((cls) => {
        if (cls.startsWith("bg-")) {
            elem.classList.remove(cls);
        }
    });
};

const removeOpacity = (elem) => {
    [...elem.classList].forEach((cls) => {
        if (cls.startsWith("opacity-")) {
            elem.classList.remove(cls);
        }
    });
};

/*
Two functions:

- One that runs when a letter is changed
- One that runs on submit or on reset or on undo


When letter is changed:
- Should only affect the current row


On submit:
- Rerender ALL rows
    - On rowIdx == guessNum: white bg, 




*/

//
const updateLetterPanelsActiveRow = () => {
    const activeRow = document.querySelector(`#game-row-${guesses.length - 1}`);
    const letterPanels = activeRow.querySelectorAll(".letter-panel");
    letterPanels.forEach((panel, idx) => {
        if (idx >= guesses.length) {
            return;
        }
        panel.innerHTML = currentGuessArr[idx];
    });
};

// overall rerender of ALL rows
const rerenderAllRows = () => {
    const rows = document.querySelectorAll(".game-row");
    const rowIdxPattern = /^game-row-(\d+)$/;
    const activeRowIdx = guesses.length;
    rows.forEach((row) => {
        const rowIdx = parseInt(row.id.match(rowIdxPattern)[1]);

        // future row, render as inactive
        if (rowIdx > activeRowIdx) {
        }
    });
};

const resetRowPanels = () => {
    const rows = document.querySelectorAll(".game-row");
    rows.forEach((row) => {
        const letterPanels = row.querySelectorAll(".letter-panel");
        letterPanels.forEach((panel) => {
            removeBgColours(panel);
            removeOpacity(panel);
            panel.classList.add(disabledColour);
            panel.innerHTML = "";
        });
    });
};

const updateRowPanels = () => {
    const rows = document.querySelectorAll(".game-row");
    const rowIdxPattern = /^game-row-(\d+)$/;
    rows.forEach((row) => {
        const letterPanels = row.querySelectorAll(".letter-panel");
        const rowIdx = parseInt(row.id.match(rowIdxPattern)[1]);
        if (rowIdx === guesses.length) {
            for (let i = 0; i < letterPanels.length; i++) {
                const panel = letterPanels[i];
                removeBgColours(panel);
                panel.classList.add("bg-gray-50");
                if (i < currentGuessArr.length) {
                    panel.innerHTML = currentGuessArr[i];
                } else {
                    panel.innerHTML = "";
                }
            }
        } else if (rowIdx === guesses.length - 1) {
            for (const elem of letterPanels) {
                removeOpacity(elem);
                elem.classList.add("opacity-50");
            }
        }
    });
};

const updateRows = (suggestions) => {
    const rowSidePanels = document.querySelectorAll(".row-side-panel");
    rowSidePanels.forEach((sidePanel, idx) => {
        // previous row
        if (idx == guesses.length - 1) {
            const undoGuessBtn = document.createElement("button");
            undoGuessBtn.id = "undo-guess-btn";
            undoGuessBtn.innerHTML = "&#9100;";
            undoGuessBtn.classList.add(
                "text-white",
                "text-2xl",
                "md:text-3xl",
                "p-2",
                "ml-1",
                "text-sm",
                "md:text-md",
                "font-bold",
                "cursor-pointer",
                "active:translate-y-0.5",
                "active:shadow-inner"
            );
            undoGuessBtn.addEventListener("click", () => {
                guessHelper.undoLastGuess();
                // updateRows(); // TODO - fix
                console.log("num guesses before:", guesses.length);
                guesses = guessHelper.guesses();
                console.log("num guesses after:", guesses.length);
            });
            sidePanel.classList.toggle("justify-start");
            sidePanel.replaceChildren(undoGuessBtn);
        } else if (idx == guesses.length) {
            const submitBtn = document.createElement("button");
            submitBtn.id = "submit-guess-btn";
            submitBtn.innerHTML = "Submit";
            submitBtn.classList.add(
                "w-16",
                "sm:w-18",
                "md:w-20",
                "h-11",
                "sm:h-14",
                "md:h-16",
                "border",
                "bg-lime-200",
                "block",
                "border",
                "border-default-medium",
                "text-sm",
                "md:text-md",
                "font-bold",
                "text-center",
                "rounded-md",
                "cursor-pointer",
                "active:translate-y-0.5",
                "active:shadow-inner"
            );
            submitBtn.addEventListener("click", (event) => {
                btn = event.currentTarget;
                if (guesses.length > 5 || currentGuessArr.length != 5) {
                    return;
                }
                const guess = currentGuessArr.join("").toLowerCase();
                if (!allValidGuessesList.has(guess)) {
                    shakeActiveLetterPanels();
                    return;
                }
                const colourPattern = getColourPattern(guesses.length);
                if (colourPattern.every((val) => val === "green")) {
                    showGameCompletePopup(
                        `Congratulations! The correct answer is <b>${guess}</b>`
                    );
                    return;
                }
                executeWithLoadingSpinner(() => {
                    suggestions = guessHelper.getSuggestions(
                        guess,
                        colourPattern
                    );
                    if (!suggestions || suggestions.length == 0) {
                        showErrorPopup(
                            "No possible answers left!<br>Are you sure you entered all the colours in correctly?"
                        );
                        return;
                    } else if (suggestions.length == 1) {
                        showGameCompletePopup(
                            `The correct answer is <b>${suggestions[0]}</b>`
                        );
                        return;
                    }
                    guesses.push(guess);
                    updateRows(suggestions, guesses.length);
                    btn.remove();
                });
            });
            const selector = document.createElement("select");
            selector.classList.add(
                "w-20",
                "sm:w-24",
                "md:w-28",
                "border",
                "h-11",
                "sm:h-14",
                "md:h-16",
                "bg-sky-300",
                "text-md",
                "md:text-lg",
                "font-bold",
                "uppercase",
                "text-center",
                "rounded-md"
            );
            for (const suggestion of suggestions) {
                selector.add(new Option(suggestion.toUpperCase(), suggestion));
            }
            const suggestionOnChange = () => {
                const selectedValue =
                    selector.options[selector.selectedIndex].value;
                currentGuessArr = selectedValue.toUpperCase().split("");
                updateRowPanels();
            };
            selector.addEventListener("change", suggestionOnChange);
            sidePanel.replaceChildren(selector, submitBtn);
            // force default value to populate the dropdown first
            if (suggestions.length > 0) {
                selector.value = suggestions[0];
                suggestionOnChange();
            }
        } else {
            sidePanel.innerHTML = "";
        }
    });
};

const getColourPattern = (rowIdx) => {
    const colourPattern = [];
    const row = document.querySelector(`#game-row-${rowIdx}`);
    const letterPanels = row.querySelectorAll(".letter-panel");
    for (const panel of letterPanels) {
        const classes = panel.classList;
        let classFound = false;
        for (const [cls, colour] of colourClasses) {
            if (classes.contains(cls)) {
                colourPattern.push(colour);
                classFound = true;
                break;
            }
        }
        if (!classFound) {
            throw new Error("panel has no class");
        }
    }
    return colourPattern;
};

const handlePressKey = (key) => {
    isLetter = charIsLetter(key);
    isBackspace = key === "Backspace";
    isEnter = key === "Enter";
    if (isLetter || isBackspace || isEnter) {
        if (isEnter) {
            document.querySelector("#submit-guess-btn").click();
            return;
        } else if (isBackspace) {
            if (currentGuessArr.length >= 1) {
                currentGuessArr.pop();
            }
        } else {
            const char = key.toUpperCase();
            if (currentGuessArr.length < 5) {
                currentGuessArr.push(char);
            }
        }
        updateRowPanels();
    }
};

document.addEventListener("keydown", (event) => {
    // prevent holding the button down
    if (event.repeat) {
        return;
    }
    handlePressKey(event.key);
});

const letterPanels = document.querySelectorAll(".letter-panel");
letterPanels.forEach((panel) => {
    panel.addEventListener("click", () => {
        for (let i = 0; i < colourClasses.length; i++) {
            const cls = colourClasses[i][0];
            if (panel.classList.contains(cls)) {
                nextCls = colourClasses[(i + 1) % colourClasses.length][0];
                panel.classList.toggle(cls);
                panel.classList.toggle(nextCls);
                break;
            }
        }
    });
});

const keyboardKeys = document.querySelectorAll(".keyboard-key-letter");
keyboardKeys.forEach((key) => {
    key.addEventListener("click", (event) => {
        const letter = event.currentTarget.innerHTML;
        handlePressKey(letter);
    });
});
const backspaceKey = document.querySelector("#keyboard-key-backspace");
backspaceKey.addEventListener("click", () => {
    handlePressKey("Backspace");
});

const restartGame = () => {
    executeWithLoadingSpinner(() => {
        resetGuessHelper();
        guesses = [];
        currentGuessArr = [];
        resetRowPanels();
        updateRows(optimalFirstGuesses);
        hideGameCompletePopup();
        hideErrorPopup();
    });
};

// const undoLastGuess = () => {
//     guessHelper.undoLastGuess();
//     guesses = guessHelper.guesses();
//     console.log("undid last guess");
// };

document.querySelector("#restart-btn").addEventListener("click", restartGame);

document
    .querySelector("#game-help-close")
    .addEventListener("click", hideGameHelpPopup);

document
    .querySelector("#help-icon-btn")
    .addEventListener("click", showGameHelpPopup);

window.mainJsInit = () => {
    // optimalFirstGuesses defined in wasm
    updateRows(optimalFirstGuesses);
    const seenHelpPopup = localStorage.getItem("seenHelpPopup");
    if (!seenHelpPopup) {
        showGameHelpPopup();
    }
    console.log("guessHelper numGuesses:", guessHelper.guesses());
};
