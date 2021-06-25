document.querySelector("#submitthread").addEventListener('click', () => {
    let title = document.getElementById("input-add-thread").value
    let content = document.getElementById("créa_thread").value
    var radios = document.getElementsByName('drone');
    for (var i = 0, length = radios.length; i < length; i++) {
        if (radios[i].checked) {
            category = radios[i].value
            break
        }
    }
    addThread(title, content, category)
})

function addThread(title, content, category) {
    console.log(title, content, category)
    const body_content = document.querySelectorAll("body")
    fetch("/thread", {
        method: "POST",
        headers: {
            "content-type": "application/json"
        },
        body: JSON.stringify({// envoyer dans le go
            title: title,
            content: content,
            category: category,
        })
    })
        .then((response) => {
            return response.json()
        })
        .then((result) => {
            console.log("dans result")
            const newDivInsideCenter = document.createElement("div")
            newDivInsideCenter.classList.add("inside-center a-thread")
            console.log(newDivInsideCenter)

            const newDivEditButtons = document.createElement("div")

            const newButtonDelete = document.createElement("button")
            newButtonDelete.setAttribute("id", "buttonDelete")
            newButtonDelete.setAttribute("type", "button")
            newButtonDelete.setAttribute("onclick", "javascript:onclickdelete(" + result.id_th + ")")
            console.log(newButtonDelete)

            const newButtonEdit = document.createElement("button")
            newButtonEdit.setAttribute("id", "buttonEdit")
            newButtonEdit.setAttribute("type", "button")
            newButtonEdit.setAttribute("onclick", "javascript:onclickedit(" + result.id_th + ")")

            newDivEditButtons.appendChild(newButtonDelete)
            newDivEditButtons.appendChild(newButtonEdit)
            newDivInsideCenter.appendChild(newDivEditButtons)

            console.log(newDivInsideCenter)

        })
        .catch((err) => {
            throw err
        })
}

document.querySelector('#button-input-add-thread').addEventListener('click', () => {
    const x = document.getElementById("input-add-thread").value
    if (x !== "") {
        document.querySelector("#namethread").value = x
        document.querySelector('#thread-background').classList.add("threadGrid")
        document.querySelector('#pop-up').classList.add("connexion-background")

        document.querySelector('#pop-up').addEventListener('click', () => {
            document.querySelector('#thread-background').classList.remove("threadGrid")
            document.querySelector('#pop-up').classList.remove("connexion-background")
        })
    }
})

function like(id_th, value) {
    const body_content = document.querySelectorAll("body")
    fetch("/like", {
        method: "POST",
        headers: {
            "content-type": "application/json"
        },
        body: JSON.stringify({
            id_th: id_th.toString(),
            value: value.toString(),
        })
    })
        .then((response) => {
            return response.json()
        })
        .then((result) => {

            const bDislike = `#btnDislike${id_th.toString()}`
            document.querySelector(bDislike).innerHTML = result[1]
            const bLike = `#btnLike${id_th.toString()}`
            document.querySelector(bLike).innerHTML = result[0]

            if (value == (-1)) {
                changeColorRed(id_th)
            } else if (value == 1) {
                changeColorGreen(id_th)
            }
        })
        .catch((err) => {
            throw err
        })
}

function onclickedit (postID) {
    document.querySelector("#submitthreadedit").setAttribute("value", postID)

    title = document.querySelector(`#title${postID}`).innerHTML
    document.querySelector("#namethreadedit").setAttribute("value", title)

    content = document.querySelector(`#content${postID}`).innerHTML
    document.getElementById("modif_thread").innerHTML = content


    document.querySelector('#thread-edit').classList.add("threadGrid")
    document.querySelector('#pop-up').classList.add("connexion-background")

    document.querySelector('#pop-up').addEventListener('click', () => {
        document.querySelector('#thread-edit').classList.remove("threadGrid")
        document.querySelector('#pop-up').classList.remove("connexion-background")
    })
}

const onclickdelete = (postID) => {
    document.querySelector("#confirmDeleteBtn").setAttribute("value", postID)
    document.querySelector('#pop-up').classList.add("connexion-background")
    document.querySelector('#thread-delete').classList.add("popUp")

    document.querySelector('#pop-up').addEventListener('click', () => {
        document.querySelector('#thread-delete').classList.remove("threadGridDelete")
        document.querySelector('#pop-up').classList.remove("connexion-background")
        document.querySelector('#thread-delete').classList.remove("popUp")
    })
}

const changeColorGreen = (id_th) => {
    console.log(id_th)
    if (document.querySelector("#testHour")) {
        let btn1 = document.querySelector(`.green${id_th}`);
        let btn2 = document.querySelector(`.red${id_th}`);

        if (btn2.classList.contains(`btnRed`)) {
            btn2.classList.remove(`btnRed`);
        }
        btn1.classList.add(`btnGreen`);
    }
}

const changeColorRed = (id_th) => {
    if (document.querySelector("#testHour")) {
        let btn1 = document.querySelector(`.green${id_th}`);
        let btn2 = document.querySelector(`.red${id_th}`);

        if (btn1.classList.contains(`btnGreen`)) {
            btn1.classList.remove(`btnGreen`);
        }
        btn2.classList.add(`btnRed`);
    }
}



document.querySelector('#profile-picture').addEventListener('click', () => {

    if (document.cookie.includes("auth")) {
        window.location.replace('http://localhost:8080/profil')
    } else {
        document.querySelector('#pop-up').classList.add("connexion-background")
        document.querySelector('#connexion').classList.add("popUp")

        document.querySelector('#sInscrire').addEventListener('click', () => {
            document.querySelector('#connexion').classList.remove("popUp")
            document.querySelector('#inscription').classList.add("popUp")
        })

        document.querySelector('#pop-up').addEventListener('click', () => {
            document.querySelector('#pop-up').classList.remove("connexion-background")
            document.querySelector('#connexion').classList.remove("popUp")
            document.querySelector('#inscription').classList.remove("popUp")
        })
    }
})

if (document.cookie.includes("auth")) {
    document.querySelector("#add-thread").classList.add("add-thread-auth")
}