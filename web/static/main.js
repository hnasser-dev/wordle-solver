let guessNum = 0;
let currentGuessArr = [];
const colourClasses = [
    ["bg-gray-50", "grey"],
    ["bg-orange-200", "yellow"],
    ["bg-green-500", "green"],
];

const charIsLetter = (char) => {
    return /^[a-z]$/i.test(char);
};

const populateRowPanels = (rowIdx) => {
    const row = document.querySelector(`#letter-row-${rowIdx}`);
    const letterPanels = row.querySelectorAll(".letter-panel");
    for (let i = 0; i < letterPanels.length; i++) {
        if (i < currentGuessArr.length) {
            letterPanels[i].innerHTML = currentGuessArr[i];
        } else {
            letterPanels[i].innerHTML = "";
        }
    }
};

const getColourPattern = (rowIdx) => {
    const colourPattern = [];
    const row = document.querySelector(`#letter-row-${rowIdx}`);
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
    isEnter = event.key === "Enter";
    if (isLetter || isBackspace || isEnter) {
        // submit guess
        if (isEnter) {
            if (currentGuessArr.length === 5) {
                console.log("pretending to submit:", currentGuessArr);
                document.querySelector("#submit-guess-btn").click();
            }
            return;
        } else if (isBackspace) {
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

document.querySelector("#submit-guess-btn").addEventListener("click", (btn) => {
    if (guessNum > 5 || currentGuessArr.length != 5) {
        btn.disabled = true;
        return;
    }
    const guess = currentGuessArr.join("").toLowerCase();
    const colourPattern = getColourPattern(guessNum);
    console.log(guess, colourPattern);
    console.log("calling go func...");

    const loadingSpinner = document.querySelector("#loading-spinner");
    loadingSpinner.classList.remove("hidden");
    // set timeout allows the removal of the loadingSpinner again
    setTimeout(() => {
        suggestions = guessHelper.getSuggestions(
            guess,
            colourPattern,
            "normal"
        );
        console.log(suggestions);
        loadingSpinner.classList.add("hidden");
        const rowSidePanels = document.querySelectorAll(".row-side-panel");
        guessNum++;
        rowSidePanels.forEach((sidePanel, idx) => {
            if (idx == guessNum) {
                const selector = document.createElement("select");
                selector.classList.add(
                    "w-40",
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
                for (suggestion of suggestions) {
                    selector.add(
                        new Option(suggestion.toUpperCase(), suggestion)
                    );
                }
                const suggestionOnChange = () => {
                    const selectedValue =
                        selector.options[selector.selectedIndex].value;
                    currentGuessArr = selectedValue.toUpperCase().split("");
                    populateRowPanels(guessNum);
                };
                selector.addEventListener("change", suggestionOnChange);
                sidePanel.replaceChildren(selector);
                // force default value to populate the dropdown first
                if (suggestions) {
                    selector.value = suggestions[0];
                    suggestionOnChange();
                }
            } else {
                sidePanel.innerHTML = "";
            }
        });
        currentGuessArr = [];
        console.log(guessNum);
    }, 10);
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

const wordDropDowns = document.querySelectorAll(".letter-panel");
