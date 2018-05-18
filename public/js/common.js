/////////////////// Check file type //////////////////////////////////////////
/**
* file : int
* ext : string array()
* Note: Extention of file with dot in front of ext.
* Ex: ext = new Array(".png", ".jpeg", ".jpg")
*/
function checkFileExt(file, ext) {
		var pathLength = file.length;
		var lastDot = file.lastIndexOf(".");
		var fileType = file.substring(lastDot, pathLength);

		for (i = 0; i < ext.length; i++) {
			if (fileType == ext[i])
			{
				return true;
			}
		}
		return false;
}
/////////////////// Sort alphabet option select list /////////////////////////
function sortAlphabet(id) {
    var prePrepend = "#";
    if (id.match("^#") == "#") prePrepend = "";
    selectedValue = $('#' + id).val();
    $(prePrepend + id).html($(prePrepend + id + " option").sort(
        function (a, b) {   
            if (a.value <= 0){
                //alert(a.value);
                return -1;
            }
            else if (b.value <= 0){
                return 1;
            }
            else 
                return a.text.toUpperCase() == b.text.toUpperCase() ? 0 : (a.text.toUpperCase() < b.text.toUpperCase() ? -1 : 1);
        }
    ));
              
    $('#' + id).val(selectedValue);
}

//Add commas to number string
function addCommas(str, addDecimal) {
    var arr,
        int,
        dec;
    str = str.replace(/,/g, '');

    str += '';
    if (str == '') {
        return '';
    }
    arr = str.split('.');

    if (str.slice(-1) == '%') {
        str = str.replace(/%/g, '');
        if (isNaN(str)) {
            return '0.00';
        }

        var decimal = 2;
        if (arr.length > 1) {
            decimal = decimal + arr[1].length - 1;
        }
        str = parseFloat(str) / 100;
        str = str.toFixed(decimal).toString();
    }
    if (isNaN(str)) {
        return '';
    }

    arr = str.split('.');
    int = parseInt(arr[0]) + '';
    if (arr[0].length == 0) {
        int = 0;
    }

    if (isNaN(int) || int == 'NaN') {
        return '';
    }

    dec = arr.length > 1 ? arr[1] : '';
    dec = dec.length > 1 ? '.' + dec : dec.length == 1 ? '.' + dec + 0 : '.00';

    if (addDecimal) {
        return int.toString().replace(/(\d)(?=(\d\d\d)+(?!\d))/g, "$1,") + dec;
    } else {
        return int.toString().replace(/(\d)(?=(\d\d\d)+(?!\d))/g, "$1,");
    }
}

