
$(function() {
	$("#logout-button").click(() => {
		$.ajax({
			url: '/session',
			type: 'DELETE',
			success: function(result) {
				alert("you've been successfully logged out");
				location.reload();
			}
		});
	});
});