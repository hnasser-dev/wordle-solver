const colourClasses = ["bg-gray-50", "bg-orange-200", "bg-green-500"];
const colourClassMapping = new Map([
    ["bg-gray-50", "grey"],
    ["bg-orange-200", "yellow"],
    ["bg-green-500", "green"],
]);
const disabledColour = "bg-gray-600";

let suggestedWords = [];

let submittedGuesses = [];
let submittedColourClasses = [];

let activeGuessArr = [];
let activeGuessColourArr = Array(5).fill("bg-gray-50");

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
        .querySelector(`#game-row-${submittedGuesses.length}`)
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

// overall rerender of ALL rows
const renderAllRows = () => {
    const rows = document.querySelectorAll(".game-row");
    const rowIdxPattern = /^game-row-(\d+)$/;
    // const activeRowIdx = guesses.length;
    rows.forEach((row) => {
        const rowIdx = parseInt(row.id.match(rowIdxPattern)[1]);
        updateRowLetterPanels(row, rowIdx);
        updateSidePanel(row, rowIdx);
    });
};

const updateRowLetterPanels = (row, rowIdx) => {
    console.log(typeof row);
    const letterPanels = row.querySelectorAll(".letter-panel");
    const isActiveRow = rowIdx === submittedGuesses.length;
    letterPanels.forEach((panel, panelIdx) => {
        removeBgColours(panel);
        removeOpacity(panel);
        if (isActiveRow) {
            panel.innerHTML =
                panelIdx < activeGuessArr.length
                    ? activeGuessArr[panelIdx]
                    : "";
            panel.classList.add(activeGuessColourArr[panelIdx]);
        } else if (rowIdx < submittedGuesses.length) {
            const word = submittedGuesses[rowIdx];
            const colourClasses = submittedColourClasses[rowIdx];
            panel.innerHTML = word[panelIdx];
            panel.classList.add("opacity-50");
            panel.classList.add(colourClasses[panelIdx]);
        } else {
            panel.innerHTML = "";
            panel.classList.add(disabledColour);
        }
    });
};

const updateSidePanel = (row, rowIdx) => {
    const rowSidePanel = row.querySelector(".row-side-panel");
    if (rowIdx === submittedGuesses.length) {
        rowSidePanel.replaceChildren(
            createSelector(row, rowIdx),
            createSubmitBtn()
        );
    } else if (rowIdx === submittedGuesses.length - 1) {
        rowSidePanel.replaceChildren(createUndoGuessBtn(rowSidePanel));
    } else {
        rowSidePanel.innerHTML = "";
    }
};

const createSelector = (row, rowIdx) => {
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
    for (const word of suggestedWords) {
        selector.add(new Option(word.toUpperCase(), word));
    }
    const selectorOnChange = () => {
        const selectedValue = selector.options[selector.selectedIndex].value;
        activeGuessArr = selectedValue.toUpperCase().split("");
        updateRowLetterPanels(row, rowIdx);
    };
    selector.addEventListener("change", selectorOnChange);
    // force default value to populate the dropdown first
    if (suggestedWords.length > 0) {
        console.log("suggested words:", suggestedWords, suggestedWords.length);
        selector.value = suggestedWords[0];
        selectorOnChange();
    }
    return selector;
};

const createSubmitBtn = () => {
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
        if (submittedGuesses.length > 5 || activeGuessArr.length != 5) {
            return;
        }
        const guess = activeGuessArr.join("").toLowerCase();
        if (!allValidGuessesList.has(guess)) {
            shakeActiveLetterPanels();
            return;
        }
        console.log("activeGuessColourArr:", activeGuessColourArr);
        const colourPattern = activeGuessColourArr.map((val) =>
            colourClassMapping.get(val)
        );
        console.log("colourPattern:", colourPattern);
        if (colourPattern.every((val) => val === "green")) {
            showGameCompletePopup(
                `Congratulations! The correct answer is <b>${guess}</b>`
            );
            return;
        }
        executeWithLoadingSpinner(() => {
            console.log("guess:", guess, "colourPattern:", colourPattern);
            suggestedWords = guessHelper.getSuggestedWords(
                guess,
                colourPattern
            );
            console.log(
                "retrieved suggested words length",
                suggestedWords.length
            );
            if (!suggestedWords || suggestedWords.length == 0) {
                showErrorPopup(
                    "No possible answers left!<br>Are you sure you entered all the colours in correctly?"
                );
                return;
            } else if (suggestedWords.length == 1) {
                showGameCompletePopup(
                    `The correct answer is <b>${suggestedWords[0]}</b>`
                );
                return;
            }
            // reset global values and re-render rows
            submittedGuesses.push(guess);
            submittedColourClasses.push(activeGuessColourArr);
            activeGuessArr = suggestedWords[0].split("");
            activeGuessColourArr = Array(5).fill("bg-gray-50");
            renderAllRows();
        });
    });
    return submitBtn;
};

const createUndoGuessBtn = (rowSidePanel) => {
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
        submittedGuesses = guessHelper.guesses();
    });
    rowSidePanel.classList.toggle("justify-start");
    return undoGuessBtn;
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
            if (activeGuessArr.length >= 1) {
                activeGuessArr.pop();
            }
        } else {
            const char = key.toUpperCase();
            if (activeGuessArr.length < 5) {
                activeGuessArr.push(char);
            }
        }
        const activeRowIdx = submittedGuesses.length;
        updateRowLetterPanels(activeRowIdx);
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
            const cls = colourClasses[i];
            if (panel.classList.contains(cls)) {
                nextCls = colourClasses[(i + 1) % colourClasses.length];
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
        submittedGuesses = [];
        activeGuessArr = [];
        resetRowPanels();
        updateRows(optimalFirstGuesses);
        hideGameCompletePopup();
        hideErrorPopup();
    });
};

document.querySelector("#restart-btn").addEventListener("click", restartGame);

document
    .querySelector("#game-help-close")
    .addEventListener("click", hideGameHelpPopup);

document
    .querySelector("#help-icon-btn")
    .addEventListener("click", showGameHelpPopup);

window.mainJsInit = () => {
    // optimalFirstGuesses defined in wasm
    suggestedWords = optimalFirstGuesses;
    console.log(suggestedWords);
    renderAllRows();
    const seenHelpPopup = localStorage.getItem("seenHelpPopup");
    if (!seenHelpPopup) {
        showGameHelpPopup();
    }
    console.log("guessHelper numGuesses:", guessHelper.guesses());
};
