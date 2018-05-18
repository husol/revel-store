$(document).ready(function() {
    $('#tableNewTransaction .table > tbody > tr').on('click', function() {
        var id = $(this).data("rel");
        if ($(this).hasClass('selected') || id == undefined) {
            return false;
        }
        $('#tableNewTransaction .table > tbody > tr').removeClass('selected');
        //Load transaction detail
        $(this).addClass("selected");
        callAjax("transactions", "detail", {id:id});
        return false;
    });

    $('.actionTransaction').on('change', function() {
        var id = $(this).closest('tr').data('rel');
        var status = $(this).val();
        window.location.href = root_url + 'home/'+id+'/update/'+status;
    });

    $('#sendmail').on('click', function() {
        //Validate
        //Check Required fields
        var requiredFields =[];
        var requiredNames = [];
        requiredFields.push('subject'); requiredNames.push('Tiêu đề');
        requiredFields.push('content'); requiredNames.push('Nội dung');
        if (!validateRequireds(requiredFields, requiredNames)) {
            return false;
        }

        var pass = false;
        if ($('#email').val().trim() != '') {
            //Check email fields
            var emailFields = [];
            var emailNames = [];
            emailFields.push('email'); emailNames.push('Email');
            if (!validateEmails(emailFields, emailNames)) {
                return false;
            }
            pass = true;
        } else {
            pass = confirm("Bạn có chắc chắn muốn gửi Email Thông báo đến tất cả khách hàng?");
        }

        if (pass) {
            $(this).closest('form').submit();
        }
        return false;
    });
});