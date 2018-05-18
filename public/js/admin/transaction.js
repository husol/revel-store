$(document).ready(function() {
    var totalPage = $("#tableTransaction").data('totalpage');
    if (totalPage > 0) {
        $('#paginationTransaction').twbsPagination({
            totalPages: totalPage,
            visiblePages: 7,
            onPageClick: function (event, page) {
                callAjax("transactions", "paging", {page:page}, displayListCallback);
                return false;
            }
        });
    } else {
        callAjax("transactions", "paging", {page:0}, displayListCallback);
        return false;
    }
});

$(document).on('click', '#tableTransaction .table > tbody > tr:not(".selected")', function() {
    var id = $(this).data("rel");
    if ($(this).hasClass('selected') || id == undefined) {
        return false;
    }
    $('#tableTransaction .table > tbody > tr').removeClass('selected');
    //Load transaction detail
    $(this).addClass("selected");
    callAjax("transactions", "detail", {id:id});
    return false;
});

function displayListCallback(result) {
    if (result !== false) {
        $('.actionTransaction').on('change', function() {
            var id = $(this).closest('tr').data('rel');
            var status = $(this).val();
            window.location.href = root_url + 'transactions/'+id+'/update/'+status;
        });
    }
}