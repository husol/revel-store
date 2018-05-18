$(document).ready(function() {
    var totalPage = $("#tableComment").data('totalpage');
    if (totalPage > 0) {
        $('#paginationComment').twbsPagination({
            totalPages: totalPage,
            visiblePages: 7,
            onPageClick: function (event, page) {
                callAjax("comments", "paging", {page:page}, displayListCallback);
                return false;
            }
        });
    } else {
        callAjax("comments", "paging", {page:0}, displayListCallback);
        return false;
    }
});

function displayListCallback(result) {
    if (result !== false) {
        $('.btn-approve').on('click', function() {
            var id = $(this).closest('tr').data('rel');
            window.location.href = root_url + 'comments/'+id+'/approve';
        });
        $('.btn-delete').on('click', function() {
            var id = $(this).closest('tr').data('rel');
            var result = confirm("Bạn có chắc chắn muốn xóa Bình luận này?");
            if (result) {
                window.location.href = root_url + 'comments/'+id+'/delete';
            }
            return false;
        });
    }
}