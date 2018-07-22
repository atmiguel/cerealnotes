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

const $createNote = function(note) {
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

const $createButtonWithText = function(text) {
    return $('<button>').addClass('mui-btn').text(text);
};

const $createHalfRowDivWithElement = function($element) {
    return $('<div>').addClass('mui-col-xs-6').append($element);
}

const $createRowWithTwoElements = function($element1, $element2) {
    const $row = $('<div>').addClass('mui-row');

    const $column1 = $createHalfRowDivWithElement($element1);
    const $column2 = $createHalfRowDivWithElement($element2);

    return $row
        .append($column1)
        .append($column2);
};

const $createGridContainer = function() {
    return $('<div>').addClass('mui-container-fluid');
};

const $createAddNoteModal = function() {
    const $modal = $('<div>').addClass('modal');

    const noteTypes = [
        'Marginalia',
        'Meta',
        'Prediction',
        'Question',
    ];

    const $buttons = noteTypes.map(noteType => {
        return $createButtonWithText(noteType).addClass('note-type-button');
    });

    return $modal
        .append($createGridContainer()
            .append($createRowWithTwoElements($buttons[0], $buttons[1]))
            .append($createRowWithTwoElements($buttons[2], $buttons[3])));
};

const activateModal = function($modal) {
    mui.overlay('on', $modal.get(0));
}

$(function() {
    const $addNoteModal = $createAddNoteModal();

    $.get('/api/user', function(usersById) {
        USERS_BY_ID = usersById;

        $.get('/api/note', function(notes) {
            const $notes = $('#notes');

            notes.forEach((note) => {
                $notes.append(
                    $createNote(note)
                );
            });
        });
    });

    $('#add-note-button').click(function() {
        activateModal($addNoteModal);
    });
});
