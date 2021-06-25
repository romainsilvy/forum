const addThreadPopUp = () => {
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
}



const like = (id_th, value) => {
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
        })
        .catch((err) => {
            throw err
        })
}

const onclickedit = (postID) => {
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

const likeColor = () => {
    if (document.querySelector("#testHour")) {
        let btn1 = document.querySelector('#green');
        let btn2 = document.querySelector('#red');

        btn1.addEventListener('click', function () {
            if (btn2.classList.contains('red')) {
                btn2.classList.remove('red');
            }
            this.classList.toggle('green');
        });

        btn2.addEventListener('click', function () {
            if (btn1.classList.contains('green')) {
                btn1.classList.remove('green');
            }
            this.classList.toggle('red');
        });
    }
}

const connexionInscriptionVoirProfilPopUp = () => {
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
}

const displayAddThread = () => {
    if (document.cookie.includes("auth")) {
        document.querySelector("#add-thread").classList.add("add-thread-auth")
    }
}





addThreadPopUp()
likeColor()
connexionInscriptionVoirProfilPopUp()
displayAddThread()
