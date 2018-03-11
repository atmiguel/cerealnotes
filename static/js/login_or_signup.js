/*jshint esversion: 6 */
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
    var displayNameField = 'displayName';
    var emailAddressField = 'emailAddress';
    var passwordField = 'password';

    var signupFieldsMetaData = {
        $form: $('#signup-form'),
        fields: [
            displayNameField,
            emailAddressField,
            passwordField,
        ],
        submitHasBeenClicked: false
    };

    var loginFieldsMetaData = {
        $form: $('#login-form'),
        fields: [
            emailAddressField,
            passwordField,
        ],
        submitHasBeenClicked: false
    };


    var installFieldValidators = function(fieldsMetaData) {
        getInputFields(fieldsMetaData.$form, fieldsMetaData.fields).forEach(
          ($field) => {
            var field = $field.attr('name');

            // continuously update validation message after failed submission
            $field.on(
                'input',
                (event) => {
                    if (fieldsMetaData.submitHasBeenClicked) {
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

    installFieldValidators(signupFieldsMetaData);
    installFieldValidators(loginFieldsMetaData);


    var installSubmitClickHandler = function(fieldsMetaData, postFunction) {
        fieldsMetaData.$form.find('button').click(
        () => {
            if (checkFormValidity(fieldsMetaData.$form, fieldsMetaData.fields)) {
                var formData = getFormData(fieldsMetaData.$form, fieldsMetaData.fields);
                var jsonData = JSON.stringify(formData);

                postFunction(jsonData)

            } else if (!fieldsMetaData.submitHasBeenClicked) {
                fieldsMetaData.submitHasBeenClicked = true;

                populateValidationMessages(fieldsMetaData.$form, fieldsMetaData.fields);
                touchAllFields(fieldsMetaData.$form, fieldsMetaData.fields);
            }
        });
    }

    installSubmitClickHandler(signupFieldsMetaData, 
        (formDataAsJsonString) => {
            $.post(
                '/user',
                formDataAsJsonString,
                (userId) => {
                    alert('Created User with id: ' + userId.value);
                },
                'json');
        });

    installSubmitClickHandler(loginFieldsMetaData, 
        (formDataAsJsonString) => {
            $.post(
                '/session', 
                formDataAsJsonString, 
                (response) => {
                alert(response);
                }, 
                'text');
        });
});
