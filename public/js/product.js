$(document).ready(function() {
    $('.btn-add').on('click', function() {
        callAjax('products', 'load_form', {}, loadFormCallback);
    });
    $('.btn-edit').on('click', function() {
        var id = $(this).closest('tr').data('rel');
        callAjax('products', 'load_form', {id: id}, loadFormCallback);
    });
    $('.btn-delete').on('click', function() {
        var id = $(this).closest('tr').data('rel');
        var result = confirm("Bạn có chắc chắn muốn xóa Sản phẩm này?");
        if (result) {
            window.location.href = root_url + 'products/'+id+'/delete';
        }
        return false;
    });
    $('#btnComment').on('click', function() {
        var id = $('#contentComment').data('id'),
            id_product = $('#contentComment').data('idproduct'),
            content = $('#contentComment').val();

        callAjax("products", "update_comment", {id:id, id_product:id_product, content:content}, updateCommentCallback);
    });
    $('#addTrolley').on('click', function() {
        var data = $(this).data("rel").split(",");
        var id = data[0], model = data[1], name = data[2], price = data[3];
        var proObj = {"model": model, "name": name, "price": price};
        if (sessionStorage.trolley == undefined) {
            sessionStorage.trolley = JSON.stringify({});
        }

        //Calculate in trolley
        var trolley = JSON.parse(sessionStorage.trolley);
        if (trolley[id] == undefined) {
            proObj.quantity = 1;
        } else {
            proObj.quantity = trolley[id].quantity + 1;
        }
        proObj.amount = proObj.quantity * price;
        trolley[id] = proObj;
        sessionStorage.trolley = JSON.stringify(trolley);

        //Update trolley view
        updateTrolleyView();
        showSuccessBubble('Thêm vào Giỏ hàng thành công.');
    });
});

$(document).on('click', '.btn-editComment', function() {
    var id = $(this).parent().data("id");
    var content = $(this).parent().next("p").html().replace(/<br *\/?>/gi, '\n');
    $('#contentComment').data('id', id);
    $('#contentComment').val(content);
    $('#contentComment').focus();
    return false;
});
$(document).on('click', '.btn-deleteComment', function() {
    var id = $(this).parent().data("id");
    var result = confirm("Bạn có chắc chắn muốn xóa Bình luận này?");
    if (result) {
        callAjax("products", "delete_comment", {id:id}, deleteCommentCallback);
    }
    return false;
});

function loadFormCallback(result) {
    if (result !== false) {
        $("#common_dialog").modal({show: true, keyboard: true, backdrop: 'static'});

        CKEDITOR.replace('description');

        $('select[multiple="multiple"]').multiselect(
            {includeSelectAllOption: true}
        );

        $("#img_cover").on("click", function () {
            $("#cover").trigger("click");
        });

        $("#cover").on('change', function(){
            previewImage(this, "img_cover");
        });

        $('#submitProduct').on('click', function () {
            //Validate
            //Check Required fields
            var requiredFields =[];
            var requiredNames = [];
            requiredFields.push('category'); requiredNames.push('Loại sản phẩm');
            requiredFields.push('model_name'); requiredNames.push('Mã sản phẩm');
            requiredFields.push('name'); requiredNames.push('Tên sản phẩm');
            if ($('#cover').data('link') == '') {
                requiredFields.push('cover'); requiredNames.push('Hình sản phẩm');
            }
            requiredFields.push('short_description'); requiredNames.push('Mô tả ngắn về sản phẩm');
            requiredFields.push('price'); requiredNames.push('Giá sản phẩm');
            if (!validateRequireds(requiredFields, requiredNames)) {
                return false;
            }

            //Validate File Type
            if ($('#cover').val().length > 0 && !isImageFile($('#cover').val())) {
                showErrorBubble('cover', "Sai định dạng hình ảnh sản phẩm.");
                return false;
            }

            if (CKEDITOR.instances.description.document.getBody().getChild(0).getText().trim() == "") {
                showErrorBubble("", "Vui lòng nhập Mô tả sản phẩm");
                return false;
            }

            $(this).closest('form').submit();
            return false;
        })
    }
}

function updateCommentCallback(result) {
    if (result !== false) {
        showSuccessBubble("Bình luận đã được cập nhật thành công.");
        $('#contentComment').val('');
    }
}

function deleteCommentCallback(result) {
    if (result !== false) {
        $('#comments > div[data-id='+result.Id+']').next("p").remove();
        $('#comments > div[data-id='+result.Id+']').remove();
        showSuccessBubble("Đã xóa bình luận thành công.");
    }
}