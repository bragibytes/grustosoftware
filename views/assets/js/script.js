$(document).ready(function () {
    const con = new BlogController();
    con.test();

    $('.dropdown-trigger').dropdown();

    // logout
    $('.logout-btn').bind('click', function(e){

        $.ajax({
            'url':'/logout',
            'type':'DELETE',
            'data':null,
            'success':()=>window.location.replace("/"),
            'error':(err)=>console.error(err)
        })
    });

    // create post
    $('#post-form').bind('submit', function (e) {
        e.preventDefault();
        console.log(e)

    })

});

class BlogController {

    constructor(){
        this.title = document.getElementById("title")
        this.body = document.getElementById("body")
    }

    test(){
        alert("at least this class is working!!")
    }

}






