var USERS_BY_ID = {};

const $createAuthor = function(authorId) {
    const user = USERS_BY_ID[authorId];

    return $('<span>', {text: user.displayName});
};

const $createType = function(type) {
    return $('<span>', {text: type});
};

const $createCreationTime = function(creationTime) {
    return $('<span>', {text: moment(creationTime).fromNow()});
};

const $createContent = function(content) {
    return $('<div>', {text: content});
};

const $createDivider = function() {
    return $('<span>', {text: ' - '});
};

const $createNote = function(noteId, note) {
    const $author = $createAuthor(note.authorId);
    const $type = $createType(note.type);
    const $creationTime = $createCreationTime(note.creationTime);
    const $content = $createContent(note.content);

    const $header = $('<div>').addClass('note-header')
        .append($author).append($createDivider())
        .append($type).append($createDivider())
        .append($creationTime);

    return $('<div>').addClass('note')
        .append($header)
        .append($content);
};

$(function() {
    $.get('/api/user', function(usersById) {
        USERS_BY_ID = usersById;

        $.get('/api/note', function(notes) {
            const $notes = $('#notes');

            for (const key of Object.keys(notes)) {
                $notes.append($createNote(key, notes[key]));
            }
        });
    });
});
