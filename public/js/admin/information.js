$(document).ready(function() {
    $('.btn-add').on('click', function() {
        callAjax('information', 'load_form', {}, loadFormCallback);
    });
    $('.btn-edit').on('click', function() {
        var id = $(this).closest('tr').data('rel');
        callAjax('information', 'load_form', {id: id}, loadFormCallback);
    });
    $('.btn-delete').on('click', function() {
        var id = $(this).closest('tr').data('rel');
        var result = confirm("Bạn có chắc chắn muốn xóa Thông tin này?");
        if (result) {
            window.location.href = root_url + 'information/'+id+'/delete';
        }
        return false;
    });
});

function loadFormCallback(result) {
    if (result !== false) {
        $("#common_dialog").modal({show: true, keyboard: true, backdrop: 'static'});

        CKEDITOR.replace('description');

        $("#img_cover").on("click", function () {
            $("#cover").trigger("click");
        });

        $("#cover").change(function(){
            previewImage(this, "img_cover");
        });

        $('#submitCategory').on('click', function () {
            //Validate
            //Check Required fields
            var requiredFields =[];
            var requiredNames = [];
            requiredFields.push('title'); requiredNames.push('Tiêu đề');
            if (!validateRequireds(requiredFields, requiredNames)) {
                return false;
            }

            //Validate File Type
            if ($('#cover').val().length > 0 && !isImageFile($('#cover').val())) {
                showErrorBubble('cover', "Sai định dạng hình ảnh Cover.");
                return false;
            }

            if (CKEDITOR.instances.description.document.getBody().getChild(0).getText().trim() == "") {
                showErrorBubble("", "Vui lòng nhập Nội dung");
                return false;
            }

            $(this).closest('form').submit();
            return false;
        })
    }
}