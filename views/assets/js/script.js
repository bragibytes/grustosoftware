M.AutoInit();

function submitComment(e, post){
    e.preventDefault();
    const form = document.getElementsByClassName('comment-form')[0];
    const cbody = form.valueOf('body');

    alert("here is the post ", JSON.stringify(post))
}

function logout() {

    const options = {
        method:'DELETE'
    };

    fetch('/logout', options)
        .then(()=>window.location.replace('/'))
        .catch(err=>console.error(err))
};

