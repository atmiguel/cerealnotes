$(function() {
    $("#logout-button").click(() => {
        $.ajax({
            url: '/api/session',
            type: 'DELETE',
            success: function() {
                alert("you've been successfully logged out");
                location.reload();
            }
        });
    });
});
