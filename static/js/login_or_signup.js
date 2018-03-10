var getInputField = function($form, field) {
    return $form.find('[name="' + field + '"]');
};

var getInputFields = function($form, fields) {
    return fields.map(
        (field) => getInputField($form, field));
};

var checkFormValidity = function($form, fields) {
    return getInputFields($form, fields).every(
        ($field) => $field.get(0).checkValidity());
};

var getFormData = function($form, fields) {
    return getInputFields($form, fields).reduce(
        (formData, $field) => {
            formData[$field.attr('name')] = $field.val();
            return formData;
        },
        {});
};

var populateValidationMessage = function($field) {
    var $validationMessage = $field.siblings('.validation-message').first();
    var validationMessage = $field.get(0).validationMessage;

    $validationMessage.text(validationMessage);
};

var populateValidationMessages = function($form, fields) {
    getInputFields($form, fields).forEach(
        ($field) => populateValidationMessage($field));
};

var touchAllFields = function($form, fields) {
    getInputFields($form, fields).forEach(
        ($field) => {
            $field.focus().blur();
        });
};

var checkKeypressIsSpace = function(event) {
    return event.which === 32;
};

$(function() {
    var $signupForm = $('#signup-form');
    var $loginForm = $('#login-form');

    var submitHasBeenClicked = false;

    var displayNameField = 'displayName';
    var emailAddressField = 'emailAddress';
    var passwordField = 'password';

    var signupFields = [
        displayNameField,
        emailAddressField,
        passwordField,
    ];

    var loginFields = [
        emailAddressField,
        passwordField,
    ];


    var installFieldValidators = function($form, fields) {
        getInputFields($form, fields).forEach(
          ($field) => {
            var field = $field.attr('name');

            // continuously update validation message after failed submission
            $field.on(
                'input',
                (event) => {
                    if (submitHasBeenClicked) {
                        populateValidationMessage($field);
                    }
                });

            if (field !== passwordField) {
                // restrict initial space character
                $field.keypress(
                    (event) => {
                        var value = $field.val();
                        var trimmedValue = $.trim(value);

                        if (checkKeypressIsSpace(event) && trimmedValue.length === 0) {
                            return false; // cancels keypress event
                        }
                    });

                // remove trailing spaces on blur
                $field.blur(
                    () => {
                        var value = $field.val();
                        var trimmedValue = $.trim(value);

                        $field.val(trimmedValue);
                    });
            }
        });
    };

    installFieldValidators($signupForm, signupFields);
    installFieldValidators($loginForm, loginFields);

    $signupForm.find('button').click(
        () => {
            if (checkFormValidity($signupForm, signupFields)) {
                var formData = getFormData($signupForm, signupFields);
                var jsonData = JSON.stringify(formData);

                $.post(
                    '/user',
                    jsonData,
                    (userId) => {
                        alert('Created User with id: ' + userId.value);
                    },
                    'json');

            } else if (!submitHasBeenClicked) {
                submitHasBeenClicked = true;

                populateValidationMessages($signupForm, signupFields);
                touchAllFields($signupForm, signupFields);
            }
        });

    $loginForm.find('button').click(function() {
    if (checkFormValidity($loginForm, loginFields)) {
            var formData = getFormData($loginForm, loginFields);
            var jsonData = JSON.stringify(formData);

            $.post('/session', jsonData, userId => {
                alert(userId);
            }, 'text');

        } else if (!submitHasBeenClicked) {
            submitHasBeenClicked = true;

            populateValidationMessages($loginForm, loginFields);
            touchAllFields($loginForm, loginFields);
        }
    });
    $loginForm.find('button').click(
        () => {
            if (checkFormValidity($loginForm, loginFields)) {
                var formData = getFormData($loginForm, loginFields);
                var jsonData = JSON.stringify(formData);
                
                $.post(
                  '/session', 
                  jsonData, 
                  (response) => {
                    alert(response);
                   }, 
                  'text');

            } else if (!submitHasBeenClicked) {
                submitHasBeenClicked = true;

                populateValidationMessages($loginForm, loginFields);
                touchAllFields($loginForm, loginFields);
            }
        });
});
