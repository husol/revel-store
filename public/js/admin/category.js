$(document).ready(function() {
    $('.btn-add').on('click', function() {
        callAjax('categories', 'load_form', {}, loadFormCallback);
    });
    $('.btn-edit').on('click', function() {
        var id = $(this).closest('tr').data('rel');
        callAjax('categories', 'load_form', {id: id}, loadFormCallback);
    });
    $('.btn-delete').on('click', function() {
        var id = $(this).closest('tr').data('rel');
        var result = confirm("Bạn có chắc chắn muốn xóa Loại sản phẩm này?");
        if (result) {
            window.location.href = root_url + 'categories/'+id+'/delete';
        }
        return false;
    });
});

function loadFormCallback(result) {
    if (result !== false) {
        $("#common_dialog").modal({show: true, keyboard: true, backdrop: 'static'});

        $('#submitCategory').on('click', function () {
            //Validate
            //Check Required fields
            var requiredFields =[];
            var requiredNames = [];
            requiredFields.push('name'); requiredNames.push('Tên loại sản phẩm');
            if (!validateRequireds(requiredFields, requiredNames)) {
                return false;
            }

            $(this).closest('form').submit();
            return false;
        })
    }
}