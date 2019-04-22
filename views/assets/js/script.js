
// Materialize Inits
$(document).ready(function(){
    $('.collapsible').collapsible();

    $('.logout-btn').bind('click', logout)
});

function submitPost(e){

}

function submitComment(e, parent){

    e.preventDefault();
    console.log("e is ", e);

    console.log("getting to the submitComment function with e as ", e, "\n", "and _parent as ", parent);

    // const body = e.target.valueOf("comment");
    //
    // const options = {
    //     method:'POST',
    //     headers:{
    //         "Content-Type":"application/json",
    //         "Accept":"application/json"
    //     },
    //     data:JSON.stringify({
    //         body:body,
    //         _parent:parent
    //     })
    // };
    //
    // fetch("/api/comments/", options)
    //     .then(res=>console.log(res))
    //     .catch(err=>console.log(err))


}

function logout() {

    const options = {
        method:'DELETE'
    }

    fetch('/logout', options)
        .then(()=>window.location.replace('/blog'))
        .catch(err=>console.error(err))
}

