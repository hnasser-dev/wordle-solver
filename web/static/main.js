let guessNum = 0;

document.querySelector("#get-suggestions-btn").onclick = (btn) => {
    if (guessNum > 5) {
        btn.disabled = true;
        return;
    }
    console.log("calling go func...");
    suggestions = guessHelper.getSuggestions(
        "hello",
        ["yellow", "yellow", "grey", "grey", "grey"],
        "normal"
    );
    console.log("suggestions", suggestions);
    const guessNumStr = guessNum.toString();
    const selectors = document.querySelectorAll(".suggestions-selector");
    selectors.forEach((selector) => {
        selector.innerHTML = "";
        if (selector.dataset.rowIdx == guessNumStr) {
            for (const suggestion of suggestions) {
                selector.add(new Option(suggestion, suggestion));
            }
            selector.disabled = false;
        } else {
            selector.disabled = true;
        }
    });
    guessNum++;
    console.log(guessNum);
};

const letterPanels = document.querySelectorAll(".letter-panel");
letterPanels.forEach((panel) => {
    panel.addEventListener("click", () => {
        colourClasses = ["bg-gray-50", "bg-orange-200", "bg-green-500"];
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

const wordDropDowns = document.querySelectorAll(".letter-panel");
