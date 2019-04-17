// if(document.getElementsByClassName('login-form')){
//     for(let form of document.getElementsByClassName('login-form')){
//         form.addEventListener('submit', login)
//     }
// }
if(document.getElementsByClassName('logout-btn')){
    for(let btn of document.getElementsByClassName('logout-btn')){
        btn.addEventListener('click', logout)
    }
}

function logout() {

    const options = {
        method:'DELETE'
    }

    fetch('/logout', options)
        .then(()=>window.location.replace('/blog'))
        .catch(err=>console.error(err))
}
//
// function login(e) {
//     e.preventDefault();
//
//     const data = {
//         name:e.target.elements[0].value,
//         password:e.target.elements[1].value
//     }
//
//     const options = {
//         method:'POST',
//         headers:{
//             'Content-Type':'application/json'
//         },
//         body:JSON.stringify(data)
//     }
//
//     fetch('/login', options)
//         .catch(err=>console.error(err))
//
//
// }
