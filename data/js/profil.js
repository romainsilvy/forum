const manageProfilPage = () => {
    alert("ses mort")
    document.querySelector('.btn-6').addEventListener('click', () => {
        function deleteCookie(name) { setCookie(name, '', -1); }
        function setCookie(name, value, days) {
            var d = new Date;
            d.setTime(d.getTime() + 24 * 60 * 60 * 1000 * days);
            document.cookie = name + "=" + value + ";path=/;expires=" + d.toGMTString();
        }
        deleteCookie("auth")
        window.location.replace('http://localhost:8080')
    })
}

alert("ses mort")
manageProfilPage()