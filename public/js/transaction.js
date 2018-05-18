$(document).on('click', 'td .removeQty', function() {
    var id = $(this).closest('tr').data('rel');
    //Calculate in trolley
    var trolley = JSON.parse(sessionStorage.trolley);
    if (trolley[id].quantity > 1) {
        trolley[id].quantity--;
        trolley[id].amount = trolley[id].quantity * trolley[id].price;
    } else {
        delete trolley[id];
    }

    sessionStorage.trolley = JSON.stringify(trolley);
    loadTrolleyTable();
    updateTrolleyView();
});
$(document).on('click', 'td .addQty', function() {
    var id = $(this).closest('tr').data('rel');
    //Calculate in trolley
    var trolley = JSON.parse(sessionStorage.trolley);
    trolley[id].quantity++;
    trolley[id].amount = trolley[id].quantity * trolley[id].price;

    sessionStorage.trolley = JSON.stringify(trolley);
    loadTrolleyTable();
    updateTrolleyView();
});

$(document).ready(function() {
    loadTrolleyTable();
    $('#btnCheckout').on('click', function () {
        //Validate
        //Check Required fields
        var requiredFields =[];
        var requiredNames = [];
        requiredFields.push('client_name'); requiredNames.push('Họ và tên người nhận');
        requiredFields.push('client_email'); requiredNames.push('Email của bạn');
        requiredFields.push('client_mobile'); requiredNames.push('Số điện thoại người nhận');
        requiredFields.push('deliver_place'); requiredNames.push('Địa điểm giao hàng');
        if (!validateRequireds(requiredFields, requiredNames)) {
            return false;
        }

        //Check email fields
        var emailFields = [];
        var emailNames = [];
        emailFields.push('client_email'); emailNames.push('Email của bạn');
        if (!validateEmails(emailFields, emailNames)) {
            return false;
        }

        //Check if trolley empty
        if (sessionStorage.trolley == undefined || sessionStorage.trolley == "{}") {
            alert('Giỏ hàng của bạn đang trống.');
            return false;
        }
        var result = confirm("Cảm ơn quý khách đã chọn Mỹ Phẩm Tây Ninh. Tiếp tục đặt hàng?");
        if (result) {
            var contactInfo = {
                name: $("#client_name").val(),
                email: $("#client_email").val(),
                mobile: $("#client_mobile").val(),
            };

            var deliverPlace = $('#deliver_place').val(), note = $('#note').val();
            var trolleyData = sessionStorage.trolley;

            contactInfo = JSON.stringify(contactInfo);
            callAjax('transactions', 'checkout', {
                "g-recaptcha-response": $("#g-recaptcha-response").val(),
                "contactInfo":contactInfo,
                "deliverPlace":deliverPlace,
                "note":note,
                "trolley":trolleyData
            }, checkoutCallback);
        }
    });
});

function checkoutCallback(result) {
    if (result !== false) {
        showSuccessBubble('Đặt hàng thành công.<br>Chúng tôi sẽ xác nhận đơn hàng của quý khách tối đa 12 tiếng.');
        sessionStorage.trolley = '{}';
        window.location.href = '/';
    }
}

function loadTrolleyTable() {
    var trolleyContent = '';
    var trolley = JSON.parse(sessionStorage.trolley);
    for(var id in trolley) {
        var row = '<tr data-rel="' +id+ '"><td>' +trolley[id]["model"]+
            '</td><td>' +trolley[id]["name"]+
            '</td><td class="text-center"><span title="Bớt" class="removeQty cursor-pointer glyphicon glyphicon-minus text-red"></span>  '+
            trolley[id]["quantity"]+
            '  <span title="Thêm" class="addQty cursor-pointer glyphicon glyphicon-plus text-blue"></span></td><td class="text-right">'+
            addCommas(trolley[id]["price"].toString(), 0)+
            ' VND</td><td class="text-right">' +addCommas(trolley[id]["amount"].toString(), 0)+
            ' VND</td></tr>';
        trolleyContent += row;
    }

    if (trolleyContent == '') {
        trolleyContent = '<tr><td colspan="5" class="text-center">Không có sản phẩm nào trong Giỏ hàng</td></tr>';
    }
    $('#tableTrolley > tbody').html(trolleyContent);
}