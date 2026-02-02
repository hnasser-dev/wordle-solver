/*
TODO
- Mobile support
    - On screen keyboard
- (Possibly?) Ability to go back a row

*/

const colourClasses = [
    ["bg-gray-50", "grey"],
    ["bg-orange-200", "yellow"],
    ["bg-green-500", "green"],
];
const disabledColour = "bg-gray-400";

let guessNum = 0;
let currentGuessArr = [];

const showErrorPopup = (msg) => {
    document.querySelector("#game-error-inner").innerHTML = msg;
    document.querySelector("#game-error-outer").classList.remove("hidden");
};

const showGameCompletePopup = (msg) => {
    document.querySelector("#game-complete-inner").innerHTML = msg;
    document.querySelector("#game-complete-outer").classList.remove("hidden");
};

const charIsLetter = (char) => {
    return /^[a-z]$/i.test(char);
};

const setInactiveColour = (elem) => {
    [...elem.classList].forEach((cls) => {
        if (cls.startsWith("bg-")) {
            elem.classList.remove(cls);
        }
    });
    elem.classList.add("bg-gray-400");
};

const populateRowPanels = (rowIdx) => {
    const rows = document.querySelectorAll(".game-row");
    rows.forEach((row) => {
        const letterPanels = row.querySelectorAll(".letter-panel");
        const sidePanel = row.querySelector(".row-side-panel");
        // update the active row
        if (row.id === `game-row-${rowIdx}`) {
            for (let i = 0; i < letterPanels.length; i++) {
                const panel = letterPanels[i];
                panel.classList.remove("bg-gray-400");
                panel.classList.add("bg-gray-50");
                if (i < currentGuessArr.length) {
                    panel.innerHTML = currentGuessArr[i];
                } else {
                    panel.innerHTML = "";
                }
            }
            sidePanel.classList.remove("bg-gray-400");
            sidePanel.classList.add("bg-purple-100");
        } else {
            for (const elem of [sidePanel, ...letterPanels]) {
                setInactiveColour(elem);
            }
        }
    });
};

const updateRows = (suggestions, guessNum) => {
    const rowSidePanels = document.querySelectorAll(".row-side-panel");
    rowSidePanels.forEach((sidePanel, idx) => {
        if (idx == guessNum) {
            const submitBtn = document.createElement("button");
            submitBtn.id = "submit-guess-btn";
            submitBtn.innerHTML = "Submit";
            submitBtn.classList.add(
                "w-20",
                "h-16",
                "bg-green-200",
                "block",
                "px-3",
                "py-2.5",
                "border",
                "border-default-medium",
                "text-md",
                "uppercase",
                "text-center",
                "rounded-md",
                "hover:cursor-pointer"
            );
            submitBtn.addEventListener("click", (event) => {
                btn = event.currentTarget;
                if (guessNum > 5 || currentGuessArr.length != 5) {
                    btn.disabled = true;
                    return;
                }
                const guess = currentGuessArr.join("").toLowerCase();
                const colourPattern = getColourPattern(guessNum);
                if (colourPattern.every((val) => val === "green")) {
                    showGameCompletePopup(
                        `Congratulations! The correct answer is <b>${guess}</b>`
                    );
                    return;
                }
                const loadingSpinner =
                    document.querySelector("#loading-spinner");
                loadingSpinner.classList.remove("hidden");
                // set timeout allows the removal of the loadingSpinner again
                setTimeout(() => {
                    suggestions = guessHelper.getSuggestions(
                        guess,
                        colourPattern
                    );
                    if (!suggestions || suggestions.length == 0) {
                        showErrorPopup(
                            "No possible answers left!<br>Are you sure you entered all the colours in correctly?"
                        );
                        loadingSpinner.classList.add("hidden");
                        return;
                    } else if (suggestions.length == 1) {
                        showGameCompletePopup(
                            `The correct answer is <b>${suggestions[0]}</b>`
                        );
                        loadingSpinner.classList.add("hidden");
                        return;
                    }
                    guessNum++;
                    updateRows(suggestions, guessNum);
                    loadingSpinner.classList.add("hidden");
                    btn.remove();
                }, 10);
            });
            const selector = document.createElement("select");
            selector.classList.add(
                "w-28",
                "h-16",
                "block",
                "px-3",
                "py-2.5",
                "border",
                "border-default-medium",
                "text-lg",
                "uppercase",
                "text-center"
            );
            for (const suggestion of suggestions) {
                selector.add(new Option(suggestion.toUpperCase(), suggestion));
            }
            const suggestionOnChange = () => {
                const selectedValue =
                    selector.options[selector.selectedIndex].value;
                currentGuessArr = selectedValue.toUpperCase().split("");
                populateRowPanels(guessNum);
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

document.addEventListener("keydown", (event) => {
    // prevent holding the button down
    if (event.repeat) {
        return;
    }
    isLetter = charIsLetter(event.key);
    isBackspace = event.key === "Backspace";
    if (isLetter || isBackspace) {
        if (isBackspace) {
            if (currentGuessArr.length >= 1) {
                currentGuessArr.pop();
            }
        } else {
            const char = event.key.toUpperCase();
            if (currentGuessArr.length < 5) {
                currentGuessArr.push(char);
            }
        }
        populateRowPanels(guessNum);
    }
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

document.querySelector("#restart-btn").addEventListener("click", () => {
    window.location.reload();
});

// optimalFirstGuesses defined in wasm

window.mainJsInit = () => {
    updateRows(optimalFirstGuesses, guessNum);
};
