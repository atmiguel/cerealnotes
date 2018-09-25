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

var checkKeypressIsEnter = function(event) {
    return event.which === 13;
}

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

    var attachFieldValidators = function(formMetadata) {
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

    attachFieldValidators(signupFormMetadata);
    attachFieldValidators(loginFormMetadata);

    var attachSubmitClickHandler = function(formMetadata, postFunction) {
        var $button = formMetadata.$form.find('button');
        $button.click(() => {
            if (checkFormValidity(formMetadata.$form, formMetadata.fields)) {
                var formData = getFormData(
                    formMetadata.$form,
                    formMetadata.fields);

                var formDataAsJsonString = JSON.stringify(formData);

                postFunction(formDataAsJsonString);

            } else if (!formMetadata.submitHasBeenClicked) {
                formMetadata.submitHasBeenClicked = true;

                populateValidationMessages(
                    formMetadata.$form,
                    formMetadata.fields);

                touchAllFields(formMetadata.$form, formMetadata.fields);
            }
        });

        formMetadata.$form.find('input:password').keydown((event) => {
            if (checkKeypressIsEnter(event)) {
                $button.click();
            }
        });
    };

    attachSubmitClickHandler(signupFormMetadata, (formDataAsJsonString) => {
        $.post('/api/user', formDataAsJsonString, (responseBody, _, $XmlHttpResponse) => {
            if ($XmlHttpResponse.status === 201) {
                mui.tabs.activate('login-form');
                alert('Successfully created user, please sign in');
            } else {
                alert('Unexpected successful status');
            }
        }).fail(($XmlHttpResponse) => {
            if ($XmlHttpResponse.status === 409) {
                alert('Email address already in use');
            } else {
                alert('Unexpected error ' + $XmlHttpResponse.responseText);
            }
        });
    });

    attachSubmitClickHandler(loginFormMetadata, (formDataAsJsonString) => {
        $.post('/api/session', formDataAsJsonString, (responseBody, _, $XmlHttpResponse) => {
            if ($XmlHttpResponse.status === 201) {
                location.reload();
            } else {
                alert('Error in logging in');
            }
        }, 'text').fail(($XmlHttpResponse) => {
            if ($XmlHttpResponse.status === 401) {
                alert('Email address and/or password was incorrect');
            } else {
                alert('Unexpected error ' + $XmlHttpResponse.responseText);
            }
        });
    });
});
