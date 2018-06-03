$(function() {
    $.get('/api/user', function(usersById) {
        $.get('/api/note', function(note) {
            $('#notes').append(
                $('<div>').addClass('note').append(
                    $('<span>', {text: usersById[note.authorId].displayName})
                ).append(
                    $('<span>', {text: ' - '})
                ).append(
                    $('<span>', {text: 'NoteType: ' + note.type})
                ).append(
                    $('<span>', {text: ' - '})
                ).append(
                    $('<span>', {text: note.creationTime})
                ).append(
                    $('<br />')
                ).append(
                    $('<span>', {text: note.content})
                )
            );
        });
    });
});
