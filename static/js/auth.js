function displayMessage(msgText, type) {
    var d = $("<div>");
    d.attr("role", "alert");
    d.addClass("alert alert-" + type);
    d.text(msgText);
    $("body").prepend(d.delay(800).fadeOut("slow"));
}

function handleHTTPError(msgText, msgStatus) {
    console.log(msgStatus);
    var d = $("<div>");
    d.attr("role", "alert");
    d.text(msgText);

    switch(msgStatus) {
    case 401:
        d.addClass("alert alert-danger")
        // redirect on 401 errors
        window.location.replace("/login");
        break;
    case 403, 500:
        d.addClass("alert alert-danger")
        break;
    default:
        d.addClass("alert alert-light")
        break;
    }

    $("body").prepend(d.delay(800).fadeOut("slow"));
}

function getToken() {
    var email = $("#emailInput").val(),
        password = $("#passwordInput").val();

    $.ajax({
        url: "/get-token",
        method: 'POST',
        data: {
            email: email,
            password: password
        }
    }).done(function(token) {
        console.log(token);
        // store in web storage
        //window.localStorage.setItem('token', token);
        window.location.replace("/");
    });
}

// jwt in web storage
//$.ajaxPrefilter(function( options ) {
//    options.beforeSend = function (xhr) {
//        xhr.setRequestHeader('Authorization', 'Bearer '+localStorage.getItem('token'));
//    }
//});
