alert("dans le script")


document.querySelector("#submitthread").addEventListener('click', ()=> {
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

function addThread (title, content, category) {
    console.log(title, content, category)
    const body_content  = document.querySelectorAll("body")
    fetch("/thread", {
    method : "POST",
    headers : {
        "content-type" : "application/json"
    },
    body : JSON.stringify({// envoyer dans le go
        title : title,
        content : content,
        category : category,
    })
    })
    .then((response)=>{
    return response.json()
    })
    .then((result)=>{
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
    .catch((err) =>{
    throw err
    })
}

document.querySelector('#button-input-add-thread').addEventListener('click', ()=>{
    const x = document.getElementById("input-add-thread").value
    if ( x !== "") {
        document.querySelector("#namethread").value = x
        document.querySelector('#thread-background').classList.add("threadGrid")
        document.querySelector('#pop-up').classList.add("connexion-background")

        document.querySelector('#pop-up').addEventListener('click', ()=>{
        document.querySelector('#thread-background').classList.remove("threadGrid")
        document.querySelector('#pop-up').classList.remove("connexion-background")
        })
    }
})

function like(id_th, value) {
    const body_content  = document.querySelectorAll("body")
    fetch("/like", {
    method : "POST",
    headers : {
        "content-type" : "application/json"
    },
    body : JSON.stringify({// envoyer dans le go
        id_th : id_th.toString(),
        value : value.toString(),
    })
    })
    .then((response)=>{
    return response.json()
    })
    .then((result)=>{


        
        const bDislike = `#btnDislike${id_th.toString()}`
        document.querySelector(bDislike).innerHTML = result[1]
        const bLike = `#btnLike${id_th.toString()}`
        document.querySelector(bLike).innerHTML = result[0]
    })
    .catch((err) =>{
    throw err
    })
}


function onclickedit (postID) {
        console.log(postID)
        document.querySelector('#thread-edit').classList.add("threadGrid")
        document.querySelector('#pop-up').classList.add("connexion-background")

        document.querySelector('#pop-up').addEventListener('click', ()=>{
        document.querySelector('#thread-edit').classList.remove("threadGrid")
        document.querySelector('#pop-up').classList.remove("connexion-background")
        })
}

function onclickdelete(postID) {
    console.log(postID)
        document.querySelector('#pop-up').classList.add("connexion-background")
        document.querySelector('#thread-delete').classList.add("popUp")

        document.querySelector('#pop-up').addEventListener('click', ()=>{
        document.querySelector('#thread-delete').classList.remove("threadGridDelete")
        document.querySelector('#pop-up').classList.remove("connexion-background")
        document.querySelector('#thread-delete').classList.remove("popUp")
    })
}

if (document.querySelector("#testHour")) {
    let btn1 = document.querySelector('#green');
    let btn2 = document.querySelector('#red');

    btn1.addEventListener('click', function() {
    
    if (btn2.classList.contains('red')) {
        btn2.classList.remove('red');
    } 
    this.classList.toggle('green');
});

btn2.addEventListener('click', function() {
    
    if (btn1.classList.contains('green')) {
        btn1.classList.remove('green');
    } 
    this.classList.toggle('red');
});
}

document.querySelector('#profile-picture').addEventListener('click', ()=>{

    if (document.cookie.includes("auth")) {
        window.location.replace('http://localhost:8080/profil')
    } else {

        document.querySelector('#pop-up').classList.add("connexion-background")
        document.querySelector('#connexion').classList.add("popUp")

        document.querySelector('#sInscrire').addEventListener('click', ()=>{
        document.querySelector('#connexion').classList.remove("popUp")
        document.querySelector('#inscription').classList.add("popUp")
    })

    document.querySelector('#pop-up').addEventListener('click', ()=>{
    document.querySelector('#pop-up').classList.remove("connexion-background")
    document.querySelector('#connexion').classList.remove("popUp")
    document.querySelector('#inscription').classList.remove("popUp")
    })
}
})

if (document.cookie.includes("auth")) {
    document.querySelector("#add-thread").classList.add("add-thread-auth")
}
