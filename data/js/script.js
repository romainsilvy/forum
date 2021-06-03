const addThread = () => {
    document.querySelector('#button-input-add-thread').addEventListener('click', ()=>{
        const divAddThread = document.getElementById('add-thread')
        const input = document.querySelector("#input-add-thread")

        addTitleFromInput(input, divAddThread)
        addInputForContent(divAddThread)
        hideElement(input)
      })
}

const addTitleFromInput = (input, box) => {
    const titleValue = input.value
    let newTitle = document.createElement("h2");
    const text = document.createTextNode(titleValue);
    newTitle.appendChild(text);
    box.append(newTitle)
}

const addInputForContent = (box) => {
    let newInput = document.createElement("input")
    newInput.setAttribute("type", "text")
    newInput.setAttribute("placeholder", "Entrez votre message ici....")
    box.append(newInput)
}

const hideElement = (el) => {
    el.classList.add("hidden")
}

addThread()