function stringToNumber(str) {
    if (typeof str === 'string') {
        return isNaN(parseFloat(str.replace(/,/g, ''))) ? 0 : parseFloat(str.replace(/,/g, ''));
    }
    return parseFloat(str);
}
/////////////////////// callAjax jQuery //////////////////////////////////////
/* Usage:

var dataSend = {id:'testID', loading:'loadingID', param_1:'paramValue_1', param_n:'paramValue_n'};
callAjax('Controller_Name', 'Action_Name', dataSend, testCallback);

function testCallback(result) {
	if (result !== false) {
        //Code here
	}
}

in which:
testID : the id of element we want to work.
loadingID : the id of loading icon we want to display.
paramValue_n: all parameters we want to post to server.
*/
function base64_decode(data) {
    return decodeURIComponent(atob(data).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
}

function fetchForm(formObj, upload) {
    if (upload == undefined) {
        upload = false;
    }

    var formData = {};
    if (upload === true) {
        formData = new FormData(formObj[0]);
    } else {
        var form = formObj.serializeArray();
        $.each(form, function() {
            formData[this.name] = this.value || '';
        });
    }

    return formData;
}

function callAjax(controller, method, args, callback) {
    var isLoading = false,
        dataType = 'json',
        isAsync = true,
        contentType = 'application/x-www-form-urlencoded; charset=UTF-8',
        processData = true;

    if (args != null && args['silent'] == null) {
        isLoading = true;
        $("#loading").show();
    }

    if (args != null && args['dataType'] != null) {
        dataType = args['dataType'];
    }

    if (args != null && args['isAsync'] != null) {
        isAsync = args['isAsync'];
    }

    var data = args;
    if (args != null && args['formData'] != null) {
        if (args['upload'] != null) {
            contentType = false;
            processData = false;
        }
        data = args['formData'];
    }

    var url = root_url + Array(controller, method).join("/");
    var urlArr = controller.split(":/");
    if (urlArr[0] == 'https' || urlArr[0] == 'http') {
        url = Array(controller, method).join("/");
    }

    objectCall = $.ajax({
        type: "POST",
        timeout: 99999999999999,
        url: url,
        data: data,
        dataType: dataType,
        async : isAsync,
        contentType: contentType,
        processData: processData,
        success: function(msg) {
            if (args['is_append'] == null) {
                $("#" + args['id']).html('');
            }
            $("#loading").hide();

            if (msg == null) {
                if (isLoading) {
                    $("#loading").hide();
                    if ($("#" + args['id']).length > 0) {
                        $("#" + args['id']).hide();
                    }
                }
                return;
            } else if (msg == "expired_session") {
                alert("Session expired.");
                location.reload();
                return;
            }

            if (msg.result == false && msg.messages != '') {
                showErrorBubble("", msg.messages);
            }

            if (typeof msg.html != 'undefined' && msg.html != null && msg.html.id != '') {
                $("#" + msg.html.id).html(base64_decode(msg.html.content));
            }

            if (callback != null) {
                if (msg.result == undefined) {
                    callback(msg, args['id']);
                } else {
                    callback(msg.result, args['id']);
                }
            }
        },
        error: function(request, status, error) {
            if (request.statusText == 'abort' || request.statusText == 'error') return;
            if (error.name == 'NS_ERROR_NOT_AVAILABLE' || request.readyState == 0) {
                if (args['silent'] == null){
                    $('#infoText').html("Request is interrupted unexpectedly");
                    setTimeout(function () {
                        $('#loading').show();
                    }, 3000);

                }
            } else {
                //default is display html return
                $('#' + args['id']).html(request.responseText);
            }
        }
    });
    return objectCall;
}
/////////////////////// End callAjax jQuery //////////////////////////////////////
function $_GET(param) {
    var vars = {};
    window.location.href.replace(location.hash, '').replace(
        /[?&]+([^=&]+)=?([^&]*)?/gi, // regexp
        function( m, key, value ) { // callback
            vars[key] = value !== undefined ? value : '';
        }
    );

    if ( param ) {
        return vars[param] ? vars[param] : null;
    }
    return vars;
}
// Show tooltip when input data is not valid
function showErrorBubble(control, error_msg, seconds) {
    var ctrl = false;
    if ($(control).length > 0) {
        ctrl = $(control);
    }
    var delay = seconds || 5000;
    jQuery.showMessage({
        thisMessage:	    [error_msg],
        className:		    'fail',
        position:		    'top',
        opacity:		    90,
        displayNavigation:	true,
        autoClose:		    true,
        delayTime:		    delay
    });

    if (ctrl !== false) {
        ctrl.focus();
    }

    return false;
}
// Show tooltip when input data is valid
function showSuccessBubble(success_msg, seconds) {
    var delay = seconds || 5000;
    jQuery.showMessage({
        thisMessage:	    [success_msg],
        className:		    'success',
        position:		    'top',
        opacity:		    90,
        displayNavigation:	true,
        autoClose:		    true,
        delayTime:		    delay
    });
    return false;
};

////////////////////// Validation Functions /////////////////////////////////////
function getExtension(filename) {
    var parts = filename.split('.');
    return parts[parts.length - 1];
}

function isImageFile(filename) {
    var ext = getExtension(filename);
    switch (ext.toLowerCase()) {
        case 'jpg':
        case 'gif':
        case 'bmp':
        case 'png':
            return true;
    }
    return false;
}

function isApplicationFile(filename) {
    var ext = getExtension(filename);
    switch (ext.toLowerCase()) {
        case 'pdf':
        case 'xlsx':
        case 'docx':
        case 'xls':
        case 'doc':
            return true;
    }
    return false;
}

function isDate(date) {
    var t = date.match(/^(\d{2})\/(\d{2})\/(\d{4})$/);
    if(t === null)
        return false;
    var d = +t[1], m = +t[2], y = +t[3];

    // Below should be a more acurate algorithm
    if(m >= 1 && m <= 12 && d >= 1 && d <= 31) {
        return true;
    }

    return false;
}

function isEmail(email) {
    var re = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(email);
}

function validateDates(fields, names) {
    var result = true;
    $.each(fields, function(index, field) {
        var obj = $('#'+field);
        var objValStr = String(obj.val());

        if (!isDate(objValStr)) {
            showErrorBubble('#'+field, "Sai định dạng " + names[index]);
            result = false;
            return false;
        }
    });
    return result;
}

function validateRequireds(fields, names) {
    var result = true;
    $.each(fields, function(index, field) {
        var obj = $('#'+field);
        var objValStr = String(obj.val());
        if (objValStr.trim() == "") {
            showErrorBubble('#'+field, "Vui lòng nhập " + names[index]);
            result = false;
            return false;
        }
    });
    return result;
}

function validateEmails(fields, names) {
    var result = true;
    $.each(fields, function(index, field) {
        var obj = $('#'+field);
        var objValStr = String(obj.val());
        var emails = [objValStr];
        if (objValStr.indexOf(',') > -1) {
            emails = objValStr.split(',');
        }

        var check = true;
        $.each(emails, function(i, email) {
            if (!isEmail(email.trim())) {
                showErrorBubble('#'+field, "Sai định dạng " + names[index]);
                check = false;
                result = false;
                return false;
            }
        });

        return check;
    });
    return result;
}
/////////////////////// End Validation Functions //////////////////////////////////////////////////////

function saveFile(url) {
    // Get file name from url.
    var filename = url.substring(url.lastIndexOf("/") + 1).split("?")[0];
    var xhr = new XMLHttpRequest();
    xhr.responseType = 'blob';
    xhr.onload = function() {
        var a = document.createElement('a');
        a.href = window.URL.createObjectURL(xhr.response); // xhr.response is a blob
        a.download = filename; // Set the file name.
        a.style.display = 'none';
        document.body.appendChild(a);
        a.click();
        delete a;
    };
    xhr.open('GET', url);
    xhr.send();
}

// Set cookie
function setCookie(c_name,value,exdays)
{
    var exdate = new Date();
    exdate.setDate(exdate.getDate() + exdays);
    var c_value = encodeURIComponent(value) + ((exdays==null) ? "" : "; expires="+exdate.toUTCString());
    document.cookie = c_name + "=" + c_value;
}

/////////////////////// auto load Data per x seconds ajax jQuery //////////////////////////////////////
/* Usage:
// On page html:
$(document).ready(function(){
			autoAjaxCall("input_url","input_html_object",time_second_unit);
		});
// On page input_url php:
echo "<p>".$row['field']."</p>"
*/
function autoAjaxCall(url, HTMLObject, jumpTime) {
	var callAjax = function(){
		$.ajax({
		  method:'POST',
		  url:url,
		  success:function(data){
			$(HTMLObject).html(data);
		  }
		});
	}
	setInterval(callAjax,jumpTime*1000);
}

function previewImage(input, imgId) {
    if (input.files && input.files[0]) {
        var reader = new FileReader();

        reader.onload = function (e) {
            $('#'+imgId).attr('src', e.target.result);
        }

        reader.readAsDataURL(input.files[0]);
    }
}

function updateTrolleyView() {
    if (sessionStorage.trolley == undefined) {
        return false;
    }
    var total = quantity = 0;
    var trolleyContent = '';
    var trolley = JSON.parse(sessionStorage.trolley);
    for(var id in trolley) {
        quantity += trolley[id]["quantity"];
        total += trolley[id]["amount"];
        if ($('#dd > ul').length == 0) {
            continue;
        }
        var row = '<li data-rel="' +id+ '"><b>' +trolley[id]["model"]+
            ':</b><span style="float:right">' +trolley[id]["quantity"]+
            ' x ' +addCommas(trolley[id]["price"].toString(), 0)+
            ' = ' +addCommas(trolley[id]["amount"].toString(), 0)+
            ' VND <span title="Bớt" class="removeQty text-red glyphicon glyphicon-minus"></span>&nbsp;&nbsp;<span title="Thêm" class="addQty text-blue glyphicon glyphicon-plus"></span></span></span></li>';
        trolleyContent += row;
    }

    if (trolleyContent == '') {
        trolleyContent = '<li>Không có sản phẩm nào trong Giỏ hàng</li>';
    }
    $('#proNum').html(quantity);
    $('#proAmount').html(addCommas(total.toString(), 0));
    if ($('#dd > ul').length > 0) {
        $('#dd > ul').html(trolleyContent);
    }

}

//Jquery code
$(document).ready(function() {
    $(document).on('blur', 'input.digit', function(event) {
        var str = addCommas($(this).val(), false);
        if (str != '') {
            $(this).val(str);
        }
    });

    $(document).on('click', '.cancelIcon', function () {
        $(this).siblings("input").val("");
    });

    $('#register').on('click', function () {
        callAjax('user', 'account', {'silent':true,'typeForm': 0}, accountCallback);
    });

    $('#account').on('click', function () {
        callAjax('user', 'account', {'silent':true,'typeForm': 1}, accountCallback);
    });

    $('#forgot_password').on('click', function () {
        callAjax('user', 'forgotpassword', {'silent':true}, forgotpasswordCallback);
    });

    $('#btnOrder').on('click', function () {
        //Check if trolley empty
        if (sessionStorage.trolley == undefined || sessionStorage.trolley == "{}") {
            alert('Giỏ hàng của bạn đang trống.');
            return false;
        }
        window.location.href = root_url+"transactions";
    });

    $('.select-text').on('click', function () {
        var el = $(this)[0];
        var range = document.createRange();
        range.selectNodeContents(el);
        var sel = window.getSelection();
        sel.removeAllRanges();
        sel.addRange(range);
    });

    $('#btnSearch').on('click', function() {
        //Validate
        if ($("#keyword").val() == '') {
            showErrorBubble("btnSearch", " Vui lòng nhập từ khóa tìm kiếm.");
            return false;
        }
        $(this).closest("form").submit();
    });

    $(".zoomer").elevateZoom({
        zoomWindowPosition: 6,
        zoomWindowWidth: 300,
        zoomWindowHeight: 300,
        tint:true,
        tintColour:'#78787A',
        tintOpacity:0.5
    });
});

function accountCallback(result) {
    if (result !== false) {
        $("#common_dialog").modal({show: true, keyboard: true, backdrop: 'static'});
        if ($('#idRecord').val() == 0) {
            var sitekey = "6LfnIDkUAAAAAP1Br-z8W5_gSqer4t_vwho3T8Wf";
            //var sitekey = "6LclTCAUAAAAALxvwNvIpveB9e09_vmY81iw84sZ";//Local Server
            grecaptcha.render('g-recaptcha', {'sitekey': sitekey});
        }
    }

    $('#submitAccount').on('click', function () {
        var formData = fetchForm($(this).closest('form'), true);
        callAjax('user', 'account_update', {upload: true, formData: formData}, accountUpdateCallback);
    });

    return false;
}

function accountUpdateCallback(result) {
    if (result === false) {
        if ($('idRecord').val() == 0) {
            grecaptcha.reset();
        }
        return false;
    }
    window.location.href = root_url;
}

function forgotpasswordCallback(result) {
    if (result !== false) {
        $("#common_dialog").modal({show: true, keyboard: true, backdrop: 'static'});
    }

    $('#submitForgotPassword').on('click', function () {
        var formData = fetchForm($(this).closest('form'));
        callAjax('user', 'forgotpassword_process', {formData: formData}, forgotpasswordProcessCallback);
    })
}

function forgotpasswordProcessCallback(result) {
    if (result !== false) {
        window.location.href = root_url;
    }
}