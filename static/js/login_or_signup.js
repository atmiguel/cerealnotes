var getInputField = function($form, field) {
    return $form.find('[name="' + field + '"]');
};

var getInputFields = function($form, fields) {
    return fields.map((field) => getInputField($form, field));
};

var checkFormValidity = function($form, fields) {
    var $inputFields = getInputFields($form, fields);

    return $inputFields.every(($field) => $field.get(0).checkValidity());
};

var getFormData = function($form, fields) {
    var $inputFields = getInputFields($form, fields);

    return $inputFields.reduce((formData, $field) => {
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
    var $inputFields = getInputFields($form, fields);
    $inputFields.forEach(($field) => populateValidationMessage($field));
};

var touchAllFields = function($form, fields) {
    var $inputFields = getInputFields($form, fields);
    $inputFields.forEach(($field) => $field.focus().blur());
};

var checkKeypressIsSpace = function(event) {
    return event.which === 32;
};

$(function() {
    var displayNameField = 'displayName';
    var emailAddressField = 'emailAddress';
    var passwordField = 'password';

    var signupFormMetadata = {
        $form: $('#signup-form'),
        fields: [displayNameField, emailAddressField, passwordField],
        submitHasBeenClicked: false,
    };

    var loginFormMetadata = {
        $form: $('#login-form'),
        fields: [emailAddressField, passwordField],
        submitHasBeenClicked: false,
    };

    var installFieldValidators = function(formMetadata) {
        var $inputFields = getInputFields(
            formMetadata.$form,
            formMetadata.fields);

        $inputFields.forEach(($field) => {
            var field = $field.attr('name');

            // continuously update validation message after failed submission
            $field.on('input', () => {
                if (formMetadata.submitHasBeenClicked) {
                    populateValidationMessage($field);
                }
            });

            if (field !== passwordField) {
                // restrict initial space character
                $field.keypress((event) => {
                    var fieldHasValue = $.trim($field.val()).length > 0;

                    if (checkKeypressIsSpace(event) && !fieldHasValue) {
                        return false; // cancels keypress event
                    }
                });

                // remove trailing spaces on blur
                $field.blur(() => {
                    var value = $field.val();
                    var trimmedValue = $.trim(value);

                    $field.val(trimmedValue);
                });
            }
        });
    };

    installFieldValidators(signupFormMetadata);
    installFieldValidators(loginFormMetadata);

    var installSubmitClickHandler = function(formMetadata, postFunction) {
        formMetadata.$form.find('button').click(() => {
            if (checkFormValidity(formMetadata.$form, formMetadata.fields)) {
                var formData = getFormData(
                    formMetadata.$form,
                    formMetadata.fields);

                var formDataAsJsonString = JSON.stringify(formData);

                postFunction(formDataAsJsonString)

            } else if (!formMetadata.submitHasBeenClicked) {
                formMetadata.submitHasBeenClicked = true;

                populateValidationMessages(
                    formMetadata.$form,
                    formMetadata.fields);

                touchAllFields(formMetadata.$form, formMetadata.fields);
            }
        });
    }

    installSubmitClickHandler(signupFormMetadata, (formDataAsJsonString) => {
        $.post('/user', formDataAsJsonString, (responseBody, _, request) => {
            if (request.status === 201) {
                alert('Successfully created user');
            } else {
                // TODO flesh out if statement possibilities
                alert('Email address already in use');
            }
        }, 'json');
    });

    installSubmitClickHandler(loginFormMetadata, (formDataAsJsonString) => {
        $.post('/session', formDataAsJsonString, (response) => {
            alert(response);
        }, 'text');
    });
});
