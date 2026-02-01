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

const updateRows = (suggestions, guessNum) => {
    console.log(suggestions, guessNum);
    const rowSidePanels = document.querySelectorAll(".row-side-panel");
    rowSidePanels.forEach((sidePanel, idx) => {
        if (idx == guessNum) {
            const submitBtn = document.createElement("button");
            submitBtn.id = "submit-guess-btn";
            submitBtn.innerHTML = "Submit";
            submitBtn.classList.add(
                "w-24",
                "h-16",
                "bg-green-200",
                "block",
                "px-3",
                "py-2.5",
                "border",
                "border-default-medium",
                "text-lg",
                "uppercase",
                "text-center",
                "rounded-md",
                "hover:cursor-pointer"
            );
            submitBtn.addEventListener("click", (event) => {
                btn = event.currentTarget;
                console.log("clicked!");
                if (guessNum > 5 || currentGuessArr.length != 5) {
                    btn.disabled = true;
                    return;
                }
                const guess = currentGuessArr.join("").toLowerCase();
                const colourPattern = getColourPattern(guessNum);
                const loadingSpinner =
                    document.querySelector("#loading-spinner");
                loadingSpinner.classList.remove("hidden");
                // set timeout allows the removal of the loadingSpinner again
                setTimeout(() => {
                    console.log("getting suggestions...");
                    suggestions = guessHelper.getSuggestions(
                        guess,
                        colourPattern,
                        "normal"
                    );
                    console.log("suggestions:", suggestions);
                    if (!suggestions || suggestions.length == 0) {
                        document.querySelector("#game-error-inner").innerHTML =
                            "No possible answers left!<br>Are you sure you entered all the colours in correctly?";
                        document
                            .querySelector("#game-error-outer")
                            .classList.remove("hidden");
                        loadingSpinner.classList.add("hidden");
                        return;
                    }
                    guessNum++;
                    updateRows(suggestions, guessNum);
                    loadingSpinner.classList.add("hidden");
                    console.log(guessNum);
                    btn.remove();
                }, 10);
            });
            const selector = document.createElement("select");
            selector.classList.add(
                "w-24",
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
                console.log("selectedValue:", selectedValue);
                currentGuessArr = selectedValue.toUpperCase().split("");
                console.log("currentGuessArr:", currentGuessArr);
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

updateRows(["raise", "thing"], guessNum);
