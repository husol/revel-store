$(document).ready(function() {
    $('.btn-add').on('click', function() {
        callAjax('users', 'load_form', {}, loadFormCallback);
    });
    $('.btn-edit').on('click', function() {
        var id = $(this).closest('tr').data('rel');
        callAjax('users', 'load_form', {id: id}, loadFormCallback);
    });
    $('.btn-delete').on('click', function() {
        var id = $(this).closest('tr').data('rel');
        var result = confirm("Bạn có chắc chắn muốn xóa User này?");
        if (result) {
            window.location.href = root_url + 'users/'+id+'/delete';
        }
        return false;
    });
});

function loadFormCallback(result) {
    if (result !== false) {
        $("#common_dialog").modal({show: true, keyboard: true, backdrop: 'static'});

        $('#submitUser').on('click', function () {
            //Validate
            //Check Required fields
            var requiredFields =[];
            var requiredNames = [];
            requiredFields.push('fullname'); requiredNames.push('Họ và tên');
            requiredFields.push('email'); requiredNames.push('Email');
            if ($('#idRecord').val() == 0) {
                requiredFields.push('password'); requiredNames.push('Mật khẩu');
            }
            requiredFields.push('mobile'); requiredNames.push('Điện thoại');
            if (!validateRequireds(requiredFields, requiredNames)) {
                return false;
            }

            //Check email fields
            var emailFields = [];
            var emailNames = [];
            emailFields.push('email'); emailNames.push('Email');
            if (!validateEmails(emailFields, emailNames)) {
                return false;
            }

            //Check logic
            if ($('#idRecord').val() == 0 && $('#password').val().length < 8) {
                showErrorBubble('password', "Mật khẩu phải có ít nhất 8 ký tự")
                return false;
            }

            if ($('#idRecord').val() == 0 && $('#password').val() != $('#confirmpassword').val()) {
                showErrorBubble('confirmpassword', "Nhập lại mật khẩu không trùng khớp")
                return false;
            }

            //Check if email existed
            callAjax('users', 'checkExistedEmail', {email:$('#email').val(), id_record: $('#idRecord').val()}, function(result){
                if (result) {
                    showErrorBubble("email", "Email đã tồn tại trong hệ thống");
                    return false;
                }

                $('#submitUser').closest('form').submit();
            });

            return false;
        })
    }
}