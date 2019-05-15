M.AutoInit();

function vote(type, id){

    console.log("type is : ", type);
    console.log("id is : ", id);

    const data = {
        value:type,
        _parent:id,
    };

    const options = {
        method:'POST',
        headers:{
            "Content-Type":"application/json"
        },
        body:JSON.stringify(data),
    };

    fetch("/api/votes/", options)
        .then(res=>{
            console.log(res);
            window.location.reload();
        })
        .catch(err => console.error(err))
}

function logout() {

    const options = {
        method:'DELETE'
    };

    fetch('/logout', options)
        .then(()=>window.location.replace('/'))
        .catch(err=>console.error(err))
};
