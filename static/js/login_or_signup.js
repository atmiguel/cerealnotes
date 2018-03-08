var getInputFields = function($form, fields) {
    return fields.map(
        field => $form.find('[name="' + field + '"]'));
};

var checkFormValidity = function($form, fields) {
    return getInputFields($form, fields).every(
        $field => $field.get(0).checkValidity());
};

var getFormData = function($form, fields) {
    return getInputFields($form, fields).reduce(
        (formData, $field) => {
            formData[$field.attr('name')] = $field.val();
            return formData;
        }, {});
};

var populateValidationMessage = function($field) {
    var $validationMessage = $field.siblings('.validation-message').first();
    var validationMessage = $field.get(0).validationMessage;

    $validationMessage.text(validationMessage);
};

var populateValidationMessages = function($form, fields) {
    getInputFields($form, fields).forEach(
        $field => populateValidationMessage($field));
};

var touchAllFields = function($form, fields) {
    getInputFields($form, fields).forEach($field => {
        $field.focus().blur();
    });
};

$(function() {
    // Signup Form
    var $signupForm = $('#signup-form');
    var submitHasBeenClicked = false;

    var fields = [
        'displayName',
        'emailAddress',
        'password',
    ];

    getInputFields($signupForm, fields).forEach($field => {
        $field.on('input', () => {
            if (submitHasBeenClicked) {
                populateValidationMessage($field);
            }
        });
    });

    $signupForm.find('button').click(function() {
        submitHasBeenClicked = true;

	if (checkFormValidity($signupForm, fields)) {
            $.post(
                '/signup',
                getFormData($signupForm, fields),
                () => {
                    console.log('done');
                });
        } else {
            populateValidationMessages($signupForm, fields);
            touchAllFields($signupForm, fields);
        }
    });
});
